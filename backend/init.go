package main

import(
  "time"
)

func initDB(){
  err := CreateDbs()
  if err != nil{
    lg.Println(err)
  }
  lg.Println("Tables are created successfully!")
  var cons ConsGorm
  cons.ChatID = "146146148"
  cons.Username = "@egrmk"
  AddConsumerToORM(cons)
  var peeer PeerGorm
  peeer.Name = "Egr_kali"
  peeer.AllowedIP = "10.0.0.2/32"
  peeer.PrivateKey = "UB4+uUtrfbhtIOZJo+gh88QcAyOL+y8rngQzK7i6kEY="
  peeer.PublicKey = "J2x2ka2YtDnSFPVXe2ze3sz5/tsbiFcPXEjSMOBOEn4="
  peeer.LatestHandshake = time.Now()
  peeer.ExpirationTime = time.Now()
  AddPeerToORM(peeer)
  peeer.Name = "Mitay_Mac"
  peeer.AllowedIP = "10.0.0.3/32"
  AddPeerToORM(peeer)
  // DeletePeerFromORM(peeer)
}

