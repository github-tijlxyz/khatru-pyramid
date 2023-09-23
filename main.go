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

var relayMaster string
var db badgern.BadgerBackend

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
	relay.Name = "a invite relay"
	relay.PubKey = ""
	relay.Contact = ""

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

	// invitedata api
	relay.Router().HandleFunc("/invitedata", inviteDataApiHandler)
	relay.Router().HandleFunc("/relaymaster", relayMasterApiHandler)
	
	// ui
	relay.Router().HandleFunc("/", embeddedUIHandler)

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
