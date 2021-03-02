package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

//Here is our waitgroup
var wg sync.WaitGroup

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

//Define function maps
var funcMap = template.FuncMap{
	"upperCase": strings.ToUpper, //upperCase is a key we can call inside of the template html file
}

type AUser struct { //Using this for Mongo
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	UserID      int    `json:"UserID"`
	Email       string `json:"Email"`
	PhoneACode  int    `json:"PhoneACode"`
	PhoneNumber int    `json:"PhoneNumber"`
	PostsMade   int    `json:"PostsMade"`
	RepliesMade int    `json:"RepliesMade"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//Parse our templates
func init() {
	/* Assign blank value to map so no nil errors occur */
	usernameMap = make(map[string]bool)                                //Clear all Usernames when loading so no problems are caused
	loadedMessagesMapHDog = make(map[int]Message)                      //Clearing this so we don't have any issues
	loadedMessagesMapHam = make(map[int]Message)                       //Clearing this so we don't have any issues
	theMessageBoardHDog.AllMessagesMap = make(map[int]Message)         //Clearing this so we don't have any issues
	theMessageBoardHDog.AllOriginalMessagesMap = make(map[int]Message) //Clearing this so we don't have any issues
	theMessageBoardHam.AllMessagesMap = make(map[int]Message)          //Clearing this so we don't have any issues
	theMessageBoardHam.AllOriginalMessagesMap = make(map[int]Message)  //Clearing this so we don't have any issues
	getbadWords()                                                      //Fill in bad words from file
	getAPICallVariables()                                              //Get all the API variables
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
	createTestMessages()
}

//Writes to the log; called from most anywhere in this program!
func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "logging.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	//Write to logger that we are handling requests
	debugMessage := "\n\nDEBUG: We are now handling requests"
	//fmt.Println(debugMessage)
	logWriter(debugMessage)
	fmt.Println(debugMessage)
	//Favicon and page spots
	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", index)
	myRouter.HandleFunc("/test", test)
	myRouter.HandleFunc("/hotdogMB", hotdogMB)
	myRouter.HandleFunc("/hamburgerMB", hamburgerMB)
	//Message update stuff
	myRouter.HandleFunc("/evaluateTenResults", evaluateTenResults).Methods("POST")
	myRouter.HandleFunc("/messageOriginalAjax", messageOriginalAjax).Methods("POST")
	myRouter.HandleFunc("/messageReplyAjax", messageReplyAjax).Methods("POST")
	//Field validation/User Creation
	myRouter.HandleFunc("/checkUsername", checkUsername).Methods("POST")
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")
	myRouter.HandleFunc("/canLogin", canLogin).Methods("POST")
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Handle Requests
	handleRequests()
}
