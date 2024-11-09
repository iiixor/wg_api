package main

func initDB() {
	err := CreateDbs()
	if err != nil {
		lg.Println(err)
	}
	lg.Println("Tables are created successfully!")
}
