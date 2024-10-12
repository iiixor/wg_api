package main

func initWgConf(){
  var inter InterfaceGorm
  inter.Name = "wg0"
  inter.PublicKey = readfile("/etc/wireguard/serverPublicKey")
  inter.PrivateKey = readfile("/etc/wireguard/serverPrivateKey")
  lg.Println("/etc/wireguard/ was read successfully!")
  AddInterfaceToORM(inter)
}

func createClientConfig(inter InterfaceGorm, peer PeerGorm) WgConfig{
  var cfg WgConfig
  cfg.Interface.Address = peer.AllowedIP
  cfg.Interface.PublicKey = inter.PublicKey
  cfg.Interface.MTU = MTU
  cfg.Interface.DNS = DNS
  cfg.Peer.PublicKey = peer.PublicKey
  cfg.Peer.AllowedIPs = "0.0.0.0/0"
  cfg.Peer.Endpoint = Endpoint
  cfg.Peer.PersistentKeepalive = "21"
  return cfg
}

func writePeersToWgConf(peers []PeerGorm){
  
}
