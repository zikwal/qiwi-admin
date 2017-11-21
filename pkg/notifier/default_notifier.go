package notifier

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty"
	"github.com/smolgu/lib/modules/setting"
)

type notifier struct {
	cmdChan chan *NotifyCmd
	client  *resty.Client
}

func (n *notifier) Start() (err error) {
	go n.run()
	return
}

func (n *notifier) Notify(cmd *NotifyCmd) {
	log.Printf("Send webhook  to %s\n", cmd.URL)
	n.cmdChan <- cmd
}

func (n *notifier) run() {
	for {
		select {
		case cmd := <-n.cmdChan:
			if cmd.Delay == NoDelay {
				err := n.hook(cmd)
				if err != nil {
					cmd.Delay = cmd.NextDelay()
					if cmd.Delay != NoDelay {
						n.Notify(cmd)
					}
				}
			} else {
				go func() {
					time.Sleep(time.Duration(cmd.Delay))
					err := n.hook(cmd)
					if err != nil {
						cmd.Delay = cmd.NextDelay()
						if cmd.Delay != NoDelay {
							n.Notify(cmd)
						}
					}
				}()
			}
		}
	}
}

func (n *notifier) hook(cmd *NotifyCmd) (err error) {
	resp, err := n.client.R().
		SetHeader("User-Agent", "Go-Qiwi-Admin/"+setting.AppVer).
		SetBody(cmd.Txn).Post(cmd.URL)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("need resend webhook, status code %d", resp.StatusCode())
	}
	return
}
