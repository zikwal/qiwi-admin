// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package exchange

import (
	"github.com/zhuharev/qiwi-admin/models"
	"github.com/zhuharev/qiwi-admin/pkg/context"
)

// Index shows exchange index page
func Index(ctx *context.Context) {
	groups, err := models.GetUserGroupsWithCounters(ctx.User.ID)
	if ctx.HasError(err, "[sql] get groups") {
		return
	}

	ctx.Data["groups"] = groups
	ctx.HTML(200, "exchange/index")
}
