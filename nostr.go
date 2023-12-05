package main

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
	sdk "github.com/nbd-wtf/nostr-sdk"
	cache_memory "github.com/nbd-wtf/nostr-sdk/cache/memory"
)

var sys = sdk.System{
	Pool:             nostr.NewSimplePool(context.Background()),
	RelaysCache:      cache_memory.New32[[]sdk.Relay](1000),
	MetadataCache:    cache_memory.New32[sdk.ProfileMetadata](1000),
	FollowsCache:     cache_memory.New32[[]sdk.Follow](1),
	RelayListRelays:  []string{"wss://purplepag.es", "wss://relay.nostr.band"},
	FollowListRelays: []string{"wss://public.relaying.io", "wss://nos.lol"},
	MetadataRelays:   []string{"wss://nostr-pub.wellorder.net", "wss://purplepag.es", "wss://relay.noswhere.com"},
	Store:            &db,
}
