package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReadWGCreds() error {
	var inter InterfaceGorm
	inter.Name = "wg0"
	inter.PublicKey = readfile(PubKeyPath)
	inter.PrivateKey = readfile(PrivateKeyPath)
	inter.PrivateKey = strings.TrimSpace(inter.PrivateKey)
	inter.PublicKey = strings.TrimSpace(inter.PublicKey)
	lg.Println("/etc/wireguard/ was read successfully!")
	err := AddInterfaceToORM(inter)
	if err != nil {
		err = fmt.Errorf("Failed to add iterface to orm: %s", err)
		return err
	}
	return nil
}

func createClientConfig(inter InterfaceGorm, peer PeerGorm) WgConfig {
	var cfg WgConfig
	cfg.FileName = peer.Name
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
	for _, peer := range peers {
		cmd := exec.Command("wg", "set", "wg0", "peer", peer.PublicKey, "allowed-ips", peer.AllowedIP)
		// Запускаем команду и возвращаем ошибку, если она произошла
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error executing command set: %v", err)
		}
		lg.Println(peer.AllowedIP)
		lg.Println(peer.PublicKey)
	}
	cmd := exec.Command("wg-quick", "save", "wg0")
	// Запускаем команду и возвращаем ошибку, если она произошла
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command save: %v", err)
	}
	lg.Println("wg-quick")
	return nil
}

func setPeer(peer PeerGorm) error {
	cmd := exec.Command("wg", "set", "wg0", "peer", peer.PublicKey, "allowed-ips", peer.AllowedIP)
	// Запускаем команду и возвращаем ошибку, если она произошла
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command set: %v", err)
	}
	lg.Println(peer.AllowedIP)
	lg.Println(peer.PublicKey)
	cmd = exec.Command("wg-quick", "save", "wg0")
	// Запускаем команду и возвращаем ошибку, если она произошла
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command save: %v", err)
	}
	lg.Println("wg-quick")
	return nil
}

func writePeersIntoWgConf(filePath string, peers []PeerGorm) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		lg.Printf("Failed to open %s:%s", filePath, err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, peer := range peers {
		str := fmt.Sprintf("[Peer]\nPublicKey = %s\nAllowedIPs = %s\n\n", peer.PublicKey, peer.AllowedIP)
		_, err := writer.Write([]byte(str))
		if err != nil {
			lg.Printf("Failed to writed data to %s:%s", filePath, err)
			return err
		}
	}
	return nil
}

func grantConsumerPeer(cons ConsGorm, month, days int) (ConsGorm, PeerGorm, error) {
	var vacantPeer PeerGorm
	vacantPeer, err := GetVacantPeerFromORM(month, days)
	if err != nil {
		return ConsGorm{}, PeerGorm{}, fmt.Errorf("Failed to get vacant peer from database: %s", err)
	}
	var resCons ConsGorm
	var resPeer PeerGorm
	resCons, resPeer, err = grantConsumerPeerInORM(cons, vacantPeer)
	if err != nil {
		err = fmt.Errorf("Failed to grant peer to consumer in database: %s", err)
		return ConsGorm{}, PeerGorm{}, err
	}
	return resCons, resPeer, nil

}

func GenAndWritePeers() error {
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

func GiveLastPaidPeer(cons ConsGorm) (ConsGorm, PeerGorm, error) {
	var resCons ConsGorm
	var resPeer PeerGorm
	resCons, resPeer, err := GiveLastPaidPeerFromORM(cons)
	if err != nil {
		err = fmt.Errorf("Failed to give last paid peer of user %s from ORM : %s", resCons.Username, err)
		return ConsGorm{}, PeerGorm{}, err
	}
	return resCons, resPeer, nil
}

func RestrictPeer(peer PeerGorm) error {
	lg.Println(peer.AllowedIP)
	peer.AllowedIP = strings.ReplaceAll(peer.AllowedIP, "10.0.0.", "0.0.0.")
	peer.Status = "Expired"
	lg.Println(peer.AllowedIP)
	cmd := exec.Command("wg", "set", "wg0", "peer", peer.PublicKey, "allowed-ips", peer.AllowedIP)
	// Запускаем команду и возвращаем ошибку, если она произошла
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command set: %v", err)
	}
	cmd = exec.Command("wg-quick", "save", "wg0")
	// Запускаем команду и возвращаем ошибку, если она произошла
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command save: %v", err)
	}
	if err := RestictPeerInORM(peer); err != nil {
		return fmt.Errorf("Failed to restrict peer %d in ORM: %s", peer.ID, err)
	}
	lg.Printf("Peer was restircted succefully %s", peer.PublicKey)
	return nil
}

func KillAndRegenPeer(oldPeer PeerGorm) error {
	newPeer, err := KillAndRegenPeerInORM(oldPeer)
	if err != nil {
		return fmt.Errorf("Failed to kill and regen peer in ORM %s: %s", oldPeer.PublicKey, err)
	}

	cmd := exec.Command("wg", "set", "wg0", "peer", newPeer.PublicKey, "allowed-ips", newPeer.AllowedIP)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command set: %v", err)
	}

	cmd = exec.Command("wg", "set", "wg0", "peer", oldPeer.PublicKey, "remove")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command set: %v", err)
	}

	cmd = exec.Command("wg-quick", "save", "wg0")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing command save: %v", err)
	}

	lg.Printf("Peer %s was killed and regened to %s", oldPeer.PublicKey, newPeer.PublicKey)
	return nil
}
