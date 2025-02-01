package main

const (
	MTU       = "1342"
	DNS       = "1.1.1.1"
	AllowedIP = "0.0.0.0/0"
)

const (
	PubKeyPath     = "/etc/wireguard/serverPublicKey"
	PrivateKeyPath = "/etc/wireguard/serverPrivateKey"
)

var (
	token         string
	preExpiredMsg string
	expiredMsg    string
	preDeadMsg    string
	deadMsg       string
	DB_HOST       string
	DB_USER       string
	DB_PASSWORD   string
	DB_NAME       string
	WG_ENDPOINT   string
)
