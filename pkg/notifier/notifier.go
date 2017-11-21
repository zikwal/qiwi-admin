package notifier

import (
	"time"

	"github.com/go-resty/resty"
	"github.com/zhuharev/qiwi-admin/models"
)

var (
	noti = &notifier{
		cmdChan: make(chan *NotifyCmd, 100),
		client:  resty.New(),
	}
)

// Notify helper for global usage
func Notify(cmd *NotifyCmd) {
	noti.Notify(cmd)
}

// NotifyTxn get account apps and notify all apps about txn
func NotifyTxn(wallet *models.Wallet, txn models.Txn) (err error) {
	apps, err := models.Apps.List(wallet.OwnerID)
	if err != nil {
		return
	}
	for _, app := range apps {
		cmd := NewCmd(app.WebHookURL, txn)
		Notify(cmd)
	}
	return
}

type Delay time.Duration

const (
	// NoDelay instant
	NoDelay  Delay = 0
	OneHour        = Delay(time.Hour)
	SixHours       = OneHour * 6
)

type Notifier interface {
	Notify(*NotifyCmd) error
	Start() error
}

func NewContext() (err error) {
	return noti.Start()
}
