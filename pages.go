package main

import (
	"context"

	. "github.com/theplant/htmlgo"
)

const buttonClass = "rounded-md bg-white p-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"

func baseHTML(inside HTMLComponent) HTMLComponent {
	navItemClass := "text-gray-600 hover:bg-gray-200 rounded-md px-3 py-2 font-medium"

	return HTML(
		Head(
			Meta().Charset("utf-8"),
			Meta().Name("viewport").Content("width=device-width, initial-scale=1"),
			Title(s.RelayName),
			Script("").Src("https://cdn.tailwindcss.com"),
			Script("").Src("https://unpkg.com/htmx.org@1.9.6"),
			Script("").Src("https://unpkg.com/hyperscript.org@0.9.12"),
		),
		Body(
			Div(
				H1(s.RelayName).Class("font-bold text-2xl"),
				P().Text(s.RelayDescription).Class("text-lg"),
			).Class("mx-auto my-6 text-center"),
			Nav(
				A().Text("information").Href("/").Class(navItemClass).Attr("hx-boost", "true", "hx-target", "main", "hx-select", "main"),
				A().Text("invite tree").Href("/users").Class(navItemClass).Attr("hx-boost", "true", "hx-target", "main", "hx-select", "main"),
				A().Text("reports").Href("/reports").Class(navItemClass).Attr("hx-boost", "true", "hx-target", "main", "hx-select", "main"),
			).Class("flex flex-1 items-center justify-center"),
			Main(inside).Class("m-4"),
		).Class("mx-4 my-6"),
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
	return Form(
		Input("pubkey").Type("text").Placeholder("npub1...").Class("w-96 rounded-md border-0 p-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600"),
		Button("invite").Class(buttonClass),
		Div(
			buildInviteTree(ctx, s.RelayPubkey),
		).Id("tree").Class("mt-3"),
	).Attr("hx-post", "/add-to-whitelist", "hx-trigger", "submit", "hx-target", "#tree")
}

func buildInviteTree(ctx context.Context, invitedBy string) HTMLComponent {
	children := make([]HTMLComponent, 0, len(whitelist))
	for _, entry := range whitelist {
		if entry.InvitedBy == invitedBy {
			user := getUserInfo(ctx, entry.PublicKey)
			children = append(children,
				Li(
					A().Href("nostr:"+user.Npub).Text(user.Name).Class("font-mono"),
					Button("remove").Class(buttonClass).Attr("hx-post", "/remove-from-whitelist", "hx-trigger", "click", "hx-target", "#tree"),
					buildInviteTree(ctx, entry.PublicKey),
				).Class("ml-3"),
			)
		}
	}
	return Ul(children...)
}
