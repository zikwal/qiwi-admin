// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"path/filepath"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/setting"
	"github.com/zhuharev/qiwi-admin/routers"
	"github.com/zhuharev/qiwi-admin/routers/auth"
	"github.com/zhuharev/qiwi-admin/routers/groups"
	"github.com/zhuharev/qiwi-admin/routers/transfers"
	"github.com/zhuharev/qiwi-admin/routers/wallets"
	macaron "gopkg.in/macaron.v1"
)

var (
	// CmdWeb starts web server
	CmdWeb = cli.Command{
		Name:   "web",
		Action: startWeb,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "data-dir",
				Value: "./data",
			},
		},
	}
)

func newMacaron(ctx *cli.Context) (m *macaron.Macaron) {

	setting.App.DataDir = ctx.String("data-dir")

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
	sessionLifeTime := 24 * 60 * 60 * 365
	m.Use(session.Sessioner(session.Options{
		CookieName:     "s",
		CookieLifeTime: sessionLifeTime,
		Provider:       "file",
		ProviderConfig: filepath.Join(setting.App.DataDir, "sessions"),
	}))

	// wrapped context
	m.Use(context.Contexter())

	return
}

func startWeb(ctx *cli.Context) {
	m := newMacaron(ctx)
	m.Get("/", auth.RedirectAutorized, routers.Index)

	m.Any("/auth", auth.RedirectAutorized, binding.Bind(models.AuthForm{}), auth.Auth)
	m.Any("/reg", auth.RedirectAutorized, binding.Bind(models.AuthForm{}), auth.Reg)
	m.Get("/logout", auth.Logout)

	m.Group("/groups", func() {
		m.Get("/:groupID", groups.Get)
	}, auth.MustAuthorized)

	m.Group("/wallets/", func() {
		m.Post("/create", wallets.Create)
		m.Get("/:id", wallets.Show)
		m.Any("/:id/setting", wallets.Setting)
	}, auth.MustAuthorized)

	m.Get("/dashboard", auth.MustAuthorized, groups.List)

	m.Any("/transfer", transfers.Transfer)

	m.Run()
}
