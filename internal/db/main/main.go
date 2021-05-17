package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	password := "password"
	secret := "secret"

	var users []user.User
	driver := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "driver",
		Password:                string(password),
		Type:                    user.EMQUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.CabEvent,
				AccessType: acl.Sub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.DriverLocation,
				AccessType: acl.Pub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.SuperappEvent,
				AccessType: acl.Sub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.GossiperLocation,
				AccessType: acl.Sub,
			},
		},
	}

	passenger := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "passenger",
		Password:                string(password),
		Type:                    user.EMQUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.CabEvent,
				AccessType: acl.Sub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.SuperappEvent,
				AccessType: acl.Sub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.PassengerLocation,
				AccessType: acl.Pub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.GossiperLocation,
				AccessType: acl.Sub,
			},
		},
	}
	box := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "box",
		Password:                string(password),
		Type:                    user.EMQUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.BoxEvent,
				AccessType: acl.Sub,
			},
		},
	}

	colonySubscriber := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "colony-subscriber",
		Password:                string(password),
		Type:                    user.EMQUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Topic:      topics.DriverLocation,
				AccessType: acl.Sub,
			},
		},
	}

	snappbox := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "snapp-box",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
		},
	}

	snappfood := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "snapp-food",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
		},
	}

	snappmarket := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "snapp-market",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
		},
	}

	snappflight := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "snapp-flight",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
		},
	}

	snappdoctor := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "snapp-doctor",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
		},
	}

	gabriel := user.User{
		MetaData: db.MetaData{
			ModelName:    "user",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		},
		Username:                "gabriel",
		Password:                string(password),
		Type:                    user.HeraldUser,
		Secret:                  string(secret),
		TokenExpirationDuration: time.Hour * 30 * 24,
		Rules: []user.Rule{
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/notification",
				AccessType: acl.Pub,
			},
			user.Rule{
				UUID:       uuid.New(),
				Endpoint:   "/event",
				AccessType: acl.Pub,
			},
		},
	}

	users = append(users, driver)
	users = append(users, passenger)
	users = append(users, box)
	users = append(users, colonySubscriber)
	users = append(users, snappbox)
	users = append(users, snappfood)
	users = append(users, snappmarket)
	users = append(users, snappflight)
	users = append(users, snappdoctor)
	users = append(users, gabriel)

	d, _ := json.Marshal(users)
	ioutil.WriteFile("db.json", d, 0644)
}

func getPublicKey(u user.Issuer) (*rsa.PublicKey, error) {
	var fileName string
	d, _ := os.Getwd()
	fmt.Println(d)
	switch u {
	case user.Passenger:
		fileName = "test/1.pem"
	case user.Driver:
		fileName = "test/0.pem"
	case user.ThirdParty:
		fileName = "test/100.pem"
	default:
		return nil, fmt.Errorf("invalid user, public key not found")
	}
	pem, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
