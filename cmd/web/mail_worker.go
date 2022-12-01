package main

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/bernie-pham/agodafake/internal/model"
	mail "github.com/xhit/go-simple-mail/v2"
)

var mailTemplate = "./templates/basic.email.html"

// this function will be listened for any sending mail request from program mainstream.
func listenForMail() {
	clientMail := newMailClient()
	go func() {
		for {
			mailReq := <-app.MailChan
			sendMail(mailReq, clientMail)
		}
	}()
}

func sendMail(m model.MailData, clientMail *mail.SMTPClient) {
	var fromEmail string
	if len(m.From) > 0 {
		fromEmail = m.From
	} else {
		fromEmail = "service@bookingroom.com"
	}
	email := mail.NewMSG()
	email.SetFrom(fromEmail).
		AddTo(m.To).
		SetSubject(m.Subject)

	mailTemplate, err := ioutil.ReadFile(mailTemplate)
	if err != nil {
		errLog.Println(err)
	}
	mailstring := string(mailTemplate)
	mailBody := strings.Replace(mailstring, "[%body%]", m.Content, 1)
	email.SetBody(mail.TextHTML, mailBody)
	err = email.Send(clientMail)
	if err != nil {
		errLog.Println(err)
	}
}

func newMailClient() *mail.SMTPClient {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errLog.Println(err)
	}
	return client
}
