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
	r.Get("/GetConsumerInfo/?={Username}", GetConsumerInfoAPI)
	r.Get("/GetVacantPeer/", GetVacantPeerAPI)
	r.Get("/GrantPeerToConsumer/{Username}+{ChatID}", GrantPeerToConsumerAPI)
	r.Get("/Init/ReadWgCreds", ReadWgCredsAPI)
	r.Get("/Init/GenAndWritePeers", GenAndWritePeersAPI)
	r.Get("/GiveLastCfg/{Username}", GiveLastPaidPeerAPI)
	r.Get("/GetTunnelList/{Username}", GetTunnelListAPI)
	r.Get("/ExtendPeer/{PeerID}", ExtendPeerAPI)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resp.Message = fmt.Sprintln("Hello World!")
		DrawJSON(w, resp, 200)
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
	if consumer.Username == "" {
		resp.Message = "Consumer not found!"
		DrawJSON(w, resp, 422)
		return
	}
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to get consumer from database: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	DrawJSON(w, consumer, 200)
}

func GetVacantPeerAPI(w http.ResponseWriter, r *http.Request) {
	var vacantPeer PeerGorm
	vacantPeer, err := GetVacantPeerFromORM()
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to get vacant peer from database: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	if vacantPeer.AllowedIP == "" {
		resp.Message = fmt.Sprintln("No vacant peers")
		DrawJSON(w, resp, 422)
		return
	}
	DrawJSON(w, vacantPeer, 200)
}

func GrantPeerToConsumerAPI(w http.ResponseWriter, r *http.Request) {
	var consumer ConsGorm
	consumer.Username = chi.URLParam(r, "Username")
	consumer.ChatID = chi.URLParam(r, "ChatID")
	var resPeer PeerGorm
	_, resPeer, err := grantConsumerPeer(consumer)
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to grant peer to consumer: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	inter, err := GetInterfaceInfoFromORM()
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to get interface info from database :%s", err)
		DrawJSON(w, resp, 422)
		return
	}
	var clientCfg WgConfig
	clientCfg = createClientConfig(inter, resPeer)
	DrawJSON(w, clientCfg, 200)
}

func ReadWgCredsAPI(w http.ResponseWriter, r *http.Request) {
	err := ReadWGCreds()
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to read wg creds: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	resp.Message = fmt.Sprintln("Wireguard creds are successfully read!")
	DrawJSON(w, resp, 200)
}

func GenAndWritePeersAPI(w http.ResponseWriter, r *http.Request) {
	err := GenAndWritePeers()
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to generate and write peers: %s", err)
		DrawJSON(w, resp, 422)
	}
	resp.Message = fmt.Sprintln("Peers generated and written successfully!")
	DrawJSON(w, resp, 200)
}

func GiveLastPaidPeerAPI(w http.ResponseWriter, r *http.Request) {
	var consumer ConsGorm
	consumer.Username = chi.URLParam(r, "Username")
	var resPeer PeerGorm
	_, resPeer, err := GiveLastPaidPeer(consumer)
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to give last paid peer: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	inter, err := GetInterfaceInfoFromORM()
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to get interface info:%s", err)
		DrawJSON(w, resp, 422)
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
		resp.Message = fmt.Sprintf("Failed to get tunnel list: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	DrawJSON(w, foundedPeers, 200)
}

func ExtendPeerAPI(w http.ResponseWriter, r *http.Request) {
	var peer PeerGorm
	strID := chi.URLParam(r, "PeerID")
	uintID, err := strconv.ParseUint(strID, 10, 32)
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to convert str to uint: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	peer.ID = uint(uintID)
	err = AddMonthToPeerExpiration(peer)
	if err != nil {
		resp.Message = fmt.Sprintf("Failed to add month to peer expiration: %s", err)
		DrawJSON(w, resp, 422)
		return
	}
	resp.Message = fmt.Sprintf("1 Month added successfully to peer ID: %d", peer.ID)
	DrawJSON(w, resp, 200)
}
