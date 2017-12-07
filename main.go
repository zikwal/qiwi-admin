// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/zhuharev/qiwi-admin/cmd"
	"github.com/zhuharev/qiwi-admin/pkg/setting"
)

// AppVer current version of app
var AppVer = "0.0.9"

func init() {
	setting.AppVer = AppVer
}

func main() {
	app := &cli.App{
		Version: AppVer,
		Commands: []cli.Command{
			cmd.CmdWeb,
		},
	}
	app.Run(os.Args)
}
