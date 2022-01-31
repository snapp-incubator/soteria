package user

// Issuer indicate issuers.
type Issuer string

const (
	Driver    Issuer = "0"
	Passenger        = "1"
	None             = "-1"
)
