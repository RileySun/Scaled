package login

import (
	"log"
	"time"
	"errors"
	"context"
	"net/http"
	"net/mail"
	"database/sql"
	
	"golang.org/x/crypto/bcrypt"
	
	"github.com/julienschmidt/httprouter"
	"github.com/golang-jwt/jwt"
	
	"github.com/RileySun/Scaled/utils"
)
type Login struct {
	DB *sql.DB
	ctx context.Context
	secretKey []byte
}

func NewLogin(ctx context.Context) (*Login, error) {
	login := &Login{
		ctx:ctx,
		secretKey:[]byte("super-secret"),
	}
	
	//Database connection
	creds := utils.LoadCredentials()
	var err error
	login.DB, err = utils.NewDB(creds.Host, creds.Port, creds.User, creds.Pass, creds.Database)
	if err != nil {
		return login, err
	}
	
	return login, nil
}

//Server Function
func (l *Login) Serve(address string) {
	//Routes
	router := httprouter.New()
	router.POST("/login", l.loginHandler)
	router.GET("/check", l.checkHandler)
	
	//Start
	go func() {utils.StartHTTPServer(l.ctx, "8080", router)}()
}

//Data Functions
func (l *Login) checkLogin(email, password string) error {
	var hashedPassword string
	row := l.DB.QueryRow("SELECT `password` FROM Users WHERE `email` = ?;", email)
	
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

func (l *Login) createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email":email,
			"expires":time.Now().Add(time.Hour * 6).Unix(),
		})
	
	tokenStr, signErr := token.SignedString(l.secretKey)
	if signErr != nil {
		return "", signErr
	}
	
	return tokenStr, nil
}

func (l *Login) checkToken(tokenStr string) error {
	//Create token refrence from token string
	token, verifyErr := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return l.secretKey, nil
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
func (l *Login) loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Parse Data
	r.ParseForm()
	email := r.PostFormValue("email")
	pass := r.PostFormValue("pass")
	
	//Validate (Saves having to check a database if it isnt even a valid email)
	_, emailErr := mail.ParseAddress(email)
	if emailErr != nil {
		log.Println(emailErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid Login Credentials"))
		return
	}
	
	//Check Login
	loginErr := l.checkLogin(email, pass)
	if loginErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(loginErr.Error()))
		return
	}
	
	//Create Token
	tokenStr, tokenErr := l.createToken(email)
	if tokenErr != nil {
		log.Println(tokenErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please Contact [ADMIN] About Token Issue"))
		return
	}
	
	//Return
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenStr))
}

func (l *Login) checkHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Bearer String
	bearerStr := r.Header.Get("Authorization")
	if bearerStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Access Token"))
		return
	}
	
	//Get Token
	tokenStr := bearerStr[len("Bearer "):]
	
	//Check Token
	checkErr := l.checkToken(tokenStr)
	if checkErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Access Token"))
		return
	}
	
	//All is well
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}