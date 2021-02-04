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
	"time"

	"github.com/gorilla/mux"
)

const sessionLength int = 180 //Length of sessions

//Here's our session struct
type theSession struct {
	username     string
	lastActivity time.Time
}

//Session Database info
var dbUsers = map[string]AUser{}         // user ID, user
var dbSessions = map[string]theSession{} // session ID, session
var dbSessionsCleaned time.Time

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

//Define function maps
var funcMap = template.FuncMap{
	"upperCase": strings.ToUpper, //upperCase is a key we can call inside of the template html file
}

//DEBUG not sure if needed
type MessageViewData struct {
	TestString  string    `json:"TestString"`
	TheMessages []Message `json:"TheMessages"`
	WhatPage    int       `json:"WhatPage"`
}

//DEBUG  Message displayed on the board
type Message struct {
	MessageID       int       `json:"MessageID"`       //ID of this Message
	UserID          int       `json:"UserID"`          //ID of the owner of this message
	Messages        []Message `json:"Messages"`        //Array of Messages under this one
	IsChild         bool      `json:"IsChild"`         //Is this message childed to another message
	HasChildren     bool      `json:"HasChildren"`     //Whether this message has children to list
	ParentMessageID int       `json:"ParentMessageID"` //The ID of this parent
	UberParentID    int       `json:"UberParentID"`    //The final parent of this parent, IF EQUAL PARENT
	Order           int       `json:"Order"`           //Order the commnet is in with it's reply tree
	RepliesAmount   int       `json:"RepliesAmount"`   //Amount of replies this message has
	TheMessage      string    `json:"TheMessage"`      //The MEssage in the post
	DateCreated     string    `json:"DateCreated"`     //When the message was created
	LastUpdated     string    `json:"LastUpdated"`     //When the message was last updated
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

/* This is the current amount of results our User is looking at
it changes as the User clicks forwards or backwards for more results */
var currentPageNumber int = 1

//Parse our templates
func init() {
	/* Assign blank value to map so no nil errors occur */
	usernameMap = make(map[string]bool) //Clear all Usernames when loading so no problems are caused
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
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
	fmt.Println(debugMessage)
	logWriter(debugMessage)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", index)
	myRouter.HandleFunc("/test", test)
	//Mongo No-SQL Stuff

	//Field validation/User Creation
	myRouter.HandleFunc("/checkUsername", checkUsername).Methods("POST")
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Handle Requests
	handleRequests()
}
