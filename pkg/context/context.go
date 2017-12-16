// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"github.com/go-macaron/session"
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/log"
	"github.com/zhuharev/qiwi-admin/pkg/setting"
	macaron "gopkg.in/macaron.v1"
)

// Context will be used in routers
type Context struct {
	*macaron.Context

	Flash   *session.Flash
	Session session.Store
	User    *models.User
}

// Contexter wrap macaron.Context
func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Flash:   f,
			Session: sess,
		}

		ctx.Data["AppVer"] = setting.AppVer

		if userIface := sess.Get("user_id"); userIface != nil {
			if userID, ok := userIface.(uint); ok {
				user, err := models.GetUser(userID)
				if err != nil {
					sess.Delete("user_id")
					c.Redirect("/auth")
					return
				}
				ctx.User = user
				c.Data["User"] = user
			}
		}

		c.Map(ctx)
	}
}

// Autorized just hellper
func (ctx *Context) Autorized() bool {
	return ctx.User != nil
}

// HTML overwrite macaron.HTML method
func (ctx *Context) HTML(code int, tmplName string, other ...interface{}) {
	layoutName := "layout"
	if !ctx.Autorized() {
		layoutName = "unauth-layout"
	}
	ctx.Context.HTML(code, tmplName, ctx.Data, macaron.HTMLOptions{Layout: layoutName})
}

// HasError check passed err and write resposne if err!=nil
func (ctx *Context) HasError(err error, args ...interface{}) bool {
	if err != nil {
		log.Trace(err, args...)
		ctx.Flash.Error(err.Error())
		if ctx.User != nil {
			ctx.Redirect("/dashboard")
		} else {
			ctx.Redirect("/")
		}

		//ctx.Error(200, err.Error())
		return true
	}
	return false
}
