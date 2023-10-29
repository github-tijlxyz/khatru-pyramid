package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/plugins/storage/badgern"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
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

	relay.Router().HandleFunc("/reports", reportsViewerHandler)
	relay.Router().HandleFunc("/add-to-whitelist", addToWhitelistHandler)
	relay.Router().HandleFunc("/remove-from-whitelist", removeFromWhitelistHandler)
	relay.Router().HandleFunc("/users", inviteTreeHandler)
	relay.Router().HandleFunc("/", homePageHandler)

	log.Info().Msg("running on http://0.0.0.0:" + s.Port)
	if err := http.ListenAndServe(":"+s.Port, relay); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}

func getLoggedUser(r *http.Request) string {
	if cookie, _ := r.Cookie("nip98"); cookie != nil {
		if evtj, err := url.QueryUnescape(cookie.Value); err == nil {
			var evt nostr.Event
			if err := json.Unmarshal([]byte(evtj), &evt); err == nil {
				return evt.PubKey
			}
		}
	}
	return ""
}
