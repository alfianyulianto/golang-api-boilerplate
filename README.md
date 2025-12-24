# NYINAUNI GOLANG (Golang Api Boilerplate)

NYINAUNI GOLANG adalah aplikasi backend berbasis Golang yang disediakan sebagai boilerplate untuk membangun layanan yang terstruktur, scalable, dan mudah dikembangkan. Aplikasi ini menerapkan Clean Architecture serta memanfaatkan library modern untuk performa dan maintainability yang optimal.

## 📋 Daftar Isi

- [Teknologi](#-teknologi)
- [Fitur](#-fitur)
- [Struktur Project](#-struktur-project)
- [Persyaratan](#-persyaratan)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Database Migration](#-database-migration)
- [API Documentation](#-api-documentation)
- [Arsitektur](#-arsitektur)

## 🚀 Teknologi

Project ini dibangun menggunakan teknologi berikut:

### Core Framework & Libraries

- **[Fiber v2](https://gofiber.io/)** - Web framework berkinerja tinggi yang terinspirasi dari Express.js
- **[GORM](https://gorm.io/)** - ORM library untuk Golang dengan fitur lengkap
- **[Viper](https://github.com/spf13/viper)** - Configuration management yang fleksibel
- **[Validator v10](https://github.com/go-playground/validator)** - Struct validation dengan tag-based rules
- **[Logrus](https://github.com/sirupsen/logrus)** - Structured logger untuk Golang
- **[Redis Go Client v9](https://github.com/redis/go-redis)** - Client Redis untuk caching dan session management

### Security & Authentication

- **[JWT (golang-jwt/jwt v5)](https://github.com/golang-jwt/jwt)** - JSON Web Token untuk autentikasi dan otorisasi

### External Services

- **SMTP** - Pengiriman email melalui SMTP (Gmail)
- **Telegram Bot** - Integrasi Telegram untuk notifikasi dan logging

### Database

- **MySQL** - Database relational menggunakan driver `gorm.io/driver/mysql`

### Additional Libraries

- **[UUID](https://github.com/google/uuid)** - Untuk generate unique identifier
- **[strcase](https://github.com/iancoleman/strcase)** - String case conversion utility

## ✨ Fitur

- **Autentikasi & Otorisasi**
  - Login dengan JWT token
  - Refresh token
  - Middleware untuk protected routes
  
- **Manajemen User**
  - CRUD operations untuk user
  - Role-based access control
  
- **File Storage**
  - Local file storage
  - Upload dan management file
  - Static file serving
  
- **Logging & Monitoring**
  - Structured logging dengan Logrus
  - Telegram hook untuk notifikasi critical errors
  - Application log file
  
- **Caching**
  - Redis integration untuk performa optimal
  - Session management
  
- **Email Notification**
  - SMTP integration untuk pengiriman email
  - Template-based email

## 📁 Struktur Project

```
golang-api-boilerplate/
├── cmd/
│   ├── setup/              # Setup utilities
│   └── website/            # Main application entry point
├── internal/
│   ├── config/             # Configuration setup (Viper, Fiber, GORM, Redis, dll)
│   ├── delivery/
│   │   └── http/           # HTTP handlers & controllers
│   │       ├── middleware/ # Custom middlewares
│   │       └── router/     # Route definitions
│   ├── entity/             # Database entities/models
│   ├── hooks/              # Custom hooks (Telegram hook untuk Logrus)
│   ├── model/              # Request/Response models & converters
│   ├── repository/         # Database operations layer
│   ├── usecase/            # Business logic layer
│   └── utils/              # Utility functions
├── migrations/             # Database migration files
├── pkg/
│   ├── auth/               # JWT service
│   ├── email/              # SMTP email service
│   ├── response/           # Standardized API response
│   ├── storage/            # File storage service
│   ├── telegram/           # Telegram bot integration
│   └── validators/         # Custom validators
├── uploads/                # Directory untuk uploaded files
├── config.json             # Configuration file
└── go.mod                  # Go module dependencies
```

### Arsitektur Layer

1. **Delivery Layer** - HTTP controllers, middlewares, dan routing
2. **Usecase Layer** - Business logic dan orchestration
3. **Repository Layer** - Database operations dan data access
4. **Entity Layer** - Database schema dan models
5. **Package Layer** - Reusable utilities dan services

## 📦 Persyaratan

Sebelum menjalankan aplikasi, pastikan Anda telah menginstal:

- **Go** >= 1.24.5
- **MySQL** >= 5.7 atau >= 8.0
- **Redis** >= 8.0
- **Git**

## 🛠️ Instalasi

1. **Clone repository**

```bash
git clone https://github.com/alfianyulianto/golang-api-boilerplate.git
cd golang-api-boilerplate
```

2. **Install dependencies**

```bash
go mod download
go mod tidy
```

3. **Setup konfigurasi**

```bash
cp config-example.json config.json
```

Edit file `config.json` sesuai dengan environment Anda.

## ⚙️ Konfigurasi

Edit file `config.json` dengan konfigurasi yang sesuai:

```json
{
  "app": {
    "name": "NYINAUNI GOLANG",
    "env": "development",
    "port": 8000,
    "base_url": "http://127.0.0.1"
  },
  "log": {
    "level": 6
  },
  "database": {
    "username": "root",
    "password": "your_password",
    "host": "localhost",
    "port": 3306,
    "name": "your_database_name",
    "pool": {
      "max_idle_conn": "10",
      "max_open_conn": "100",
      "conn_max_life_time": 60,
      "conn_max_idle_time": 5
    }
  },
  "storage": {
    "driver": "local"
  },
  "mail": {
    "host": "smtp.gmail.com",
    "port": 587,
    "username": "your-email@gmail.com",
    "password": "your_app_password",
    "to_address": ["recipient@example.com"],
    "from_address": "your-email@gmail.com"
  },
  "jwt": {
    "secret_key": "your_secret_key",
    "expire_duration": 600,
    "refresh_expire_duration": 604800
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "db": 0
  },
  "telegram": {
    "bot_token": "your_telegram_bot_token",
    "chat_id": "your_telegram_chat_id"
  }
}
```

### Penjelasan Konfigurasi

- **app.env**: Environment mode (`development` atau `production`)
- **app.port**: Port aplikasi akan berjalan
- **log.level**: Level logging (6 = Trace, 5 = Debug, 4 = Info, 3 = Warn, 2 = Error, 1 = Fatal, 0 = Panic)
- **database.pool**: Connection pool settings untuk optimasi koneksi database
- **jwt.expire_duration**: Durasi token dalam detik (600 = 10 menit)
- **jwt.refresh_expire_duration**: Durasi refresh token dalam detik (604800 = 7 hari)

### Setup Gmail SMTP (Optional)

Untuk menggunakan Gmail SMTP:

1. Enable 2-Factor Authentication di akun Gmail
2. Generate App Password di [Google Account Security](https://myaccount.google.com/security)
3. Gunakan App Password sebagai `mail.password` di config.json

### Setup Telegram Bot (Optional)

1. Buat bot baru melalui [@BotFather](https://t.me/botfather)
2. Dapatkan bot token dari BotFather
3. Dapatkan chat ID Anda melalui [@userinfobot](https://t.me/userinfobot)
4. Masukkan bot_token dan chat_id ke config.json

## 🚀 Menjalankan Aplikasi

### Development Mode

```bash
go run cmd/website/main.go
```

### Build & Run

```bash
# Build aplikasi
go build -o golang-api-boilerplate-service cmd/website/main.go

# Jalankan aplikasi
./golang-api-boilerplate  # Linux/Mac
golang-api-boilerplate.exe  # Windows
```

Aplikasi akan berjalan pada `http://127.0.0.1:8000`

## 🗄️ Database Migration

Project ini menggunakan migration files yang tersimpan di folder `migrations/`.

### Menjalankan Migration

Anda dapat menggunakan migration tool seperti [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
# Install migrate CLI
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Make migration files
migrate create -ext sql -dir migrations your_migration_name

# Run migrations
migrate -path migrations -database "mysql://username:password@tcp(localhost:3306)/your_migration_name" up

# Rollback migrations
migrate -path migrations -database "mysql://username:password@tcp(localhost:3306)/your_migration_name" down
```

## 📚 API Documentation

### Base URL

```
http://127.0.0.1:8000
```

### Authentication

Untuk endpoint yang memerlukan autentikasi, sertakan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

### Endpoints

#### Auth

- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login dan dapatkan JWT token
- `POST /api/auth/refresh-token` - Refresh JWT token
- `POST /api/auth/request-reset-password` - Request reset password
- `POST /api/auth/reset-password` - Reset password


#### Account

- `GET /api/auth/_current` - Get account info (Protected)
- `PUT /api/auth/account` - Update account (Protected)
- `POST /api/auth/logout` - Logout user (Protected)
- `PUT /api/auth/update-password` - Update password (Protected)

#### Users

- `GET /api/users` - Get all users (Protected)
- `GET /api/users/:id` - Get user by ID (Protected)
- `POST /api/users` - Create new user (Protected)
- `PUT /api/users/:id` - Update user (Protected)
- `DELETE /api/users/:id` - Delete user (Protected)

### Response Format

Semua response menggunakan format standar:

**Success Response:**

```json
{
  "code": 200,
  "status": "success",
  "data": {
    // response data
  }
}
```

**Error Response:**

```json
{
  "code": 400,
  "status": "error",
  "errors": "Error message or validation errors"
}
```

## 🏗️ Arsitektur

Project ini menggunakan **Clean Architecture** dengan prinsip:

- **Separation of Concerns** - Setiap layer memiliki tanggung jawab yang jelas
- **Dependency Rule** - Dependencies hanya mengalir ke dalam (dari outer ke inner layer)
- **Testability** - Setiap layer dapat di-test secara independen
- **Independence** - Business logic tidak bergantung pada framework atau tools eksternal

### Layer Flow

```
HTTP Request → Router → Middleware → Controller (Delivery) 
    → UseCase (Business Logic) → Repository (Data Access) 
    → Database
```

### Dependency Injection

Aplikasi menggunakan manual dependency injection yang di-setup di `internal/config/app.go`

## 📝 Custom Validators

Project ini menyediakan custom validators:

- **exists** - Validasi apakah record ada di database
- **unique** - Validasi unique constraint
- **image** - Validasi format dan size image
- **match_password** - Validasi password confirmation
- **size** - Validasi ukuran file

Contoh penggunaan:

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=255"`
    Email string `json:"email" validate:"required,email,unique=schools.email"`
}
```

## 🔐 Security Best Practices

- JWT token dengan expiration time
- Password hashing (implementasikan bcrypt)
- Input validation menggunakan validator
- SQL injection prevention dengan GORM
- CORS configuration
- Environment-based configuration

## 📊 Monitoring & Logging

### Log Levels

- **Trace** (6) - Informasi sangat detail
- **Debug** (5) - Informasi debugging
- **Info** (4) - Informasi umum
- **Warn** (3) - Warning message
- **Error** (2) - Error message
- **Fatal** (1) - Fatal error (aplikasi akan berhenti)
- **Panic** (0) - Panic level

### Telegram Notification

Error log level Error dan Fatal akan dikirimkan ke Telegram jika hook sudah dikonfigurasi.

### Log File

Log disimpan di file `application.log`

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 👤 Author

**Alfian Yulianto**

- GitHub: [@alfianyulianto](https://github.com/alfianyulianto)

## 🙏 Acknowledgments

- [Fiber Framework](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Viper](https://github.com/spf13/viper)
- [Logrus](https://github.com/sirupsen/logrus)
- [Go Playground Validator](https://github.com/go-playground/validator)
- [Redis Go Client](https://github.com/redis/go-redis)

## 📞 Support

Jika Anda memiliki pertanyaan atau menemukan bug, silakan buat issue di [GitHub Issues](https://github.com/alfianyulianto/golang-api-boilerplate/issues)

---

**Happy Coding! 🚀**

