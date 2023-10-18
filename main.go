package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/plugins/storage/badgern"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port             string `envconfig:"PORT" default:"3334"`
	RelayName        string `envconfig:"RELAY_NAME" required:"true"`
	RelayPubkey      string `envconfig:"RELAY_PUBKEY" required:"true"`
	RelayDescription string `envconfig:"RELAY_DESCRIPTION"`
	RelayContact     string `envconfig:"RELAY_CONTACT"`
}

var (
	db  badgern.BadgerBackend
	s   Settings
	log = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func main() {
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig")
		return
	}

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

	relay.Name = s.RelayName
	relay.PubKey = s.RelayPubkey
	relay.Description = s.RelayDescription
	relay.Contact = s.RelayContact

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
	relay.Router().HandleFunc("/", homePageHandler)

	log.Info().Msg("running on http://0.0.0.0:" + s.Port)
	if err := http.ListenAndServe(":"+s.Port, relay); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}

// save whitelist on shutdown
func handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	saveWhitelist()
	os.Exit(0)
}
