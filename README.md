# IF2211 Tugas Kecil 3: Ice Sliding Puzzle Solver

## Deskripsi

Aplikasi desktop GUI untuk menyelesaikan *Ice Sliding Puzzle* menggunakan algoritma pencarian **UCS**, **GBFS**, dan **A\***. Pemain bergerak di atas papan es dan meluncur hingga menabrak dinding atau batas papan. Solver mencari jalur optimal dari posisi awal menuju tujuan, dengan visualisasi langkah-langkah solusi secara interaktif.

## Struktur Proyek

```
src/
├── main.go                # Entry point aplikasi GUI
├── go.mod / go.sum        # Dependensi Go
├── core/
│   ├── models/            # Struct data (MapData, GameState, SolverResult, dll)
│   ├── parser/            # Parser file konfigurasi papan
│   └── solver/            # Mesin solver (UCS, GBFS, A*)
└── gui/                   # Layer tampilan (View)
    ├── window.go          # MainWindow
    ├── leftpanel.go       # Panel kiri (kontrol)
    ├── boardrenderer.go   # Renderer papan permainan
    ├── viewstate.go       # State & playback logic
    └── *_test.go          # Unit test
```

## Persyaratan

- [Go](https://go.dev/) versi 1.26 atau lebih baru
- Library grafis untuk Fyne. Jika terjadi error build grafis, instal dependensi berikut sesuai distro:

  ```bash
  # Debian/Ubuntu
  sudo apt-get install libgl1-mesa-dev xorg-dev

  # Fedora
  sudo dnf install libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libGL-devel

  # Arch
  sudo pacman -S mesa libxrandr libxinerama libxcursor libxi
  ```

## Instalasi

```bash
cd src
go mod tidy
```

## Cara Menjalankan

```bash
cd src
go run main.go
```

## Cara Menggunakan

1. Klik **Import** untuk memuat file konfigurasi papan (format `.txt`).
2. Pilih **algoritma** pencarian (UCS, GBFS, atau A\*) dan **heuristik** jika diperlukan.
3. Klik **Run** untuk menjalankan solver.
4. Gunakan tombol **Next Step** / **Prev Step** atau slider untuk memutar ulang animasi solusi.
5. Statistik waktu eksekusi dan jumlah node yang dievaluasi ditampilkan di panel kiri.

### Format File Konfigurasi

File `.txt` berisi:
1. Baris pertama: `N M` (dimensi papan N baris × M kolom)
2. `N` baris berikutnya: grid karakter (`X` = dinding, `*` = jalan, `K` = posisi awal, `T` = tujuan)
3. `N` baris berikutnya: matriks biaya (integer)

Contoh file tersedia di `test/input/`.

## Cara Menjalankan Test

```bash
cd src
go test ./gui -v
```

## Author

| Nama | NIM |
|------|-----|
| Nicholas Wise Saragih Sumbayak | 13524037 |
| Reinhard Alfonzo Hutabarat | 13524056 |
