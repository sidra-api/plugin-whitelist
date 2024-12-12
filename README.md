# Plugin Whitelist

## Deskripsi
Plugin Whitelist digunakan untuk memfilter permintaan berdasarkan alamat IP pada Sidra Api. Plugin ini memastikan hanya IP yang diizinkan dalam daftar whitelist yang dapat mengakses layanan backend.

---

## Cara Kerja
1. **Pemeriksaan IP**
   - Plugin memeriksa IP klien dari header berikut:
     - `X-Real-Ip`
     - `X-Forwarded-For`
     - `Remote-Addr`
   - Jika IP klien cocok dengan salah satu IP dalam daftar whitelist, akses akan diizinkan.

2. **Respon**
   - Jika IP diizinkan:
     - Status: `200 OK`
     - Body: "Allowed"
   - Jika IP tidak diizinkan:
     - Status: `403 Forbidden`
     - Body: "IP not allowed"

---

## Konfigurasi
- **Daftar Whitelist**
  - Dapat dikonfigurasi langsung pada file `main.go`:
    ```go
    var allowedIPs = map[string]bool{
        "192.168.1.1": true,
        "192.168.1.2": true,
    }
    ```

---

## Cara Menjalankan
1. Pastikan Anda sudah menginstal **Sidra Api**.
2. Tambahkan plugin ini ke direktori `plugins/whitelist/main.go` pada Sidra Api.
3. Kompilasi dan jalankan Sidra Api.
4. Plugin akan otomatis terhubung melalui UNIX socket pada path `/tmp/whitelist.sock`.

---

## Pengujian

### Endpoint
- **URL**: Endpoint mana saja yang dikonfigurasi untuk melewati plugin Whitelist.

### Langkah Pengujian
1. Kirim request dari IP yang diizinkan atau tidak diizinkan.
2. Respons yang diharapkan:
   - Jika IP diizinkan:
     - Status: `200 OK`
     - Body: "Allowed"
   - Jika IP tidak diizinkan:
     - Status: `403 Forbidden`
     - Body: "IP not allowed"

---

## Catatan Penting
- **Keamanan**: Pastikan hanya IP yang sah dimasukkan dalam whitelist.
- **Header IP**: Pastikan header IP seperti `X-Real-Ip` atau `X-Forwarded-For` dikirim oleh proxy atau load balancer sebelum mencapai Sidra Gateway.

---

## Lisensi
MIT License
