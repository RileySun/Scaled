package utils

import(
	"os"
	"log"
	"time"
	"net/http"
	"math/rand"
	
	"github.com/joho/godotenv"
)

//Init
func init() {
	rand.Seed(time.Now().UnixNano())
}

//Database credentials
type Credentials struct {
	User string `json:"user"`			//Database Username
	Pass string `json:"pass"`			//Database Pass
	Host string `json:"host"`			//Database Host
	Port string `json:"port"`			//Database Port
	Database string `json:"database"`	//Database Table
}

func LoadCredentials() *Credentials {
	envErr := godotenv.Load("../utils/.env") //Always relative to all folders in project
	if envErr != nil {
		log.Println("Utils: Error loading .env file - ", envErr)
		log.Println("This may be caused by running in docker")
	}
	
	creds := &Credentials{
		User:os.Getenv("DB_USER"),
		Pass:os.Getenv("DB_PASS"),
		Host:os.Getenv("DB_HOST"),
		Port:os.Getenv("DB_PORT"),
		Database:os.Getenv("DB_DATABASE"),
	}
	
	return creds
}

func LoadEnv() {
	envErr := godotenv.Load("./.env")
	if envErr != nil {
		log.Println("Utils: Error loading .env file - ", envErr)
		log.Println("This may be caused by running in docker")
	}
}

//HTTP Server
func StartHTTPServer(r http.Handler, port string) *http.Server {
	srv := &http.Server{
		Handler: r,
		Addr: ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			//prob intentional close
		} else {
			log.Printf("Httpserver: ListenAndServe() closing...")
		}
	}()

	//returned for Shutdown()
	return srv
}

//Get Random Int (Min Inclusive->Max Exclusive)
func GetRandomInt(min int, max int) int {
    return min + rand.Intn(max-min)
}