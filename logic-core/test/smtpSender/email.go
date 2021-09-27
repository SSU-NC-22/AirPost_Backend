package main
 
import (
    "log"
	"fmt"
    "net/smtp"
    "strings"
)
 
const (
    // Gmail SMTP Server
    GoogleSMTPServer = "smtp.gmail.com"
)
 
type smtpSender struct {
    senderEmail string
    password    string
}
 
func NewSender(senderEmail string, password string) smtpSender {
    return smtpSender{senderEmail: senderEmail, password: password}
}
 
func (sender *smtpSender) SendMail(Dest []string, Subject string, Message string) error {
    msg := "From: " + sender.senderEmail + "\n" +
        "To: " + strings.Join(Dest, ",") + "\n" +
        "Subject: " + Subject + "\n" + Message
 
    err := smtp.SendMail(GoogleSMTPServer+":587",
        smtp.PlainAuth("", sender.senderEmail, sender.password, GoogleSMTPServer),
        sender.senderEmail, Dest, []byte(msg))
 
    if err != nil {
        fmt.Printf("smtp error: %s", err)
        return err
    }
 
    fmt.Println("Mail sent successfully!")
    return nil
}

func main() {
 
    const (
        senderEmail = "REDACTED"
        password    = ""
    )
 
    receiver := []string{"eunseo@q.ssu.ac.kr"}
    subject := "테스트 메일 발송입니다!! 메일 타이틀!"
    message := "이건 본문이구요~~~!! 메일 본문!!"
 
    smtpSender := NewSender(senderEmail, password)
    if err := smtpSender.SendMail(receiver, subject, message); err != nil {
        log.Panicln("smtp send error: ", err)
    } else {
        log.Println("smtp send ok")
    }
}
