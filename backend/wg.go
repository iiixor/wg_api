package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
  "os/exec"
)


func ReadWGCreds() error {
  var inter InterfaceGorm
  inter.Name = "wg0"
  inter.PublicKey = readfile("/etc/wireguard/serverPublicKey")
  inter.PrivateKey = readfile("/etc/wireguard/serverPrivateKey")
  inter.PrivateKey = strings.TrimSpace(inter.PrivateKey)
  inter.PublicKey = strings.TrimSpace(inter.PublicKey)
  lg.Println("/etc/wireguard/ was read successfully!")
  err := AddInterfaceToORM(inter)
  if err != nil{
    lg.Printf("Failed to add iterface to orm: %s", err)
    return err
  }
  return nil
}

func createClientConfig(inter InterfaceGorm, peer PeerGorm) WgConfig{
  var cfg WgConfig
  cfg.Interface.Address = peer.AllowedIP
  cfg.Interface.PrivateKey = peer.PrivateKey
  cfg.Interface.MTU = MTU
  cfg.Interface.DNS = DNS
  cfg.Peer.PublicKey = inter.PublicKey
  cfg.Peer.AllowedIPs = "0.0.0.0/0"
  cfg.Peer.Endpoint = Endpoint
  cfg.Peer.PersistentKeepalive = "21"
  return cfg
}

func setPeers(peers []PeerGorm) error {
  for _,peer := range peers {
    cmd := exec.Command("wg", "set", "wg0", "peer", peer.PublicKey, "allowed-ips", peer.AllowedIP)

    // Запускаем команду и возвращаем ошибку, если она произошла
    if err := cmd.Run(); err != nil {
      return fmt.Errorf("error executing command: %v", err)
    }
  }
  return nil
}

func writePeersIntoWgConf(filePath string, peers []PeerGorm) error {
  file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil{
    lg.Printf("Failed to open %s:%s", filePath, err)
    return err
  }
  defer file.Close()

  writer := bufio.NewWriter(file)
  defer writer.Flush()

  for _,peer := range peers{
    str := fmt.Sprintf("[Peer]\nPublicKey = %s\nAllowedIPs = %s\n\n", peer.PublicKey, peer.AllowedIP)
    _, err := writer.Write([]byte(str))
    if err != nil{
      lg.Printf("Failed to writed data to %s:%s", filePath, err)
      return err
    }
  }
  return nil
}
func grantConsumerPeer(cons ConsGorm) (ConsGorm,PeerGorm, error){
  var vacantPeer PeerGorm
  vacantPeer, err := GetVacantPeerFromORM()
  if err != nil{
    lg.Printf("Failed to get vacant peer: %s", err)
    return ConsGorm{},PeerGorm{}, err
  }
  var resCons ConsGorm
  var resPeer PeerGorm
  resCons,resPeer, err = grantConsumerPeerInORM(cons, vacantPeer)
  if err != nil{
    lg.Printf("Failed to grant vacant peer to the comsumer: %s", err)
    return ConsGorm{},PeerGorm{}, err
  }
  return resCons,resPeer, nil
}

func GenAndWritePeers() error{
  peers := generatePeers()
  err := setPeers(peers)
  if err != nil {
    lg.Printf("Failed to write peers into wg conf: %s", err)
    return err
  }
  err = writePeersToORM(peers)
  if err != nil {
    lg.Printf("Failed to write peers into ORM: %s", err)
    return err
  }
  return nil
}


func GiveLastPaidPeer(cons ConsGorm) (ConsGorm,PeerGorm, error){
  var resCons ConsGorm
  var resPeer PeerGorm
  resCons, resPeer , err := GiveLastPaidPeerFromORM(cons)
  if err != nil{
    lg.Printf("Failed to find last given peer to the comsumer(%s): %s", resCons.ChatID, err)
    return ConsGorm{},PeerGorm{}, err
  }
  return resCons,resPeer, nil
}
