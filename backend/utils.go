package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/curve25519"
)

func readPeersInfoJSON(filePath string){
	data, err := os.ReadFile(filePath)
	if err != nil {
		lg.Fatalf("Failed to open JSON file: %v", err)
	}

  err = json.Unmarshal(data, &response)
  if err != nil{
    lg.Printf("Failed to unmarshal %s: %s", filePath, err)
  }
}

func GetAvailablePeer() Peer {
  for _, peer := range response.Data.ConfigurationPeers{
    if peer.LatestHandshake == "No Handshake"{
      // Отметить, что пир выдан
      return peer
    }
  }
  return Peer{}
}

func AddUserPeers(userID uint64){
  var consumer Consumer
  consumer.ChatID = userID
  consumer.Peers = make([]Peer, 6)
  peer := GetAvailablePeer()
  consumer.Peers = append(consumer.Peers, peer)
  lg.Println(consumer)
  // Добавить пира пользователю в бд
}

func DeletePeer(userID uint64, peerName string){
  var consumer Consumer
  consumer.ChatID = userID
  // прочитать пира пользователя из бд
  for _, peer := range consumer.Peers{
    if peer.Name == peerName{
      peer = Peer{}
      break
    }
  }
  lg.Printf("Peer %s belonged to user %s was deleted!", peerName, userID)
}
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

func readfile(filePath string) string{
	data, err := os.ReadFile(filePath)
	if err != nil {
		lg.Fatalf("Failed to open %s: %v", filePath, err)
	}
  return string(data)
}

func generatePeers() []PeerGorm{
  peersArray := make([]PeerGorm, 254)
  for i:=0;i<254;i++{
    privateKey, publicKey, err := generateKeys()
    if err != nil{
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

