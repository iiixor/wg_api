package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/curve25519"
)

func generateKeys() (string, string, error) {
	// Создаем 32 байта для приватного ключа
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return "", "", err
	}

	// Генерация публичного ключа на основе приватного
	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	// Преобразуем в base64
	privKeyStr := base64.StdEncoding.EncodeToString(privateKey[:])
	pubKeyStr := base64.StdEncoding.EncodeToString(publicKey[:])

	return privKeyStr, pubKeyStr, nil
}

func readfile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		lg.Fatalf("Failed to open %s: %v", filePath, err)
	}
	return string(data)
}

func generatePeers() []PeerGorm {
	peersArray := make([]PeerGorm, 254)
	for i := 0; i < 254; i++ {
		privateKey, publicKey, err := generateKeys()
		if err != nil {
			lg.Printf("Failed to generate keys: %s", err)
		}
		peersArray[i].Name = publicKey
		peersArray[i].PublicKey = publicKey
		peersArray[i].PrivateKey = privateKey
		peersArray[i].AllowedIP = fmt.Sprintf("10.0.0.%d/32", i+2)
		peersArray[i].Status = "Virgin"
		peersArray[i].InterfaceID = 1
		// lg.Println(peersArray[i])
		// lg.Printf("Name:%s", peersArray[i].Name)
		// lg.Printf("PubKey:%s", peersArray[i].PublicKey)
		// lg.Printf("PrivateKey:%s", peersArray[i].PrivateKey)
		// lg.Printf("AllowedIP:%s", peersArray[i].AllowedIP)
		// lg.Printf("Status:%s", peersArray[i].Status)
		// lg.Printf("InterfaceID:%d\n\n", peersArray[i].InterfaceID)
		// lg.Printf("Name: %s\nPublicKey:%s\nPrivateKey:%s\nAllowedIP:%s\nStatus:%s\nInterfaceID:%s", pe)
		// lg.Println(peersArray[i])
	}
	return peersArray
}

func MakePeerIDArray(cons []ConsGorm) []int {
	res := make([]int, len(cons), len(cons))
	for i, con := range cons {
		res[i] = int(con.PeerID)
	}
	return res
}

func AddMonthToExpire(currentTime time.Time) time.Time {
	return currentTime.AddDate(0, 1, 0)
}
