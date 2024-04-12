package main

import (
	"boardcamp/cmd/api"
	"boardcamp/internal/api/platform/database"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

////
/////

func init() {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	//err := godotenv.Load(".env")

	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}

}

func main() {
	// INIT DB
	db := database.Init()
	if db != nil {
		fmt.Println("Conexão com o banco de dados estabelecida com sucesso.")
	} else {
		fmt.Println("Falha ao estabelecer conexão com o banco de dados.")
	}

	// INIT API
	app := api.InitMicroService(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Listening server at %s\n", port)
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
