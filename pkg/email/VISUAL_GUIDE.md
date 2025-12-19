# 📖 Visual Guide - Struktur Template Email

Panduan visual untuk memahami struktur template email Nyinauni Golang.

---

## 🎨 Struktur Visual Template

```
┌─────────────────────────────────────────────────┐
│                                                 │
│         🐹 LOGO NYINAUNI                       │
│            NYINAUNI                             │
│      Belajar Golang Bersama                     │
│                                                 │
└─────────────────────────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│                                                 │
│  Halo, {{.Name}}! 👋                           │
│                                                 │
│  {{.Message}}                                   │
│  (Pesan utama email)                           │
│                                                 │
│  ┌─────────────────────────────────┐           │
│  │     {{.ButtonText}}             │ (Optional)│
│  └─────────────────────────────────┘           │
│                                                 │
│  ┌────────────────────────────────────────┐    │
│  │ 💡 {{.InfoTitle}}                      │    │
│  │    {{.InfoContent}}                    │    │
│  └────────────────────────────────────────┘    │
│                    (Optional)                   │
│                                                 │
│  ┌────────────────────────────────────────┐    │
│  │ ⚡ {{.HighlightText}}                  │    │
│  └────────────────────────────────────────┘    │
│                    (Optional)                   │
│                                                 │
│  ─────────────────────────────────────────     │
│                                                 │
│  📌 Catatan Penting:                           │
│  {{.Note}}                                      │
│  (Default text jika tidak diisi)               │
│                                                 │
│  Salam hangat,                                  │
│  Tim Nyinauni Golang                           │
│  🐹 Go  💻 Programming  🚀 Learn               │
│                                                 │
└─────────────────────────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│              FOOTER (Dark)                      │
│                                                 │
│    🔗 GitHub  🌐 Website  ✉️ Email             │
│                                                 │
│    © {{.Year}} Nyinauni Golang.                │
│    All rights reserved.                         │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## 📊 Field Mapping

| Section | Field | Required | Default | Preview |
|---------|-------|----------|---------|---------|
| **Header** | (Logo) | - | Auto | Logo Nyinauni |
| **Header** | (Title) | - | Auto | "NYINAUNI" |
| **Header** | (Subtitle) | - | Auto | "Belajar Golang Bersama" |
| **Greeting** | `Name` | ✅ | - | "Halo, Alfian!" |
| **Message** | `Message` | ✅ | - | Konten pesan utama |
| **Button** | `ButtonText` | ⚪ | Hidden | "Verifikasi Akun" |
| **Button** | `ButtonURL` | ⚪ | Hidden | "https://..." |
| **Info Box** | `InfoTitle` | ⚪ | Hidden | "Kenapa Perlu Verifikasi?" |
| **Info Box** | `InfoContent` | ⚪ | Hidden | Konten info |
| **Highlight** | `HighlightText` | ⚪ | Hidden | "Link berlaku 24 jam!" |
| **Note** | `Note` | ⚪ | Default text | Custom note atau default |
| **Footer** | `Year` | ✅ | - | "2025" |

---

## 🎯 Field Behavior

### Required Fields (Wajib)
```go
Name    string  // Jika kosong → Error/tidak tampil baik
Subject string  // Untuk subject email
Message string  // Jika kosong → Body kosong
Year    int     // Jika 0 → Tampil "0" (aneh)
```

### Optional Fields (Akan Hide Jika Kosong)
```go
ButtonText     string  // Jika kosong → Button tidak muncul
ButtonURL      string  // Jika kosong → Button tidak muncul
InfoTitle      string  // Jika kosong → Info box tidak muncul
InfoContent    string  // Jika kosong → Info box tidak muncul
HighlightText  string  // Jika kosong → Highlight tidak muncul
Note           string  // Jika kosong → Pakai default text
```

---

## 📝 Contoh Kombinasi Field

### 1. Minimal (Hanya Required)
```go
{
    Name:    "John",
    Subject: "Hello",
    Message: "Welcome!",
    Year:    2025,
}
```
**Hasil:** Email sederhana tanpa button/info/highlight

---

### 2. Dengan Button
```go
{
    Name:       "John",
    Subject:    "Verify Email",
    Message:    "Please verify your email",
    Year:       2025,
    ButtonText: "Verify Now",
    ButtonURL:  "https://example.com/verify",
}
```
**Hasil:** Email + button biru cantik

---

### 3. Dengan Info Box
```go
{
    Name:        "John",
    Subject:     "Info",
    Message:     "Here's some information",
    Year:        2025,
    InfoTitle:   "Did You Know?",
    InfoContent: "This is useful info",
}
```
**Hasil:** Email + box info dengan border biru

---

### 4. Full Features
```go
{
    Name:          "John",
    Subject:       "Important",
    Message:       "Read this carefully",
    Year:          2025,
    ButtonText:    "Take Action",
    ButtonURL:     "https://example.com",
    InfoTitle:     "Why This Matters",
    InfoContent:   "Explanation here",
    HighlightText: "Urgent! Act now",
    Note:          "Custom note here",
}
```
**Hasil:** Email lengkap dengan semua elemen

---

## 🎨 Design Preview

### Colors Used
```
Header Background: Linear gradient #4facfe → #00f2fe (Blue)
Button: Linear gradient #667eea → #764ba2 (Purple)
Info Box: Linear gradient #f5f7fa → #c3cfe2 (Gray)
Highlight Box: #fff9e6 background, #ffd700 border (Yellow)
Footer: #2c3e50 (Dark gray)
Text: #333333 (Dark), #555555 (Medium), #666666 (Light)
```

### Typography
```
Header Title: 26px, bold
Greeting: 24px, bold
Message: 16px, normal
Button: 16px, bold
Info Title: 18px
Footer: 14px
```

### Spacing
```
Content Padding: 40px (top/bottom), 30px (left/right)
Button Padding: 16px (vertical), 45px (horizontal)
Info Box Padding: 25px
Margins: 20-35px between sections
```

---

## 📱 Responsive Design

Template otomatis responsive untuk mobile:

```
Desktop (>600px)         Mobile (≤600px)
─────────────────        ───────────────
Logo: 180px              Logo: 150px
Title: 26px              Title: 22px
Greeting: 24px           Greeting: 20px
Button: 16px             Button: 15px
Padding: 40px            Padding: 30px
```

---

## 🔧 Customization Points

Jika ingin customize template, edit file `template.gohtml`:

### 1. Ubah Warna
```css
/* Di section <style> */
.header {
    background: linear-gradient(135deg, #YOUR_COLOR_1, #YOUR_COLOR_2);
}

.button {
    background: linear-gradient(135deg, #YOUR_COLOR_3, #YOUR_COLOR_4);
}
```

### 2. Ubah Logo
Ganti file: `pkg/email/logo-nyinauni-golang.png`

### 3. Ubah Font
```css
body {
    font-family: 'Your-Font', Arial, sans-serif;
}
```

### 4. Tambah Section Baru
```html
{{if .YourNewField}}
<div class="your-custom-class">
    {{.YourNewField}}
</div>
{{end}}
```

---

## 🎭 Template Variants

### Standard Email (Default)
- Full width
- Gradient header
- White content area
- Dark footer

### Notification Email (Suggested)
- Minimal button
- Focus on message
- Optional highlight

### Action Email (Suggested)
- Prominent button
- Clear CTA
- Time-sensitive highlight

### Info Email (Suggested)
- Large info box
- Multiple sections
- Detailed content

---

## ⚡ Performance

### Email Size
```
HTML: ~10KB (template only)
With Logo: ~25KB (total with embedded image)
Load Time: <1s (email clients)
```

### Compatibility
```
✅ Gmail
✅ Outlook
✅ Apple Mail
✅ Yahoo Mail
✅ Mobile apps
⚠️ Very old email clients (may not show gradients)
```

---

## 📚 Best Practices

### ✅ DO
- Keep message concise (3-5 paragraphs max)
- Use button for one primary action
- Test on multiple email clients
- Use descriptive button text
- Keep file size under 100KB

### ❌ DON'T
- Don't use too many colors
- Don't make text too small
- Don't add too many buttons (max 2)
- Don't forget mobile users
- Don't send without testing

---

**Happy Designing! 🎨**

