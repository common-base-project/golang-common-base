package email

import (
	"bytes"
	"text/template"
)

type MIMEContent struct {
	From      string
	To        string
	Organizer string
	ReplyTo   string
	Now       string
	Subject   string
	Content   string
	FromT     string
	ToT       string
	Status    string // CONFIRMED CANCELLED
	Method    string // REQUEST  CANCEL
	Attendees []Attendee
	Location  string
}

type PlainTextContent struct {
	From    string // 发送者
	To      string // 接受者
	CC      string // 抄送
	Subject string // 主题
	Content string // 消息体
}

type Attendee struct {
	Email string
}

const MIMEContentT = ` 
Content-Type: multipart/mixed; boundary="===============7605613248146577809=="
MIME-Version: 1.0
Reply-To: {{.ReplyTo}}
Date: {{.Now}}
Subject: {{.Subject}}
From: {{.From}}
To: {{.To}}

--===============7605613248146577809==
Content-Type: multipart/alternative; boundary="===============4337741840263426948=="
MIME-Version: 1.0

--===============4337741840263426948==
Content-Type: text/html; charset="utf-8"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit

{{.Content}}
--===============4337741840263426948==
MIME-Version: 1.0
Content-Type: text/calendar; method="{{.Method}}"; charset="utf-8"
Content-Transfer-Encoding: 7bit

BEGIN:VCALENDAR
PRODID:-//xxx corporation//eim//EN
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:{{.Method}}
BEGIN:VTIMEZONE
TZID:China Standard Time
BEGIN:STANDARD
DTSTART:16010101T000000
TZOFFSETFROM:+0800
TZOFFSETTO:+0800
END:STANDARD
BEGIN:DAYLIGHT
DTSTART:16010101T000000
TZOFFSETFROM:+0800
TZOFFSETTO:+0800
END:DAYLIGHT
END:VTIMEZONE
BEGIN:VEVENT
DTSTART;TZID=China Standard Time:{{.FromT}}
DTEND;TZID=China Standard Time:{{.ToT}}
DTSTAMP:{{.Now}}
ORGANIZER;CN=organiser:mailto:{{.Organizer}}
UID:FIXMEUID{{.Now}}
{{range .Attendees -}}
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-    PARTICIPANT;PARTSTAT=ACCEPTED;RSVP=TRUE
 ;CN={{.Email}};X-NUM-GUESTS=0:
 mailto:{{.Email}}
{{end -}}
CREATED:{{.Now}}
DESCRIPTION: {{.Subject}}
LAST-MODIFIED:{{.Now}}
LOCATION:{{.Location}}
SEQUENCE:0
STATUS:{{.Status}}
SUMMARY:{{.Subject}}
TRANSP:OPAQUE
END:VEVENT
END:VCALENDAR

--===============4337741840263426948==--

--===============7605613248146577809==--

`

const PlainText = ` 
From: {{.From}}
To: {{.To}}
CC: {{.CC}}
Subject: {{.Subject}}
Content-Type: text/html; charset=UTF-8
MIME-Version: 1.0

<html>
<body>
<p>{{.Content}}</p>
</body>
</html>
`

func NewTemplate(kind string) *template.Template {
	switch kind {
	case MailKindInvitation:
		return template.Must(template.New("MIMEContent").Parse(MIMEContentT))
	case MailKindPlainText:
		return template.Must(template.New("PlainTextContent").Parse(PlainText))
	default:
		return nil
	}
}

func RenderTemplate(tpl *template.Template, content interface{}) (strc string, err error) {
	var buff bytes.Buffer
	err = tpl.Execute(&buff, content)
	if err == nil {
		strc = buff.String()
	}
	return
}
