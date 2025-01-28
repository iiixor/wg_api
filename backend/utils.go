package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/curve25519"
)

func setLogger() {
	// Открываем файл для записи логов
	logFile, err := os.OpenFile("wg_api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		lg.Fatalf("Ошибка при открытии файла: %v", err)
	}
	// Настраиваем логер на запись в файл
	lg.SetOutput(logFile)
	lgWG.SetOutput(logFile)
	lgORM.SetOutput(logFile)
	lgError.SetOutput(logFile)
	lgAPI.SetOutput(logFile)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	token = os.Getenv("BOT_TOKEN")
	preExpiredMsg = os.Getenv("PRE_EXPIRED_MSG")
	expiredMsg = os.Getenv("EXPIRED_MSG")
	preDeadMsg = os.Getenv("PRE_DEAD_MSG")
	deadMsg = os.Getenv("DEAD_MSG")

	setLogger()
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
	}
	return peersArray
}

func RegenOnePeer(oldPeer PeerGorm) PeerGorm {
	var newPeer PeerGorm
	privateKey, publicKey, err := generateKeys()
	if err != nil {
		lg.Printf("Failed to generate keys: %s", err)
	}
	newPeer = oldPeer
	newPeer.Name = publicKey
	newPeer.PublicKey = publicKey
	newPeer.PrivateKey = privateKey
	newPeer.AllowedIP = turnOnPeer(oldPeer.AllowedIP)
	newPeer.Status = "Virgin"
	newPeer.InterfaceID = 1
	lg.Printf("New peer: %s pubKey: %s allowed_ip %s was generated! Old name: %s old allowed_ip: %s", newPeer.Name, newPeer.PublicKey, newPeer.AllowedIP, oldPeer.Name, oldPeer.AllowedIP)
	return newPeer
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

func Expired(expDate time.Time) bool {
	currentDate := time.Now()
	if currentDate.After(expDate) {
		return true
	}
	return false
}

func StartExpirationChecker(interval time.Duration) {
	for {
		lg.Println("Started to check expiration...")
		err := CheckExpiration()
		if err != nil {
			lg.Printf("Failed to check expiration: %s", err)
		}
		// Ожидание до следующей проверки
		lg.Printf("Next checking will be in %s...\n", interval)
		time.Sleep(interval)
	}
}

func escapeMarkdownV2(text string) string {
	// Полный список специальных символов для MarkdownV2
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`",
		">", "#", "+", "-", "=", "|", "{", "}",
		".", "!",
	}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func turnOffPeer(peer string) string {
	if strings.HasPrefix(peer, "10.") {
		return strings.Replace(peer, "10.", "0.", 1)
	}
	return peer
}

func turnOnPeer(peer string) string {
	if strings.HasPrefix(peer, "0.") {
		return strings.Replace(peer, "0.", "10.", 1)
	}
	return peer
}
