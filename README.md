# GG Store

Aplikasi CLI (Command Line Interface) untuk simulasi toko penjualan game digital berbasis *license key*. Dibangun dengan Go menggunakan arsitektur berlapis (layered architecture) dan MySQL sebagai database.

## Fitur

### Customer
- Registrasi & Login
- Melihat katalog game
- Menambahkan game ke keranjang belanja
- Melihat isi keranjang belanja
- Checkout / membeli game di keranjang
- Melihat riwayat pesanan
- Melihat detail pesanan (termasuk license key yang didapat)
- Top-up saldo wallet

### Admin
- Melihat katalog game
- Menambahkan stok license key baru untuk sebuah game

## Arsitektur

Proyek ini mengikuti pola **layered architecture**:

```
cli/        -> Menu & interaksi terminal (Auth, Admin, Customer)
handler/    -> Menjembatani CLI dengan service, memformat I/O
service/    -> Business logic
repository/ -> Akses data ke database (MySQL)
entity/     -> Struct model data (User, Game, Cart, Order, dst.)
config/     -> Koneksi database
database/   -> Skrip SQL (schema & seed)
```

Alur request: `CLI -> Handler -> Service -> Repository -> Database`

## Struktur Direktori

```
game-store/
├── cli/
│   ├── admin.go          # Menu & handler untuk role admin
│   ├── auth.go           # Menu login/register
│   ├── cli.go            # Entry point CLI, routing berdasarkan role
│   └── customer.go       # Menu & handler untuk role customer
├── config/
│   └── database.go       # Koneksi ke MySQL via DB_URL
├── database/
│   ├── schema.sql        # Struktur tabel database
│   └── seed.sql          # Data awal (seed) untuk testing
├── entity/
│   └── models.go         # Struct: User, Game, Cart, Order, dll.
├── handler/
│   ├── cart_handler.go
│   ├── game_handler.go
│   ├── order_handler.go
│   └── user_handler.go
├── repository/
│   ├── cart_repository.go
│   ├── game_repository.go
│   ├── order_repository.go
│   └── user_repository.go
├── service/
│   ├── cart_service.go
│   ├── game_service.go
│   ├── order_service.go
│   └── user_service.go
├── main.go               # Wiring dependency & start aplikasi
└── go.mod
```

## Persyaratan

- Go 1.21 atau lebih baru
- MySQL (lokal atau remote)

## Instalasi & Menjalankan

1. Clone repository dan masuk ke direktori proyek:
   ```bash
   git clone <https://github.com/BillsNorth/game-store.git>
   cd game-store
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Siapkan database MySQL, lalu jalankan skrip berikut untuk membuat tabel dan (opsional) data awal:
   ```bash
   mysql -u <user> -p < database/schema.sql
   mysql -u <user> -p < database/seed.sql
   ```

4. Buat file `.env` di root proyek dengan koneksi database MySQL :
   ```env
   DB_URL="user:password@tcp(host:port)/game_store?parseTime=true"
   ```

5. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## Menjalankan Test

Test unit tersedia untuk seluruh handler (menggunakan mock service, tanpa koneksi database):

```bash
go test ./...
```

Untuk melihat detail tiap test case:

```bash
go test ./handler/... -v
```

## Alur Penggunaan Singkat

1. Jalankan aplikasi, lalu **Register** akun baru atau **Login** dengan akun yang sudah ada.
2. Role `customer` akan diarahkan ke dashboard customer: browse game, tambah ke keranjang, checkout, top-up saldo, dan cek riwayat pesanan.
3. Role `admin` akan diarahkan ke dashboard admin: browse game dan menambahkan license key baru.
4. Saat checkout, saldo wallet akan dipakai untuk membayar total keranjang dan license key yang tersedia akan dialokasikan ke order.


