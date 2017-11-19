// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package transfers

import (
	"github.com/zhuharev/qiwi-admin/pkg/context"
	"github.com/zhuharev/qiwi-admin/pkg/qiwi"
)

// TransferFromGroup transfer money
func TransferFromGroup(ctx *context.Context) {

	var (
		//from = ctx.Query("from")
		to      = ctx.Query("to")
		groupID = uint(ctx.ParamsInt(":groupID"))
	)

	errs := qiwi.TransferFromGroup(groupID, ctx.User.ID, to, 300)
	if errs != nil {
		ctx.Data["errs"] = errs
		ctx.HTML(200, "transfers/errors")
		return
	}

	ctx.Redirect("/dashboard")
}
