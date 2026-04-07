# 🧹 Git-Cleaner Pro

![Go Version](https://img.shields.io/github/go-mod/go-version/Nando-suka/git-cleaner-pro)
![Build Status](https://img.shields.io/github/actions/workflow/status/Nando-suka/git-cleaner-pro/ci.yml?branch=main)
![License](https://img.shields.io/github/license/Nando-suka/git-cleaner-pro)

**Alat CLI yang aman dan interaktif untuk membersihkan cabang Git yang sudah di-merge.**

Bosan melihat puluhan cabang lama yang sudah di-merge namun masih berserakan di repositori lokal dan remote Anda?  
Git-Cleaner Pro hadir untuk merapikan repositori Anda dengan satu perintah singkat, lengkap dengan konfirmasi interaktif agar Anda tidak salah hapus.

---

## ✨ Fitur Utama

- 🔍 **Deteksi Otomatis** – Menemukan cabang dasar (`main`/`master`) secara cerdas tanpa konfigurasi tambahan.
- 🗑️ **Hapus Lokal & Remote** – Bersihkan cabang di mesin lokal **dan** remote (`origin`) sekaligus.
- 🛡️ **Konfirmasi Interaktif** – Setiap cabang yang akan dihapus ditampilkan terlebih dahulu dan meminta persetujuan Anda.
- ⚡ **Native Go** – Dibangun dengan pustaka [go-git](https://github.com/go-git/go-git) murni, tidak bergantung pada instalasi `git` di sistem.
- 🎨 **Tampilan Rapi** – Output berwarna dan simbol yang memudahkan pembacaan.

---

## 📦 Instalasi

### Via `go install` (Rekomendasi)

```bash
go install github.com/Nando-suka/git-cleaner-pro@latest

Pastikan $GOPATH/bin sudah ada di PATH Anda.

Build dari Source

git clone https://github.com/<username-anda>/git-cleaner-pro.git
cd git-cleaner-pro
go build -o git-cleaner-pro .
sudo mv git-cleaner-pro /usr/local/bin   # opsional

🚀 Cara Penggunaan

Masuk ke direktori repositori Git mana pun, lalu jalankan:
bash

git-cleaner-pro clean

Contoh Sesi Interaktif
text

🔍 Memeriksa cabang yang sudah di-merge ke 'main'...

📋 Ditemukan 3 cabang yang sudah di-merge:
  - feature/login (lokal)
  - bugfix/header (lokal)
  - chore/update-deps (remote)

Hapus 3 cabang yang terdaftar? (y/N) y

🗑️  Menghapus cabang...
  ✅ Berhasil menghapus feature/login
  ✅ Berhasil menghapus bugfix/header
  ✅ Berhasil menghapus origin/chore/update-deps

🎉 Pembersihan selesai!

Opsi Lanjutan
Flag	Singkatan	Deskripsi
--target <branch>	-t	Tentukan cabang target secara manual (contoh: develop)
--remote	-r	Sertakan cabang remote dalam pembersihan
--yes	-y	Lewati konfirmasi (otomatis hapus semua)
bash

# Hapus semua cabang lokal & remote yang sudah di-merge ke 'develop' tanpa konfirmasi
git-cleaner-pro clean -t develop -r -y

🧪 Menjalankan Pengujian
bash

go test -v ./...

🤝 Kontribusi

Kontribusi sangat diterima! Jika Anda ingin membantu meningkatkan Git-Cleaner Pro, silakan ikuti langkah berikut:

    Fork repositori ini.

    Buat branch fitur (git checkout -b fitur-keren).

    Commit perubahan Anda (git commit -m 'Menambahkan fitur keren').

    Push ke branch (git push origin fitur-keren).

    Buat Pull Request baru.

Pastikan kode Anda lolos go test dan golangci-lint run sebelum mengirim PR.