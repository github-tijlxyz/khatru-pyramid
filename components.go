package main

import (
	"context"

	sdk "github.com/nbd-wtf/go-nostr/sdk"
	. "github.com/theplant/htmlgo"
)

func inviteTreeComponent(ctx context.Context, inviter string, loggedUser string) HTMLComponent {
	children := make([]HTMLComponent, 0, len(whitelist)/2)
	for pubkey, invitedBy := range whitelist {
		if invitedBy == inviter {
			profile := sys.FetchProfileMetadata(ctx, pubkey)
			children = append(children, userRowComponent(ctx, profile, loggedUser))
		}
	}
	return Ul(children...)
}

func userRowComponent(ctx context.Context, profile sdk.ProfileMetadata, loggedUser string) HTMLComponent {
	button := Span("")
	if isAncestorOf(loggedUser, profile.PubKey) && loggedUser != "" {
		button = Button("remove").
			Class(buttonClass+" px-2 bg-red-100 hover:bg-red-300").
			Attr(
				"hx-post", "/remove-from-whitelist",
				"hx-trigger", "click",
				"hx-target", "#tree",
				"hx-vals", `{"pubkey": "`+profile.PubKey+`"}`,
			)
	}

	return Li(
		userNameComponent(profile),
		button,
		inviteTreeComponent(ctx, profile.PubKey, loggedUser),
	).Class("ml-6")
}

func userNameComponent(profile sdk.ProfileMetadata) HTMLComponent {
	return A().Href("https://nosta.me/" + profile.Npub()).Target("_blank").Children(
		Span(profile.ShortName()).Attr("title", profile.Npub()),
	).Class("font-mono py-1")
}
