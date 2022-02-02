package user

// Issuer indicate issuers.
type Issuer string

const (
	Driver    Issuer = "0"
	Passenger Issuer = "1"
	None      Issuer = "-1"
)
