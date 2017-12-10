// Copyright 2017 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	confFile = "conf/app.yaml"

	//DevMode if true, run in dev mode
	DevMode bool
	// Verbose enable debug output
	Verbose bool
	// AppVer current app version
	AppVer string
	// App main config
	App struct {
		DataDir string `yaml:"data_dir"`
		// Db struct {
		// 	Driver string
		// 	Config string
		// }

		Reg struct {
			Disabled bool
		}
	}
)

// NewContext create new context
func NewContext(ops ...func()) (err error) {

	for _, v := range ops {
		v()
	}

	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &App)
	if err != nil {
		return
	}

	// iniFile, err = ini.Load(confFile)
	// if err != nil {
	// 	return
	// }
	// iniFile.NameMapper = mapper
	// err = iniFile.MapTo(&App)

	return
}

func CustomLocation(path string) func() {
	return func() {
		confFile = path
	}
}
