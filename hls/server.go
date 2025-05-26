package main

import(
	"io"
	"os"
	"log"
	"strings"
	"context"
	"net/http"
	"path/filepath"
	
	"github.com/julienschmidt/httprouter"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/RileySun/Scaled/utils"
)

type VideoServer struct {
	router *httprouter.Router
}

//Create
func NewVideoServer() *VideoServer {
	server := &VideoServer{
		router:httprouter.New(),
	}
	
	//Routes
	server.router.POST("/upload", server.Upload)
	
	//Serve Files
	server.router.ServeFiles("/stream/*filepath", http.Dir("stream"))
	
	return server
}

//Server
func (s *VideoServer) Start(ctx context.Context) {
	utils.StartHTTPServer(ctx, "8080", s.router)
}

func (s *VideoServer) Stop(cancel func()) {
	//Cleanup
	
	//Graceful Shutdown
	cancel()
}

//Actions
func (s *VideoServer) Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Set max size
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512)
	
	//Get Reader
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	for {
		part, partErr := reader.NextPart()
		if partErr != nil {
			if partErr == io.EOF {
				//Done
				return
			} else {
				//Other Error
				log.Println(partErr)
				http.Error(w, "Error uploading file", http.StatusBadRequest)
				continue
			}
		}
		
		//Is an uploaded file
		if part.FormName() != "files" {
			continue
		}
		
		//Filename
		fullfile := strings.ReplaceAll(part.FileName(), " ", "_")
		filename := strings.Split(fullfile, ".")[0]
		
		//Path
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			log.Println(wdErr)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			continue
		}
		path := wd + "/stream/" + filename
		fpath := filepath.Join(path, "/", fullfile)
		
		//Mkdir
		if dirErr := os.Mkdir(path, os.ModePerm); dirErr != nil {
			log.Println(dirErr)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			continue
		}
		
		//Create File
		f, fileErr := os.Create(fpath)
		if fileErr != nil {
			log.Println(fileErr)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			continue
		}
		
		//Read to file
		_, readErr := f.ReadFrom(part)
		if readErr != nil {
			log.Println(readErr)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			continue
		}
		
		//Convert to hls
		convErr := s.Convert(fpath, filename)
		if convErr != nil {
			//log.Println(convErr)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			continue
		}
		
		//Remove original file
		delErr := os.Remove(fpath)
		if convErr != nil {
			log.Println(delErr)
			http.Error(w, "Error converting file", http.StatusBadRequest)
			continue
		}
		
		serverPath := "http://localhost:8080/stream/" + filename + "/" + filename + ".m3u8"
		w.Write([]byte(serverPath))
	}
}

func (s *VideoServer) Convert(fpath string, filename string) error {	
	//Get Paths
	dir := filepath.Dir(fpath)
	outFile := filepath.Join(dir, filename+".m3u8")
	
	//Run ffmpeg
	output := ffmpeg.Input(fpath).Output(outFile, ffmpeg.KwArgs{"c:a":"copy","hls_time":"10","hls_list_size":"0","f":"hls",})
	return output.OverWriteOutput().Silent(true).Run()
}