package main

func main(){
  readPeersInfoJSON("./peers.json")
  lg.Println("All peers are read succesfully!")
  initDB()
  // for i:=0; i < 3;i++{
  //   AddUserPeers(145145145)
  // }
  // vacantPeer := GetAvailablePeer()
  // if vacantPeer.AllowedIP == ""{
  //   lg.Println("No vacant peers!")
  //   return
  // }
  // lg.Println(vacantPeer)

  // initDB()
  //
  privateKey, publicKey,_ := generateKeys()
  lg.Printf("%s\n%s", privateKey, publicKey)
  // initWgConf()
  initServer()
}
