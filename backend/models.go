package main

import(
  "log"
  "os"
)

type DataUsage struct {
    Receive float64 `json:"Receive"`
    Sent    float64 `json:"Sent"`
    Total   float64 `json:"Total"`
}

type ConfigurationInfo struct {
    Address      string    `json:"Address"`
    ConnectedPeers int      `json:"ConnectedPeers"`
    DataUsage    DataUsage `json:"DataUsage"`
    ListenPort   string    `json:"ListenPort"`
    Name         string    `json:"Name"`
    PostDown     string    `json:"PostDown"`
    PostUp       string    `json:"PostUp"`
    PreDown      string    `json:"PreDown"`
    PreUp        string    `json:"PreUp"`
    PrivateKey   string    `json:"PrivateKey"`
    PublicKey    string    `json:"PublicKey"`
    SaveConfig   bool      `json:"SaveConfig"`
    Status       bool      `json:"Status"`
}

type Job struct {
    Action        string  `json:"Action"`
    Configuration string  `json:"Configuration"`
    CreationDate  string  `json:"CreationDate"`
    ExpireDate    *string `json:"ExpireDate"`
    Field         string  `json:"Field"`
    JobID         string  `json:"JobID"`
    Operator      string  `json:"Operator"`
    Peer          string  `json:"Peer"`
    Value         string  `json:"Value"`
}



type Peer struct {
    DNS                 string          `json:"DNS"`
    ShareLink           []interface{}   `json:"ShareLink"`
    AllowedIP           string          `json:"allowed_ip"`
    Configuration       ConfigurationInfo `json:"configuration"`
    CumuData            float64         `json:"cumu_data"`
    CumuReceive         float64         `json:"cumu_receive"`
    CumuSent            float64         `json:"cumu_sent"`
    Endpoint            string          `json:"endpoint"`
    EndpointAllowedIP   string          `json:"endpoint_allowed_ip"`
    ID                  string          `json:"id"`
    Jobs                []Job           `json:"jobs"`
    Keepalive           int             `json:"keepalive"`
    LatestHandshake     string          `json:"latest_handshake"`
    MTU                 int             `json:"mtu"`
    Name                string          `json:"name"`
    PresharedKey        string          `json:"preshared_key"`
    PrivateKey          string          `json:"private_key"`
    RemoteEndpoint      string          `json:"remote_endpoint"`
    Status              string          `json:"status"`
    TotalData           float64         `json:"total_data"`
    TotalReceive        float64         `json:"total_receive"`
    TotalSent           float64         `json:"total_sent"`
}

type ConfigurationData struct {
    ConfigurationInfo         ConfigurationInfo `json:"configurationInfo"`
    ConfigurationPeers        []Peer            `json:"configurationPeers"`
    ConfigurationRestrictedPeers []Peer          `json:"configurationRestrictedPeers"`
}

type Response struct {
    Data    ConfigurationData `json:"data"`
    Message *string           `json:"message"`
    Status  bool              `json:"status"`
}

type Consumer struct{
  ChatID uint64
  Peers []Peer
}

type PeerCfg struct{
  PublicKey string
  AllowedIPs string
  Endpoint string
  PersistentKeepalive string
}

type InterfaceCfg struct{
  PrivateKey string
  Address string
  MTU string
  DNS string
}

type WgConfig struct{
  Interface InterfaceCfg
  Peer PeerCfg
}

var response Response

var lg *log.Logger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lmicroseconds)
