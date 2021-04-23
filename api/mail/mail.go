package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/bson/primitive"
	config "thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
)

//SendWithGmail with gmail mail server
func SendWithGmail(fromMail string, toMail string, subject string, body string, token string) {
	e := email.NewEmail()
	e.From = "Runex Support<" + fromMail + ">"
	e.To = []string{toMail}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = subject
	//e.Text = []byte("Welcome to Runex\n")
	e.HTML = []byte(body)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com"))
	log.Println(err)
}

// SendConfirmRegMail for send email confirm account
func SendConfirmRegMail(name string, email string, userID primitive.ObjectID, pf string, role string) {
	log.Println("Send Email")
	token, err := GenarateToken(userID, role, pf)
	if err == nil {
		subject := "Welcome to Runex"
		body := "<h1>Welcome to Runex! Please confirm your email.\n<h1>" +
			"<p>Welcome, " + name + " !\n</p>" +
			"<p>Please confirm your email (" + email + ") so you don't miss out on anything important. We will also use this email as reference when you contact our support.\n</p>" +
			"<p>Let's get it on!\n</p>" +
			"<a href='https://runex.co/users/confirmation/" + token + "'>Confirm</a>"
		// "<a href='https://runex.co/users/confirmation/" + token + "'>Confirm</a>"

		admin := "runex.contact@gmail.com"
		SendWithGmail(admin, email, subject, body, token)
	}
}

// SendForgotMail for send email confirm account
func SendForgotMail(user model.UserForgot) error {
	log.Println("Send Email")
	token, err := GenarateToken(user.UserID, user.Role, user.PF)
	if err != nil {
		log.Println(err)
	}
	user.Token = token
	log.Printf("[info] data.Email %s", user.Email)
	r := NewRequest([]string{user.Email}, "Forgot password", "")
	err = r.ParseTemplate("./templates/ForgotPassword.html", user)
	if err := r.ParseTemplate("./templates/ForgotPassword.html", user); err == nil {
		e := email.NewEmail()
		e.From = "Runex Support<" + "runex.contact@gmail.com" + ">"
		e.To = []string{user.Email}
		//e.Bcc = []string{"test_bcc@example.com"}
		//e.Cc = []string{"test_cc@example.com"}
		e.Subject = "ตั้งค่ารหัสผ่านใหม่"
		//e.Text = []byte("Welcome to Runex\n")
		e.HTML = []byte(r.body)
		err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com"))
		//ok, _ := r.SendEmail()
		//log.Printf("[info] ok %s", ok)
	}
	log.Printf("[info] err %s", err)
	return err
}

// GenarateToken send to confirm mail
func GenarateToken(userID primitive.ObjectID, role string, pf string) (string, error) {
	// Create the token
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// Set some claims
	claims[config.ID_KEY] = userID
	claims[config.ROLE_KEY] = role
	claims[config.PF] = pf
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	// log.Println(tokenString)
	// log.Println(err)
	return tokenString, err

}

// SendRegRaceRun for send email register race run event
func SendRegRaceRun(data model.EmailTemplateData2) {
	log.Println("Send Email for Register Race Run")
	r := NewRequest([]string{data.Email}, "ลงทะเบียน", "")
	err := r.ParseTemplate("./templates/Order.html", data)
	// if err := r.ParseTemplate("./templates/Order.html", data); err == nil {
	// 	ok, _ := r.SendEmail()
	// 	log.Printf("[info] ok %s", ok)
	// }
	// log.Printf("[info] err %s", err)
	if err := r.ParseTemplate("./templates/Order.html", data); err == nil {
		e := email.NewEmail()
		e.From = "Runex Support<" + "runex.contact@gmail.com" + ">"
		e.To = []string{data.Email}
		//e.Bcc = []string{"test_bcc@example.com"}
		//e.Cc = []string{"test_cc@example.com"}
		e.Subject = "ลงทะเบียน " + data.EventName
		//e.Text = []byte("Welcome to Runex\n")
		e.HTML = []byte(r.body)
		err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com"))
		//ok, _ := r.SendEmail()
		//log.Printf("[info] ok %s", ok)
	}
	log.Printf("[info] send mail err %s", err)
}

func SendRegEventMail(data model.EmailTemplateData) {

	log.Printf("[info] data.Email %s", data.Email)
	r := NewRequest([]string{data.Email}, "Confirm register event!", "")
	err := r.ParseTemplate("./templates/RegisterEventTemplate.html", data)
	if err := r.ParseTemplate("./templates/RegisterEventTemplate.html", data); err == nil {
		ok, _ := r.SendEmail()
		log.Printf("[info] ok %s", ok)
	}
	log.Printf("[info] err %s", err)
}

func SendRegFreeEventMail(data model.EmailTemplateData) {

	log.Printf("[info] data.Email %s", data.Email)
	r := NewRequest([]string{data.Email}, "Confirm register event!", "")
	err := r.ParseTemplate("./templates/RegisterFreeEventTemplate.html", data)
	if err := r.ParseTemplate("./templates/RegisterFreeEventTemplate.html", data); err == nil {
		ok, _ := r.SendEmail()
		log.Printf("[info] ok %s", ok)
	}
	log.Printf("[info] err %s", err)
}

func SendRegEventMail2(data model.EmailTemplateData2) {

	// log.Printf("[info] data.Email %s", data.Email)
	// r := NewRequest([]string{data.Email}, "Confirm register event!", "")
	// err := r.ParseTemplate("./templates/RegisterEventTemplate2.html", data)
	// if err := r.ParseTemplate("./templates/RegisterEventTemplate2.html", data); err == nil {
	// 	ok, _ := r.SendEmail()
	// 	log.Printf("[info] ok %s", ok)
	// }
	// log.Printf("[info] err %s", err)
	// 	ประเภทการแข่งขัน  Chinjang Shiro Run 60 KM  500 THB
	// หมายเลขผู้สมัคร  00001
	// ticket_nameChinjang Shiro Run 60 KM
	// StatusPAYMENT_WAITING
	// Payment TypePAYMENT_TRANSFER
	// ที่อยู่ 987 หมู่ 8 City บ้านเป็ด District เมืองขอนแก่น Province ขอนแก่น 40000
	subject := "ขอบคุณที่สมัครเข้าร่วม " + data.TicketName
	body := "<h1>ประเภทการแข่งขัน " + data.TicketName + "\n<h1>" +
		"<p>ชื่อ : " + data.Name + "\n</p>" +
		"<p>เบอร์โทรศัพท์ : " + data.Phone + "\n</p>" +
		"<p>หมายเลขอ้างอิง : " + data.RefID + "\n</p>" +
		"<p>E BIB : " + data.RegisterNumber + "\n</p>" +
		"<p>Status : " + data.Status + " !\n</p>" +
		"<p>Email :" + data.Email + "\n</p>" +
		"<p>Payment Type : " + data.PaymentType + "\n</p>" +
		"<p>ราคา : " + strconv.FormatFloat(data.Price, 'f', 2, 64) + "\n</p>" +
		"<p>ที่อยู่จัดส่ง : " + data.ShipingAddress + "\n</p>" +
		"<p>Web site : <a href='https://runex.co'>runex.co</a>\n</p>"

	//admin := "runex.contact@gmail.com"
	e := email.NewEmail()
	e.From = "Runex Support<" + "runex.contact@gmail.com" + ">"
	e.To = []string{data.Email}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = subject
	//e.Text = []byte("Welcome to Runex\n")
	e.HTML = []byte(body)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com"))
	log.Println(err)
}

func SendRegFreeEventMail2(data model.EmailTemplateData2) {

	// log.Printf("[info] data.Email %s", data.Email)
	// r := NewRequest([]string{data.Email}, "Confirm register event!", "")
	// err := r.ParseTemplate("./templates/RegisterEventTemplate2.html", data)
	// if err := r.ParseTemplate("./templates/RegisterEventTemplate2.html", data); err == nil {
	// 	ok, _ := r.SendEmail()
	// 	log.Printf("[info] ok %s", ok)
	// }
	// log.Printf("[info] err %s", err)
	// 	ประเภทการแข่งขัน  Chinjang Shiro Run 60 KM  500 THB
	// หมายเลขผู้สมัคร  00001
	// ticket_nameChinjang Shiro Run 60 KM
	// StatusPAYMENT_WAITING
	// Payment TypePAYMENT_TRANSFER
	// ที่อยู่ 987 หมู่ 8 City บ้านเป็ด District เมืองขอนแก่น Province ขอนแก่น 40000
	subject := "ขอบคุณที่สมัครเข้าร่วม " + data.TicketName
	body := "<h1>ประเภทการแข่งขัน " + data.TicketName + "\n<h1>" +
		"<p>ชื่อ : " + data.Name + "\n</p>" +
		"<p>เบอร์โทรศัพท์ : " + data.Phone + "\n</p>" +
		"<p>หมายเลขอ้างอิง : " + data.RefID + "\n</p>" +
		"<p>E BIB : " + data.RegisterNumber + "\n</p>" +
		"<p>Status : " + data.Status + " !\n</p>" +
		"<p>Email :" + data.Email + "\n</p>" +
		"<p>Web site : <a href='https://runex.co'>runex.co</a>\n</p>"

	//admin := "runex.contact@gmail.com"
	e := email.NewEmail()
	e.From = "Runex Support<" + "runex.contact@gmail.com" + ">"
	e.To = []string{data.Email}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = subject
	//e.Text = []byte("Welcome to Runex\n")
	e.HTML = []byte(body)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com"))
	log.Println(err)
}

func TestMailTemplate() {

	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
	}
	r := NewRequest([]string{"techitblue@gmail.com"}, "Confirm register event", "")
	err := r.ParseTemplate("./templates/RegisterEventTemplate.html", templateData)
	if err := r.ParseTemplate("./templates/template.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
		log.Printf("[info] ok %s", ok)
	}
	log.Printf("[info] err %s", err)
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	from := "Runex Support <runex.contact@gmail.com>"
	//subject := "Subject: " + r.subject + ""
	msg := []byte(r.body)
	addr := "smtp.gmail.com:587"
	var auth smtp.Auth

	//auth = smtp.PlainAuth("", "support@runex.co", "Think@2019", "smtp.gmail.com")
	auth = smtp.PlainAuth("", "runex.contact@gmail.com", "zkwzkunlbfmesyvd", "smtp.gmail.com")
	if err := smtp.SendMail(addr, auth, from, r.to, msg); err != nil {
		log.Printf("[info] err SendEmail %s", err)
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
