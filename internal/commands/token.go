package commands

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db/redis"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/emq"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

var (
	ErrInvalidToken = errors.New("invalid token type")
	ErrNoToken      = errors.New("must be at least one token type")
)

type ErrInvalidInput struct {
	Message string
}

func (err ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input. %s", err.Message)
}

var Token = &cobra.Command{
	Use:   "token",
	Short: "token issues token based on type of user",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return ErrNoToken
		}

		t := args[0]
		if t == "herald" || t == "superuser" {
			return nil
		}

		return ErrInvalidInput{Message: fmt.Sprintf("token with type %s is not valid", t)}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.InitConfig()

		pk100, err := config.ReadPrivateKey(user.ThirdParty)
		if err != nil {
			return fmt.Errorf("cannot load private keys %w", err)
		}

		// nolint: exhaustivestruct
		app.GetInstance().SetAuthenticator(&authenticator.Authenticator{
			PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
				user.ThirdParty: pk100,
			},
		})

		switch args[0] {
		case "herald":
			return heraldToken(cmd, args)
		case "superuser":
			return superuserToken(config, cmd, args)
		default:
			return ErrInvalidToken
		}
	},
}

func heraldToken(cmd *cobra.Command, args []string) error {
	r := bufio.NewScanner(os.Stdin)

	cmd.Print("username? >>> ")
	r.Scan()
	username := r.Text()

	cmd.Print("duration? (in hours) >>> ")
	r.Scan()
	duration, err := time.ParseDuration(fmt.Sprintf("%sh", r.Text()))
	if err != nil {
		return fmt.Errorf("cannot parse duration %w", err)
	}

	var (
		allowedTopics    []acl.Topic
		allowedEndpoints []acl.Endpoint
	)

	heraldEndpoints := map[string]string{
		"1": "/event",
		"2": "/events",
		"3": "/notification",
	}
	topicTypes := map[string]topics.Type{
		"1": topics.CabEvent,
		"2": topics.BoxEvent,
		"3": topics.SuperappEvent,
		"4": topics.DriverLocation,
		"5": topics.PassengerLocation,
		"6": topics.SharedLocation,
		"7": topics.Chat,
	}
	accessTypes := map[string]acl.AccessType{
		"1": acl.Sub,
		"2": acl.Pub,
		"3": acl.PubSub,
	}

	for {
		cmd.Println("do you want to give permissions? [y/n]")
		cmd.Print(">>> ")
		r.Scan()

		switch r.Text() {
		case "n":
			token, err := app.GetInstance().Authenticator.HeraldToken(username, allowedEndpoints, allowedTopics, duration)
			if err != nil {
				return fmt.Errorf("token creation failed %w", err)
			}

			cmd.Println(token)

			return nil
		case "y":
			cmd.Println("which one do you want to grant access to? \n\t1. topic\n\t2. endpoint")
			cmd.Print(">>> ")
			r.Scan()

			switch r.Text() {
			case "1":
				cmd.Println("which topic do you want to grant access to?\n" +
					"\t1. Snapp Cab Events" +
					"\t2. Snapp Box Events" +
					"\t3. Snapp Super App Events" +
					"\t4. Snapp Driver Locations")
				cmd.Print(">>> ")
				r.Scan()
				topicType, ok := topicTypes[r.Text()]
				if !ok {
					return ErrInvalidInput{Message: "selected topic does not exist"}
				}

				cmd.Println("which access type do you want to grant?\n" +
					"\t1. Subscribe" +
					"\t2. Publish" +
					"\t3. Publish-Subscribe")
				cmd.Print(">>> ")
				r.Scan()
				at, ok := accessTypes[r.Text()]
				if !ok {
					return ErrInvalidInput{Message: "invalid input. selected access type does not exist"}
				}
				topic := acl.Topic{
					Type:       topicType,
					AccessType: at,
				}
				allowedTopics = append(allowedTopics, topic)
			case "2":
				cmd.Println("which endpoint do you want to grant access to?\n" +
					"\t1. event" +
					"\t2. events" +
					"\t3. notification")
				cmd.Print(">>> ")
				r.Scan()
				endpoint, ok := heraldEndpoints[r.Text()]
				if !ok {
					return ErrInvalidInput{Message: "invalid input. selected endpoint does not exist"}
				}
				e := acl.Endpoint{Name: endpoint}
				allowedEndpoints = append(allowedEndpoints, e)
			default:
				return ErrInvalidInput{Message: ""}
			}
		default:
			return ErrInvalidInput{Message: "you should enter y or n"}
		}
	}
}

func superuserToken(config config.AppConfig, cmd *cobra.Command, _ []string) error {
	cli, err := redis.NewRedisClient(config.Redis)
	if err != nil {
		return fmt.Errorf("redis connection failed %w", err)
	}

	app.GetInstance().SetEMQStore(emq.Store{
		Client: cli,
	})

	r := bufio.NewScanner(os.Stdin)

	cmd.Print("username? >>> ")
	r.Scan()
	username := r.Text()

	cmd.Print("duration? (in hours) >>> ")
	r.Scan()
	duration, err := time.ParseDuration(fmt.Sprintf("%sh", r.Text()))
	if err != nil {
		return fmt.Errorf("cannot parse duration %w", err)
	}

	cmd.Print("password? >>> ")
	r.Scan()
	password := r.Text()

	token, err := app.GetInstance().Authenticator.SuperuserToken(username, duration)
	if err != nil {
		return fmt.Errorf("token creation failed %w", err)
	}

	if err := app.GetInstance().EMQStore.Save(cmd.Context(), emq.User{
		Username:    token,
		Password:    password,
		IsSuperuser: true,
	}); err != nil {
		return fmt.Errorf("cannot save token to redis %w", err)
	}

	cmd.Println(token)

	return nil
}
