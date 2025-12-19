# 🚀 Quick Start - Email Template Nyinauni

Panduan cepat untuk langsung mulai menggunakan email template.

---

## ⚡ Cara Paling Mudah (Recommended)

### 1. Setup Email Service (Hanya Sekali)

```go
import "github.com/alfianyulianto/pds-service/pkg/email"

// Di main.go atau config
emailService := email.NewEmailService(&email.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password", // App Password dari Gmail
    From:     "noreply@nyinauni.com",
})
```

### 2. Kirim Email (Pilih Salah Satu)

#### ✅ Verifikasi Email
```go
err := email.QuickSendVerification(
    emailService,
    "user@example.com",           // Email tujuan
    "Alfian Yulianto",            // Nama user
    "https://nyinauni.com/verify?token=abc123", // Link verifikasi
)
```

#### 🔐 Reset Password
```go
err := email.QuickSendResetPassword(
    emailService,
    "user@example.com",
    "Alfian Yulianto",
    "https://nyinauni.com/reset?token=xyz789",
)
```

#### 👋 Welcome Email
```go
err := email.QuickSendWelcome(
    emailService,
    "user@example.com",
    "Alfian Yulianto",
)
```

#### 🔔 Login Notification
```go
err := email.QuickSendLoginNotification(
    emailService,
    "user@example.com",
    "Alfian Yulianto",
    "17 Des 2025, 14:30 WIB",     // Waktu login
    "Chrome on Windows",           // Device
    "Jakarta, Indonesia",          // Lokasi
)
```

**Selesai!** ✨

---

## 🎨 Cara Custom (Lebih Fleksibel)

### Gunakan Builder Pattern

```go
// Buat email custom dengan builder
emailData := email.NewEmailTemplate().
    SetName("John Doe").
    SetSubject("Email Custom").
    SetMessage("Ini adalah pesan custom saya.").
    AddButton("Click Me", "https://example.com").
    AddHighlight("Penting!").
    Build()

// Kirim
err := emailService.SendTemplateEmail(
    []string{"user@example.com"},
    emailData.Subject,
    emailData,
)
```

### Atau Gunakan Struct Langsung

```go
emailData := email.EmailTemplateData{
    Name:    "Jane Smith",
    Subject: "Test Email",
    Message: "Hello World!",
    Year:    2025,
    ButtonText: "Visit",
    ButtonURL:  "https://example.com",
}

err := emailService.SendTemplateEmail(
    []string{"user@example.com"},
    emailData.Subject,
    emailData,
)
```

---

## 📦 Contoh Real World: Di Use Case

### Contoh 1: Saat Register User

```go
// Di auth_usecase.go
func (u *authUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
    // ... buat user di database ...

    // Generate token
    token := generateVerificationToken()
    verifyURL := fmt.Sprintf("https://nyinauni.com/verify?token=%s", token)

    // Kirim email verifikasi (async agar tidak block)
    go func() {
        if err := email.QuickSendVerification(
            u.SMTPService,
            request.Email,
            request.Name,
            verifyURL,
        ); err != nil {
            u.Log.WithError(err).Error("Failed to send verification email")
        }
    }()

    return userResponse, nil
}
```

### Contoh 2: Saat Forgot Password

```go
func (u *authUseCase) ForgotPassword(ctx context.Context, email string) error {
    // ... cari user ...

    // Generate reset token
    token := generateResetToken()
    resetURL := fmt.Sprintf("https://nyinauni.com/reset-password?token=%s", token)

    // Kirim email reset password
    go func() {
        if err := email.QuickSendResetPassword(
            u.SMTPService,
            user.Email,
            user.Name,
            resetURL,
        ); err != nil {
            u.Log.WithError(err).Error("Failed to send reset password email")
        }
    }()

    return nil
}
```

### Contoh 3: Email Custom dengan Builder

```go
func (u *authUseCase) SendCustomEmail(ctx context.Context, userEmail, userName string) error {
    // Buat email custom menggunakan builder
    emailData := email.NewEmailTemplate().
        SetName(userName).
        SetSubject("Promo Spesial Untuk Anda!").
        SetMessage(`Halo! Ada promo spesial bulan ini.
        
Dapatkan diskon 50% untuk semua course premium. Buruan ambil sekarang!`).
        AddButton("Lihat Promo", "https://nyinauni.com/promo").
        AddHighlight("Promo berakhir 31 Desember!").
        AddInfoBox("Syarat & Ketentuan", "Promo hanya berlaku untuk user terdaftar. Satu akun satu promo.").
        Build()

    // Kirim async
    go func() {
        if err := u.SMTPService.SendTemplateEmail(
            []string{userEmail},
            emailData.Subject,
            emailData,
        ); err != nil {
            u.Log.WithError(err).Error("Failed to send promo email")
        }
    }()

    return nil
}
```

---

## 🔧 Setup Gmail (Langkah-langkah)

### Jika Menggunakan Gmail:

1. **Buka Google Account**: https://myaccount.google.com/
2. **Security** → **2-Step Verification** → Aktifkan
3. **App passwords** → Pilih "Mail" → Generate
4. **Copy password** yang muncul (16 karakter)
5. **Gunakan di config**:

```go
emailService := email.NewEmailService(&email.SMTPConfig{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "abcd efgh ijkl mnop", // Password 16 karakter dari step 4
    From:     "noreply@nyinauni.com",
})
```

---

## 📝 Template Tersedia

| Template | Function | Keterangan |
|----------|----------|------------|
| ✅ Verification | `QuickSendVerification()` | Email verifikasi akun |
| 🔐 Reset Password | `QuickSendResetPassword()` | Email reset password |
| 👋 Welcome | `QuickSendWelcome()` | Email selamat datang |
| 🔔 Login Notification | `QuickSendLoginNotification()` | Notifikasi login |
| 💳 Payment Success | `PaymentSuccessEmailTemplate()` | Konfirmasi pembayaran |
| 🔑 Password Changed | `PasswordChangedEmailTemplate()` | Konfirmasi ubah password |
| 📅 Event Invitation | `EventInvitationEmailTemplate()` | Undangan event/webinar |
| 🗑️ Account Deleted | `AccountDeletedEmailTemplate()` | Konfirmasi hapus akun |

---

## ⚠️ Tips Penting

### 1. Selalu Kirim Async (Jangan Block User)
```go
// ❌ JANGAN seperti ini (blocking)
err := emailService.SendTemplateEmail(...)
if err != nil {
    return err // User menunggu email selesai dikirim
}

// ✅ LAKUKAN seperti ini (non-blocking)
go func() {
    if err := emailService.SendTemplateEmail(...); err != nil {
        log.Printf("Error: %v", err)
    }
}()
return nil // User langsung dapat response
```

### 2. Handle Error dengan Bijak
```go
go func() {
    if err := emailService.SendTemplateEmail(...); err != nil {
        // Log error tapi jangan panic
        u.Log.WithError(err).Error("Failed to send email")
        
        // Optional: save ke database untuk retry nanti
        // saveToEmailQueue(emailData)
    }
}()
```

### 3. Gunakan Rate Limiting
```go
// Jangan spam user dengan banyak email
// Implementasi rate limiting jika perlu
```

---

## 🐛 Troubleshooting Cepat

### Email Tidak Terkirim?

**Cek ini:**
1. ✅ SMTP credentials benar?
2. ✅ Sudah pakai App Password (bukan password Gmail biasa)?
3. ✅ Port 587 tidak diblock firewall?
4. ✅ File logo ada di `pkg/email/logo-nyinauni-golang.png`?

**Debug:**
```go
err := emailService.SendTemplateEmail(...)
if err != nil {
    fmt.Printf("Error detail: %+v\n", err)
}
```

---

## 📚 Next Steps

- Lihat `README.md` untuk dokumentasi lengkap
- Lihat `example_usage.go` untuk contoh lebih banyak
- Lihat `helpers.go` untuk semua template tersedia

---

**Happy Coding! 🚀**

