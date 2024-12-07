package main

import (
	"fmt"
	"strconv"
	"strings"
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

func AddMonthToPeerExpiration(peer PeerGorm) error {
	db := OpenDB()
	var resPeer PeerGorm
	db.Where("id = ?", peer.ID).First(&resPeer)
	if resPeer.ID == 0 {
		return fmt.Errorf("Failed to find peer with id %d in database", peer.ID)
	}
	resPeer.ExpirationTime = resPeer.ExpirationTime.AddDate(0, 1, 0)
	resPeer.Status = "Paid"
	resPeer.AllowedIP = strings.ReplaceAll(resPeer.AllowedIP, "0.0.0.", "10.0.0.")
	err := setPeer(resPeer)
	if err != nil {
		return fmt.Errorf("Failed to set Peer %s new info", peer.AllowedIP)
	}
	db.Save(&resPeer)
	return nil
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
	peer.Name = fmt.Sprintf("%d-%s-%s", peer.ID, cons.Username, time.Now().Format("2006-01-02-15-04-05"))
	lg.Println(peer.Name)
	cons.PeerID = uint32(peer.ID)
	db.Save(&cons)
	db.Save(&peer)
	return cons, peer, nil
}

func GiveLastPaidPeerFromORM(cons ConsGorm) (ConsGorm, PeerGorm, error) {
	db := OpenDB()
	var resCons ConsGorm
	var resPeer PeerGorm
	db.Where("username = ?", cons.Username).Last(&resCons)
	if resCons.PeerID == 0 {
		err := fmt.Errorf("Failed to find consumer with username %s in database", cons.Username)
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

func FindChatIDsByPeerIDs(PeerIDs uint) (int64, error) {
	db := OpenDB()
	var resCons ConsGorm
	db.Where("peer_id = ?", uint32(PeerIDs)).First(&resCons)
	lg.Println(resCons)
	if resCons.ChatID == "" {
		return int64(0), fmt.Errorf("Consumer with peer_id %d was not found in database!", PeerIDs)
	}
	intChatID, err := strconv.Atoi(resCons.ChatID)
	if err != nil {
		return int64(0), fmt.Errorf("Failed to convert string %s to int", resCons.ChatID)
	}
	return int64(intChatID), nil
}

func FindPeers(peerIDs []int) ([]PeerGorm, error) {
	db := OpenDB()
	var resPeers []PeerGorm
	if len(peerIDs) > 0 {
		db.Find(&resPeers, peerIDs)
		if len(resPeers) == 0 {
			err := fmt.Errorf("Peers: %d were not found in database!", peerIDs)
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

func changePeerStatusInORM(peer PeerGorm, status string) error {
	db := OpenDB()
	var resPeer PeerGorm
	db.Where("public_key = ?", peer.PublicKey).First(&resPeer)
	if resPeer.PublicKey == "" {
		return fmt.Errorf("Failed to find peer with public_key %s", peer.PublicKey)
	}
	resPeer.Status = status
	db.Save(&resPeer)
	return nil
}

func CheckExpiration() error {
	db := OpenDB()
	var resPeers []PeerGorm
	db.Find(&resPeers, "status != ?", "Virgin")
	if len(resPeers) == 0 {
		return fmt.Errorf("Failed to find peers without status 'Virgin'")
	}

	for _, peer := range resPeers {
		days := time.Since(peer.ExpirationTime).Hours() / 24
		switch {
		case days > float64(-1) && days < float64(0) && peer.Status == "Paid":
			lg.Printf("PRE_EXPIRED %s DAYS: %f", peer.AllowedIP, days)
			ChatID, err := FindChatIDsByPeerIDs(peer.ID)
			if err != nil {
				return fmt.Errorf("Failed to find ChatID of Peer %d: %s", peer.ID, err)
			}
			err = changePeerStatusInORM(peer, "Pre_Expired")
			if err != nil {
				return fmt.Errorf("Failed to change Peer status %s  status: %s", peer.Name, err)
			}
			msg := fmt.Sprintf(escapeMarkdownV2(preExpiredMsg), peer.Name)
			go sendMessage(ChatID, msg)

		case days >= float64(0) && days < float64(6) && peer.Status == "Pre_Expired":
			lg.Printf("EXPIRED %s DAYS: %f", peer.AllowedIP, days)
			err := RestrictPeer(peer)
			if err != nil {
				return fmt.Errorf("Failed to restrict Peer %d: %s", peer.ID, err)
			}
			ChatID, err := FindChatIDsByPeerIDs(peer.ID)
			if err != nil {
				return fmt.Errorf("Failed to find ChatID of Peer %d: %s", peer.ID, err)
			}
			msg := fmt.Sprintf(escapeMarkdownV2(expiredMsg), peer.Name)
			go sendMessage(ChatID, msg)

		case days >= float64(6) && days < float64(7) && peer.Status == "Expired":
			lg.Printf("PRE_DEAD %s DAYS: %f", peer.AllowedIP, days)
			ChatID, err := FindChatIDsByPeerIDs(peer.ID)
			if err != nil {
				return fmt.Errorf("Failed to find ChatID of Peer %d: %s", peer.ID, err)
			}
			err = changePeerStatusInORM(peer, "Pre_Dead")
			if err != nil {
				return fmt.Errorf("Failed to change Peer %s  status: %s", peer.Name, err)
			}
			msg := fmt.Sprintf(escapeMarkdownV2(preDeadMsg), peer.Name)
			go sendMessage(ChatID, msg)

		case days >= float64(7) && days < float64(8) && peer.Status == "Pre_Dead":
			lg.Printf("KILLING %s DAYS: %f", peer.AllowedIP, days)
			err := KillAndRegenPeer(peer)
			if err != nil {
				return fmt.Errorf("Failed to kill peer %d regen new: %s", peer.ID, err)
			}
			ChatID, err := FindChatIDsByPeerIDs(peer.ID)
			if err != nil {
				return fmt.Errorf("Failed to find ChatIDs of Peers %d: %s", peer.ID, err)
			}
			msg := fmt.Sprintf(escapeMarkdownV2(deadMsg), peer.Name)
			go sendMessage(ChatID, msg)

		}
	}

	return nil
}

func RestictPeerInORM(peer PeerGorm) error {
	db := OpenDB()
	var resPeer PeerGorm
	db.Find(&resPeer, "public_key = ?", peer.PublicKey)
	if resPeer.ID == 0 {
		return fmt.Errorf("Failed to find peer %s", peer.PublicKey)
	}
	resPeer.AllowedIP = peer.AllowedIP
	resPeer.Status = "Expired"
	db.Save(&resPeer)
	return nil
}

func KillAndRegenPeerInORM(oldPeer PeerGorm) (PeerGorm, error) {
	db := OpenDB()
	db.Find(&oldPeer, "public_key = ?", oldPeer.PublicKey)
	if oldPeer.ID == 0 {
		return PeerGorm{}, fmt.Errorf("Failed to find peer %s", oldPeer.PublicKey)
	}
	oldPeer = RegenOnePeer(oldPeer)

	db.Save(&oldPeer)
	return oldPeer, nil
}
