package main

import (
	"net/http"
	"os"
)

//DEBUG not sure if needed
type MessageViewData struct {
	TestString  string    `json:"TestString"`
	TheMessages []Message `json:"TheMessages"`
	WhatPage    int       `json:"WhatPage"`
}

//This is for data that can be used for posting messages
type ViewData struct {
	Username       string    `json:"Username"`       //The Username
	UserID         int       `json:"UserID"`         //The UserID
	TheMessages    []Message `json:"TheMessages"`    //The Messages we need to display
	MessageDisplay int       `json:"MessageDisplay"` //This is IF we need a message displayed
	WhatPage       int       `json:"WhatPage"`       //What pageNumber we are on
	WhatBoard      string    `json:"WhatBoard"`      //This is what board we are posting from
}

//Handles the Index requests
func index(w http.ResponseWriter, r *http.Request) {
	usernameMap = loadUsernames()
	//fmt.Printf("DEBUG: here we are in index: \n")
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err1)
}

//Handles the Hotdog Messageboard Requests
func hotdogMB(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	/* First, we need to query for this messageboard, in case other users
	made comments while this other user was on another page */
	refreshDatabases() //Refresh our DBS/Messageboards Maps
	/* Second, we need to get 10 results based off of what page number it is */
	ourMessages, _ := getTenResults(currentPageNumHotDog, "hotdog")
	vd := ViewData{
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		TheMessages:    ourMessages,
		MessageDisplay: 0,
		WhatPage:       currentPageNumHotDog,
		WhatBoard:      "hotdog",
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "hotdogmsb.gohtml", vd)
	HandleError(w, err1)
}

//Handles the Hotdog Messageboard Requests
func hamburgerMB(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	/* First, we need to query for this messageboard, in case other users
	made comments while this other user was on another page */
	refreshDatabases() //Refresh our DBS/Messageboards Maps
	/* Second, we need to get 10 results based off of what page number it is */
	ourMessages, _ := getTenResults(currentPageNumHamburger, "hamburger")
	vd := ViewData{
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		TheMessages:    ourMessages,
		MessageDisplay: 0,
		WhatPage:       currentPageNumHamburger,
		WhatBoard:      "hamburger",
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "hamburgermsb.gohtml", vd)
	HandleError(w, err1)
}

//Handles the documentation page
func documentation(w http.ResponseWriter, req *http.Request) {
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "8080"
	}

	err1 := template1.ExecuteTemplate(w, "documentation.gohtml", nil)
	HandleError(w, err1)
}

//Handles the Contact page

func contact(w http.ResponseWriter, r *http.Request) {
	//Serve the template normally
	err1 := template1.ExecuteTemplate(w, "contact.gohtml", nil)
	HandleError(w, err1)
}

//Handles the test page
func test(w http.ResponseWriter, r *http.Request) {
	userStuff := MessageViewData{
		TestString:  "bootyhole",
		TheMessages: []Message{},
		WhatPage:    0,
	}

	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "test.gohtml", userStuff)
	HandleError(w, err1)
}
