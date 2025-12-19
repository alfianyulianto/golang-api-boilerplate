# Telegram Package

Package ini digunakan untuk mengirim pesan ke Telegram menggunakan Telegram Bot API.

## Setup

### 1. Membuat Bot Telegram

1. Buka [@BotFather](https://t.me/botfather) di Telegram
2. Ketik `/newbot` untuk membuat bot baru
3. Ikuti instruksi untuk memberi nama dan username bot
4. Simpan **Bot Token** yang diberikan

### 2. Mendapatkan Chat ID

**Untuk Personal Chat:**
1. Mulai chat dengan bot Anda
2. Kirim pesan apa saja ke bot
3. Buka browser dan akses: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Cari `"chat":{"id":` untuk mendapatkan chat ID Anda

**Untuk Group Chat:**
1. Tambahkan bot ke group
2. Kirim pesan di group
3. Buka: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Cari chat ID group (biasanya angka negatif)

### 3. Konfigurasi

Tambahkan konfigurasi Telegram di `config.json`:

```json
{
  "telegram": {
    "bot_token": "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz",
    "chat_id": "123456789"
  }
}
```

## Penggunaan

Package ini sudah terintegrasi dengan logrus hooks. Setiap log dengan level error, fatal, atau panic akan otomatis dikirim ke Telegram.

### Contoh Penggunaan Manual

```go
package main

import (
    "github.com/alfianyulianto/pds-service/pkg/telegram"
)

func main() {
    client := telegram.NewTelegramClient("YOUR_BOT_TOKEN", "YOUR_CHAT_ID")
    
    err := client.SendMessage("Hello from Go!")
    if err != nil {
        // handle error
    }
}
```

### Format Pesan

Pesan yang dikirim mendukung format HTML. Contoh:

```go
message := "<b>Bold Text</b>\n<i>Italic Text</i>\n<code>Code</code>"
client.SendMessage(message)
```

## Logrus Hook

Hook telegram sudah dikonfigurasi di `internal/config/logrus.go` dan akan mengirim notifikasi untuk:
- `logrus.PanicLevel`
- `logrus.FatalLevel`
- `logrus.ErrorLevel`

Jika ingin menambahkan `WarnLevel`, edit file `internal/hooks/telegram_hook.go` pada method `Levels()`.

## Troubleshooting

### Bot tidak mengirim pesan
- Pastikan bot token benar
- Pastikan chat ID benar
- Pastikan bot sudah di-start (kirim `/start` ke bot)
- Jika menggunakan group, pastikan bot sudah ditambahkan sebagai member

### Error "Forbidden: bot was blocked by the user"
- User harus melakukan `/start` ke bot terlebih dahulu

### Error "Bad Request: chat not found"
- Chat ID salah atau bot belum pernah memulai chat dengan user tersebut

