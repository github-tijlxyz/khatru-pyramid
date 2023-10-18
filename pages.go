package main

import (
	"context"

	. "github.com/theplant/htmlgo"
)

func baseHTML(inside HTMLComponent) HTMLComponent {
	navItemClass := "text-gray-300 hover:bg-gray-700 hover:text-white rounded-md px-3 py-2 font-medium"

	return HTML(
		Head(
			Meta().Charset("utf-8"),
			Meta().Name("viewport").Content("width=device-width, initial-scale=1"),
			Title(s.RelayName),
			Script("").Src("https://cdn.tailwindcss.com"),
		),
		Body(
			Div(
				H1(s.RelayName).Class("font-bold text-2xl"),
				P().Text(s.RelayDescription).Class("text-lg"),
			).Class("mx-auto my-6 text-center"),
			Nav(
				A().Text("information").Href("/").Class(navItemClass),
				A().Text("invite tree").Href("/users").Class(navItemClass),
				A().Text("reports").Href("/reports").Class(navItemClass),
			).Class("flex flex-1 items-center justify-center"),
			Div(inside).Class("m-4"),
		).Class("bg-gray-800 mx-4 my-6 text-white"),
	)
}

type HomePageParams struct {
	RelayOwnerInfo SimpleUserInfo
}

func homePageHTML(ctx context.Context, params HomePageParams) HTMLComponent {
	contact := Div()
	if s.RelayContact != "" {
		contact = Div().Text("alternative contact: " + s.RelayContact)
	}

	description := Div()
	if s.RelayDescription != "" {
		description = Div().Text("description: " + s.RelayDescription)
	}

	return Div(
		Div().Text("name: "+s.RelayName),
		description,
		contact,
		Div(
			Text("relay master: "),
			A().Text(params.RelayOwnerInfo.Name).Href("nostr:"+params.RelayOwnerInfo.Npub),
		),
		Br(),
		Div(
			Text("this relay uses"),
			A().Target("_blank").Href("https://github.com/github-tijlxyz/khatru-invite").Text("Khatru Invite"),
			Text(" which is built with "),
			A().Target("_blank").Href("https://github.com/fiatjaf/khatru").Text("Khatru"),
		),
	)
}

type InviteTreePageParams struct{}

func inviteTreePageHTML(ctx context.Context, params InviteTreePageParams) HTMLComponent {
	return Div(
		Input("").Type("text").Placeholder("npub1..."),
		Button("invite"),
		buildInviteTree(ctx, ""),
	)
}

func buildInviteTree(ctx context.Context, invitedBy string) HTMLComponent {
	tree := Ul()
	for _, entry := range whitelist {
		if entry.InvitedBy == invitedBy {
			user := getUserInfo(ctx, entry.PublicKey)
			tree = tree.Children(
				Li(
					A().Href("nostr:"+user.Npub).Text(user.Name),
					A().Text("remove"),
					buildInviteTree(ctx, entry.PublicKey),
				),
			)
		}
	}
	return tree
}
