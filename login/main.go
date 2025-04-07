package main

import (
	"os"
	"log"
	"time"
	"errors"
	"context"
	"syscall"
	"net/http"
	"net/mail"
	"os/signal"
	"database/sql"
	
	"golang.org/x/crypto/bcrypt"
	
	"github.com/julienschmidt/httprouter"
	"github.com/golang-jwt/jwt"
	
	//Internal Utils
	"github.com/RileySun/Scaled/utils"
)

var DB *sql.DB
var server *http.Server
var secretKey = []byte("super-secret")

//Data Functions
func checkLogin(email, password string) error {
	var hashedPassword string
	row := DB.QueryRow("SELECT `password` FROM Users WHERE `email` = ?;", email)
	
	scanErr := row.Scan(&hashedPassword)
	if scanErr != nil {
		return errors.New("Invalid Login Credentials") //Same invalid response each time, for security
	}
	
	//Check Password
	cryptErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if cryptErr != nil {
		return errors.New("Invalid Login Credentials") //Same invalid response each time, for security
	} else {
		return nil
	}
}

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email":email,
			"expires":time.Now().Add(time.Hour * 6).Unix(),
		})
	
	tokenStr, signErr := token.SignedString(secretKey)
	if signErr != nil {
		return "", signErr
	}
	
	return tokenStr, nil
}

func checkToken(tokenStr string) error {
	//Create token refrence from token string
	token, verifyErr := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if verifyErr != nil {
		return verifyErr
	}
	
	//Check if token is valid
	if !token.Valid {
		return errors.New("Invalid Access Token")
	} else {
		return nil
	}
}

//Route Handlers
func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	//Parse Data
	r.ParseForm()
	email := r.PostFormValue("email")
	pass := r.PostFormValue("pass")
	
	//Validate (Saves having to check a database if it isnt even a valid email)
	_, emailErr := mail.ParseAddress(email)
	if emailErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid Login Credentials"))
	}
	
	//Check Login
	loginErr := checkLogin(email, pass)
	if loginErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(loginErr.Error()))
	}
	
	//Create Token
	tokenStr, tokenErr := createToken(email)
	if tokenErr != nil {
		log.Println(tokenErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please Contact [ADMIN] About Token Issue"))
	}
	
	//Return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenStr))
}

func checkHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Bearer String
	bearerStr := r.Header.Get("Authorization")
	if bearerStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Access Token"))
	}
	
	//Get Token
	tokenStr := bearerStr[len("Bearer "):]
	
	//Check Token
	checkErr := checkToken(tokenStr)
	if checkErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Access Token"))
	}
	
	//All is well
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

//Run Server
func main() {
	//Create DB
	creds := utils.LoadCredentials()
	DB = utils.NewDB(creds.Host, creds.Port, creds.User, creds.Pass, creds.Database)

	//create http router
	router := httprouter.New()
	
	//set routes
	router.POST("/login", loginHandler) //Logins
	router.GET("/check", checkHandler) //Is logged in checker
	
	//Create server (can close gracefully with Shutdown())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	server = utils.StartHTTPServer(router, "8080")
	<-done
	
	//Context for shutting fown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		//Graceful shutdown functions here
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}