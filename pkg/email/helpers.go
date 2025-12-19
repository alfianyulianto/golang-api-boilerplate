package email

import (
	"fmt"
	"time"
)

// EmailTemplateBuilder adalah builder pattern untuk membuat EmailTemplateData dengan mudah
type EmailTemplateBuilder struct {
	data EmailTemplateData
}

// NewEmailTemplate membuat builder baru untuk email template
func NewEmailTemplate() *EmailTemplateBuilder {
	return &EmailTemplateBuilder{
		data: EmailTemplateData{
			Year: time.Now().Year(), // Default ke tahun sekarang
		},
	}
}

// SetName mengatur nama penerima (REQUIRED)
func (b *EmailTemplateBuilder) SetName(name string) *EmailTemplateBuilder {
	b.data.Name = name
	return b
}

// SetSubject mengatur subject email (REQUIRED)
func (b *EmailTemplateBuilder) SetSubject(subject string) *EmailTemplateBuilder {
	b.data.Subject = subject
	return b
}

// SetMessage mengatur pesan utama (REQUIRED)
func (b *EmailTemplateBuilder) SetMessage(message string) *EmailTemplateBuilder {
	b.data.Message = message
	return b
}

// SetYear mengatur tahun untuk footer (optional, default: tahun sekarang)
func (b *EmailTemplateBuilder) SetYear(year int) *EmailTemplateBuilder {
	b.data.Year = year
	return b
}

// AddButton menambahkan tombol call-to-action
func (b *EmailTemplateBuilder) AddButton(text, url string) *EmailTemplateBuilder {
	b.data.ButtonText = text
	b.data.ButtonURL = url
	return b
}

// AddInfoBox menambahkan info box
func (b *EmailTemplateBuilder) AddInfoBox(title, content string) *EmailTemplateBuilder {
	b.data.InfoTitle = title
	b.data.InfoContent = content
	return b
}

// AddHighlight menambahkan text highlight dengan emoji petir
func (b *EmailTemplateBuilder) AddHighlight(text string) *EmailTemplateBuilder {
	b.data.HighlightText = text
	return b
}

// AddNote menambahkan catatan khusus
func (b *EmailTemplateBuilder) AddNote(note string) *EmailTemplateBuilder {
	b.data.Note = note
	return b
}

// Build menghasilkan EmailTemplateData yang siap digunakan
func (b *EmailTemplateBuilder) Build() EmailTemplateData {
	return b.data
}

// Validate memvalidasi data sebelum dikirim
func (b *EmailTemplateBuilder) Validate() error {
	if b.data.Name == "" {
		return fmt.Errorf("name is required")
	}
	if b.data.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if b.data.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// SendWith mengirim email menggunakan EmailService yang diberikan
func (b *EmailTemplateBuilder) SendWith(service *EmailService, to []string) error {
	if err := b.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	data := b.Build()
	return service.SendTemplateEmail(to, data.Subject, data)
}

// ============================================================
// Pre-built Templates untuk use case umum
// ============================================================

// VerificationEmailTemplate membuat template email verifikasi
func VerificationEmailTemplate(name, verifyURL string) EmailTemplateData {
	return NewEmailTemplate().
		SetName(name).
		SetSubject("Verifikasi Email Anda - Nyinauni Golang").
		SetMessage(`Terima kasih telah mendaftar di platform Nyinauni Golang!

Untuk mengaktifkan akun Anda, silakan klik tombol verifikasi di bawah ini. Link verifikasi akan kadaluarsa dalam 24 jam.`).
		AddButton("Verifikasi Email", verifyURL).
		AddInfoBox("Kenapa Perlu Verifikasi?", "Verifikasi email membantu kami memastikan keamanan akun Anda dan mencegah penyalahgunaan platform.").
		AddHighlight("Link berlaku selama 24 jam!").
		AddNote("Jika Anda tidak mendaftar di platform kami, abaikan email ini. Akun tidak akan dibuat tanpa verifikasi.").
		Build()
}

// ResetPasswordEmailTemplate membuat template email reset password
func ResetPasswordEmailTemplate(name, resetURL string) EmailTemplateData {
	return NewEmailTemplate().
		SetName(name).
		SetSubject("Reset Password - Nyinauni Golang").
		SetMessage(`Kami menerima permintaan untuk mereset password akun Anda.

Jika Anda yang melakukan permintaan ini, silakan klik tombol di bawah untuk membuat password baru. Link ini akan kadaluarsa dalam 1 jam.`).
		AddButton("Reset Password", resetURL).
		AddInfoBox("Tips Password Aman", "Gunakan kombinasi huruf besar, huruf kecil, angka, dan simbol. Minimal 8 karakter untuk keamanan maksimal.").
		AddHighlight("Link reset password berlaku 1 jam!").
		AddNote("Jika Anda tidak meminta reset password, abaikan email ini atau hubungi kami jika Anda khawatir tentang keamanan akun.").
		Build()
}

// WelcomeEmailTemplate membuat template email selamat datang
func WelcomeEmailTemplate(name string) EmailTemplateData {
	return NewEmailTemplate().
		SetName(name).
		SetSubject("Selamat Datang di Nyinauni Golang! 🎉").
		SetMessage(`Selamat! Akun Anda telah berhasil diverifikasi dan aktif.

Sekarang Anda dapat mengakses semua fitur pembelajaran Golang yang kami sediakan. Mari mulai perjalanan belajar Anda bersama kami!`).
		AddInfoBox("Langkah Selanjutnya", fmt.Sprint("1. Lengkapi profil Anda untuk personalisasi pengalaman\n2. Mulai dengan course \"Golang Fundamentals\"\n3. Bergabung dengan komunitas Discord kami\n4. Eksplorasi project-project menarik")).
		AddNote("Ada pertanyaan? Tim support kami siap membantu Anda 24/7. Jangan ragu untuk menghubungi kami kapan saja!").
		Build()
}

// LoginNotificationEmailTemplate membuat template notifikasi login
func LoginNotificationEmailTemplate(name, loginTime, device string) EmailTemplateData {
	message := fmt.Sprintf(`Kami mendeteksi login baru ke akun Anda pada:

⏰ Waktu: %s
💻 Perangkat: %s

Jika ini adalah Anda, tidak perlu melakukan apa-apa. Jika bukan, segera ubah password Anda.`, loginTime, device)

	return NewEmailTemplate().
		SetName(name).
		SetSubject("Aktivitas Login Terdeteksi - Nyinauni Golang").
		SetMessage(message).
		AddHighlight("Jika bukan Anda, segera amankan akun!").
		AddNote("Untuk keamanan akun, kami merekomendasikan mengaktifkan 2-Factor Authentication (2FA).").
		Build()
}

// EventInvitationEmailTemplate membuat template undangan event
func EventInvitationEmailTemplate(name, eventTitle, eventDate, eventTime, topic, registerURL string) EmailTemplateData {
	message := fmt.Sprintf(`Anda diundang untuk mengikuti webinar spesial kami!

📅 Tanggal: %s
⏰ Waktu: %s
🎯 Topik: %s

Jangan lewatkan kesempatan ini untuk belajar dan berinteraksi dengan expert!`, eventDate, eventTime, topic)

	return NewEmailTemplate().
		SetName(name).
		SetSubject(fmt.Sprintf("Undangan: %s - Nyinauni Golang", eventTitle)).
		SetMessage(message).
		AddButton("Daftar Sekarang", registerURL).
		AddInfoBox("Yang Akan Anda Pelajari", "Materi lengkap, live coding, Q&A session, dan sertifikat kehadiran untuk peserta.").
		AddHighlight("Kuota terbatas! Daftar sekarang").
		AddNote("Link Zoom akan dikirim 1 hari sebelum acara. Simpan tanggalnya!").
		Build()
}

// PasswordChangedEmailTemplate membuat template konfirmasi perubahan password
func PasswordChangedEmailTemplate(name, changeTime string) EmailTemplateData {
	return NewEmailTemplate().
		SetName(name).
		SetSubject("Password Berhasil Diubah - Nyinauni Golang").
		SetMessage(fmt.Sprintf(`Password akun Anda telah berhasil diubah pada %s.

Jika ini bukan Anda, segera hubungi tim support kami untuk mengamankan akun Anda.`, changeTime)).
		AddHighlight("Jika bukan Anda, segera hubungi kami!").
		AddInfoBox("Tips Keamanan", "Jangan bagikan password Anda kepada siapapun. Gunakan password yang unik untuk setiap layanan.").
		Build()
}

// AccountDeletedEmailTemplate membuat template konfirmasi penghapusan akun
func AccountDeletedEmailTemplate(name string) EmailTemplateData {
	return NewEmailTemplate().
		SetName(name).
		SetSubject("Akun Anda Telah Dihapus - Nyinauni Golang").
		SetMessage(`Akun Anda telah berhasil dihapus dari sistem kami sesuai permintaan Anda.

Semua data personal Anda telah dihapus secara permanen. Kami sedih melihat Anda pergi, namun pintu kami selalu terbuka jika Anda ingin bergabung kembali.`).
		AddNote("Terima kasih telah menjadi bagian dari komunitas Nyinauni Golang. Kami berharap dapat bertemu lagi di masa depan!").
		Build()
}

// PaymentSuccessEmailTemplate membuat template konfirmasi pembayaran
func PaymentSuccessEmailTemplate(name, orderID, amount, item string) EmailTemplateData {
	message := fmt.Sprintf(`Pembayaran Anda telah berhasil diproses! 🎉

📦 Order ID: %s
💰 Total: %s
📋 Item: %s

Terima kasih atas pembelian Anda. Akses ke course sudah aktif di akun Anda.`, orderID, amount, item)

	return NewEmailTemplate().
		SetName(name).
		SetSubject("Pembayaran Berhasil - Nyinauni Golang").
		SetMessage(message).
		AddButton("Mulai Belajar", "https://nyinauni.com/my-courses").
		AddInfoBox("Akses Course", "Course Anda sudah aktif dan dapat diakses kapan saja. Selamat belajar!").
		Build()
}

// ============================================================
// Helper Functions
// ============================================================

// QuickSendVerification shortcut untuk mengirim email verifikasi
func QuickSendVerification(service *EmailService, to, name, verifyURL string) error {
	data := VerificationEmailTemplate(name, verifyURL)
	return service.SendTemplateEmail([]string{to}, data.Subject, data)
}

// QuickSendResetPassword shortcut untuk mengirim email reset password
func QuickSendResetPassword(service *EmailService, to, name, resetURL string) error {
	data := ResetPasswordEmailTemplate(name, resetURL)
	return service.SendTemplateEmail([]string{to}, data.Subject, data)
}

// QuickSendWelcome shortcut untuk mengirim email welcome
func QuickSendWelcome(service *EmailService, to, name string) error {
	data := WelcomeEmailTemplate(name)
	return service.SendTemplateEmail([]string{to}, data.Subject, data)
}

// QuickSendLoginNotification shortcut untuk mengirim notifikasi login
func QuickSendLoginNotification(service *EmailService, to, name, loginTime, device string) error {
	data := LoginNotificationEmailTemplate(name, loginTime, device)
	return service.SendTemplateEmail([]string{to}, data.Subject, data)
}

// QuickPasswordChangedEmail shortcut untuk mengirim notifikasi terhadap perubahan password
func QuickPasswordChangedEmail(service *EmailService, to, name, changeTime string) error {
	data := PasswordChangedEmailTemplate(name, changeTime)
	return service.SendTemplateEmail([]string{to}, data.Subject, data)
}
