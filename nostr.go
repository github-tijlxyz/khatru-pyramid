package main

import (
	"context"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type SimpleUserInfo struct {
	Npub    string
	Name    string
	Picture string
	Time    time.Time
}

var (
	userInfoCache = make(map[string]SimpleUserInfo)
)

func getUserInfo(ctx context.Context, hexpubkey string) SimpleUserInfo {
	// check if in cache
	v, o := userInfoCache[hexpubkey]
	if o {
		if !(time.Since(v.Time) > 2*time.Hour) { // use cache for 2 hours
			return v
		}
	}

	npub, _ := nip19.EncodePublicKey(hexpubkey)
	var name string = string(npub)
	var picture string = ""

	evts, err := db.QueryEvents(ctx, nostr.Filter{
		Authors: []string{hexpubkey},
		Kinds:   []int{0},
		Limit:   1,
	})
	if err != nil {
		return SimpleUserInfo{npub, name, picture, time.Now()}
	}
	for ev := range evts {
		name, picture = getProfileInfoFromJson(ev.Content)
	}

	userInfoCache[hexpubkey] = SimpleUserInfo{npub, name, picture, time.Now()}
	return SimpleUserInfo{npub, name, picture, time.Now()}
}
