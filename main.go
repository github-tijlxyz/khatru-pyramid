package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/fiatjaf/eventstore/badger"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/plugins"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port             string `envconfig:"PORT" default:"3334"`
	Domain           string `envconfig:"DOMAIN" required:"true"`
	RelayName        string `envconfig:"RELAY_NAME" required:"true"`
	RelayPubkey      string `envconfig:"RELAY_PUBKEY" required:"true"`
	RelayDescription string `envconfig:"RELAY_DESCRIPTION"`
	RelayContact     string `envconfig:"RELAY_CONTACT"`
	RelayIcon        string `envconfig:"RELAY_ICON"`
	DatabasePath     string `envconfig:"DATABASE_PATH" default:"./db"`
}

var (
	s         Settings
	db        = badger.BadgerBackend{}
	log       = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	whitelist = make(Whitelist)
	relay     = khatru.NewRelay()
)

func main() {
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig")
		return
	}

	// load db
	db.Path = s.DatabasePath
	if err := db.Init(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
		return
	}
	log.Debug().Str("path", db.Path).Msg("initialized database")

	// init relay
	relay.Name = s.RelayName
	relay.PubKey = s.RelayPubkey
	relay.Description = s.RelayDescription
	relay.Contact = s.RelayContact
	relay.IconURL = s.RelayIcon

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.RejectEvent = append(relay.RejectEvent,
		rejectEventsFromUsersNotInWhitelist,
		plugins.RestrictToSpecifiedKinds(supportedKinds...),
		validateAndFilterReports,
	)
	relay.OverwriteFilter = append(relay.OverwriteFilter,
		plugins.RemoveAllButKinds(supportedKinds...),
		removeAuthorsNotWhitelisted,
	)
	relay.RejectFilter = append(relay.RejectFilter,
		plugins.NoSearchQueries,
		discardFiltersWithTooManyAuthors,
	)

	// load users registry
	if err := loadWhitelist(); err != nil {
		log.Fatal().Err(err).Msg("failed to load whitelist")
		return
	}

	// http routes
	relay.Router().HandleFunc("/add-to-whitelist", addToWhitelistHandler)
	relay.Router().HandleFunc("/remove-from-whitelist", removeFromWhitelistHandler)
	relay.Router().HandleFunc("/reports", reportsViewerHandler)
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
				if tag := evt.Tags.GetFirst([]string{"domain", ""}); tag != nil && (*tag)[1] == s.Domain {
					if ok, _ := evt.CheckSignature(); ok {
						return evt.PubKey
					}
				}
			}
		}
	}
	return ""
}
