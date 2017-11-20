package notifier

import "github.com/zhuharev/qiwi-admin/models"

type NotifyCmd struct {
	URL   string
	Delay Delay
	Txn   models.Txn
}

// NewCmd returns new cmd with zero delay
func NewCmd(url string, txn models.Txn) *NotifyCmd {
	return &NotifyCmd{
		URL:   url,
		Txn:   txn,
		Delay: NoDelay,
	}
}

func (nc NotifyCmd) NextDelay() Delay {
	switch nc.Delay {
	case NoDelay:
		return OneHour
	case OneHour:
		return SixHours
	default:
		return NoDelay
	}
}
