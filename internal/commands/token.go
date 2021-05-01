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
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
)

var Token = &cobra.Command{
	Use:   "token",
	Short: "token issues token based on type of user",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you should set the token type")
		}

		t := args[0]
		if t == "herald" || t == "superuser" {
			return nil
		}

		return fmt.Errorf("token with type %s is not valid", t)
	},
	PreRunE: tokenPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "herald":
			return heraldToken(cmd, args)
		case "superuser":
			return superuserToken(cmd, args)
		default:
			panic("invalid token type")
		}
		return nil
	},
}

func tokenPreRun(cmd *cobra.Command, args []string) error {
	config := config.InitConfig()

	pk100, err := config.ReadPrivateKey(user.ThirdParty)
	if err != nil {
		return err
	}

	app.GetInstance().SetAuthenticator(&authenticator.Authenticator{
		PrivateKeys: map[user.Issuer]*rsa.PrivateKey{
			user.ThirdParty: pk100,
		},
	})

	return nil
}

func heraldToken(cmd *cobra.Command, args []string) error {
	r := bufio.NewScanner(os.Stdin)

	fmt.Print("username? >>> ")
	r.Scan()
	username := r.Text()

	fmt.Print("duration? (in hours) >>> ")
	r.Scan()
	duration, err := time.ParseDuration(fmt.Sprintf("%sh", r.Text()))
	if err != nil {
		return err
	}
	var allowedTopics []acl.Topic
	var allowedEndpoints []acl.Endpoint

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
	}
	accessTypes := map[string]acl.AccessType{
		"1": acl.Sub,
		"2": acl.Pub,
		"3": acl.PubSub,
	}
	for {
		fmt.Println("do you want to give permissions? [y/n]")
		fmt.Print(">>> ")
		r.Scan()
		in := r.Text()
		switch in {
		case "n":
			token, err := app.GetInstance().Authenticator.HeraldToken(username, allowedEndpoints, allowedTopics, duration)
			if err != nil {
				return err
			}
			fmt.Println(token)
			return nil
		case "y":
			fmt.Println("which one do you want to grant access to? \n\t1. topic\n\t2. endpoint")
			fmt.Print(">>> ")
			r.Scan()
			switch r.Text() {
			case "1":
				fmt.Println("which topic do you want to grant access to?\n" +
					"\t1. Snapp Cab Events" +
					"\t2. Snapp Box Events" +
					"\t3. Snapp Super App Events" +
					"\t4. Snapp Driver Locations")
				fmt.Print(">>> ")
				r.Scan()
				topicType, ok := topicTypes[r.Text()]
				if !ok {
					return fmt.Errorf("invaid input. selected topic does not exist")
				}
				fmt.Println("which access type do you want to grant?\n" +
					"\t1. Subscribe" +
					"\t2. Publish" +
					"\t3. Publish-Subscribe")
				fmt.Print(">>> ")
				r.Scan()
				at, ok := accessTypes[r.Text()]
				if !ok {
					return fmt.Errorf("invalid input. selected access type does not exist")
				}
				topic := acl.Topic{
					Type:       topicType,
					AccessType: at,
				}
				allowedTopics = append(allowedTopics, topic)
			case "2":
				fmt.Println("which endpoint do you want to grant access to?\n" +
					"\t1. event" +
					"\t2. events" +
					"\t3. notification")
				fmt.Print(">>> ")
				r.Scan()
				endpoint, ok := heraldEndpoints[r.Text()]
				if !ok {
					return fmt.Errorf("invalid input. selected endpoint does not exist")
				}
				e := acl.Endpoint{Name: endpoint}
				allowedEndpoints = append(allowedEndpoints, e)
			default:
				return fmt.Errorf("invalid input")
			}
		default:
			return fmt.Errorf("invalid answer. you should enter y or n")
		}
	}
}

func superuserToken(cmd *cobra.Command, args []string) error {
	r := bufio.NewScanner(os.Stdin)

	fmt.Print("username? >>> ")
	r.Scan()
	username := r.Text()

	fmt.Print("duration? (in hours) >>> ")
	r.Scan()
	duration, err := time.ParseDuration(fmt.Sprintf("%sh", r.Text()))
	if err != nil {
		return err
	}

	token, err := app.GetInstance().Authenticator.SuperuserToken(username, duration)
	if err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}
