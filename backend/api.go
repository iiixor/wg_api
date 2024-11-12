package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func initServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/GetConsumerInfo/?={ChatID}", GetConsumerInfoAPI)
	r.Get("/GetVacantPeer/", GetVacantPeerAPI)
	r.Get("/GrantPeerToConsumer/{Username}", GrantPeerToConsumerAPI)
	r.Get("/Init/ReadWgCreds", ReadWgCredsAPI)
	r.Get("/Init/GenAndWritePeers", GenAndWritePeersAPI)
	r.Get("/GiveLastCfg/{ChatID}", GiveLastPaidPeerAPI)
	r.Get("/GetTunnelList/{Username}", GetTunnelListAPI)
	r.Get("/ExtendPeer/{PeerID}", ExtendPeerAPI)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":3000", r)
}

func DrawJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

func GetConsumerInfoAPI(w http.ResponseWriter, r *http.Request) {
	var consumer ConsGorm
	consumer.ChatID = chi.URLParam(r, "ChatID")
	consumer, err := GetConsumerInfoDB(consumer)
	// DrawJSON(w, consumer, 200)
	if consumer.Username == "" {
		w.WriteHeader(404)
		w.Write([]byte("Consumer not found!"))
		// lg.Printf("Consumer %s not found!", consumer.ChatID)
		return
	}
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
		return
	}
	DrawJSON(w, consumer, 200)
}

func GetVacantPeerAPI(w http.ResponseWriter, r *http.Request) {
	var vacantPeer PeerGorm
	vacantPeer, err := GetVacantPeerFromORM()
	if err != nil {
		w.WriteHeader(422)
		lg.Printf("Failed to get vacant peer %s", err)
		return
	}
	if vacantPeer.AllowedIP == "" {
		w.WriteHeader(404)
		lg.Println("No vacant peers")
		return
	}
	DrawJSON(w, vacantPeer, 200)
	lg.Printf("Vacant peer: %s", vacantPeer.Name)
}

func GrantPeerToConsumerAPI(w http.ResponseWriter, r *http.Request) {
	var consumer ConsGorm
	consumer.Username = chi.URLParam(r, "Username")
	// var resCons ConsGorm
	var resPeer PeerGorm
	_, resPeer, err := grantConsumerPeer(consumer)
	if err != nil {
		lg.Println("Failed to grant peer to consumer!")
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
		return
	}
	inter, err := GetInterfaceInfoFromORM()
	if err != nil {
		lg.Printf("Failed to get interface info:%s", err)
		return
	}
	var clientCfg WgConfig
	clientCfg = createClientConfig(inter, resPeer)

	DrawJSON(w, clientCfg, 200)
	// w.WriteHeader(200)
	// w.Write([]byte(fmt.Sprintf("UserID: %s was granted peer %d", consumer.ChatID, peerID)))
}

func ReadWgCredsAPI(w http.ResponseWriter, r *http.Request) {
	err := ReadWGCreds()
	if err != nil {
		lg.Printf("Failed to read wg creds: %s", err)
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Wireguard creds are successfully read!"))
}

func GenAndWritePeersAPI(w http.ResponseWriter, r *http.Request) {
	err := GenAndWritePeers()
	if err != nil {
		lg.Printf("Failed to generate and write peers: %s", err)
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
	w.Write([]byte("Peers generated and written successfully!"))
}

func GiveLastPaidPeerAPI(w http.ResponseWriter, r *http.Request) {
	var consumer ConsGorm
	consumer.ChatID = chi.URLParam(r, "ChatID")
	var resPeer PeerGorm
	_, resPeer, err := GiveLastPaidPeer(consumer)
	if err != nil {
		msg := fmt.Sprintf("Failed to give last paid peer: ", err)
		lg.Println(msg)
		w.WriteHeader(422)
		w.Write([]byte(msg))
		return
	}
	inter, err := GetInterfaceInfoFromORM()
	if err != nil {
		msg := fmt.Sprintf("Failed to get interface info:%s", err)
		lg.Printf(msg)
		w.WriteHeader(422)
		w.Write([]byte(msg))
		return
	}
	var clientCfg WgConfig
	clientCfg.FileName = resPeer.Name
	clientCfg = createClientConfig(inter, resPeer)
	DrawJSON(w, clientCfg, 200)
}

func GetTunnelListAPI(w http.ResponseWriter, r *http.Request) {
	var cons ConsGorm
	cons.Username = chi.URLParam(r, "Username")
	foundedPeers, err := getTunnelList(cons)
	if err != nil {
		lg.Printf("Failed to get tunnel list: %s", err)
		w.WriteHeader(422)
		w.Write([]byte(fmt.Sprintf("Failed to get tunnel list: %s", err)))
		return
	}
	DrawJSON(w, foundedPeers, 200)
}

func ExtendPeerAPI(w http.ResponseWriter, r *http.Request) {
	var peer PeerGorm
	strID := chi.URLParam(r, "PeerID")
	uintID, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Failed to convert str to uint: %s", err)
		w.WriteHeader(422)
		w.Write([]byte(msg))
		return
	}
	peer.ID = uint(uintID)
	err = AddMonthToPeerExpiration(peer)
	if err != nil {
		msg := fmt.Sprintf("Failed to add month to peer expiration: %s", err)
		w.WriteHeader(422)
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(200)
	msg := fmt.Sprintf("1 Month added successfully to peer ID: %d", peer.ID)
	w.Write([]byte(msg))
}
