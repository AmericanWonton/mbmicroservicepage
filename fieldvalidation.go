package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

/* DEFINED SLURS */
var slurs []string = []string{}

//This gets the slur words we check against in our username and
//text messages
func getbadWords() {
	file, err := os.Open("security/badphrases.txt")

	if err != nil {
		fmt.Printf("DEBUG: Trouble opening bad word text file: %v\n", err.Error())
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	slurs = text
}

//Runs a mongo query to get all Usernames, then puts it in a map to return
func loadUsernames() map[string]bool {
	mapOusernameToReturn := make(map[string]bool)
	usernameMap = make(map[string]bool) //Clear all Usernames when loading so no problems are caused

	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	email := theUser.Email
	areacode := theUser.PhoneACode
	phonenum := theUser.PhoneNumber

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
		Email:       email,
		PhoneACode:  areacode,
		PhoneNumber: phonenum,
		PostsMade:   0,
		RepliesMade: 0,
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
	/* AT THIS POINT, WE SHOULD BE GOOD TO CONCURRENTLY SEND THE USER AN EMAIL
	AND TEXT MESSAGE */
	//Format variable to send in functions
	aMessage := ""
	aSubject := ""
	if theSuccMessage.SuccessNum != 0 {
		//Create failure messages
		aMessage = "Sorry, " + username + ", there was an error creating your profile. Please see below and " +
			"see if you can adjust your information for a successful submission. Otherwise, you can email me at johnnycowboy39@gmail.com " +
			"for further support. \n Thanks and have a great day.\n " + "Error: " + otherReturnedMessage.ResultMsg
		aSubject = "Profile creation failure"
	} else {
		//Create success messages
		aMessage = "Hello, " + username + ", thank you for creating a profile on this test site! Feel free to read the comments and add some of " +
			"your own! If you have any quesitons, please email me at johnnycowboy39@gmail.com.\n Thanks again,\nJoseph Keller"
		aSubject = "Profile Created"
	}

	//Send Text Message
	wg.Add(1)
	go sendText(aMessage, areacode, phonenum)
	//Send Email
	wg.Add(1)
	go sendEmail(aMessage, email, aSubject)

	wg.Wait() //Not sure if this is needed or in the right spot...DEBUG
}

//Sends a text message with a go routine, no response needed, it's just logged
func sendText(theMessage string, areacode int, phonenumber int) {
	//Declare DataType for JSON
	type TextInfo struct {
		TextMessage string `json:"TextMessage"`
		PhoneACode  int    `json:"PhoneACode"`
		PhoneNumber int    `json:"PhoneNumber"`
	}

	dataTextMessage := TextInfo{
		TextMessage: theMessage,
		PhoneACode:  areacode,
		PhoneNumber: phonenumber,
	}

	//Marshal into JSON to use
	theJSONMessage, err := json.Marshal(dataTextMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send the text message info and wait for a response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.188.234.83:80/sendTextMessage", payload)
	if err != nil {
		theErr := "There was an error sending a text in sendText: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for sendText: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Declare return information for JSON
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	var otherReturnedMessage ReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)

	//Log the results
	logMessage := ""
	if otherReturnedMessage.SuccOrFail != 0 {
		logMessage = "Failure to send text message: " + otherReturnedMessage.TheErr + "\n" + otherReturnedMessage.ResultMsg
	} else {
		area := strconv.Itoa(areacode)
		phonnum := strconv.Itoa(phonenumber)
		logMessage = "Message successfully created for " + area + "-" + phonnum
	}
	logWriter(logMessage)

	wg.Done() //For GoRoutines
}

//Sends an email with a go routine, no response needed it's just logged
func sendEmail(theMessage string, emailAddress string, subject string) {
	//Declare DataType for JSON
	//Declare DataType from JSON
	type EmailInfo struct {
		EmailMessage string `json:"EmailString"`
		EmailAddress string `json:"EmailAddressString"`
		EmailSubject string `json:"EmailSubject"`
	}

	dataEmail := EmailInfo{
		EmailMessage: theMessage,
		EmailAddress: emailAddress,
		EmailSubject: subject,
	}

	fmt.Printf("DEBUG: Here is our email: %v\n", dataEmail)

	//Marshal into JSON to use
	theJSONMessage, err := json.Marshal(dataEmail)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send the text message info and wait for a response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.188.234.83:80/sendEmail", payload)
	if err != nil {
		theErr := "There was an error sending an email in sendEmail: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response in sendEmail " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Declare return information for JSON
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
	}
	var otherReturnedMessage ReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)

	//Log the results
	logMessage := ""
	if otherReturnedMessage.SuccOrFail != 0 {
		logMessage = "Failure to email: " + otherReturnedMessage.TheErr + "\n" + otherReturnedMessage.ResultMsg
	} else {
		logMessage = "Email successfully created for " + emailAddress
	}
	logWriter(logMessage)

	wg.Done() //For GoRoutines
}

//User Sign In Check
func canLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: We reached canLogin\n")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Send a response back to Ajax after session is made
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{}

	//Declare DataType from Ajax
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	//Marshal the user data into our type
	var dataForLogin LoginData
	json.Unmarshal(bs, &dataForLogin)

	fmt.Printf("DEBUG: We got this username: %v and this password: %v\n", dataForLogin.Username,
		dataForLogin.Password)

	bsString := []byte(dataForLogin.Password)     //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2
	dataForLogin.Password = encodedString
	//Check to see if the login is legit
	//Query Mongo for those username and password
	//Send to CRUD OPERATIONS API
	theJSONMessage, err := json.Marshal(dataForLogin)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.191.212.197:8080/userLogin", payload)
	if err != nil {
		theErr := "There was an error pinging userLogin in canLogin: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	fmt.Printf("DEBUG: We got stuff back from the API\n")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for userLogin in canLogin: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	fmt.Printf("DEBUG: Here is the body: %v\n", string(body))

	//Marshal the returned response from Create User
	type ReturnMessage struct {
		TheErr     string `json:"TheErr"`
		ResultMsg  string `json:"ResultMsg"`
		SuccOrFail int    `json:"SuccOrFail"`
		TheUser    AUser  `json:"TheUser"`
	}
	var returnMessage ReturnMessage
	json.Unmarshal(body, &returnMessage)

	//Check for a null User returned
	if returnMessage.SuccOrFail != 0 {
		theSuccMessage = SuccessMSG{
			Message:    returnMessage.ResultMsg,
			SuccessNum: returnMessage.SuccOrFail,
		}
	} else {
		//Log User in and give session cookie, if needed
		theUser := returnMessage.TheUser
		dbUsers[dataForLogin.Username] = theUser
		// create session
		uuidWithHyphen := uuid.New().String()

		cookie := &http.Cookie{
			Name:  "session",
			Value: uuidWithHyphen,
		}
		cookie.MaxAge = sessionLength
		http.SetCookie(w, cookie)
		dbSessions[cookie.Value] = theSession{dataForLogin.Username, time.Now()}

		theSuccMessage = SuccessMSG{
			Message:    returnMessage.ResultMsg,
			SuccessNum: returnMessage.SuccOrFail,
		}
	}

	theJSONMessage, err = json.Marshal(theSuccMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}
