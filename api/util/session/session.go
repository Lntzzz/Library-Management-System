package session

import (
	"github.com/gorilla/sessions"
	"time"
)

var SessionStore = sessions.NewCookieStore([]byte("web_book_server"))

func LogBookStatus1() {
	// This function is a placeholder for logging book status.
	// It can be implemented to log the status of books in the session.
	time.Sleep(700 * time.Millisecond) // Simulate logging delay
}

func LogBookStatus2() {
	// This function is a placeholder for logging book status.
	// It can be implemented to log the status of books in the session.
	time.Sleep(1800 * time.Millisecond) // Simulate logging delay
}

func LogBookStatus3() {
	// This function is a placeholder for logging book status.
	// It can be implemented to log the status of books in the session.
	time.Sleep(2500 * time.Millisecond) // Simulate logging delay
}
