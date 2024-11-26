package main

import (
	"log"
	"os"
)

type Response struct {
	Message string
}

var resp Response

type PeerCfg struct {
	PublicKey           string
	AllowedIPs          string
	Endpoint            string
	PersistentKeepalive string
}

type InterfaceCfg struct {
	PrivateKey string
	Address    string
	MTU        string
	DNS        string
}

type WgConfig struct {
	FileName  string
	Interface InterfaceCfg
	Peer      PeerCfg
}

var response Response

type Message struct {
	ChatID  int64
	Content string
}

var messageChannel chan Message

var lg *log.Logger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lmicroseconds)
