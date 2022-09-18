package utils

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	// Email Server Host
	MailHost = os.Getenv("EMAIL_HOST") //"smtp.163.com"
	// Email Port
	MailPort = os.Getenv("EMAIL_PORT")
	// Email User
	MailUser = os.Getenv("EMAIL_USER")
	// Eamil Authorization Code
	MailPwd = os.Getenv("EMAIL_PASSWORD")
)

/*
title send email
@param []string mailAddress
@param string subject
@param string body
@return error
*/
func SendMail(mailAddress []string, subject string, body string) error {
	m := gomail.NewMessage()
	nickname := "PVA Opeator"
	m.SetHeader("From", nickname+"<"+MailUser+">")
	m.SetHeader("To", mailAddress...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	mail_port, _ := strconv.Atoi(MailPort)
	d := gomail.NewDialer(MailHost, mail_port, MailUser, MailPwd)
	err := d.DialAndSend(m)
	return err
}
