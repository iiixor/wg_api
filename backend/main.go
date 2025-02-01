package main

import (
	"time"
)

func main() {
	loadEnv()
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	// DB, err := dbConstruct(dsn)
	// if err != nil {
	// 	lgORM.Fatalf("Failed to construct Db %v", err)
	// }

	// if err := DB.CreateTestLines(); err != nil {
	// 	lgORM.Fatalf("Failed to create test line: %v", err)
	// }

	go initBot()
	time.Sleep(time.Second * 2)
	go StartExpirationChecker(24 * time.Hour)
	initServer()
}
