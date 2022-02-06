package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// data structure for visitor entries
type VisitorEntry struct {
	Timestamp time.Time
	Source    string
	FName     string
	LName     string
	Country   string
	City      string
	State     string
	Message   string
}

// date formats
const (
	layoutISO   = "2006-01-02"
	layoutUS    = "January 2, 2006"
	journalPath = "../journal/"
	logPath     = "../logs/"
	defaultPort = "80"
)

var policy = bluemonday.StrictPolicy()

// format entry as CSV
func (p *VisitorEntry) toCSV() string {
	return p.Timestamp.UTC().Format(time.UnixDate) + ",\"" + p.Source + "\"," + p.FName + "," + p.LName + "," + p.City + "," + p.State + "," + p.Country + ",\"" + p.Message + "\""
}

// format entry for printing
func (p *VisitorEntry) toString() string {
	return p.Timestamp.UTC().Format(layoutUS) + " (" + p.Source + "): " + p.FName + " " + p.LName + " from " + p.City + ", " + p.State + " (" + p.Country + ") wrote: \"" + p.Message + "\""
}

// save entry to disk
func (p *VisitorEntry) save() {
	// filename is based on the date of message creation
	filename := journalPath + "visitor-log-" + p.Timestamp.Format(layoutISO) + ".csv"

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
		if _, err := f.Write([]byte("Time,Source,First Name,Last Name,City,State,Country,Message\n")); err != nil {
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
		Source:    r.RemoteAddr + ";" + r.UserAgent(),
		FName:     strings.TrimSpace(policy.Sanitize(r.FormValue("fname"))),
		LName:     strings.TrimSpace(policy.Sanitize(r.FormValue("lname"))),
		Country:   strings.TrimSpace(policy.Sanitize(r.FormValue("country"))),
		City:      strings.TrimSpace(policy.Sanitize(r.FormValue("city"))),
		State:     strings.TrimSpace(policy.Sanitize(r.FormValue("state"))),
		Message:   strings.TrimSpace(policy.Sanitize(r.FormValue("message"))),
	}
	p.save()
	fmt.Println(p.toString())
	// display confirmationß
	http.Redirect(w, r, "/confirmation", http.StatusFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("Received from " + r.RemoteAddr + "; UserAgent: " + r.UserAgent())
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
	// create log and journal folders, if do not exist
	os.Mkdir(journalPath, 0755)
	os.Mkdir(logPath, 0755)

	// create log folder, create daily log file if needed
	file, err := os.OpenFile(logPath+"log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("Application starting...")

	var port = defaultPort

	// process command-line parameters
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--port":
			if (i + 1) < len(os.Args) {
				port = os.Args[i+1]
			} else {
				fmt.Println("Port wasn't specified. Exiting.")
				log.Fatal("Port wasn't specified. Exiting.")
			}
		default:
			log.Println("Unknown parameter: " + os.Args[i])
		}
	}

	/* Testing
	p1 := &VisitorEntry{Timestamp: time.Now(), Source: "Test", FName: "John", LName: "Smith", City: "New York", State: "NY", Country: "US", Message: "This is a sample message."}
	p1.save()
	fmt.Println(p1.toString())
	*/

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/confirmation", confirmation)
	err = http.ListenAndServe(":"+port, nil)
	fmt.Println(err)
	log.Fatal(err)
}
