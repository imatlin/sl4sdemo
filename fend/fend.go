package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// data structure for visitor entries
type VisitorEntry struct {
	Timestamp time.Time
	FName     string
	LName     string
	Country   string
	City      string
	State     string
	Message   string
}

// date formats
const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
	relPath   = "./logs/"
	port      = ":8081"
)

var policy = bluemonday.StrictPolicy()

// format entry as CSV
func (p *VisitorEntry) toCSV() string {
	return p.Timestamp.UTC().Format(time.UnixDate) + "," + p.FName + "," + p.LName + "," + p.City + "," + p.State + "," + p.Country + ",\"" + p.Message + "\""
}

// format entry for printing
func (p *VisitorEntry) toString() string {
	return p.Timestamp.UTC().Format(layoutUS) + ": " + p.FName + " " + p.LName + " from " + p.City + "," + p.State + " (" + p.Country + ") wrote: \"" + p.Message + "\""
}

// save entry to disk
func (p *VisitorEntry) save() {
	// filename is based on the date of message creation
	filename := relPath + "visitor-log-" + p.Timestamp.Format(layoutISO) + ".txt"

	// open or create file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// if file is empty, write out header
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if _, err := f.Write([]byte("Time,First Name,Last Name,City,State,Country,Message\n")); err != nil {
			log.Fatal(err)
		}
	}

	// write out visitor entry as CSV
	if _, err := f.Write([]byte(p.toCSV() + "\n")); err != nil {
		log.Fatal(err)
	}

	// close the file
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	// save submitted values to a file
	p := &VisitorEntry{
		Timestamp: time.Now(),
		FName:     policy.Sanitize(r.FormValue("fname")),
		LName:     policy.Sanitize(r.FormValue("lname")),
		Country:   policy.Sanitize(r.FormValue("country")),
		City:      policy.Sanitize(r.FormValue("city")),
		State:     policy.Sanitize(r.FormValue("state")),
		Message:   policy.Sanitize(r.FormValue("message")),
	}
	p.save()
	fmt.Println(p.toString())
	// display confirmationß
	http.Redirect(w, r, "/confirmation", http.StatusFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //get request method
	if r.Method == http.MethodGet {
		// serve home page
		http.ServeFile(w, r, "templates/home.html")
	} else {
		// save form values
		saveHandler(w, r)
	}
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	// display confirmation, wait 5 seconds, reditect to home
	http.ServeFile(w, r, "templates/confirmation.html")
}

func main() {
	// create log folder, if does not exist
	os.Mkdir(relPath, 0755)

	/* Testing
	p1 := &VisitorEntry{Timestamp: time.Now(), FName: "John", LName: "Smith", City: "New York", State: "NY", Country: "US", Message: "This is a sample message."}
	p1.save()
	fmt.Println(p1.toString())
	*/

	http.HandleFunc("/", home)
	http.HandleFunc("/confirmation", confirmation)
	log.Fatal(http.ListenAndServe(port, nil))
}
