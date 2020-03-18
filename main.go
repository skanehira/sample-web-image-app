package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const imageDir = "images"

type File struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Content string `json:"content"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	outFile, err := os.Create(filepath.Join(imageDir, header.Filename))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fs := []File{}
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(imageDir, file.Name()))
		if err != nil {
			continue
		}

		f := File{
			Name:    file.Name(),
			Size:    file.Size(),
			Content: base64.StdEncoding.EncodeToString(b),
		}

		fs = append(fs, f)
	}

	if err := json.NewEncoder(w).Encode(&fs); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getFiles(w, r)
		case http.MethodPost:
			uploadFile(w, r)
		}

	})
	log.Println("start http server :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
