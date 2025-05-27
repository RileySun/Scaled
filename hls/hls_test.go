package main

import(
	"os"
	"time"
	"errors"
	"context"
	"testing"
	
	"github.com/RileySun/Scaled/utils"
)


var mainCtx context.Context
var cancel func()
var server *VideoServer

func TestMain(m *testing.M) {	
	exit := m.Run()
	cancel()
	os.Exit(exit)
}

func TestLoad(t *testing.T) {
	mainCtx, cancel = context.WithCancel(context.Background())
	
	server = NewVideoServer()
	go func() {server.Start(mainCtx)}()
}

func TestUpload(t *testing.T) {
	ctx, cancel := context.WithTimeout(mainCtx, time.Duration(time.Second * 10))
	defer cancel()
	
	err := utils.UploadFile(ctx, "http://localhost:8080/upload", "RAW/Example.mp4", "files")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestExists(t *testing.T) {
	if _, err := os.Stat("RAW/Example.mp4"); errors.Is(err, os.ErrNotExist) {
		t.Error("Test File Does Not Exist")
		t.Fail()
	}
}

func TestCleanup(t *testing.T) {
	os.RemoveAll("./stream/")
	os.Mkdir("./stream", 0777)
}