// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mail

import (
    "fmt"
    "net"
    "net/mail"
    "net/smtp"
    "os"
    "strings"
    "prosnav.com/wxserver/modules/setting"
    "crypto/tls"
    "github.com/gogits/gogs/modules/log"
)

type Message struct {
    To      []string
    From    string
    Subject string
    Body    string
    Type    string
    Massive bool
    Info    string
}

// create mail content
func (m Message) Content() string {
    // set mail type
    contentType := "text/plain; charset=UTF-8"
    if m.Type == "html" {
        contentType = "text/html; charset=UTF-8"
    }

    // create mail content
    content := "From: " + m.From + "\r\nSubject: " + m.Subject + "\r\nContent-Type: " + contentType + "\r\n\r\n" + m.Body
    return content
}

var mailQueue chan *Message


// Direct Send mail message
func Send(msg *Message) (int, error) {
    log.Info("Sending mails to: %s", strings.Join(msg.To, "; "))
    // get message body
    content := msg.Content()

    if len(msg.To) == 0 {
        return 0, fmt.Errorf("empty receive emails")
    } else if len(msg.Body) == 0 {
        return 0, fmt.Errorf("empty email body")
    }
    if msg.Massive {
        // send mail to multiple emails one by one
        num := 0
        for _, to := range msg.To {
            body := []byte("To: " + to + "\r\n" + content)
            err := sendMail(setting.MailService, []string{to}, body)
            if err != nil {
                return num, err
            }
            num++
        }
        return num, nil
    } else {
        body := []byte("To: " + strings.Join(msg.To, ";") + "\r\n" + content)

        // send to multiple emails in one message
        err := sendMail(setting.MailService, msg.To, body)
        if err != nil {
            return 0, err
        } else {
            return 1, nil
        }
    }

}

func processmailQueue() {
    for {
        select {
        case msg := <- mailQueue:
            num, err := Send(msg)
            tos := strings.Join(msg.To, "; ")
            info := ""
            if err != nil {
                if len(msg.Info) > 0 {
                    info = ", info: " + msg.Info
                }
                log.Error(4, fmt.Sprintf("Async sent email %d succeed, not send emails: %s%s err: %s", num, tos, info, err))
            } else {
                log.Info(fmt.Sprintf("Async sent email %d succeed, sent emails: %s%s", num, tos, info))
            }
        }
    }
}


func NewMailerContext() {
    mailQueue= make(chan *Message, setting.Cfg.Section("mail").Key("SEND_BUFFER_LEN").MustInt(10))
    go processmailQueue()
}

// Async Send mail message
func SendAsync(msg *Message) {
    go func() {
        mailQueue <- msg
    }()

}
// sendMail allows mail with self-signed certificates.
func sendMail(settings *setting.Mailer, recipients []string, msgContent []byte) error {
    host, port, err := net.SplitHostPort(settings.Host)
    if err != nil {
        return err
    }

    tlsconfig := &tls.Config{
        InsecureSkipVerify: settings.SkipVerify,
        ServerName:         host,
    }

    if settings.UseCertificate {
        cert, err := tls.LoadX509KeyPair(settings.CertFile, settings.KeyFile)
        if err != nil {
            return err
        }
        tlsconfig.Certificates = []tls.Certificate{cert}
    }

    conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
    if err != nil {
        return err
    }
    defer conn.Close()

    isSecureConn := false
    // Start TLS directly if the port ends with 465 (SMTPS protocol)
    if strings.HasSuffix(port, "465") {
        conn = tls.Client(conn, tlsconfig)
        isSecureConn = true
    }

    client, err := smtp.NewClient(conn, host)
    if err != nil {
        return err
    }

    hostname, err := os.Hostname()
    if err != nil {
        return err
    }

    if err = client.Hello(hostname); err != nil {
        return err
    }

    // If not using SMTPS, alway use STARTTLS if available
    hasStartTLS, _ := client.Extension("STARTTLS")
    if !isSecureConn && hasStartTLS {
        if err = client.StartTLS(tlsconfig); err != nil {
            return err
        }
    }

    canAuth, options := client.Extension("AUTH")

    if canAuth && len(settings.User) > 0 {
        var auth smtp.Auth

        if strings.Contains(options, "CRAM-MD5") {
            auth = smtp.CRAMMD5Auth(settings.User, settings.Passwd)
        } else if strings.Contains(options, "PLAIN") {
            auth = smtp.PlainAuth("", settings.User, settings.Passwd, host)
        }

        if auth != nil {
            if err = client.Auth(auth); err != nil {
                return err
            }
        }
    }

    if fromAddress, err := mail.ParseAddress(settings.From); err != nil {
        return err
    } else {
        if err = client.Mail(fromAddress.Address); err != nil {
            return err
        }
    }

    for _, rec := range recipients {
        if err = client.Rcpt(rec); err != nil {
            return err
        }
    }

    w, err := client.Data()
    if err != nil {
        return err
    }
    if _, err = w.Write([]byte(msgContent)); err != nil {
        return err
    }

    if err = w.Close(); err != nil {
        return err
    }

    return client.Quit()
}




// Create html mail message
func NewHtmlMessage(To []string, From, Subject, Body string) Message {
        return Message{
                    To:      To,
                    From:    From,
                    Subject: Subject,
                    Body:    Body,
                    Type:    "html",
                }
}
