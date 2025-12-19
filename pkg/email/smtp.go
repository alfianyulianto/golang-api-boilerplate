package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailService struct {
	*SMTPConfig
}

// EmailTemplateData adalah struct yang berisi data untuk template email
type EmailTemplateData struct {
	Name          string // Nama penerima (WAJIB)
	Subject       string // Subject email (WAJIB)
	Message       string // Pesan utama (WAJIB)
	Year          int    // Tahun untuk footer (WAJIB)
	ButtonText    string // Text tombol (OPSIONAL)
	ButtonURL     string // URL tujuan tombol (OPSIONAL)
	InfoTitle     string // Judul info box (OPSIONAL)
	InfoContent   string // Isi info box (OPSIONAL)
	HighlightText string // Text highlight dengan emoji petir (OPSIONAL)
	Note          string // Catatan khusus (OPSIONAL, ada default di template)
}

func NewEmailService(SMTPConfig *SMTPConfig) *EmailService {
	return &EmailService{SMTPConfig: SMTPConfig}
}

// SendEmail mengirim email HTML sederhana tanpa attachment
func (s *EmailService) SendEmail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n\r\n"+
		"%s\r\n", s.From, to[0], subject, body))

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.From, to, msg)
}

// SendTemplateEmail mengirim email menggunakan template HTML dengan logo embedded
func (s *EmailService) SendTemplateEmail(to []string, subject string, data interface{}) error {
	// Parse template
	tmpl, err := template.ParseFiles("pkg/email/template.gohtml")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var bodyBuf bytes.Buffer
	if err = tmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Baca logo file
	logoPath := "pkg/email/logo-nyinauni-golang.png"
	logoData, err := os.ReadFile(logoPath)
	if err != nil {
		return fmt.Errorf("error reading logo file: %w", err)
	}

	// Kirim email dengan logo embedded
	return s.sendEmailWithEmbeddedImage(to, subject, bodyBuf.String(), logoData, "logo", "logo-nyinauni-golang.png")
}

// sendEmailWithEmbeddedImage mengirim email dengan gambar embedded (inline)
func (s *EmailService) sendEmailWithEmbeddedImage(to []string, subject, htmlBody string, imageData []byte, contentID, filename string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// Boundary untuk multipart
	boundary := generateBoundary()

	// Build message
	var msg bytes.Buffer

	// Write headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=\"%s\"\r\n", boundary))
	msg.WriteString("\r\n")

	// HTML part
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)
	msg.WriteString("\r\n\r\n")

	// Image part (embedded)
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString(fmt.Sprintf("Content-Type: image/png; name=\"%s\"\r\n", filename))
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString(fmt.Sprintf("Content-ID: <%s>\r\n", contentID))
	msg.WriteString(fmt.Sprintf("Content-Disposition: inline; filename=\"%s\"\r\n", filename))
	msg.WriteString("\r\n")

	// Encode image to base64
	encoded := base64.StdEncoding.EncodeToString(imageData)
	// Split into 76 character lines (SMTP standard)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		msg.WriteString(encoded[i:end] + "\r\n")
	}
	msg.WriteString("\r\n")

	// End boundary
	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Send email
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.From, to, msg.Bytes())
}

// generateBoundary generates a unique boundary string
func generateBoundary() string {
	return fmt.Sprintf("boundary_%d", time.Now().UnixNano())
}
