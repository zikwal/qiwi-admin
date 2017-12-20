// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/setting"
	"github.com/zhuharev/qiwi-admin/routers"
	"github.com/zhuharev/qiwi-admin/routers/accounts"
	"github.com/zhuharev/qiwi-admin/routers/apps"
	"github.com/zhuharev/qiwi-admin/routers/auth"
	"github.com/zhuharev/qiwi-admin/routers/exchange"
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
			cli.BoolFlag{
				Name: "prod",
			},
		},
	}
)

func newMacaron(ctx *cli.Context) (m *macaron.Macaron) {

	setting.App.DataDir = ctx.String("data-dir")
	setting.Verbose = ctx.GlobalBool("verbose")
	if setting.Verbose {
		log.Println("Запуск в режиме расширенного логгирования")
	}

	err := routers.GlobalInit(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	m = macaron.Classic()
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["startTime"] = time.Now()
	})

	// html templates
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Layout: "layout",
		Funcs: []template.FuncMap{map[string]interface{}{
			"LoadTimes": func(startTime time.Time) string {
				return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "мс"
			},
		}},
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
	m.Get("/logout", auth.MustAuthorized, auth.Logout)

	m.Group("/groups", func() {
		m.Get("/:groupID", groups.Get)
		m.Post("/create", groups.Create)
		m.Any("/:groupID/setting", groups.Setting)
		m.Get("/:id/delete", groups.Delete)
	}, auth.MustAuthorized)

	m.Group("/wallets/", func() {
		m.Post("/create", wallets.Create)
		m.Get("/:id", wallets.Show)
		m.Any("/:id/setting", wallets.Setting)
		m.Get("/:id/delete", wallets.Delete)
	}, auth.MustAuthorized)

	m.Get("/dashboard", auth.MustAuthorized, groups.List)
	m.Group("/dashboard", func() {
		m.Get("/apps", apps.Apps)
		m.Post("/apps/:appID/webhook", apps.SaveWebHookURL)
		m.Post("/apps/create", apps.Create)
		m.Post("/apps/test", apps.Test)
	}, auth.MustAuthorized)

	m.Group("/exchange", func() {
		m.Get("/", exchange.Index)
		m.Post("/wallets/:id", exchange.Wallet)
		m.Get("/trades", exchange.Trades)
	}, auth.MustAuthorized)

	m.Any("/transfer", transfers.Transfer, auth.MustAuthorized)
	m.Post("/transfer/groups/:groupID", transfers.TransferFromGroup, auth.MustAuthorized)

	m.Get("/setting", accounts.Setting, auth.MustAuthorized)
	m.Post("/account/setting", auth.MustAuthorized, accounts.SaveSetting)

	m.Group("/users", func() {
		m.Post("/create", auth.MustAuthorized, binding.Bind(models.AuthForm{}), accounts.Create)
	})

	m.Run()
}
