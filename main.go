package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/plugins/storage/badgern"
	"github.com/joho/godotenv"
)

var (
	relayMaster      string
	db               badgern.BadgerBackend
	relayName        string = ""
	relayPubkey      string = ""
	relayDescription string = "none"
	relayContact     string = "none"
)

func main() {
	// save whitelist on shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSignals()
	}()

	// backup whitelist every hour
	go func() {
		for {
			time.Sleep(time.Hour)
			saveWhitelist()
		}
	}()

	// init env config
	godotenv.Load(".env")

	// init relay
	relay := khatru.NewRelay()

	relayMaster = os.Getenv("INVITE_RELAY_MASTER")

	// add information here!
	relayName = os.Getenv("RELAY_NAME")
	relayPubkey = os.Getenv("RELAY_PUBKEY")
	relayDescription = os.Getenv("RELAY_DESCRIPTION")
	relayContact = os.Getenv("RELAY_CONTACT")

	relay.Name = relayName
	relay.PubKey = relayPubkey
	relay.Description = relayDescription
	relay.Contact = relayContact

	// load whitelist storage
	if err := loadWhitelist(); err != nil {
		panic(err)
	}

	// load db
	db = badgern.BadgerBackend{Path: "./khatru-badgern-db"}
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

	relay.RejectEvent = append(relay.RejectEvent, whitelistRejecter)

	// ui
	relay.Router().HandleFunc("/reports", reportsViewerHandler)
	relay.Router().HandleFunc("/users", inviteTreeHandler)
	relay.Router().HandleFunc("/", redirectHandler)

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}

// save whitelist on shutdown
func handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	saveWhitelist()
	os.Exit(0)
}
