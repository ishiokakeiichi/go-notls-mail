package go-notls-mail

import (
	"errors"
	"net/smtp"
	"net/http"
	"fmt"
	"strings"
	"bytes"
	"encoding/base64"
)

// Login Auth
type smtpLoginAuth struct {
	username, password string
}
func LoginAuth(username, password string) smtp.Auth {
	return &smtpLoginAuth{username, password}
}
func (a *smtpLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}
func (a *smtpLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

type Recipients struct{
	To []string
	Cc []string
	Bcc []string
}

// Make Maile Message
func utf8Split(utf8string string, length int) []string {
	var resultString []string
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}

func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}
type MailAttachement struct{
	Data bytes.Buffer
	Filename string
}

func MakeMessage(from string, recipents Recipients, subject string, body string, delimeter string, attachment *MailAttachement) []byte {
	var bufBody bytes.Buffer
	if len(delimeter) == 0{
		delimeter = "**=SimpleMail Delimeter"
	}
	bufBody.WriteString(body)


	var header bytes.Buffer
	header.WriteString("From: " + from + "\r\n")
	header.WriteString("To: " + strings.Join(recipents.To, ",") + "\r\n")
	header.WriteString(encodeSubject(subject))
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimeter))
	header.WriteString(fmt.Sprintf("\r\n--%s\r\n", delimeter))
	header.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	header.WriteString("Content-Transfer-Encoding: base64\r\n")

	var message bytes.Buffer
	message = header
	message.WriteString("\r\n")
	message.WriteString(add76crlf(base64.StdEncoding.EncodeToString(bufBody.Bytes())))
	message.WriteString("\r\n")

	if attachment != nil{
		contentType := http.DetectContentType(attachment.Data.Bytes())
		message.WriteString(fmt.Sprintf("\r\n--%s\r\n", delimeter))
		message.WriteString("Content-Type: " + contentType + "; charset=\"utf-8\"\r\n")
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString("Content-Disposition: attachment;filename=\"" + attachment.Filename + "\"\r\n")

		message.WriteString("\r\n" + base64.StdEncoding.EncodeToString(attachment.Data.Bytes()))
	}


	return []byte(message.String())
}

// Send MAIL(No TLS)
func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}
	if a != nil{
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

