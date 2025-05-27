package utils

import(
	"os"
	"os/exec"
	"os/signal"
	
	"io"
	"fmt"
	"log"
	"net"
	"time"
	"bytes"
	"context"
	"strings"
	"strconv"
	"syscall"
	"net/http"
	"math/rand"
	"path/filepath"
	"mime/multipart"
	
	"golang.org/x/sync/errgroup"
	
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

func LoadCredentials(path string) *Credentials {
	envErr := godotenv.Load(path) //Always relative to all folders in project
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
func StartHTTPServer(ctx context.Context, port string, r http.Handler) {
	//Get Signal Context
	mainCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	//Server Config
	srv := &http.Server{
		Handler: r,
		Addr: ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	//Create Errgroup, Launch Server, && Setup Shutdown
	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return srv.ListenAndServe()
	}) //Launch Server
	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	}) //Graceful Shutdown

	//Wait for Errgroup
	if err := g.Wait(); err != nil {
		//fmt.Printf("Shutdown: %s \n", err)
	}
}

//Get Random Int (Min Inclusive->Max Exclusive)
func GetRandomInt(min int, max int) int {
    return min + rand.Intn(max-min)
}

//Format Duration for Uptime
func FormatUptime(dur time.Duration) string {
	dur = dur.Round(time.Minute)
	
	d := dur / (time.Hour * 24)
	dur -= d * (time.Hour * 24)
	
	h := dur / time.Hour
	dur -= h * time.Hour
    
	m := dur / time.Minute
	dur -= m * time.Minute
	
    return fmt.Sprintf("%01dd %01dh %01dm", d, h, m)
}

//Get Memory of running golang process
func GetMemory() string {
	//Get PID and use for ps
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	
	cmd := exec.Command("bash", "-c",  "ps -p " + pidStr +" -o rss=")
	
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return "Error"
	}
	
	formatted := strings.TrimSpace(string(out[:]))
	mem, _ := strconv.ParseFloat(formatted, 32)
	memStr := fmt.Sprintf("%.1f", mem/1000)
	
	return memStr
}

//Basic Server Data
type ServerData struct {
	Name string `json:"Name"`
	Memory string `json:"Memory"`
	Uptime string `json:"Uptime"`
}

//Upload File to URL
func UploadFile(ctx context.Context, url string, fpath string, formfield string) (error) {	
	filename := filepath.Base(fpath)
	b, w := multipartFormData(fpath, filename, formfield)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	
	client := &http.Client{}
	_, reqErr := client.Do(req)
	return reqErr
}

func multipartFormData(fpath string, filename string, formfield string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	var fw io.Writer
	
	w := multipart.NewWriter(&b)
	
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal("Error opening file", fpath, err)
	}
	
	fw, err = w.CreateFormFile(formfield, filename)
	if err != nil {
		log.Fatal("Error creating writer:", err)
	}
	
	_, err = io.Copy(fw, file)
	if err != nil {
		log.Fatal("Error copying file:", err)
	}
	w.Close()
	
	return b, w
}