package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/* This is the current amount of results our User is looking at
it changes as the User clicks forwards or backwards for more results */
var currentPageNumHotDog int = 1
var currentPageNumHamburger int = 1

var loadedMessagesMapHDog map[int]Message
var loadedMessagesMapHam map[int]Message
var theMessageBoardHDog MessageBoard //The board containing all our hotdog messages
var theMessageBoardHam MessageBoard  //The board containing all our hamburger messages

type Message struct {
	MessageID       int       `json:"MessageID"`       //ID of this Message
	UserID          int       `json:"UserID"`          //ID of the owner of this message
	PosterName      string    `json:"PosterName"`      //Username of the poster of this message
	Messages        []Message `json:"Messages"`        //Array of Messages under this one
	IsChild         bool      `json:"IsChild"`         //Is this message childed to another message
	HasChildren     bool      `json:"HasChildren"`     //Whether this message has children to list
	ParentMessageID int       `json:"ParentMessageID"` //The ID of this parent
	UberParentID    int       `json:"UberParentID"`    //The final parent of this parent, IF EQUAL PARENT
	Order           int       `json:"Order"`           //Order the commnet is in with it's reply tree
	RepliesAmount   int       `json:"RepliesAmount"`   //Amount of replies this message has
	TheMessage      string    `json:"TheMessage"`      //The Message in the post
	WhatBoard       string    `json:"WhatBoard"`       //The board this message is apart of
	DateCreated     string    `json:"DateCreated"`     //When the message was created
	LastUpdated     string    `json:"LastUpdated"`     //When the message was last updated
}

type MessageBoard struct {
	MessageBoardID         int             `json:"MessageBoardID"`
	BoardName              string          `json:"BoardName"`              //The Name of the board
	AllMessages            []Message       `json:"AllMessages"`            //All the IDs listed
	AllMessagesMap         map[int]Message `json:"AllMessagesMap"`         //A map of ALL messages
	AllOriginalMessages    []Message       `json:"AllOriginalMessages"`    //All the messages that AREN'T replies
	AllOriginalMessagesMap map[int]Message `json:"AllOriginalMessagesMap"` //Map of original Messages
	LastUpdated            string          `json:"LastUpdated"`            //Last time this messageboard was updated
	DateCreated            string          `json:"DateCreated"`            //Date this board was created
}

/*Creates a list of test messages, (if messageboard isn't created).
Initially called in the init function */
func createTestMessages() {
	//Ping our CRUD Microservice to see if Messageboards are already created
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080/isMessageBoardCreated", nil)
	if err != nil {
		theErr := "There was an error reaching out to isMessageBoardCreated: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for seeing if a messageboard is created " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr      []string     `json:"TheErr"`
		ResultMsg   []string     `json:"ResultMsg"`
		SuccOrFail  int          `json:"SuccOrFail"`
		GivenHDogMB MessageBoard `json:"GivenHDogMB"`
		GivenHamMB  MessageBoard `json:"GivenHamMB"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	/*Assign the messageboards for hotdog and hamburger off of the response, (if it is 0)
	Also fill those loadedMessage */
	if returnedMessage.SuccOrFail != 0 {
		//Log the failure to get the database
		message := "Failure to get the hotdog and hamburger messageboards"
		logWriter(message)
		theMessageBoardHDog = MessageBoard{}
		theMessageBoardHam = MessageBoard{}
	} else {
		fmt.Printf("DEBUG: Got the messageboards for hotdog and hamburger\n")
		theMessageBoardHDog = returnedMessage.GivenHDogMB
		theMessageBoardHam = returnedMessage.GivenHamMB
		//Fill the hotdog Messagemap
		for g := 0; g < len(theMessageBoardHDog.AllOriginalMessages); g++ {
			loadedMessagesMapHDog[g+1] = theMessageBoardHDog.AllOriginalMessages[g]
		}
		//Fill the hamburger MessageMap
		for z := 0; z < len(theMessageBoardHam.AllOriginalMessages); z++ {
			loadedMessagesMapHam[z+1] = theMessageBoardHam.AllOriginalMessages[z]
		}
	}
}

/* This refreshes messages from the database anytime a User
pings one of our web pages */
func refreshDatabases() {
	//Nullify all maps to not cause issues
	loadedMessagesMapHDog = make(map[int]Message)
	loadedMessagesMapHam = make(map[int]Message)
	/* Ping our CRUD API to get our most recent databases */
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080/isMessageBoardCreated", nil)
	if err != nil {
		theErr := "There was an error getting a random id in getRandomID: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for refreshDatabases: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	type ReturnMessage struct {
		TheErr      []string     `json:"TheErr"`
		ResultMsg   []string     `json:"ResultMsg"`
		SuccOrFail  int          `json:"SuccOrFail"`
		GivenHDogMB MessageBoard `json:"GivenHDogMB"`
		GivenHamMB  MessageBoard `json:"GivenHamMB"`
	}
	var otherReturnedMessage ReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)

	//Validate our responses; if successful response, update our db
	if otherReturnedMessage.SuccOrFail != 0 {
		//Failure, log error
		message := ""
		for j := 0; j < len(otherReturnedMessage.TheErr); j++ {
			message = message + otherReturnedMessage.TheErr[j] + "\n"
		}
		logWriter(message)
		fmt.Println(message)
	} else {
		//Fill our databases with response
		theMessageBoardHDog = otherReturnedMessage.GivenHDogMB
		theMessageBoardHam = otherReturnedMessage.GivenHamMB
		//Fill our messages with go routine
		wg.Add(1)
		go fillMessageMaps("hotdog")
		wg.Add(1)
		go fillMessageMaps("hamburger")
		wg.Wait()
	}
}

/* DEBUG: This goroutine is supposed to quickly fill our messageMaps
concurrently....not sure how useful this is yet, or if it can be converted into a
channel */
func fillMessageMaps(whichMap string) {
	switch whichMap {
	case "hotdog":
		for x := 0; x < len(theMessageBoardHDog.AllOriginalMessages); x++ {
			loadedMessagesMapHDog[x+1] = theMessageBoardHDog.AllOriginalMessages[x]
		}
		break
	case "hamburger":
		for x := 0; x < len(theMessageBoardHam.AllOriginalMessages); x++ {
			loadedMessagesMapHam[x+1] = theMessageBoardHam.AllOriginalMessages[x]
		}
		break
	default:
		err := "Wrong 'whichMap' entered in fillMessageMaps: " + whichMap
		fmt.Println(err)
		logWriter(err)
		break
	}
	wg.Done()
}

//This gets 10 results for display on a messageboard page
func getTenResults(whatPageNum int, whatBoard string) ([]Message, bool) {
	giveMessages := []Message{}
	topResult := whatPageNum * 10              //Last Result to add to map
	minResult := ((whatPageNum * 10) - 10) + 1 //First result to add to map
	okayResult := true                         //The result returned if we have messages to return

	switch whatBoard {
	case "hotdog":
		//Initial check to see if a map exists
		if _, ok := loadedMessagesMapHDog[minResult]; ok {
			//top value exists, get the message in range
			for x := minResult; x <= topResult; x++ {
				//Check to see if top result exists to value of ten; if not, add nothing
				if _, ok := loadedMessagesMapHDog[x]; ok {
					giveMessages = append(giveMessages, loadedMessagesMapHDog[x])
				} else {
					//Do nothing, there's no message here
				}
			}
			okayResult = true
		} else {
			fmt.Printf("DEBUG: Page value does not exist! The Value: %v\n", minResult)
			fmt.Printf("DEBUG: Here is our map currently: \n\n%v\n\n", loadedMessagesMapHDog)
			okayResult = false
		}

		//Reversing order of slice for 'MessageBoard Display' purposes
		giveMessages = reverseSlice(giveMessages)

		for q := 0; q < len(giveMessages); q++ {
			//fmt.Printf("giveMessages results %v: %v\n", q, giveMessages[q])
		}
		break
	case "hamburger":
		//Initial check to see if a map exists
		if _, ok := loadedMessagesMapHam[minResult]; ok {
			//top value exists, get the message in range
			for x := minResult; x <= topResult; x++ {
				//Check to see if top result exists to value of ten; if not, add nothing
				if _, ok := loadedMessagesMapHam[x]; ok {
					giveMessages = append(giveMessages, loadedMessagesMapHam[x])
				} else {
					//Do nothing, there's no message here
				}
			}
			okayResult = true
		} else {
			fmt.Printf("DEBUG: Page value does not exist! The Value: %v\n", minResult)
			fmt.Printf("DEBUG: Here is our map currently: \n\n%v\n\n", loadedMessagesMapHam)
			okayResult = false
		}

		//Reversing order of slice for 'MessageBoard Display' purposes
		giveMessages = reverseSlice(giveMessages)

		for q := 0; q < len(giveMessages); q++ {
			//fmt.Printf("giveMessages results %v: %v\n", q, giveMessages[q])
		}
		break
	default:
		theMessage := "Error, wrong board entered: " + whatBoard
		logWriter(theMessage)
		fmt.Printf("DEBUG: %v\n", theMessage)
		okayResult = false
		break
	}

	return giveMessages, okayResult
}

/* Called in Ajax from Javascript everytime User clicks left or right or submits a page
with results they'd like to see. If it's successful, it returns a number of JSON formatted Messages
for the page to update with. If not, it returns an error, which can be put in the pageNumber field. */
func evaluateTenResults(w http.ResponseWriter, r *http.Request) {
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Declare datatype from Ajax
	type PageData struct {
		ThePage  int    `json:"ThePage"`
		WhatPage string `json:"WhatPage"`
	}
	//Unmarshal JSON
	var pageDataPosted PageData
	json.Unmarshal(bs, &pageDataPosted)
	//Attempt to get data from loaded message map
	someMessages, goodMessageFind := getTenResults(pageDataPosted.ThePage, pageDataPosted.WhatPage)
	//Declare data to return
	type ReturnMessage struct {
		Messages   []Message `json:"Messages"`
		ResultMsg  string    `json:"ResultMsg"`
		SuccOrFail int       `json:"SuccOrFail"`
	}
	if goodMessageFind == true {
		//Set the current page number server side in case User refreshes
		switch pageDataPosted.WhatPage {
		case "hotdog":
			currentPageNumHotDog = pageDataPosted.ThePage
			//Return failure message
			theReturnMessage := ReturnMessage{
				Messages:   someMessages,
				ResultMsg:  "Page Found",
				SuccOrFail: 0,
			}
			theJSONMessage, err := json.Marshal(theReturnMessage)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(w, string(theJSONMessage))
			break
		case "hamburger":
			currentPageNumHamburger = pageDataPosted.ThePage
			//Return failure message
			theReturnMessage := ReturnMessage{
				Messages:   someMessages,
				ResultMsg:  "Page Found",
				SuccOrFail: 0,
			}
			theJSONMessage, err := json.Marshal(theReturnMessage)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(w, string(theJSONMessage))
			break
		default:
			//Return failure message
			theReturnMessage := ReturnMessage{
				Messages:   someMessages,
				ResultMsg:  "Error finding page...wrong data posted: " + pageDataPosted.WhatPage,
				SuccOrFail: 1,
			}
			theJSONMessage, err := json.Marshal(theReturnMessage)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(w, string(theJSONMessage))
			break
		}
	} else {
		//Return failure message
		theReturnMessage := ReturnMessage{
			Messages:   someMessages,
			ResultMsg:  "Error finding page...",
			SuccOrFail: 1,
		}
		theJSONMessage, err := json.Marshal(theReturnMessage)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	}
}

//This is for reversing the order of a Message array for display
func reverseSlice(orderedSlice []Message) []Message {
	last := len(orderedSlice) - 1
	for i := 0; i < len(orderedSlice)/2; i++ {
		orderedSlice[i], orderedSlice[last-i] = orderedSlice[last-i], orderedSlice[i]
	}

	return orderedSlice
}

/* Called in ajax from Javascript when a User submits an original message to a thread.
This calls the CRUD Microservice to perform the updates */
func messageOriginalAjax(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type OriginalMessage struct {
		TheMessage string `json:"TheMessage"`
		PosterName string `json:"PosterName"`
		WhatBoard  string `json:"WhatBoard"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedMessage OriginalMessage
	json.Unmarshal(bs, &postedMessage)

	newestMessage := Message{}

	//Declare return data and inform Ajax
	type DataReturn struct {
		SuccessMsg     string  `json:"SuccessMsg"`
		SuccessBool    bool    `json:"SuccessBool"`
		SuccessInt     int     `json:"SuccessInt"`
		CreatedMessage Message `json:"CreatedMessage"`
		ThePageNow     int     `json:"ThePageNow"`
	}

	theDataReturn := DataReturn{
		SuccessMsg:     "You created a new, original message",
		SuccessBool:    true,
		SuccessInt:     0,
		CreatedMessage: newestMessage,
		ThePageNow:     0,
	}

	theTimeNow := time.Now() //Needed for setting time values
	theOrder := 1            //Needed for setting the order value in messages
	//Set the 'order' for the newest message
	switch postedMessage.WhatBoard {
	case "hotdog":
		theOrder = len(theMessageBoardHDog.AllOriginalMessages) + 1
		//Format the new Original Message
		newestMessage = Message{
			MessageID:       getRandomID(),
			UserID:          getRandomID(),
			PosterName:      postedMessage.PosterName,
			Messages:        []Message{},
			IsChild:         false,
			HasChildren:     false,
			ParentMessageID: 0,
			UberParentID:    0,
			Order:           theOrder,
			RepliesAmount:   0,
			TheMessage:      postedMessage.TheMessage,
			WhatBoard:       postedMessage.WhatBoard,
			DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
			LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Insert new Message into database and update on server
		wg.Add(1)
		go insertOneMessage(newestMessage)
		//Update our hotdog messageboard
		theMessageBoardHDog.AllMessages = append(theMessageBoardHDog.AllMessages, newestMessage)
		theMessageBoardHDog.AllMessagesMap[newestMessage.MessageID] = newestMessage
		theMessageBoardHDog.AllOriginalMessages = append(theMessageBoardHDog.AllOriginalMessages, newestMessage)
		theMessageBoardHDog.AllOriginalMessagesMap[newestMessage.MessageID] = newestMessage
		theMessageBoardHDog.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		wg.Add(1)
		go updateMongoMessageBoard(theMessageBoardHDog)
		//Update the messagemap as well
		loadedMessagesMapHDog[len(loadedMessagesMapHDog)+1] = newestMessage
		//Set value for return data
		theDataReturn.ThePageNow = currentPageNumHotDog
		break
	case "hamburger":
		theOrder = len(theMessageBoardHam.AllOriginalMessages) + 1
		//Format the new Original Message
		newestMessage = Message{
			MessageID:       getRandomID(),
			UserID:          getRandomID(),
			PosterName:      postedMessage.PosterName,
			Messages:        []Message{},
			IsChild:         false,
			HasChildren:     false,
			ParentMessageID: 0,
			UberParentID:    0,
			Order:           theOrder,
			RepliesAmount:   0,
			TheMessage:      postedMessage.TheMessage,
			WhatBoard:       postedMessage.WhatBoard,
			DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
			LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Insert new Message into database and update on server
		wg.Add(1)
		go insertOneMessage(newestMessage)
		//Update our hotdog messageboard
		theMessageBoardHam.AllMessages = append(theMessageBoardHam.AllMessages, newestMessage)
		theMessageBoardHam.AllMessagesMap[newestMessage.MessageID] = newestMessage
		theMessageBoardHam.AllOriginalMessages = append(theMessageBoardHam.AllOriginalMessages, newestMessage)
		theMessageBoardHam.AllOriginalMessagesMap[newestMessage.MessageID] = newestMessage
		theMessageBoardHam.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		wg.Add(1)
		go updateMongoMessageBoard(theMessageBoardHam)
		//Update the messagemap as well
		loadedMessagesMapHam[len(loadedMessagesMapHam)+1] = newestMessage
		//Set value for return data
		theDataReturn.ThePageNow = currentPageNumHamburger
		break
	default:
		message := "Incorrect 'whatboard' is put in messageOriginalAjax: " + postedMessage.WhatBoard
		fmt.Println(message)
		logWriter(message)
		theDataReturn.SuccessBool = false
		theDataReturn.SuccessInt = 1
		theDataReturn.SuccessMsg = message
		break
	}

	wg.Wait() //Not sure if this is needed or in the right spot...DEBUG

	dataJSON, err := json.Marshal(theDataReturn)
	if err != nil {
		fmt.Println("There's an error marshalling this data")
	}
	fmt.Fprintf(w, string(dataJSON))
}

/* Called in Ajax from Javascript when a User submits a reply to any message/
reply */
func messageReplyAjax(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type MessageReply struct {
		ParentMessage Message `json:"ParentMessage"`
		ChildMessage  Message `json:"ChildMessage"`
		CurrentPage   int     `json:"CurrentPage"`
		PosterName    string  `json:"PosterName"`
		WhatBoard     string  `json:"WhatBoard"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedMessageReply MessageReply
	json.Unmarshal(bs, &postedMessageReply)

	//Declare return data and inform Ajax
	newestMessage := Message{}
	type ReturnData struct {
		SuccessMsg     string  `json:"SuccessMsg"`
		SuccessBool    bool    `json:"SuccessBool"`
		SuccessInt     int     `json:"SuccessInt"`
		CreatedMessage Message `json:"CreatedMessage"`
		ParentMessage  Message `json:"ParentMessage"`
	}
	theReturnData := ReturnData{
		SuccessMsg:     "You updated the messages",
		SuccessBool:    true,
		SuccessInt:     0,
		CreatedMessage: newestMessage,
		ParentMessage:  postedMessageReply.ParentMessage,
	}

	/* Format the parent message; we grab this becuase it might have been updated
	from Ajax, but the page hasn't refreshed to give us a new parent value we might
	have from a previous refresh */
	switch postedMessageReply.WhatBoard {
	case "hotdog":
		formattedParent := theMessageBoardHDog.AllMessagesMap[postedMessageReply.ParentMessage.MessageID]

		//Determine if this parent is an UberParent
		theMessageUberID := 0
		if formattedParent.IsChild == false {
			theMessageUberID = postedMessageReply.ParentMessage.MessageID
		} else {
			theMessageUberID = formattedParent.UberParentID
		}

		theTimeNow := time.Now()

		//Format the newestMessage
		newestMessage := Message{
			MessageID:       getRandomID(),
			UserID:          getRandomID(),
			PosterName:      postedMessageReply.PosterName,
			Messages:        []Message{},
			IsChild:         true,
			HasChildren:     false,
			ParentMessageID: formattedParent.MessageID,
			UberParentID:    theMessageUberID,
			Order:           len(formattedParent.Messages) + 1,
			RepliesAmount:   0,
			TheMessage:      postedMessageReply.ChildMessage.TheMessage,
			WhatBoard:       postedMessageReply.WhatBoard,
			DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
			LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
		}

		//Update the Message
		success := uberUpdate(newestMessage, formattedParent, postedMessageReply.WhatBoard)
		if success == true {
			theReturnData.SuccessBool = true
			theReturnData.SuccessInt = 0
			theReturnData.SuccessMsg = "Successful datatbase update"
		} else {
			theReturnData.SuccessBool = false
			theReturnData.SuccessInt = 1
			theReturnData.SuccessMsg = "Un-Successful datatbase update"
		}
		break
	case "hamburger":
		formattedParent := theMessageBoardHDog.AllMessagesMap[postedMessageReply.ParentMessage.MessageID]

		//Determine if this parent is an UberParent
		theMessageUberID := 0
		if formattedParent.IsChild == false {
			theMessageUberID = postedMessageReply.ParentMessage.MessageID
		} else {
			theMessageUberID = formattedParent.UberParentID
		}

		theTimeNow := time.Now()

		//Format the newestMessage
		newestMessage := Message{
			MessageID:       getRandomID(),
			UserID:          getRandomID(),
			PosterName:      postedMessageReply.PosterName,
			Messages:        []Message{},
			IsChild:         true,
			HasChildren:     false,
			ParentMessageID: formattedParent.MessageID,
			UberParentID:    theMessageUberID,
			Order:           len(formattedParent.Messages) + 1,
			RepliesAmount:   0,
			TheMessage:      postedMessageReply.ChildMessage.TheMessage,
			WhatBoard:       postedMessageReply.WhatBoard,
			DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
			LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
		}

		//Update the Message
		success := uberUpdate(newestMessage, formattedParent, postedMessageReply.WhatBoard)
		if success == true {
			theReturnData.SuccessBool = true
			theReturnData.SuccessInt = 0
			theReturnData.SuccessMsg = "Successful datatbase update"
		} else {
			theReturnData.SuccessBool = false
			theReturnData.SuccessInt = 1
			theReturnData.SuccessMsg = "Un-Successful datatbase update"
		}
		break
	default:
		message := "Incorrect 'whatboard' in messageReplyAjax: " + postedMessageReply.WhatBoard
		logWriter(message)
		fmt.Println(message)
		theReturnData.SuccessBool = false
		theReturnData.SuccessInt = 1
		theReturnData.SuccessMsg = message
		break
	}

	dataJSON, err := json.Marshal(theReturnData)
	if err != nil {
		fmt.Println("There's an error marshalling this data")
	}
	fmt.Fprintf(w, string(dataJSON))
}

/* Calls the UberUpdate in the CRUD Microservice. Updates our messageboards
and maps */
func uberUpdate(newestMessage Message, parentMessage Message, whatBoard string) bool {
	success := true
	//Declare JSON to send
	type UberUpdateMessages struct {
		TheNewestMessage Message         `json:"TheNewestMessage"`
		TheParentMessage Message         `json:"TheParentMessage"`
		WhatBoard        string          `json:"WhatBoard"`
		HotdogMB         MessageBoard    `json:"HotdogMB"`
		HamburgerMB      MessageBoard    `json:"HamburgerMB"`
		LoadedMapHDog    map[int]Message `json:"LoadedMapHDog"`
		LoadedMapHam     map[int]Message `json:"LoadedMapHam"`
	}
	uberUpdateSend := UberUpdateMessages{
		TheNewestMessage: newestMessage,
		TheParentMessage: parentMessage,
		WhatBoard:        whatBoard,
		HotdogMB:         theMessageBoardHDog,
		HamburgerMB:      theMessageBoardHam,
		LoadedMapHDog:    loadedMessagesMapHDog,
		LoadedMapHam:     loadedMessagesMapHam,
	}

	theJSONMessage, err := json.Marshal(uberUpdateSend)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send to CRUD OPERATIONS API
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.191.212.197:8080/uberUpdate", payload)
	if err != nil {
		theErr := "There was an error inserting a message in uberUpdate: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for updating MessageBoard in uberUpdate: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the returned response from Create User
	type ReturnMessage struct {
		TheErr             []string        `json:"TheErr"`
		ResultMsg          []string        `json:"ResultMsg"`
		SuccOrFail         int             `json:"SuccOrFail"`
		GivenHDogMB        MessageBoard    `json:"GivenHDogMB"`
		GivenHamMB         MessageBoard    `json:"GivenHamMB"`
		GivenLoadedMapHDog map[int]Message `json:"GivenLoadedMapHDog"`
		GivenLoadedMapHam  map[int]Message `json:"GivenLoadedMapHam"`
	}
	var otherReturnedMessage ReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)

	/* If given a successful response, update our messageboard appropriatley;
	if no successful response, log appropriatly */
	if otherReturnedMessage.SuccOrFail != 0 {
		//Log the error
		message := ""
		for j := 0; j < len(otherReturnedMessage.TheErr); j++ {
			message = message + "\n" + otherReturnedMessage.TheErr[j]
		}
		logWriter(message)
		fmt.Println(message)
		success = false
	} else {
		message := ""
		for j := 0; j < len(otherReturnedMessage.ResultMsg); j++ {
			message = message + otherReturnedMessage.ResultMsg[j]
		}
		logWriter(message)
		fmt.Println(message)
		//Update all of our message variables
		loadedMessagesMapHDog = otherReturnedMessage.GivenLoadedMapHDog
		loadedMessagesMapHam = otherReturnedMessage.GivenLoadedMapHam
		theMessageBoardHDog = otherReturnedMessage.GivenHDogMB
		theMessageBoardHam = otherReturnedMessage.GivenHamMB
	}

	return success
}

/* Insert one message,(calls the CRUD Microservice, can be used with Go Routines) */
func insertOneMessage(theMessage Message) {
	theJSONMessage, err := json.Marshal(theMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send to CRUD OPERATIONS API
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.191.212.197:8080/insertOneNewMessage", payload)
	if err != nil {
		theErr := "There was an error inserting a message in insertOneMessage: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
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

	if otherReturnedMessage.SuccOrFail != 0 {
		theMessage := otherReturnedMessage.TheErr
		logWriter(theMessage)
		fmt.Println(theMessage)
	} else {
		theMessage := otherReturnedMessage.ResultMsg
		logWriter(theMessage)
		fmt.Println(theMessage)
	}

	wg.Done()
}

/* Update one messageboard, (calls the CRUD Microservice, can be used with GO Routines) */
func updateMongoMessageBoard(updatedMessageBoard MessageBoard) {

	type UpdatedMongoBoard struct {
		UpdatedMessageBoard MessageBoard `json:"UpdatedMessageBoard"`
	}
	theboardUpdate := UpdatedMongoBoard{
		UpdatedMessageBoard: updatedMessageBoard,
	}

	theJSONMessage, err := json.Marshal(theboardUpdate)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	//Send to CRUD OPERATIONS API
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", "http://18.191.212.197:8080/updateMongoMessageBoard", payload)
	if err != nil {
		theErr := "There was an error inserting a message in updateMongoMessageBoard: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for updating message in updateMongoMessageBoard: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the returned response from Create User
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	var theReturnedMessage ReturnMessage
	json.Unmarshal(body, &theReturnedMessage)

	if theReturnedMessage.SuccOrFail != 0 {
		theMessage := ""
		for i := 0; i < len(theReturnedMessage.TheErr); i++ {
			theMessage = "\n" + theReturnedMessage.TheErr[i]
		}
		fmt.Println(theMessage)
		logWriter(theMessage)
	} else {
		theMessage := ""
		for i := 0; i < len(theReturnedMessage.ResultMsg); i++ {
			theMessage = "\n" + theReturnedMessage.ResultMsg[i]
		}
		fmt.Println(theMessage)
		logWriter(theMessage)
	}

	wg.Done()
}

/* This function calls the CRUD Microservice to get a random ID */
func getRandomID() int {
	theID := 0

	//Send to CRUD OPERATIONS API
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080/randomIDCreationAPI", nil)
	if err != nil {
		theErr := "There was an error getting a random id in getRandomID: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for random ID in getRandomID: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		RandomID   int      `json:"RandomID"`
	}
	var otherReturnedMessage ReturnMessage
	json.Unmarshal(body, &otherReturnedMessage)

	theID = otherReturnedMessage.RandomID
	if otherReturnedMessage.SuccOrFail != 0 {
		for j := 0; j < len(otherReturnedMessage.TheErr); j++ {
			logWriter(otherReturnedMessage.TheErr[j])
		}
	}
	return theID
}
