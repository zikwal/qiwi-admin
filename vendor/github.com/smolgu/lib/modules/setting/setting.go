package setting

import (
	"github.com/smolgu/lib/modules/base"
	"github.com/zhuharev/menu"
	"gopkg.in/ini.v1"
)

var (
	AppVer string

	BuildHash string

	DataDir   string
	UploadDir string
	LogDir    string

	RunMode  string
	HostName string

	VkAccessToken string
	VkGroupID     int

	MainMenu menu.Menus

	Users struct {
		Secret        string
		DbDriver      string
		DbSetting     string
		KvPath        string
		AdminLogin    string
		AdminPassword string

		LogFile string
	}

	DbDriver  string
	DbSetting string

	ItemsInPage = 15
)

func init() {
	BuildHash = base.GetRandomString(6)
}

func NewContext(mode, configLocation string) {
	RunMode = mode

	iniFile, e := ini.Load(configLocation)
	if e != nil {
		panic(e)
	}
	iniFile.NameMapper = ini.TitleUnderscore

	as := iniFile.Section(mode)
	DataDir = as.Key("data_dir").String()

	e = iniFile.Section(mode + ".users").MapTo(&Users)
	if e != nil {
		panic(e)
	}

	VkAccessToken = as.Key("vk_access_token").String()
	VkGroupID = as.Key("vk_group_id").MustInt(0)

	LogDir = as.Key("log_dir").String()

	DbDriver = as.Key("db.driver").String()
	DbSetting = as.Key("db.setting").String()

	MainMenu, e = menu.NewFromFile("conf/menu/main.conf")
	if e != nil {
		panic(e)
	}
}
