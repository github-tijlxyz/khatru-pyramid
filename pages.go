package main

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
	. "github.com/theplant/htmlgo"
)

const buttonClass = "rounded-md text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300"

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
				A().Text("invite tree").Href("/").Class(navItemClass).Attr("hx-boost", "true", "hx-target", "main", "hx-select", "main"),
				A().Text("browse").Href("/browse").Class(navItemClass),
				A().Text("reports").Href("/reports").Class(navItemClass).Attr("hx-boost", "true", "hx-target", "main", "hx-select", "main"),
				A().Text("").Href("#").Class(navItemClass).
					Attr("_", `
on click if my innerText is equal to "login" get window.nostr.signEvent({created_at: Math.round(Date.now()/1000), kind: 27235, tags: [['domain', "`+s.Domain+`"]], content: ''}) then get JSON.stringify(it) then set cookies['nip98'] to it otherwise call cookies.clear('nip98') end then call location.reload()

on load get cookies['nip98'] then if it is undefined set my innerText to "login" otherwise set my innerText to "logout"`),
			).Class("flex flex-1 items-center justify-center"),
			Main(inside).Class("m-4"),
			P(
				Text("powered by "),
				A().Href("https://github.com/github-tijlxyz/khatru-pyramid").Text("khatru-pyramid").Class("hover:underline cursor-pointer text-blue-500"),
			).Class("text-end my-4 text-sm"),
		).Class("my-6 mx-auto max-w-min min-w-96"),
	)
}

type InviteTreePageParams struct {
	loggedUser string
}

func inviteTreePageHTML(ctx context.Context, params InviteTreePageParams) HTMLComponent {
	inviteForm := Div()

	if params.loggedUser != "" && (params.loggedUser == s.RelayPubkey || !hasInvitedAtLeast(params.loggedUser, s.MaxInvitesPerPerson)) {
		inviteForm = Form(
			Input("pubkey").Type("text").Placeholder("npub1...").Class("w-96 rounded-md border-0 p-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600"),
			Button("invite").Class(buttonClass+" ml-2 p-2 bg-white hover:bg-gray-50"),
		).Attr(
			"hx-post", "/add-to-whitelist",
			"hx-trigger", "submit",
			"hx-target", "#tree",
			"_", "on htmx:afterRequest(elt, successful) if successful and elt is I call I.reset()",
		).Class("flex")
	}

	return Div(
		inviteForm,
		Div(
			inviteTreeComponent(ctx, "", params.loggedUser),
		).Id("tree").Class("mt-3"),
	)
}

type ReportsPageParams struct {
	reports    chan *nostr.Event
	loggedUser string
}

func reportsPageHTML(ctx context.Context, params ReportsPageParams) HTMLComponent {
	items := make([]HTMLComponent, 0, 52)
	for report := range params.reports {
		var primaryType string
		var secondaryType string
		var relatedContent HTMLComponent

		if e := report.Tags.GetFirst([]string{"e", ""}); e != nil {
			// event report
			res, _ := sys.StoreRelay.QuerySync(ctx, nostr.Filter{IDs: []string{(*e)[1]}})
			if len(res) == 0 {
				sys.Store.DeleteEvent(ctx, report)
				continue
			}

			if len(*e) >= 3 {
				primaryType = (*e)[2]
			}

			relatedEvent := res[0]
			relatedContent = Div(
				Text("event reported: "),
				Div().Text(relatedEvent.String()).Class("text-mono"),
			)
		} else if p := report.Tags.GetFirst([]string{"p", ""}); p != nil {
			// pubkey report
			if !isPublicKeyInWhitelist((*p)[1]) {
				sys.Store.DeleteEvent(ctx, report)
				continue
			}

			if len(*p) >= 3 {
				primaryType = (*p)[2]
			}

			relatedProfile := sys.FetchProfileMetadata(ctx, (*p)[1])
			relatedContent = Div(
				Text("profile reported: "),
				userNameComponent(relatedProfile),
			)
		} else {
			continue
		}

		reporter := sys.FetchProfileMetadata(ctx, report.PubKey)
		report := Div(
			Div(Span(primaryType).Class("font-semibold"), Text(" report")).Class("font-lg"),
			Div().Text(secondaryType),
			Div(Text("by "), userNameComponent(reporter)),
			Div().Text(report.Content).Class("p-3"),
			relatedContent,
		)

		items = append(items, report)
	}
	return baseHTML(
		Div(
			H1("reports received").Class("text-xl p-4"),
			Div(items...),
		),
	)
}
