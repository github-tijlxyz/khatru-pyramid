package main

import (
	"fmt"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/theplant/htmlgo"
)

func inviteTreeHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)
	content := inviteTreePageHTML(r.Context(), InviteTreePageParams{
		loggedUser: loggedUser,
	})
	htmlgo.Fprint(w, baseHTML(content, loggedUser), r.Context())
}

func addToWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)

	pubkey := r.PostFormValue("pubkey")
	if pfx, value, err := nip19.Decode(pubkey); err == nil && pfx == "npub" {
		pubkey = value.(string)
	}

	if loggedUser != s.RelayPubkey && hasInvitedAtLeast(loggedUser, s.MaxInvitesPerPerson) {
		http.Error(w, fmt.Sprintf("cannot invite more than %d", s.MaxInvitesPerPerson), 403)
		return
	}

	if err := addToWhitelist(pubkey, loggedUser); err != nil {
		http.Error(w, "failed to add to whitelist: "+err.Error(), 500)
		return
	}

	content := inviteTreeComponent(r.Context(), "", loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

func removeFromWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)
	pubkey := r.PostFormValue("pubkey")
	if err := removeFromWhitelist(pubkey, loggedUser); err != nil {
		http.Error(w, "failed to remove from whitelist: "+err.Error(), 500)
		return
	}
	content := inviteTreeComponent(r.Context(), "", loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

// this deletes all events from users not in the relay anymore
func cleanupStuffFromExcludedUsersHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)
	if loggedUser != s.RelayPubkey {
		http.Error(w, "unauthorized, only the relay owner can do this", 403)
		return
	}

	oldLimit := db.MaxLimit
	db.MaxLimit = 999999
	ch, err := db.QueryEvents(r.Context(), nostr.Filter{Limit: db.MaxLimit})
	if err != nil {
		http.Error(w, "failed to query", 500)
		return
	}
	db.MaxLimit = oldLimit

	count := 0

	for evt := range ch {
		if isPublicKeyInWhitelist(evt.PubKey) {
			continue
		}

		if err := db.DeleteEvent(r.Context(), evt); err != nil {
			http.Error(w, fmt.Sprintf(
				"failed to delete %s: %s -- stopping, %d events were deleted before this error", evt, err, count), 500)
			return
		}
		count++
	}

	fmt.Fprintf(w, "deleted %d events", count)
}

func reportsViewerHandler(w http.ResponseWriter, r *http.Request) {
	events, err := db.QueryEvents(r.Context(), nostr.Filter{
		Kinds: []int{1984},
		Limit: 52,
	})
	if err != nil {
		http.Error(w, "failed to query reports: "+err.Error(), 500)
		return
	}

	content := reportsPageHTML(r.Context(), ReportsPageParams{
		reports:    events,
		loggedUser: getLoggedUser(r),
	})
	htmlgo.Fprint(w, content, r.Context())
}

func joubleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>pyramid</title>
    <script>
window.relayGroups = [{
  groupName: 'pyramid',
  relayUrls: [location.href.replace('http', 'ws').replace('/browse', '')],
  isActive: true,
}]
window.hideRelaySettings = true
</script>
    <script type="module" crossorigin src="https://unpkg.com/jouble/dist/index.js"></script>
    <link rel="stylesheet" crossorigin href="https://unpkg.com/jouble/dist/index.css">
  </head>

  <body>
    <div id="root"></div>
    <script src="https://unpkg.com/window.nostr.js/dist/window.nostr.js"></script>
  </body>
</html>
`)
}
