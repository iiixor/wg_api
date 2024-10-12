package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// _ "github.com/lib/pq"
)

type ConsGorm struct{
  gorm.Model
  ChatID string
  Username string
  PeerID uint32
}

type PeerGorm struct{
  gorm.Model
  Name string
  InterfaceID uint32
  PrivateKey string
  PublicKey string
  PresharedKey string
  AllowedIP string
  Status string // Выдан или нет
  LatestHandshake time.Time
  ExpirationTime time.Time
}

type InterfaceGorm struct{
  gorm.Model
  Name string
  PrivateKey string
  PublicKey string
}

func CreateDbs() error {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  err = db.AutoMigrate(&ConsGorm{}, &PeerGorm{}, &InterfaceGorm{})
  if err != nil {
    lg.Printf("Failed to migrate schema %s", err)
    return err
  }
  return nil
}

func AddConsumerToORM(consumer ConsGorm) error {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  db.Create(&consumer)
  lg.Printf("%s was successfully added to %s", consumer.Username, db.Name())
  return nil
}

func AddPeerToORM(peer PeerGorm) error {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  db.Create(&peer)
  lg.Printf("%s was successfully added to %s", peer.Name, db.Name())
  return nil
}

func AddInterfaceToORM(inter InterfaceGorm) error{
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  db.Create(&inter)
  lg.Printf("%s was successfully added to %s", inter.Name, db.Name())
  return nil
}

func DeletePeerFromORM(peer PeerGorm) error {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  db.Delete(&peer, "name = ?", peer.Name)
  lg.Printf("%s was deleted successfully from %s", peer.Name, db.Name())
  return nil
}

func GetConsumerInfoDB(consumer ConsGorm) (ConsGorm, error) {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  var res ConsGorm
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return ConsGorm{},err
  }
  db.Where("chat_id=?", consumer.ChatID).Find(&res)
  return res,nil
}

func GetVacantPeerFromORM() (PeerGorm, error) {
  var vacantPeer PeerGorm
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return PeerGorm{}, err
  }
  db.Where("status <> ?", "Paid").First(&vacantPeer)
  vacantPeer.Status = "Paid"
  vacantPeer.ExpirationTime = time.Now().AddDate(0, 1, 0)
  db.Save(&vacantPeer)
  return vacantPeer,nil
}

func GetInterfaceInfoFromORM() (InterfaceGorm, error){
  var inter InterfaceGorm
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return InterfaceGorm{}, err
  }
  db.Last(&inter)
  return inter, nil
}

func writePeersToORM(peers []PeerGorm) error {
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return err
  }
  db.Create(peers)
  return nil
}

func grantConsumerPeerInORM(cons ConsGorm, peer PeerGorm) (ConsGorm,PeerGorm, error){
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return ConsGorm{},PeerGorm{}, err
  }
  cons.PeerID = uint32(peer.ID)
  db.Save(&cons)
  return cons, peer, nil
}

func GiveLastPaidPeerFromORM(cons ConsGorm) (ConsGorm,PeerGorm, error){
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil{
    lg.Printf("Failed to open db %s: %s", dbName, err)
    return ConsGorm{},PeerGorm{}, err
  }
  var resCons ConsGorm
  var resPeer PeerGorm
  db.Where("chat_id = ?", cons.ChatID).Last(&resCons)
  db.Where("id = ?", resCons.PeerID).First(&resPeer)
  return resCons, resPeer, nil
}
