package main

import (
	"encoding/json"
	"os"
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
