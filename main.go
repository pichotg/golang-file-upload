package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const indexPage = "/index.html"

//UploadFile upload files on server :
func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(10 << 20) // maximum upload 10 MB file size

	file, handler, err := r.FormFile("file_name")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("tmp", "upload-*.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	fmt.Fprintf(w, fmt.Sprintf("Successfully Uploaded File \n name :%s", tempFile.Name()))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/upload", UploadFile)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
	fmt.Println("Server start on port :8080")
}
