package handler

import (
	"fmt"
	"net/http"
	"net/mail"
	"open-platform/utils"

	"strings"

	sasl "github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
)

type sender struct {
	To      string   `json:"to"`
	ToList  []string `json:"toList"`
	Name    string   `json:"name"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}

// SendMailHandler is a func to handle send email template requests
func SendMailHandler(c *gin.Context) {
	var data sender
	c.BindJSON(&data)

	to := data.To
	toList := data.ToList
	name := data.Name
	content := data.Content
	subject := data.Subject

	if ((to == "") && (len(toList) == 0)) || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Post data missing parameter"})
		return
	}

	renderContent, err := utils.RenderHTML(name, content)

	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusBadRequest})
	}

	err = SendToMail(
		utils.AppConfig.SMTP.Sender,
		utils.AppConfig.SMTP.Password,
		utils.AppConfig.SMTP.Host,
		subject, renderContent, "html", utils.RemoveDuplicate(append(toList, to)))

	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusBadRequest})
	} else {
		fmt.Println("Send mail success!")
		c.JSON(http.StatusOK, gin.H{"message": "OK", "code": http.StatusOK})
	}
}

// SendToMail is a function to handle send email smtp requests
func SendToMail(user, password, host, subject, body, mailtype string, to []string) error {
	auth := sasl.NewPlainClient("", user, password)
	fromName := "联创团队"
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}
	msg := strings.NewReader("To: " + strings.Join(to, ",") + "\r\nReply-To: " + "contact@hustunique.com" + "\r\nFrom: " + fromName + " <" + user + ">\r\nSubject: " + encodeRFC2047(subject) + "\r\n" + contentType + "\r\n\r\n" + body)

	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: String, Address: ""}
	return strings.Trim(addr.String(), "<@>")
}
