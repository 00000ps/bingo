package utils

import (
	"bingo/pkg/log"
	"fmt"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

func sendEmail() {
	ch := make(chan *email.Email, 10)
	p, err := email.NewPool(
		"smtp.126.com:25",
		4,
		smtp.PlainAuth("", "leedarjun@126.com", "358942617ldj", "smtp.126.com"),
	)

	if err != nil {
		log.Fatal("failed to create pool:", err)
	}

	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}

	for i := 0; i < 10; i++ {
		e := email.NewEmail()
		e.From = "dj <leedarjun@126.com>"
		e.To = []string{"935653229@qq.com"}
		e.Subject = "Awesome web"
		e.Text = []byte(fmt.Sprintf("Awesome Web %d", i+1))
		ch <- e
	}

	close(ch)
	wg.Wait()
}
