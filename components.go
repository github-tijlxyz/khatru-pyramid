package main

import (
	"context"

	sdk "github.com/nbd-wtf/nostr-sdk"
	. "github.com/theplant/htmlgo"
)

func inviteTreeComponent(ctx context.Context, inviter string, loggedUser string) HTMLComponent {
	children := make([]HTMLComponent, 0, len(whitelist)/2)
	for pubkey, invitedBy := range whitelist {
		if invitedBy == inviter {
			profile := fetchAndStoreProfile(ctx, pubkey)
			children = append(children, userRowComponent(ctx, profile, loggedUser))
		}
	}
	return Ul(children...)
}

func userRowComponent(ctx context.Context, profile sdk.ProfileMetadata, loggedUser string) HTMLComponent {
	button := Span("")
	if isAncestorOf(loggedUser, profile.PubKey) && loggedUser != profile.PubKey {
		button = Button("remove").
			Class(buttonClass+" px-2 bg-red-100 hover:bg-red-300").
			Attr(
				"hx-post", "/remove-from-whitelist",
				"hx-trigger", "click",
				"hx-target", "#tree",
				"hx-vals", `{"profile.pubkey": "`+profile.PubKey+`"}`,
			)
	}

	return Li(
		A().Href("nostr:"+profile.Npub()).Children(
			Span(profile.ShortName()).Attr(
				"npub", profile.Npub(),
				"name", profile.ShortName(),
				"_", `
on mouseenter set my innerText to @npub then hide the next <button />
on mouseleave set my innerText to @name then show the next <button />`,
			),
		).Class("font-mono py-1"),
		button,
		inviteTreeComponent(ctx, profile.PubKey, loggedUser),
	).Class("ml-6")
}
