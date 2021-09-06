package mailer

import (
	"io"
	"mime"
	"net/mail"
	"net/textproto"
	"strings"
	"time"
)

type Email struct {
	// ReplyTo     []string
	From    string
	To      string
	Tos     []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte // Plaintext message (optional)
	HTML    []byte // Html message (optional)
	Sender  string // override From as SMTP envelope sender (optional)
	// Attachments []*Attachment
	// ReadReceipt []string
}

func (e *Email) buildHeaders() (body []byte, headers textproto.MIMEHeader, err error) {
	res := make(textproto.MIMEHeader, 6)

	// check content type
	if len(e.HTML) > 0 {
		res.Set("Content-Type", "text/html; charset=UTF-8")
		body = e.HTML
	} else {
		res.Set("Content-Type", "text/plain; charset=UTF-8")
		body = e.Text
	}
	res.Set("Content-Transfer-Encoding", "quoted-printable")
	// Set headers if there are values.
	if e.To != "" {
		e.Tos = append(e.Tos, strings.Split(e.To, ",")...)
	}
	if len(e.Tos) > 0 {
		res.Set("To", strings.Join(e.Tos, ", "))
	}
	if len(e.Cc) > 0 {
		res.Set("Cc", strings.Join(e.Cc, ", "))
	}
	if e.Subject != "" {
		res.Set("Subject", e.Subject)
	}
	// Date and From are required headers.
	res.Set("From", e.From)
	res.Set("Date", time.Now().Format(time.RFC1123Z))
	res.Set("MIME-Version", "1.0")
	return body, res, nil
}

// get encoded email
func (e *Email) Bytes(buf io.Writer) error {
	// -- build headers
	if body, headers, err := e.buildHeaders(); err == nil {
		headerToBytes(buf, &headers)
		_, err := buf.Write(body)
		return err
	} else {
		return err
	}
}

func headerToBytes(buff io.Writer, header *textproto.MIMEHeader) {
	for field, vals := range *header {
		for _, subval := range vals {
			io.WriteString(buff, field)
			io.WriteString(buff, ": ")
			// Write the encoded header if needed
			switch {
			case field == "Content-Type" || field == "Content-Disposition":
				io.WriteString(buff, subval)
			case field == "From" || field == "To" || field == "Cc" || field == "Bcc":
				participants := strings.Split(subval, ",")
				for i, v := range participants {
					addr, err := mail.ParseAddress(v)
					if err != nil {
						continue
					}
					participants[i] = addr.String()
				}
				io.WriteString(buff, strings.Join(participants, ", "))
			default:
				io.WriteString(buff, mime.QEncoding.Encode("UTF-8", subval))
			}
			io.WriteString(buff, "\r\n")
		}
	}
}
