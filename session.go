package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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

func getUser(w http.ResponseWriter, req *http.Request) AUser {
	// get cookie
	cookie, err := req.Cookie("session")
	//If there is no session cookie, create a new session cookie
	if err != nil {
		uuidWithHyphen := uuid.New().String()
		cookie = &http.Cookie{
			Name:  "session",
			Value: uuidWithHyphen,
		}

	}
	//Set the cookie age to the max length again.
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie) //Set the cookie to our grabbed cookie,(or new cookie)

	// if the user exists already, get user
	var theUser AUser
	if session, ok := dbSessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
		theUser = dbUsers[session.username]
	}
	return theUser
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false //If there is an error getting the cookie, return false
	}
	//if session is found, we update the session with the newest time since activity!
	session, ok := dbSessions[cookie.Value]
	if ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
	}
	/* Check to see if the Username exists from this Session Username. If not, we return false. */
	_, ok = dbUsers[session.username]
	// refresh session
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
	return ok
}
