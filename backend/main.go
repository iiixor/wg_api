package main

import (
	"time"
)

func main() {
	loadEnv()
	go initBot()
	time.Sleep(time.Second * 2)
	go StartExpirationChecker(24 * time.Hour)
	initServer()
}
