// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package semail
 
import (
 	"net/mail"
	"encoding/base64"
	"net/smtp"
	"fmt"
	"strings"
)

type Email struct {
	host string
	from mail.Address
	passwd string
}

func NewEmail(host, from, passwd, fromName string) *Email {
	e := &Email {
		host: host,
		from: mail.Address{fromName, from},
		passwd: passwd,
	}

	return e
}

func (m *Email) auth() smtp.Auth {
	return smtp.PlainAuth(
		"",
		m.from.Address,
		m.passwd,
		m.host,
	)
}

func (m *Email) encodeRFC2047(s string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{s, ""}
	return strings.Trim(addr.String(), " <>")
}

func (m *Email) rmdup(addrs []string) []string {

	ck := make(map[string]bool)
	for _, a := range addrs {
		if len(a) > 0 {
			ck[a] = true
		}
	}

	rv := make([]string, 0)

	for k, _ := range ck {
		rv = append(rv, k)
	}

	return rv
}


func (m *Email) Html(to, cc []string, subject, body string) error {
	to = m.rmdup(to)
	cc = m.rmdup(cc)

	header := make(map[string]string)
	header["From"] = m.from.String()
	header["To"] = strings.Join(to, ";")
	header["Cc"] = strings.Join(cc, ";")

	header["Subject"] = m.encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	//header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"


	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

 
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one ste
	tolist := append(to, cc...)
	tolist = m.rmdup(tolist)
	//fmt.Println(tolist)
	err := smtp.SendMail(
		m.host + ":25",
		m.auth(),
		m.from.Address,
		tolist,
		[]byte(message),
	)


	return err
}


