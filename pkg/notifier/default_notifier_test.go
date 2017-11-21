package notifier

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/zhuharev/qiwi-admin/models"
)

func startTestServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Printf("Method not a post %s", r.Method)
			return
		}

		bts, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while reading request body %s", err)
			return
		}
		log.Printf("%s", bts)
	}

	http.HandleFunc("/handler", handler)

	log.Fatal(http.ListenAndServe(":22228", nil))
}

// TestDefaultNotifier make http post request
func TestDefaultNotifier(t *testing.T) {
	go startTestServer()
	err := NewContext()
	if err != nil {
		log.Fatalln(err)
	}
	txn := models.Txn{
		QiwiTxnID: 228,
	}

	cmd := NewCmd("http://localhost:22228/handler", txn)
	Notify(cmd)

	time.Sleep(2 * time.Second)

}
