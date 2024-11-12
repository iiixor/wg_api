package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// _ "github.com/lib/pq"
)

type ConsGorm struct {
	gorm.Model
	ChatID   string
	Username string
	PeerID   uint32
}

type PeerGorm struct {
	gorm.Model
	Name            string
	InterfaceID     uint32
	PrivateKey      string
	PublicKey       string
	PresharedKey    string
	AllowedIP       string
	Status          string // Выдан или нет
	LatestHandshake time.Time
	ExpirationTime  time.Time
}

type InterfaceGorm struct {
	gorm.Model
	Name       string
	PrivateKey string
	PublicKey  string
}

func initDB() {
	err := CreateDbs()
	if err != nil {
		lg.Println(err)
	}
	lg.Println("Tables are created successfully!")
}

func CreateDbs() error {
	db := OpenDB()
	err := db.AutoMigrate(&ConsGorm{}, &PeerGorm{}, &InterfaceGorm{})
	if err != nil {
		lg.Printf("Failed to migrate schema %s", err)
		return err
	}
	return nil
}

func OpenDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		lg.Printf("Failed to open db %s: %s", dbName, err)
		return nil
	}
	return db
}

func AddConsumerToORM(consumer ConsGorm) error {
	db := OpenDB()
	db.Create(&consumer)
	lg.Printf("%s was successfully added to %s", consumer.Username, db.Name())
	return nil
}

func AddPeerToORM(peer PeerGorm) error {
	db := OpenDB()
	db.Create(&peer)
	lg.Printf("%s was successfully added to %s", peer.Name, db.Name())
	return nil
}

func AddInterfaceToORM(inter InterfaceGorm) error {
	db := OpenDB()
	db.Create(&inter)
	lg.Printf("%s was successfully added to %s", inter.Name, db.Name())
	return nil
}

func DeletePeerFromORM(peer PeerGorm) error {
	db := OpenDB()
	db.Delete(&peer, "name = ?", peer.Name)
	lg.Printf("%s was deleted successfully from %s", peer.Name, db.Name())
	return nil
}

func GetConsumerInfoDB(consumer ConsGorm) (ConsGorm, error) {
	db := OpenDB()
	var res ConsGorm
	db.Where("chat_id=?", consumer.ChatID).Find(&res)
	return res, nil
}

func GetVacantPeerFromORM() (PeerGorm, error) {
	var vacantPeer PeerGorm
	db := OpenDB()
	db.Where("status <> ?", "Paid").First(&vacantPeer)
	vacantPeer.Status = "Paid"
	vacantPeer.ExpirationTime = time.Now().AddDate(0, 1, 0)
	db.Save(&vacantPeer)
	return vacantPeer, nil
}

func GetInterfaceInfoFromORM() (InterfaceGorm, error) {
	var inter InterfaceGorm
	db := OpenDB()
	db.Last(&inter)
	if inter.ID == 0 {
		err := fmt.Errorf("Failed to find Interface %s in the database", inter.Name)
		return InterfaceGorm{}, err
	}
	return inter, nil
}

func writePeersToORM(peers []PeerGorm) error {
	db := OpenDB()
	db.Create(peers)
	return nil
}

func grantConsumerPeerInORM(cons ConsGorm, peer PeerGorm) (ConsGorm, PeerGorm, error) {
	db := OpenDB()
	cons.PeerID = uint32(peer.ID)
	db.Save(&cons)
	return cons, peer, nil
}

func GiveLastPaidPeerFromORM(cons ConsGorm) (ConsGorm, PeerGorm, error) {
	db := OpenDB()
	var resCons ConsGorm
	var resPeer PeerGorm
	db.Where("chat_id = ?", cons.ChatID).Last(&resCons)
	if resCons.PeerID == 0 {
		err := fmt.Errorf("Failed to find consumer with ChatID %s in database", cons.ChatID)
		return ConsGorm{}, PeerGorm{}, err
	}
	db.Where("id = ?", resCons.PeerID).First(&resPeer)
	if resPeer.ID == 0 {
		err := fmt.Errorf("Failed to find peer with PeerID %d in database", resCons.PeerID)
		return ConsGorm{}, PeerGorm{}, err
	}
	return resCons, resPeer, nil
}

func FindCons(cons ConsGorm) ([]ConsGorm, error) {
	db := OpenDB()
	var resCons []ConsGorm
	db.Find(&resCons, "username = ?", cons.Username)
	if len(resCons) == 0 {
		err := fmt.Errorf("Consumer %s was not found in database!", cons.Username)
		return []ConsGorm{}, err
	}
	return resCons, nil
}

func FindPeers(peerIDs []int) ([]PeerGorm, error) {
	db := OpenDB()
	var resPeers []PeerGorm
	if len(peerIDs) > 0 {
		db.Find(&resPeers, peerIDs)
		if len(resPeers) == 0 {
			err := fmt.Errorf("Peers: %s were not found in database!", peerIDs)
			return []PeerGorm{}, err
		}
		return resPeers, nil
	}
	return []PeerGorm{}, nil
}

func getTunnelList(cons ConsGorm) ([]PeerGorm, error) {
	foundedCons, err := FindCons(cons)
	if err != nil {
		err = fmt.Errorf("Failed to find consumer: %s", err)
		return []PeerGorm{}, err
	}
	peerIDs := MakePeerIDArray(foundedCons)
	foundedPeers, err := FindPeers(peerIDs)
	if err != nil {
		err = fmt.Errorf("Failed to find peer: %s", err)
		return []PeerGorm{}, err
	}
	return foundedPeers, nil
}
