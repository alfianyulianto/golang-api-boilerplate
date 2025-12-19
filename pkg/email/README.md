# 📧 Email Service - Nyinauni Golang

Dokumentasi lengkap untuk menggunakan Email Service dengan template HTML yang cantik dan responsive.

## 📋 Daftar Isi
- [Setup](#-setup)
- [Struktur Data Template](#-struktur-data-template)
- [Cara Penggunaan](#-cara-penggunaan)
- [Contoh Lengkap](#-contoh-lengkap)
- [Tips & Best Practices](#-tips--best-practices)

---

## 🚀 Setup

### 1. Konfigurasi SMTP

Buat konfigurasi SMTP (biasanya di `config.json` atau environment variables):

```json
{
  "smtp": {
    "host": "smtp.gmail.com",
    "port": 587,
    "username": "your-email@gmail.com",
    "password": "your-app-password",
    "from": "noreply@nyinauni.com"
  }
}
```

**Catatan untuk Gmail:**
- Gunakan **App Password**, bukan password Gmail biasa
- Cara mendapatkan App Password:
  1. Buka Google Account → Security
  2. Aktifkan 2-Step Verification
  3. Cari "App passwords"
  4. Generate password untuk "Mail"

### 2. Inisialisasi Service

```go
import "github.com/alfianyulianto/pds-service/pkg/email"

// Setup SMTP Config
smtpConfig := &email.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password",
    From:     "noreply@nyinauni.com",
}

// Buat Email Service
emailService := email.NewEmailService(smtpConfig)
```

---

## 📊 Struktur Data Template

Template menerima struct dengan field berikut:

| Field | Type | Required | Deskripsi |
|-------|------|----------|-----------|
| `Name` | string | ✅ Wajib | Nama penerima email |
| `Subject` | string | ✅ Wajib | Subject/judul email |
| `Message` | string | ✅ Wajib | Pesan utama email (bisa multi-line) |
| `Year` | int | ✅ Wajib | Tahun untuk footer copyright |
| `ButtonText` | string | ⚪ Opsional | Text tombol call-to-action |
| `ButtonURL` | string | ⚪ Opsional | URL tujuan tombol |
| `InfoTitle` | string | ⚪ Opsional | Judul info box |
| `InfoContent` | string | ⚪ Opsional | Konten info box |
| `HighlightText` | string | ⚪ Opsional | Text yang di-highlight (dengan emoji ⚡) |
| `Note` | string | ⚪ Opsional | Catatan khusus (ada default jika kosong) |

---

## 💡 Cara Penggunaan

### Penggunaan Dasar

```go
package main

import (
    "fmt"
    "time"
    "github.com/alfianyulianto/pds-service/pkg/email"
)

func main() {
    // 1. Setup Email Service
    emailService := email.NewEmailService(&email.SMTPConfig{
        Host:     "smtp.gmail.com",
        Port:     587,
        Username: "your-email@gmail.com",
        Password: "your-app-password",
        From:     "noreply@nyinauni.com",
    })

    // 2. Siapkan data untuk template
    emailData := email.EmailTemplateData{
        Name:    "Alfian Yulianto",
        Subject: "Selamat Datang!",
        Message: "Terima kasih telah bergabung dengan Nyinauni Golang.",
        Year:    time.Now().Year(),
    }

    // 3. Kirim email
    err := emailService.SendTemplateEmail(
        []string{"user@example.com"},  // Penerima (bisa multiple)
        emailData.Subject,              // Subject
        emailData,                      // Data template
    )

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Email berhasil dikirim!")
    }
}
```

---

## 🎯 Contoh Lengkap

### 1. Email Verifikasi Akun

```go
emailData := email.EmailTemplateData{
    Name:    "John Doe",
    Subject: "Verifikasi Akun Anda",
    Message: `Terima kasih telah mendaftar di Nyinauni Golang!
    
Untuk mengaktifkan akun Anda, silakan klik tombol verifikasi di bawah ini. Link verifikasi akan kadaluarsa dalam 24 jam.`,
    Year:          2025,
    ButtonText:    "Verifikasi Akun",
    ButtonURL:     "https://nyinauni.com/verify?token=abc123",
    InfoTitle:     "Kenapa Perlu Verifikasi?",
    InfoContent:   "Verifikasi email membantu kami memastikan keamanan akun Anda dan mencegah spam.",
    HighlightText: "Link berlaku selama 24 jam!",
    Note:          "Jika Anda tidak mendaftar, abaikan email ini.",
}

err := emailService.SendTemplateEmail(
    []string{"john@example.com"},
    emailData.Subject,
    emailData,
)
```

**Hasil:** Email dengan tombol biru cantik, info box, dan highlight warning.

---

### 2. Email Reset Password

```go
emailData := email.EmailTemplateData{
    Name:    "Jane Smith",
    Subject: "Reset Password Anda",
    Message: `Kami menerima permintaan reset password untuk akun Anda.
    
Klik tombol di bawah untuk membuat password baru. Link ini akan kadaluarsa dalam 1 jam.`,
    Year:          2025,
    ButtonText:    "Reset Password",
    ButtonURL:     "https://nyinauni.com/reset?token=xyz789",
    HighlightText: "Link berlaku 1 jam saja!",
    InfoTitle:     "Tips Password Aman",
    InfoContent:   "Gunakan minimal 8 karakter dengan kombinasi huruf, angka, dan simbol.",
}

emailService.SendTemplateEmail([]string{"jane@example.com"}, emailData.Subject, emailData)
```

---

### 3. Email Welcome (Tanpa Tombol)

```go
emailData := email.EmailTemplateData{
    Name:    "Budi Santoso",
    Subject: "Selamat Datang di Nyinauni! 🎉",
    Message: `Halo Budi! Selamat datang di komunitas Nyinauni Golang.
    
Kami sangat senang Anda bergabung. Mulai perjalanan belajar Golang Anda sekarang!`,
    Year: 2025,
    InfoTitle: "Apa Selanjutnya?",
    InfoContent: `1. Lengkapi profil Anda
2. Pilih course pertama
3. Bergabung di Discord
4. Mulai coding!`,
}

emailService.SendTemplateEmail([]string{"budi@example.com"}, emailData.Subject, emailData)
```

---

### 4. Email Notifikasi Sederhana

```go
emailData := email.EmailTemplateData{
    Name:    "Ahmad",
    Subject: "Login Terdeteksi",
    Message: fmt.Sprintf(`Kami mendeteksi login ke akun Anda pada:

Waktu: %s
Device: Chrome on Windows
Lokasi: Jakarta, Indonesia`, time.Now().Format("02 Jan 2006 15:04")),
    Year:          2025,
    HighlightText: "Jika bukan Anda, segera ubah password!",
}

emailService.SendTemplateEmail([]string{"ahmad@example.com"}, emailData.Subject, emailData)
```

---

### 5. Penggunaan di Use Case (Real World)

```go
// Di file auth_usecase.go
func (u *authUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
    // ... create user di database ...

    // Generate verification token
    token := generateToken()
    verifyURL := fmt.Sprintf("https://nyinauni.com/verify?token=%s", token)

    // Kirim email verifikasi
    emailData := email.EmailTemplateData{
        Name:    request.Name,
        Subject: "Verifikasi Email - Nyinauni Golang",
        Message: `Terima kasih telah mendaftar! 
        
Untuk melanjutkan, silakan verifikasi email Anda dengan klik tombol di bawah.`,
        Year:          time.Now().Year(),
        ButtonText:    "Verifikasi Email",
        ButtonURL:     verifyURL,
        HighlightText: "Link berlaku 24 jam",
        InfoTitle:     "Perlu Bantuan?",
        InfoContent:   "Jika tombol tidak berfungsi, copy link ini: " + verifyURL,
    }

    // Kirim email async (non-blocking)
    go func() {
        if err := u.SMTPService.SendTemplateEmail(
            []string{request.Email},
            emailData.Subject,
            emailData,
        ); err != nil {
            u.Log.WithError(err).Error("Failed to send verification email")
        }
    }()

    return userResponse, nil
}
```

---

## 📝 Tips & Best Practices

### 1. **Kirim Email Secara Async**
Jangan block request user, kirim email di goroutine:

```go
go func() {
    if err := emailService.SendTemplateEmail(to, subject, data); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}()
```

### 2. **Handle Error dengan Baik**
```go
if err := emailService.SendTemplateEmail(to, subject, data); err != nil {
    // Log error tapi jangan gagalkan proses utama
    log.WithError(err).Error("Email sending failed")
    // Bisa save ke queue untuk retry
}
```

### 3. **Gunakan Multi-line String**
Untuk message yang panjang, gunakan backtick:

```go
Message: `Ini baris pertama.

Ini paragraf kedua.

- Poin 1
- Poin 2`,
```

### 4. **Validasi Email**
Pastikan format email valid sebelum kirim:

```go
import "net/mail"

func isValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}
```

### 5. **Rate Limiting**
Jangan spam! Batasi jumlah email:

```go
// Contoh: max 10 email per menit
limiter := rate.NewLimiter(rate.Every(time.Minute/10), 1)
if !limiter.Allow() {
    return errors.New("too many emails")
}
```

### 6. **Testing**
Untuk development, gunakan service seperti [Mailtrap](https://mailtrap.io) atau [MailHog](https://github.com/mailhog/MailHog):

```go
// Development config
smtpConfig := &email.SMTPConfig{
    Host:     "smtp.mailtrap.io",
    Port:     2525,
    Username: "your-mailtrap-username",
    Password: "your-mailtrap-password",
    From:     "test@nyinauni.com",
}
```

---

## 🎨 Customization

### Menambah Field Baru

1. **Update struct di `smtp.go`** (atau buat struct terpisah):
```go
type EmailTemplateData struct {
    // Existing fields...
    CustomField string
}
```

2. **Update template di `template.gohtml`**:
```html
{{if .CustomField}}
<div class="custom-section">
    {{.CustomField}}
</div>
{{end}}
```

---

## 🐛 Troubleshooting

### Email Tidak Terkirim

**Cek:**
1. ✅ SMTP credentials benar?
2. ✅ Port benar? (587 untuk TLS, 465 untuk SSL)
3. ✅ Gmail: sudah pakai App Password?
4. ✅ Firewall tidak block koneksi SMTP?

**Debug:**
```go
err := emailService.SendTemplateEmail(to, subject, data)
if err != nil {
    log.Printf("Error detail: %+v", err)
}
```

### Logo Tidak Muncul

- Pastikan file `logo-nyinauni-golang.png` ada di `pkg/email/`
- Cek path relatif sesuai dengan working directory

### Template Error

```go
// Pastikan semua field required diisi
if data.Name == "" || data.Subject == "" || data.Message == "" {
    return errors.New("required fields missing")
}
```

---

## 📚 Referensi

- [Go SMTP Documentation](https://pkg.go.dev/net/smtp)
- [HTML Email Best Practices](https://www.campaignmonitor.com/css/)
- [Go Templates](https://pkg.go.dev/html/template)

---

## 👨‍💻 Author

**Alfian Yulianto**  
- GitHub: [@alfianyulianto](https://github.com/alfianyulianto)
- Email: alfianyulianto36@gmail.com
- Website: [alfiansite.my.id](http://alfiansite.my.id/)

---

**Happy Coding! 🚀**

