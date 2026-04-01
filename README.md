# Mindful Journal Bot

Aplikasi bot Telegram sumber terbuka (*open-source*) untuk mencatat jurnal harian, memantau produktivitas, dan memberikan dukungan psikologis ringan. Proyek ini dibangun menggunakan **Golang**, didukung kecerdasan buatan **Gemini AI**, dan menggunakan basis data PostgreSQL dari **Supabase**.

Sistem ini dirancang dengan pendekatan **Empathetic Conversational AI**, dioptimalkan khusus untuk mendampingi individu yang sedang berada pada fase *quarter-life crisis*. Bot dapat di-*hosting* secara mandiri (*self-hosting*) menggunakan layanan tingkat gratis (*free-tier*) dari penyedia pihak ketiga.

## Fitur Utama
* **Jurnal Harian:** Ruang aman untuk mencatat keluh kesah, progres, atau sekadar bercerita tanpa adanya tuntutan produktivitas.
* **AI Companion (Empathetic Mode):** Merespons layaknya teman pendengar yang hangat, suportif, dan netral secara gender. AI telah diinstruksikan untuk memvalidasi perasaan pengguna dan menghindari *toxic positivity*.
* **Sapaan Terjadwal Dinamis (Context-Aware):**
  * **Pagi (05:00 WIB):** AI membaca riwayat kegiatan hari sebelumnya untuk menyusun ucapan selamat pagi yang relevan dan menyemangati.
  * **Malam (20:00 WIB):** AI mengevaluasi catatan hari ini untuk memberikan apresiasi atau sekadar menanyakan kabar.
* **Rekapitulasi Jangka Panjang:** Mampu menyusun wawasan (*insight*) perkembangan secara mingguan dan menyimpan arsip rangkuman secara permanen.
* **Manajemen Penyimpanan Otomatis:** Sistem akan memperingatkan apabila basis data mendekati batas kapasitas (contoh: batas 500MB pada *Free Tier*). Dilengkapi *inline keyboard* untuk membersihkan log pesan mentah lama tanpa menghapus wawasan perkembangan utama.
* **Privasi Penuh:** Fitur pemusnahan data (*nuclear option*) melalui perintah `/stop` yang akan menghapus seluruh identitas dan riwayat pengguna secara permanen dari basis data.

## Prasyarat
Untuk menjalankan bot ini secara mandiri, diperlukan tiga buah kredensial (API Key):
1. **Telegram Bot Token**
2. **Gemini API Key**
3. **Supabase Database URL**

---

## Panduan Instalasi dan Konfigurasi

### 1. Konfigurasi Telegram
1. Buka Telegram dan cari **@BotFather**.
2. Kirim perintah `/newbot` dan ikuti petunjuk untuk menentukan nama dan *username* bot.
3. Salin **HTTP API Token** yang diberikan.

### 2. Konfigurasi Supabase (Database)
1. Buat akun dan proyek baru di [Supabase](https://supabase.com).
2. Setelah proyek aktif, masuk ke menu **SQL Editor**.
3. Salin isi dari file `schema.sql` pada repositori ini, tempelkan ke SQL Editor, lalu jalankan (*Run*).
4. Masuk ke **Project Settings -> Database**.
5. Salin **Connection String** (Pilih format *URI* dan pastikan mode *Connection Pooling* aktif).

### 3. Konfigurasi Gemini AI
1. Kunjungi [Google AI Studio](https://aistudio.google.com/).
2. Buat proyek baru dan hasilkan **API Key**.

### 4. Menjalankan Aplikasi secara Lokal
1. Lakukan kloning repositori ini.
2. Salin file `.env.example` menjadi `.env` dan masukkan seluruh kredensial yang telah dikumpulkan.
3. Unduh modul dependensi Golang:
   ```bash
   go mod tidy
   ```
4. Jalankan aplikasi:
   ```bash
   go run cmd/bot/main.go
   ```

### 5. Panduan *Deployment* (Render / Cloud Services)
Repositori ini telah dilengkapi dengan `Dockerfile` dan server HTTP *dummy* untuk memfasilitasi *deployment* pada layanan awan agar bot tetap aktif 24/7.
1. Hubungkan repositori GitHub ke platform seperti [Render](https://render.com).
2. Buat layanan **Web Service** baru.
3. Pilih *runtime* **Docker**.
4. Masukkan seluruh variabel lingkungan (`TELEGRAM_BOT_TOKEN`, `GEMINI_API_KEY`, `SUPABASE_DB_URL`) pada menu **Environment Variables**.
5. Lakukan *Deploy*.

## Daftar Perintah Bot
* `/start` - Menampilkan informasi sistem dan daftar perintah.
* `/mulai_percakapan` - Memulai obrolan atau pencatatan hari ini.
* `/hari_ini` - Melihat seluruh catatan yang telah dimasukkan pada hari berjalan.
* `/rekap_mingguan` - Menyusun rangkuman progres selama satu minggu terakhir.
* `/rekap_bulanan` - Melihat *insight* jangka panjang yang disusun dari arsip permanen.
* `/stop` - Menghapus seluruh data pengguna dari sistem secara permanen.

## Manajemen Data
Aplikasi ini menggunakan `Telegram_ID` sebagai pengidentifikasi utama. Melalui relasi *Cascade* pada PostgreSQL, pencabutan akses pengguna akan memastikan seluruh entri yang berkaitan (pesan mentah maupun rangkuman) terhapus tanpa sisa.