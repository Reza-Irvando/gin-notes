# Notes Management System - Backend

REST API modern untuk mengelola catatan dengan autentikasi pengguna dan otorisasi, dibangun dengan framework **Gin** dan database **MySQL**.

## Fitur Utama

- ✅ Pendaftaran dan Login pengguna dengan JWT
- ✅ Logout pengguna
- ✅ Manajemen catatan (Create, Read, Update, Delete) dengan soft delete
- ✅ Kategori untuk mengorganisir catatan
- ✅ Tag untuk kategorisasi yang lebih fleksibel
- ✅ Favorite notes untuk menandai catatan penting
- ✅ Activity log untuk tracking semua aktivitas pengguna
- ✅ Pencarian dan filter catatan
- ✅ Pagination untuk daftar catatan
- ✅ Keamanan password dengan bcrypt
- ✅ Isolasi data per pengguna

## Stack Teknologi

- **Gin** - Framework web
- **GORM** - ORM untuk database
- **MySQL** - Database
- **JWT** - Autentikasi berbasis token
- **bcrypt** - Hashing password

---

## Setup & Konfigurasi Awal

### Prerequisites

- Go 1.16+
- MySQL 5.7+
- Git

### 1. Clone Repository

```bash
git clone https://github.com/Reza-Irvando/gin-notes
cd gin-notes
```

### 2. Setup Database

#### Login ke MySQL:

```bash
mysql -u root -p
```

#### Jalankan schema:

```bash
source schema.sql
```

#### Verifikasi database:

```bash
USE go-notes;
SHOW TABLES;
```

### 3. Konfigurasi Database Connection

Edit file `configs/config.go` dan sesuaikan connection string:

```go
// Ubah nilai berikut dengan kredensial MySQL Anda
"root:YOUR_PASSWORD@tcp(127.0.0.1:3306)/go-notes?charset=utf8&parseTime=True&loc=Local"
```

**Parameter:**

- `root` - Username MySQL
- `YOUR_PASSWORD` - Password MySQL Anda
- `127.0.0.1:3306` - Host dan port MySQL
- `go-notes` - Nama database

### 4. Install Dependencies

```bash
go mod download
go mod tidy
```

### 5. Build dan Run

#### Build executable:

```bash
go build -o gin-notes.exe
./gin-notes.exe
```

#### Atau jalankan langsung:

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

---

## Konfigurasi Postman

### Import Collection & Environment

1. Buka Postman
2. Import collection: `postman/collections/Notes_Management_System.postman_collection.json`
3. Atau import dari link url _[Postman Notes Collection](https://red-firefly-423695.postman.co/workspace/gin-notes~a96d39dc-ee1d-4baf-a1d2-38c733eab583/collection/56061881-282b8046-983e-46ca-8862-e98f3d5c1afb?action=share&source=copy-link&creator=56061881)_.

### Workflow Testing

1. **Register** - Buat akun baru
2. **Login** - Dapatkan JWT token
3. **Copy token** ke environment variable
4. **Test endpoints** - Gunakan token untuk akses endpoint yang terlindungi

### API Documentation

Dapat dilihat pada Postman Overview -> Documentation
