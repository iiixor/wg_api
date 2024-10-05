package main

func main(){
  readPeersInfoJSON("./peers.json")
  lg.Println("All peers are read succesfully!")
  // for i:=0; i < 3;i++{
  //   AddUserPeers(145145145)
  // }
  vacantPeer := GetAvailablePeer()
  if vacantPeer.AllowedIP == ""{
    lg.Println("No vacant peers!")
    return
  }
  lg.Println(vacantPeer)

  // initDB()
  initServer()
}
