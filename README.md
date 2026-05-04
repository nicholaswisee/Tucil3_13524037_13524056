# Tucil 3: Ice Sliding Puzzle Solver

## Deskripsi

Aplikasi desktop GUI untuk menyelesaikan *Ice Sliding Puzzle* menggunakan algoritma pencarian (UCS, GBFS, A*). GUI dibuat dengan [Fyne](https://fyne.io/) dan terpisah sepenuhnya dari mesin solver sehingga mudah dihubungkan dengan controller.

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
- Library grafis untuk Fyne (biasanya sudah tersedia di Linux Desktop). Jika terjadi error build grafis, instal:
  ```bash
  # Debian/Ubuntu
  sudo apt-get install libgl1-mesa-dev xorg-dev

  # Fedora
  sudo dnf install libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libGL-devel

  # Arch
  sudo pacman -S mesa libxrandr libxinerama libxcursor libxi
  ```

## Setup

```bash
cd src
go mod tidy
```

## Cara Menjalankan Aplikasi

```bash
cd src
go run main.go
```

Akan muncul jendela GUI dengan:
- **Panel kiri**: pemilih algoritma, pemilih heuristic, tombol import/export konfigurasi, tombol run, kontrol playback (step forward/backward), dan statistik (waktu & iterasi).
- **Panel kanan**: visualisasi papan permainan dengan token pemain, trail jalur, dan angka.

## Cara Menjalankan Test

```bash
cd src
go test ./gui -v
```

Test mencakup:
- `TestViewState_SetMap` & `TestViewState_Playback` — playback logic
- `TestBoardRenderer_Creation` & `TestBoardRenderer_Draw*` — rendering papan
- `TestLeftPanel_Creation`, `TestLeftPanel_SetStats`, `TestLeftPanel_SetStepLabel` — widget panel kiri

## Menghubungkan Controller

Semua widget di `LeftPanel` dan method di `ViewState` bersifat **publik** sehingga controller dapat mengaksesnya langsung:

```go
mw := gui.NewMainWindow()

// 1. Callback Import
mw.LeftPanel.ImportBtn.OnTapped = func() {
    // Buka file dialog, parse ke MapData, lalu:
    mw.State.SetMap(mapData)
    mw.BoardRenderer.Refresh()
}

// 2. Callback Run
mw.LeftPanel.RunBtn.OnTapped = func() {
    // Jalankan solver, lalu:
    mw.State.SetResult(solverResult)
    mw.LeftPanel.SetStats(solverResult.TimeMs, solverResult.NodesEval)
    mw.LeftPanel.SetStepLabel(0, len(solverResult.PathHistory)-1)
    mw.BoardRenderer.Refresh()
}

// 3. Callback Playback Step Forward
mw.LeftPanel.NextStepBtn.OnTapped = func() {
    if mw.State.StepForward() {
        mw.LeftPanel.SetStepLabel(mw.State.CurrentStep, len(mw.State.Result.PathHistory)-1)
        mw.BoardRenderer.Refresh()
    }
}
```

Lihat field publik di `gui/leftpanel.go` dan method di `gui/viewstate.go` untuk seluruh API yang tersedia.

## Format File Konfigurasi

File `.txt` berisi:
1. Baris pertama: `N M` (dimensi papan)
2. `N` baris berikutnya: grid karakter (`X`, `*`, `L`, `Z`, `O`, `0..9`)
3. `N` baris berikutnya: matriks biaya (integer)

Contoh ada di `test/input/input1.txt`.
