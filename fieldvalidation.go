package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

/* DEFINED SLURS */
var slurs []string = []string{"penis", "vagina", "dick", "cunt", "asshole", "fag", "faggot",
	"nigglet", "nigger", "beaner", "wetback", "wet back", "chink", "tranny", "bitch", "slut",
	"whore", "fuck", "damn",
	"shit", "piss", "cum", "jizz"}

//Runs a mongo query to get all Usernames, then puts it in a map to return
func loadUsernames() map[string]bool {
	mapOusernameToReturn := make(map[string]bool)
	usernameMap = make(map[string]bool) //Clear all Usernames when loading so no problems are caused

	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080/giveAllUsernames", nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	fmt.Printf("DEBUG: Here is our body we recieved: %v\n", string(body))

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign the map,(log failures if there are any)
	if returnedMessage.SuccOrFail != 0 {
		mapOusernameToReturn = returnedMessage.ReturnedUserMap
		bigLogMessage := ""
		for i := 0; i < len(returnedMessage.ResultMsg); i++ {
			bigLogMessage = bigLogMessage + returnedMessage.ResultMsg[i]
		}
		logWriter(bigLogMessage)
	} else {
		mapOusernameToReturn = returnedMessage.ReturnedUserMap
		bigLogMessage := ""
		for i := 0; i < len(returnedMessage.ResultMsg); i++ {
			bigLogMessage = bigLogMessage + returnedMessage.ResultMsg[i]
		}
		logWriter(bigLogMessage)
	}

	return mapOusernameToReturn
}

//Checks the Usernames after every keystroke
func checkUsername(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 20 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, usernameMap[sbs])
	}
}

//Checks to see if the Username contains language
func containsLanguage(theText string) bool {
	hasLanguage := false
	textLower := strings.ToLower(theText)
	for i := 0; i < len(slurs); i++ {
		if strings.Contains(textLower, slurs[i]) {
			hasLanguage = true
			return hasLanguage
		}
	}
	return hasLanguage
}

//User creation
func createUser(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var theUser AUser
	json.Unmarshal(bs, &theUser)

	// get form values
	username := theUser.UserName
	password := theUser.Password

	// create session
	uuidWithHyphen := uuid.New().String()
	newCookie := &http.Cookie{
		Name:  "session",
		Value: uuidWithHyphen,
	}
	newCookie.MaxAge = sessionLength
	http.SetCookie(w, newCookie)
	dbSessions[newCookie.Value] = theSession{username, time.Now()}

	//Begin to add User to Mongo
	bsString := []byte(password)                  //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2
	theTimeNow := time.Now()

	//Get RandomID for User
	//Call our crudOperations Microservice in order to get our new UserID
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080/randomIDCreationAPI", nil)
	if err != nil {
		theErr := "There was an error getting a random id in createUser: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for random ID in createUser: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		RandomID   int      `json:"RandomID"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	userSend := AUser{
		UserName:    username,
		Password:    encodedString,
		UserID:      returnedMessage.RandomID,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}
	theJSONMessage, err := json.Marshal(userSend)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send to CRUD OPERATIONS API
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err = http.NewRequest("POST", "http://18.191.212.197:8080/addUser", payload)
	if err != nil {
		theErr := "There was an error getting a random id in createUser: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err = http.DefaultClient.Do(req.WithContext(ctx))

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for random ID in createUser: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the returned response from Create User
	type OtherReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	var otherReturnedMessage OtherReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)
	/* Send the response back to Ajax */
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    otherReturnedMessage.ResultMsg,
		SuccessNum: otherReturnedMessage.SuccOrFail,
	}

	theJSONMessage, err = json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
