package acl

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.StandardClaims
	Topics    []Topic    `json:"topics"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Topic struct {
	Type       string     `json:"type"`
	AccessType AccessType `json:"access_type"`
}

type Endpoint struct {
	Name string `json:"name"`
}

type SuperuserClaims struct {
	jwt.StandardClaims
	IsSuperuser bool `json:"is_superuser"`
}
