package login

import(
	"os"
	"testing"
	
	"strings"
	"context"
	"net/url"
	"net/http"
	"io/ioutil"
)

var mainCtx context.Context
var login *Login
var token string

func TestMain(m *testing.M) {
	//Context
	var cancel func()
	mainCtx, cancel = context.WithCancel(context.Background())
	defer cancel()
	
	//Run Tests
	exit := m.Run()
	
	//Exit code
	os.Exit(exit)
}

func TestLoad(t *testing.T) {
	login, err := NewLogin(mainCtx)
	if err != nil {
		t.Error("Could not connect to database")
		t.Fail()
	}
	
	login.Serve("8080")
}

func TestLogin(t *testing.T) {
	data := url.Values{}
    data.Set("email", "riley@example.com")
    data.Set("pass", "potato")
    
	req, err := http.NewRequestWithContext(mainCtx, http.MethodPost, "http://localhost:8080/login", strings.NewReader(data.Encode()))
	if err != nil {
		t.Error("loginHandler", "Could not create request")
		t.Fail()
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		t.Error("loginHandler", "Could not complete request")
		t.Fail()
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		t.Error("loginHandler", "Could not rad body")
		t.Fail()
	}
	
	if string(byt) == "Invalid Login Credentials" {
		t.Error("loginHandler", "Login information somehow wrong, check database")
	}
	
	token = string(byt)
}

func TestCheck(t *testing.T) {
	req, err := http.NewRequestWithContext(mainCtx, http.MethodGet, "http://localhost:8080/check", nil)
	if err != nil {
		t.Error("loginHandler", "Could not create request")
		t.Fail()
	}
	req.Header.Add("Authorization", "BEARER:"+token)
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		t.Error("checkHandler", "Could not complete request")
		t.Fail()
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		t.Error("checkHandler", "Could not read body")
		t.Fail()
	}
	
	if string(byt) != "OK" {
		t.Error("checkHandler", "Token value somehow wrong, serious issue")
	}
}
