package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

//Runs a mongo query to get all Usernames, then puts it in a map to return
func loadUsernames() map[string]bool {
	mapOusernameToReturn := make(map[string]bool)
	usernameMap = make(map[string]bool) //Clear all Usernames when loading so no problems are caused

	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", "http://18.191.212.197:8080", nil)
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

	return mapOusernameToReturn
}
