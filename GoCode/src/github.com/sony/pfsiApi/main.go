package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
)

type imagesHash1 struct {
	ID string `json:"ID" `
}
type imagesHash struct {
	ID   string `json:"ID" `
	Name string `json:"Name" `
}

var imagesHashes []imagesHash

var cid, err string

var tmpl = template.Must(template.New("tmpl").ParseFiles("ipfs.html"))

func addData(w http.ResponseWriter, r *http.Request) {
	sh := shell.NewShell("localhost:5001")

	cid, err := sh.Add(strings.NewReader("Hello World!!"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Println("added ", cid)
	imagesHashes = append(imagesHashes, imagesHash{ID: "1", Name: cid})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imagesHashes)
}
func readData(w http.ResponseWriter, r *http.Request) {
	sh := shell.NewShell("localhost:5001")
	params := mux.Vars(r)
	err := sh.Get(params["id"], "C:/GoCode/src/github.com/sony/pfsiApi/storeFile")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Println("Saved ")

}

func addDir(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addDir")

	sh := shell.NewShell("localhost:5001")

	cid, err := sh.AddDir("C:/GoCode/src/github.com/sony/pfsiApi/file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	fmt.Println(cid)
	HomePageVars := imagesHash1{ID: cid}
	t, err := template.ParseFiles("hash.html") //parse the html file homepage.html
	if err != nil {                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func deleteData(w http.ResponseWriter, r *http.Request) {
	sh := shell.NewShell("localhost:5001")
	params := mux.Vars(r)

	err := sh.Unpin(params["id"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	fmt.Println("hello world12")
	cmd := exec.Command("ipfs repo gc")
	errd := cmd.Run()
	if errd != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errd)
	}

	fmt.Println("hello worlddd1")
}

func main() {

	fmt.Println("hello world")
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "ipfs.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.HandleFunc("/readData/{id}", readData).Methods("GET")
	r.HandleFunc("/addData", addData).Methods("POST")
	r.HandleFunc("/addDir", addDir).Methods("POST")
	r.HandleFunc("/deleteData/{id}", deleteData).Methods("DELETE")
	fmt.Println("ListenAndServe")
	log.Fatal(http.ListenAndServe(":8080", r))

}
