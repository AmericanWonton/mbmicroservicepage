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
