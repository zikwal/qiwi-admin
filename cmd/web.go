// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/routers"
	"github.com/zhuharev/qiwi-admin/routers/auth"
	"github.com/zhuharev/qiwi-admin/routers/groups"
	"github.com/zhuharev/qiwi-admin/routers/wallets"
	macaron "gopkg.in/macaron.v1"
)

var (
	// CmdWeb starts web server
	CmdWeb = cli.Command{
		Name:   "web",
		Action: startWeb,
	}
)

func newMacaron() (m *macaron.Macaron) {

	err := routers.GlobalInit()
	if err != nil {
		log.Fatalln(err)
	}

	m = macaron.Classic()

	// html templates
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Layout: "layout",
	}))

	// sessions, auth, cookies
	m.Use(session.Sessioner(session.Options{
		CookieName:     "s",
		Provider:       "file",
		ProviderConfig: "data/sessions",
	}))

	// wrapped context
	m.Use(context.Contexter())

	return
}

func startWeb(ctx *cli.Context) {
	m := newMacaron()
	m.Get("/", auth.RedirectAutorized, routers.Index)

	m.Any("/auth", auth.RedirectAutorized, binding.Bind(models.AuthForm{}), auth.Auth)
	m.Any("/reg", auth.RedirectAutorized, binding.Bind(models.AuthForm{}), auth.Reg)
	m.Get("/logout", auth.Logout)

	m.Group("/groups", func() {
		m.Get("/:groupID", wallets.List)
	}, auth.MustAuthorized)

	m.Group("/wallets/", func() {
		m.Post("/create", wallets.Create)
		m.Get("/:id", wallets.Show)
	}, auth.MustAuthorized)

	m.Get("/dashboard", auth.MustAuthorized, groups.List)

	m.Run()
}
