package main

import (
	"context"

	"github.com/fiatjaf/khatru/plugins"
	"github.com/nbd-wtf/go-nostr"
)

func rejectEventsFromUsersNotInWhitelist(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
	if isPublicKeyInWhitelist(event.PubKey) {
		return false, ""
	}
	if event.Kind == 1985 {
		// we accept reports from anyone (will filter them for relevance in the next function)
		return false, ""
	}
	return true, "not authorized"
}

var restrictToKinds = plugins.RestrictToSpecifiedKinds(
	0, 1, 3, 5, 6, 8, 16, 1063, 1985, 9735, 10000, 10001, 10002, 30008, 30009, 30311, 31922, 31923, 31924, 31925)

func validateAndFilterReports(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
	if event.Kind == 1985 {
		if e := event.Tags.GetFirst([]string{"e", ""}); e != nil {
			// event report: check if the target event is here
			res, _ := sys.StoreRelay().QuerySync(ctx, nostr.Filter{IDs: []string{(*e)[1]}})
			if len(res) == 0 {
				return true, "we don't know anything about the target event"
			}
		} else if p := event.Tags.GetFirst([]string{"p", ""}); p != nil {
			// pubkey report
			if !isPublicKeyInWhitelist((*p)[1]) {
				return true, "target pubkey is not a user of this relay"
			}
		} else {
			return true, "invalid report"
		}
	}

	return false, ""
}
