package main

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

func rejectEventsFromUsersNotInWhitelist(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
	if isPublicKeyInWhitelist(event.PubKey) {
		return false, ""
	}
	return true, "not authorized"
}
