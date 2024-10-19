package main

func main(){
  initDB()
  ReadWGCreds()
  GenAndWritePeers()

  // err := CreateDbs()
  // if err != nil{
  //   lg.Fatalf("Failed to create tables in psql: %s", err)
  // }
  // lg.Println("Tables are created successfully!")

  // var testCons ConsGorm
  //
  // testCons.ChatID = "146146146"
  // testCons.Username = "@egrmk"
  // grantConsumerPeer(testCons)
  initServer()
}
