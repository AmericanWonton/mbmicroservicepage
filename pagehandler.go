package main

import (
	"fmt"
	"net/http"
)

//Handles the Index requests
func index(w http.ResponseWriter, r *http.Request) {

	usernameMap = loadUsernames()
	fmt.Printf("DEBUG: here we are in index: \n")
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
	fmt.Printf("DEBUG: here we are in hotdogMB: \n")
	type ViewData struct {
		TheUser        AUser  `json:"TheUser"`        //The User we use
		MessageDisplay int    `json:"MessageDisplay"` //This is IF we need a message displayed
		WhatPage       string `json:"WhatPage"`       //What messageboard is displayed
	}
	vd := ViewData{
		TheUser:        aUser,
		MessageDisplay: 0,
		WhatPage:       "hotdogmb",
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
	fmt.Printf("DEBUG: here we are in hamburgerMB: \n")
	type ViewData struct {
		TheUser        AUser  `json:"TheUser"`        //The User we use
		MessageDisplay int    `json:"MessageDisplay"` //This is IF we need a message displayed
		WhatPage       string `json:"WhatPage"`       //What messageboard is displayed
	}
	vd := ViewData{
		TheUser:        aUser,
		MessageDisplay: 0,
		WhatPage:       "hamburgermb",
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "hamburgermsb.gohtml", vd)
	HandleError(w, err1)
}

//Handles the test page
func test(w http.ResponseWriter, r *http.Request) {
	userStuff := MessageViewData{
		TestString:  "bootyhole",
		TheMessages: []Message{},
		WhatPage:    currentPageNumber,
	}

	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "test.gohtml", userStuff)
	HandleError(w, err1)
}
