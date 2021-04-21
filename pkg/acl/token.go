package acl

import (
	"github.com/dgrijalva/jwt-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
)

type Claims struct {
	jwt.StandardClaims
	Topics    []Topic    `json:"topics"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Topic struct {
	Type       topics.Type `json:"type"`
	AccessType AccessType  `json:"access_type"`
}

type Endpoint struct {
	Name string `json:"name"`
}

type SuperuserClaims struct {
	jwt.StandardClaims
	IsSuperuser bool `json:"is_superuser"`
}
