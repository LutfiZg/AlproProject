package main

import (
	"encoding/gob" // Menyimpan dan memuat data dalam bentuk binary (GOB format) ke/dari file.
	"fmt" // Menampilkan teks ke terminal dan menerima input dari pengguna.
	"os" // Berinteraksi dengan sistem file: buat, buka, dan tutup file.
	"os/exec" // Menjalankan perintah sistem (terminal/shell) dari dalam program Go.
	//"math" // Operasi matematika yang lebih kompleks.
	"strings" // Manipulasi dan pencarian string.
)

const NMAX = 1000

const (
	Purple     = "\033[35m"
	Reset      = "\033[0m"
	fileAkun   = "akun.gob"
	fileBarang = "Barang.gob"
)

type Akun struct {
    Username string
    Password string
}

type Barang struct {
    Kode  string
    Nama  string
    Stok  int
    Harga int
}

type tabBarang [NMAX]Barang
type tabAkun [NMAX]Akun

func main() {
	akun, totAcc := loadAkun(fileAkun)      // ← muat akun
	data, jumBrg := loadBarang(fileBarang)  // ← muat Barang

	if authSystem(&akun, &totAcc) {         // ← login
		menu(&data, &jumBrg)               // ← menu Stok-Barang
	}

	saveAkun(fileAkun, akun, totAcc)        // ← simpan akun
	saveBarang(fileBarang, data, jumBrg)    // ← simpan Barang
}





// Function authSystem untuk melakukan autentikasi login pengguna.
func authSystem(A *tabAkun, n *int) bool {
	end := false
	for !end {
		clearScreen()
		fmt.Println(Purple + "┌────────────────────────────┐" + Reset)
		fmt.Println(Purple + "│          MENU LOGIN        │" + Reset)
		fmt.Println(Purple + "├────────────────────────────┤" + Reset)
		fmt.Println(Purple + "│ 1. Login                   │" + Reset)
		fmt.Println(Purple + "│ 2. Register                │" + Reset)
		fmt.Println(Purple + "│ 3. Lupa Kata Sandi         │" + Reset)
		fmt.Println(Purple + "│ 4. Keluar                  │" + Reset)
		fmt.Println(Purple + "└────────────────────────────┘" + Reset)
		fmt.Print("Pilih (1-4): ")
		var pilih int
		fmt.Scan(&pilih)
		fmt.Println()

		switch pilih {
		case 1: // LOGIN
			var u, p string
			fmt.Print("Username : ")
			fmt.Scan(&u)
			fmt.Print("Password : ")
			fmt.Scan(&p)
			if login(*A, *n, u, p) {
				fmt.Println("✔️  Login berhasil!")
				fmt.Println("Tekan ENTER...")
				fmt.Scanln(); fmt.Scanln()
				return true
			}
			fmt.Println("❌  Username / Password salah.")
			fmt.Println("Tekan ENTER...")
			fmt.Scanln(); fmt.Scanln()

		case 2: // REGISTER
			if *n == NMAX {
				fmt.Println("⚠️  Database akun penuh.")
			} else {
				var u, p string
				fmt.Print("Buat Username : ")
				fmt.Scan(&u)
				
				if isTaken(*A, *n, u) {
					fmt.Println("⚠️  Username sudah dipakai.")
				} else {
					fmt.Print("Buat Password : ")
					fmt.Scan(&p)
					(*A)[*n] = Akun{u, p}
					(*n)++
					saveAkun(fileAkun, *A, *n)
					fmt.Println("🎉  Registrasi berhasil!")
				}
			}
			fmt.Println("Tekan ENTER...")
			fmt.Scanln(); fmt.Scanln()

		case 3: // LUPA Password
			var u string
			fmt.Print("Masukkan Username : ")
			fmt.Scan(&u)
			idx := findUser(*A, *n, u)
			if idx == -1 {
				fmt.Println("⚠️  Username tidak ditemukan.")
			} else {
				var newPass string
				fmt.Print("Password baru : ")
				fmt.Scan(&newPass)
				(*A)[idx].Password = newPass
				saveAkun(fileAkun, *A, *n)
				fmt.Println("🔑  Password berhasil di-reset.")
			}
			fmt.Println("Tekan ENTER...")
			fmt.Scanln(); fmt.Scanln()

		case 4: // KELUAR
			end = true

		default:
			fmt.Println("Pilihan tidak valid.")
			fmt.Println("Tekan ENTER...")
			fmt.Scanln(); fmt.Scanln()
		}
	}
	return false
}

// Fungsi login menggunakan Algoritma SequentialSearch untuk cek sandi dan username.
func login(A tabAkun, n int, u, p string) bool {
	for i := 0; i < n; i++ {
		if A[i].Username == u && A[i].Password == p {
			return true
		}
	}
	return false
}

// fungsi isTaken untuk mendeteksi username yang sama.
func isTaken(A tabAkun, n int, u string) bool {
	return findUser(A, n, u) != -1
}

// Fungsi findUser untuk mencari Username
func findUser(A tabAkun, n int, u string) int {
	var i int
	for i = 0; i<n; i++ {
		if A[i].Username == u {
			return i
		}
	}
	return -1
}

// Procedure save akun menggunakan encoding/gob
func saveAkun(path string, A tabAkun, n int) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	enc.Encode(n)
	enc.Encode(A[:n])
}

// Fungsi loadAkun menggunakan encoding/gob dan os untuk membuka file.gob
func loadAkun(path string) (tabAkun, int) {
	var A tabAkun
	var n int
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		dec := gob.NewDecoder(f)
		dec.Decode(&n)
		var slice []Akun
		dec.Decode(&slice)
		copy(A[:], slice)
	}
	return A, n
}

// Procedure save akun menggunakan encoding/gob
func saveBarang(path string, A tabBarang, n int) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	enc.Encode(n)
	enc.Encode(A[:n])
}

// Fungsi loadBarang menggunakan encoding/gob dan os untuk membuka file.gob
func loadBarang(path string) (tabBarang, int) {
	var A tabBarang
	var n int
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		dec := gob.NewDecoder(f)
		dec.Decode(&n)
		var slice []Barang
		dec.Decode(&slice)
		copy(A[:], slice)
	}
	return A, n
}






// Procedure menampilkan menu
func menu(A *tabBarang, n *int) {
	var pilih int

	for {
		clearScreen()
		fmt.Println(Purple + "┌────────────────────────────────────────┐" + Reset)
		fmt.Println(Purple + "│        APLIKASI Stok Barang TOKO       │" + Reset)
		fmt.Println(Purple + "│        by: 1. Lutfi Ghifari Hibban     │" + Reset)
		fmt.Println(Purple + "│            2. Muhamad Gyan Kausal      │" + Reset)
		fmt.Println(Purple + "├────────────────────────────────────────┤" + Reset)
		fmt.Println(Purple + "│ 1. Tambah Data Barang                  │" + Reset)
		fmt.Println(Purple + "│ 2. Tampilkan Data Barang               │" + Reset)
		fmt.Println(Purple + "│ 3. Cari Barang                         │" + Reset)
		fmt.Println(Purple + "│ 4. Edit Data Barang                    │" + Reset)
		fmt.Println(Purple + "│ 5. Hapus Data Barang                   │" + Reset)
		fmt.Println(Purple + "│ 6. Urutkan Barang                      │" + Reset)
		fmt.Println(Purple + "│ 7. Hitung Total Nilai Stok             │" + Reset)
		fmt.Println(Purple + "│ 8. Statistik Barang                    │" + Reset)
		fmt.Println(Purple + "│ 9. Keluar                              │" + Reset)
		fmt.Println(Purple + "└────────────────────────────────────────┘" + Reset)
		fmt.Print("Pilih (1-8): ")
		fmt.Scan(&pilih)

		switch pilih {
		case 1:
			tambah_data(A, n)
			saveBarang(fileBarang, *A, *n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 2:
			cetak_data(*A, *n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 3:
			menuCariBarang(A, n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 4:
			var Kode string
			fmt.Print("Kode Barang apa? ")
			fmt.Scan(&Kode)
			edit_data(A, n, Kode)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 5:
			var Kode string
			fmt.Print("Masukan data yang ingin dihapus: ")
			fmt.Scan(&Kode)
			hapus_data(A, n, Kode)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 6:
			menuPengurutan(A, n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 7:
			hitung_totStok(*A, *n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 8:
			statistik_Barang(*A, *n)
			fmt.Println("\nTekan ENTER untuk kembali ke menu...")
			fmt.Scanln(); fmt.Scanln()
		case 9:
			fmt.Println("                                   @@@@#*#@@@@@@@@@:%@@ @@@@@@@@@@      @@@+@@@       ")    
			fmt.Println("                            @@@@@@@@@@++++%#*++*#@@@@@@@@@@@@@@@@@@@@@@@@@**+@@@@@@@@ ") 
			fmt.Println("                   @@@@@@@@@@-..:+@@@@*++++++++++++@@@@+++++++++++@@@++++++@*+++**#@@@") 
			fmt.Println("         @@@    @@@@@++++++@@@++++#@@#+++++#@@@@*+++++@@@@@@@@@#+++@@@@@#*@%@@*@@@    ") 
			fmt.Println("     @@@@@@@@@@@@@%++++++++@@%++++%@@#++++++@@+++++++*@@@%@@@@*+++@@@@@@++#@%++#@@@   ") 
			fmt.Println("    @@@@*++*#+#@@@@++++++#@@@@@@@@@@@+++++++*%++%*+++@@@@+:@@@++*@@@@@@*@@@@@@@*@@@   ") 
			fmt.Println(" @@@@@@*+++++++@@@@@@@+++*@@@@@@@@@@@@@@@++++@%++++@@@@@#@%:%@@@@@@  @@@@@   @@@@@@   ") 
			fmt.Println("@@@%+++++++++++++@@@@@+++*@@@@@@%++*@@@@@%+++%@@@@@@@@@@@@*%+.@@@                     ") 
			fmt.Println(" @@@*#@@@@++++++%@@@@@++*%@@@+++++++@@@@@@@@@@@@@@=:%@@@@@@#*+:-@@                    ") 
			fmt.Println("  @@@@@@++++++++#+++@@@@@@@@@*+**#@@@@@@@@@@@@@@=@*@-%#@@@@@@++-.@@@                  ") 
			fmt.Println("      @@@@@@*+++++++@@@@@@@@@@@@@@@@@@@@@@@@@@@%*@*@@=%+@@@@@@%=+.@@@                 ") 
			fmt.Println("       @@@*+++++++@@@@@@@@@@@@@:#@@@:@@@@@@@@@@#@@*@@@-@-%@@@@@@+=:.@@@               ") 
			fmt.Println("      @@@*++%@++++@@@@@@@@@@@@::@@@=-@@@@@@@@@#:@@-#@@@:-:%@@@#%@%-+:.@@@  @@@@@      ") 
			fmt.Println("       @@%++++++#@@@%@@@@@@@@#*=%%*:*##@@@@@@@-*@%#@@@@@#=#*@@@@*%%-++:.#@@@#.@@      ") 
			fmt.Println("       @@@@@@@@@@@#-@@@@@%@@@+#+@@=+@@@@@@@@%##@@#@@@@@@@@@--@@@@%*%+-+++=-:.@@@      ") 
			fmt.Println("         @@@@@@%@#+:@@@@@%@@@-@+@#%@@@@@@@@@@@@@@@@@@@@@@@@-+-%@@@@**+:..=@@@@        ") 
			fmt.Println("         @@--%%%@#+-%@@@@#@@##@@@@@@@@@@@@@@@@@@@@@@@@@@@@@.-++=%@@@#=+.@@@           ") 
			fmt.Println("         @@.+%#%@#+=%@@@@#@@#@@-. -**=-@@@@@@@@@@@@@@@@@@@@-%-*+-%@@@@=-.@@           ") 
			fmt.Println("         @@.+#*%@*+-@@@@@#@@#-.-+:---:@@@@@@@@@@@@@@*...--:%##++++=@@@@=.@@           ") 
			fmt.Println("         @@.+#*#@*-@*%@@@%#%*:=#@@#+=#:@@@@@@@@@@@%*@@@@@@@@+#+=+++-+@@@:@@@          ") 
			fmt.Println("        @@@.+**#-#@@**@@@@=%+:@*%%*:+*+#@@@@@@@@@@@@@@@@@@@@@:++=++++:%@#.@@          ") 
			fmt.Println("        @@.:=:*@@%=-+=#@@@%-%+#@-*@@@@+@@@@@@@@@@@@@@@@@@%@@@-++:=:-+=:%%.@@          ") 
			fmt.Println("        @@@.+++=@#++++-@@@@*+*#@@:%@@*@@@@@@@@@@@@@@@@@@@@@@#+++=..:.*#:@:@@          ") 
			fmt.Println("         @@=:++=*%+=-=+:%@@@#=-:@@@@@@@@@@@%*+**++=@@@@@@@@%=++++:#@@@@.+.@@          ") 
			fmt.Println("          @@.=++-#*=@@==:#@@@%=+@@@@@@@@:%%%%%%%%%%#@@@@@@.:=-...@@@@@@+.@@@          ") 
			fmt.Println("          @@@.++:::+#@@=#=:%@@@@%***+*=%@%%%%%%%%%*@@@@-:::.....-.#..*@@@@            ") 
			fmt.Println("           @@@.:*@@@@@@=**+%:+@@@@@%=:=-%@+@%%%%+@@-.=:::-:#+@@@@@@@@*.@@             ") 
			fmt.Println("           @@@.@@@@@@@@-:-:@@@@@@@@@@#@=*#=*=.:==:--%+::+-::@@@=...:=+*@@             ") 
			fmt.Println("            @@+.##+..:=+%+::=%@@@@@@@@@+***=.-::::.*=+:=::::.=@@@@@@@@@@     @@@@@    ") 
			fmt.Println("             @@@@@@@@@=.:.+.:-=..#-+*@@+%%%@:*%@@*#@##*@=:=+#-.=+#@@@@@@@@@@@@..#@@   ") 
			fmt.Println("         @@@@ @@@@@  @@@@@@@@@@*.==*==#%@@@#=+%%@@==%@@**%@@@*-=#%+....@@@#:..==.@@   ") 
			fmt.Println("        @@@@@@@@%@@@@@@@@@@@@@@@@#.+-::-::::.#%@@@@@---*@@@@=--=--@*=%......=%+:+@@   ") 
			fmt.Println("       @@@@@@  @@@@@@@@%@@@@@@@@@@@:.:::....:#+@#@@@%#+--*@@--:::-%@@#-@%@@#-@-.@@@   ") 
			fmt.Println("        @@@%@@@@@@@@@@@@@@@@@@@%@@@@@@@@.:#%%@@@#@@@+%@@+----+--:*#=..+@#@*..#@@@@    ") 
			fmt.Println("          @@@@@@@@@@@@@@@@@@@@@@@     @@:=@@*%@@+*=--*##%@@@@+#=@:..-.:.+@@@@@@       ") 
			fmt.Println("             @@@@@@@@@@@@@@@@          @@%.+==*#+=-==*######*++=%...+@@@@@           ") 
			fmt.Println("░▒▓███████▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░  ")
			fmt.Println("░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░")
			fmt.Println("░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ")
			fmt.Println(" ░▒▓██████▓▒░░▒▓████████▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░░▒▓████████▓▒░ ")
			fmt.Println("       ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ")
			fmt.Println("       ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ")
			fmt.Println("░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░    ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░")
			fmt.Println("")
		}
		if pilih == 9 {
			break
		}
	}
}

// Procedure menambahkan data Array barang
func tambah_data(A *tabBarang, n *int) {
	fmt.Println(Purple + "========================================" + Reset)
	fmt.Println(Purple + "         TAMBAH DATA Barang BARU        " + Reset)
	fmt.Println(Purple + "========================================" + Reset)
	if *n < NMAX {
		fmt.Printf("%s%-22s%s", Purple, "Masukkan Kode Barang: ", Reset)
		fmt.Scan(&A[*n].Kode)
		
		fmt.Printf("%s%-22s%s", Purple, "Masukkan Nama Barang: ", Reset)
		fmt.Scan(&A[*n].Nama)

		fmt.Printf("%s%-22s%s", Purple, "Masukkan Stok       : ", Reset)
		fmt.Scan(&A[*n].Stok)

		fmt.Printf("%s%-22s%s", Purple, "Masukkan Harga      : ", Reset)
		fmt.Scan(&A[*n].Harga)

		(*n)++
		fmt.Println(Purple + "----------------------------------------" + Reset)
		fmt.Println("   → Data Barang berhasil ditambahkan!   ")
		fmt.Println(Purple + "----------------------------------------" + Reset)
	} else {
		fmt.Println(Purple + "Data telah penuh" + Reset)
	}
}

// Procedure untuk mencetak data
func cetak_data(A tabBarang, n int) {
	fmt.Println("Daftar Barang")
	fmt.Println("================================================================================")
	fmt.Printf("%-10s | %-20s | %6s | %12s | %15s\n",
		"Kode", "Nama Barang", "Stok", "Harga", "Total Nilai")
	fmt.Println("--------------------------------------------------------------------------------")
	for i := 0; i < n; i++ {
		total := A[i].Stok * A[i].Harga
		fmt.Printf("%-10s | %-20s | %6d | %12d | %15d\n",
			A[i].Kode, A[i].Nama, A[i].Stok, A[i].Harga, total)
	}
	fmt.Println("================================================================================")
}






// Procedure untuk memunculkan menu Pencarian Barang
func menuCariBarang(A *tabBarang, n *int) {
	for {
		clearScreen()
		fmt.Println(Purple + "┌────────────────────────────────────────┐" + Reset)
		fmt.Println(Purple + "│              PENCARIAN BARANG          │" + Reset)
		fmt.Println(Purple + "├────────────────────────────────────────┤" + Reset)
		fmt.Println(Purple + "│ 1. Berdasarkan Nama  (parsial)         │" + Reset)
		fmt.Println(Purple + "│ 2. Berdasarkan Kode  (tepat)           │" + Reset)
		fmt.Println(Purple + "│ 3. Berdasarkan Rentang Harga           │" + Reset)
		fmt.Println(Purple + "│ 4. Berdasarkan Rentang Stok            │" + Reset)
		fmt.Println(Purple + "│ 5. Kembali                             │" + Reset)
		fmt.Println(Purple + "└────────────────────────────────────────┘" + Reset)
		fmt.Print("Pilih (1-5): ")
		var pilih int
		fmt.Scan(&pilih)
		
		switch pilih {
		case 1:
			var kata string
			fmt.Print("Masukkan sebagian nama barang: ")
			fmt.Scan(&kata)
			cari_NamaParsial(*A, *n, kata)
		case 2:
			var kode string
			fmt.Print("Masukkan kode barang: ")
			fmt.Scan(&kode)
			cari_kodeBarang(*A, *n, kode)
		case 3:
			var minH, maxH int
			fmt.Print("Harga minimum : ")
			fmt.Scan(&minH)
			fmt.Print("Harga maksimum: ")
			fmt.Scan(&maxH)
			cari_HargaRange(*A, *n, minH, maxH)
		case 4:
			var minS, maxS int
			fmt.Print("Stok minimum : ")
			fmt.Scan(&minS)
			fmt.Print("Stok maksimum: ")
			fmt.Scan(&maxS)
			cari_StokRange(*A, *n, minS, maxS)
		case 5:
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
		fmt.Println("\nTekan ENTER untuk kembali ...")
		fmt.Scanln(); fmt.Scanln()
	}
}

// Procedure mencari nama secara parsial (melalui substring)
func cari_NamaParsial(A tabBarang, n int, kunci string) {
	k := strings.ToLower(kunci)
	ditemukan := false

	fmt.Println("\nHasil pencarian nama:", kunci)
	fmt.Println("+----------+----------------------+--------+--------------+")
	fmt.Printf("| %-8s | %-20s | %4s | %10s |\n",
		"Kode", "Nama", "Stok", "Harga")
	fmt.Println("+----------+----------------------+--------+--------------+")
	for i := 0; i < n; i++ {
		if strings.Contains(strings.ToLower(A[i].Nama), k) {
			fmt.Printf("| %-8s | %-20s | %4d | %10d |\n",
				A[i].Kode, A[i].Nama, A[i].Stok, A[i].Harga)
			ditemukan = true
		}
	}
	if !ditemukan {
		fmt.Println("|    (tidak ada data cocok)              |")
	}
	fmt.Println("+----------+----------------------+--------+--------------+")
}

// Procedure cari barang berdasarkan kode menggunakan Algoritma SequentialSearch.
func cari_kodeBarang(A tabBarang, n int, x string) {
    // diberikan array A tabBarang, banyaknya data(n), dan kode barang yang ingin dicari(x) terdefinisi. Kode akan dicari menggunakan Algoritma SequentialSearch
    var idx, i int
	idx = -1
    for i < n && idx == -1 {
        if A[i].Kode == x {
            idx = i
        }
		i++
    }
    if idx != -1 {
        // cetak header
        fmt.Println("Data barang ditemukan!")
        fmt.Println("+----------+----------------------+--------+--------------+")
        fmt.Printf("| %-8s | %-20s | %4s | %10s |\n",
            "Kode", "Nama", "Stok", "Harga")
        fmt.Println("+----------+----------------------+--------+--------------+")
        // cetak data
        fmt.Printf("| %-8s | %-20s | %4d | %10d |\n",
            A[idx].Kode, A[idx].Nama, A[idx].Stok, A[idx].Harga)
        fmt.Println("+----------+----------------------+--------+--------------+")
    } else {
        fmt.Println("Data barang tidak ditemukan.")
    }
}

// Procedure mencari barang berdasarkan range harga
func cari_HargaRange(A tabBarang, n, minH, maxH int) {
	ditemukan := false
	fmt.Printf("\nBarang dengan harga %d – %d\n", minH, maxH)
	fmt.Println("+----------+----------------------+--------+--------------+")
	fmt.Printf("| %-8s | %-20s | %4s | %10s |\n",
		"Kode", "Nama", "Stok", "Harga")
	fmt.Println("+----------+----------------------+--------+--------------+")
	for i := 0; i < n; i++ {
		if A[i].Harga >= minH && A[i].Harga <= maxH {
			fmt.Printf("| %-8s | %-20s | %4d | %10d |\n",
				A[i].Kode, A[i].Nama, A[i].Stok, A[i].Harga)
			ditemukan = true
		}
	}
	if !ditemukan {
		fmt.Println("|    (tidak ada data cocok)              |")
	}
	fmt.Println("+----------+----------------------+--------+--------------+")
}

// Procedure mencari data berdasarkan Rentang Stok
func cari_StokRange(A tabBarang, n int, minS, maxS int) {
	var i, left, right, mid int
	var start, end int

	if n == 0 {
		fmt.Println("⚠️  Belum ada data.")
		return
	}
	pengurutan_StokMenaik(&A, n)
	left = 0
	right = n - 1
	start = n // Inisialisasi di luar jangkauan, jika tidak ada yg memenuhi
	for left <= right {
		mid = (right + left) / 2 
		if A[mid].Stok >= minS {
			start = mid       // Temukan kandidat, simpan.
			right = mid - 1   // Coba cari lagi di sebelah kiri.
		} else {
			left = mid + 1
		}
	}
	left = 0
	right = n - 1
	end = -1 // Inisialisasi di luar jangkauan, jika tidak ada yg memenuhi
	for left <= right {
		mid = (right + left) / 2
		if A[mid].Stok <= maxS {
			end = mid         // Temukan kandidat, simpan.
			left = mid + 1    // Coba cari lagi di sebelah kanan.
		} else {
			right = mid - 1
		}
	}
	fmt.Printf("\nBarang dengan stok %d – %d (binary search)\n", minS, maxS)
	fmt.Println("+----------+----------------------+--------+--------------+")
	fmt.Printf("| %-8s | %-20s | %-6s | %-12s |\n", "Kode", "Nama", "Stok", "Harga")
	fmt.Println("+----------+----------------------+--------+--------------+")
	if start > end || start == n { // Kondisi jika tidak ada data yang cocok
		fmt.Println("|              (tidak ada data cocok)                |")
	} else {
		for i = start; i <= end; i++ {
			fmt.Printf("| %-8s | %-20s | %-6d | %-12d |\n",
				A[i].Kode, A[i].Nama, A[i].Stok, A[i].Harga)
		}
	}
	fmt.Println("+----------+----------------------+--------+--------------+")
}

// Procedure edit data.
func edit_data(A *tabBarang, n *int, x string) {
	fmt.Println("========================================")
	fmt.Println("            EDIT DATA Barang            ")
	fmt.Println("========================================")

	k := seqSearchKode(A, *n, x)
	if k == -1 {
		fmt.Printf("⚠️  Data dengan Kode \"%s\" tidak ditemukan.\n", x)
	} else {
		fmt.Println("Data saat ini:")
		fmt.Printf("  Kode       : %s\n", (*A)[k].Kode)
		fmt.Printf("  Nama Barang: %s\n", (*A)[k].Nama)
		fmt.Printf("  Stok       : %d\n", (*A)[k].Stok)
		fmt.Printf("  Harga      : %d\n", (*A)[k].Harga)
		fmt.Println("----------------------------------------")

		fmt.Printf("Masukkan Kode baru   : ")
		fmt.Scan(&(*A)[k].Kode)
		fmt.Printf("Masukkan Nama baru   : ")
		fmt.Scan(&(*A)[k].Nama)
		fmt.Printf("Masukkan Stok baru   : ")
		fmt.Scan(&(*A)[k].Stok)
		fmt.Printf("Masukkan Harga baru  : ")
		fmt.Scan(&(*A)[k].Harga)

		saveBarang(fileBarang, *A, *n)
		fmt.Println("✅  Data Barang berhasil diperbarui.")
	}
}

// Procedure hapus data.
func hapus_data(A *tabBarang, n *int, x string) {
	k := seqSearchKode(A, *n, x)
	if k != -1 {
		for i := k; i < *n-1; i++ {
			(*A)[i] = (*A)[i+1]
		}
		(*n)--
		fmt.Println("Data telah berhasil dihapus!")
	} else {
		fmt.Println("Data tidak ditemukan.")
	}
	saveBarang(fileBarang, *A, *n)
}

// Fungsi mencari barang berdasarkan kode menggunakan Algoritma SequentialSearch.
func seqSearchKode(A *tabBarang, n int, KodeBarang string) int {
	for i := 0; i < n; i++ {
		if (*A)[i].Kode == KodeBarang {
			return i
		}
	}
	return -1
}






// Menampilkan menu Sorting
func menuPengurutan(A *tabBarang, n *int) {
	for {
		clearScreen()
		fmt.Println(Purple + "┌────────────────────────────────────────┐" + Reset)
		fmt.Println(Purple + "│             PENGURUTAN BARANG          │" + Reset)
		fmt.Println(Purple + "├────────────────────────────────────────┤" + Reset)
		fmt.Println(Purple + "│ 1. Stok  ↓  (selection)                │" + Reset)
		fmt.Println(Purple + "│ 2. Stok  ↑  (insertion)                │" + Reset)
		fmt.Println(Purple + "│ 3. Harga ↓  (selection)                │" + Reset)
		fmt.Println(Purple + "│ 4. Harga ↑  (insertion)                │" + Reset)
		fmt.Println(Purple + "│ 5. Total ↓ (selection)                 │" + Reset)
		fmt.Println(Purple + "│ 6. Total ↑ (insertion)                 │" + Reset)
		fmt.Println(Purple + "│ 7. Kode ↑ (insertion)                  │" + Reset)
		fmt.Println(Purple + "│ 8. Kode ↓ (selection)                  │" + Reset)
		fmt.Println(Purple + "│ 9. Kembali                             │" + Reset)
		fmt.Println(Purple + "└────────────────────────────────────────┘" + Reset)
		fmt.Print("Pilih (1-9): ")
		var pilih int
		fmt.Scan(&pilih)

		switch pilih {
        case 1: // Stok menaik  (selection sort)
            pengurutan_StokMenaik(A, *n)
            cetak_data(*A, *n)
        case 2: // Stok menurun (insertion sort)
            pengurutan_StokMenurun(A, *n)
            cetak_data(*A, *n)
        case 3: // Harga menaik (selection sort)
            pengurutan_HargaMenaik(A, *n)
            cetak_data(*A, *n)
        case 4: // Harga menurun (insertion sort)
            pengurutan_HargaMenurun(A, *n)
            cetak_data(*A, *n)
        case 5: // Total (stok×harga) menaik
            pengurutan_TotalMenaik(A, *n)
            cetak_data(*A, *n)
        case 6: // Total (stok×harga) menurun
            pengurutan_TotalMenurun(A, *n)
            cetak_data(*A, *n)
		case 7: // kode menaik (selection sort)
			pengurutan_KodeMenaik(A, *n)
			cetak_data(*A, *n)
		case 8: // kode menaik (selection sort)
			pengurutan_KodeMenurun(A, *n)
			cetak_data(*A, *n)
		case 9: // kembali ke menu utama
			return
        default:
            fmt.Println("Pilihan tidak valid.")
        }
		fmt.Println("\nTekan ENTER untuk kembali ...")
		fmt.Scanln(); fmt.Scanln()
	}
}

// Procedure pengurutan berdasarkan stok (ascending) menggunakan algoritma selectionsort.
func pengurutan_StokMenaik(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		idx := pass - 1
		for i := pass; i < n; i++ {
			if (*A)[i].Stok < (*A)[idx].Stok { // pilih stok terkecil
				idx = i
			}
		}
		temp := (*A)[idx]
		(*A)[idx] = (*A)[pass-1]
		(*A)[pass-1] = temp
	}
}

// Procedure pengurutan berdasarkan stok (descending) menggunakan algoritma selectionsort.
func pengurutan_StokMenurun(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		temp := (*A)[pass]
		idx := pass - 1
		for idx >= 0 && (*A)[idx].Stok < temp.Stok {
			(*A)[idx+1] = (*A)[idx]
			idx--
		}
		(*A)[idx+1] = temp
	}
}

// Procedure pengurutan berdasarkan harga (ascending) menggunakan algoritma selectionsort.
func pengurutan_HargaMenaik(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		idx := pass - 1
		for i := pass; i < n; i++ {
			if (*A)[i].Harga < (*A)[idx].Harga {
				idx = i
			}
		}
		temp := (*A)[idx]
		(*A)[idx] = (*A)[pass-1]
		(*A)[pass-1] = temp
	}
}

// Procedure pengurutan berdasarkan harga (descending) menggunakan algoritma insertionsort.
func pengurutan_HargaMenurun(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		temp := (*A)[pass]
		idx := pass - 1
		for idx >= 0 && (*A)[idx].Harga < temp.Harga {
			(*A)[idx+1] = (*A)[idx]
			idx--
		}
		(*A)[idx+1] = temp
	}
}

// Procedure pengurutan berdasarkan total harga (ascending) menggunakan algoritma selectionsort.
func pengurutan_TotalMenaik(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		idx := pass - 1
		for i := pass; i < n; i++ {
			if (*A)[i].Stok*(*A)[i].Harga < (*A)[idx].Stok*(*A)[idx].Harga {
				idx = i
			}
		}
		temp := (*A)[idx]
		(*A)[idx] = (*A)[pass-1]
		(*A)[pass-1] = temp
	}
}

// Procedure pengurutan berdasarkan stok (descending) menggunakan algoritma insertionsort.
func pengurutan_TotalMenurun(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		temp := (*A)[pass]
		idx := pass - 1
		for idx >= 0 && (*A)[idx].Stok*(*A)[idx].Harga < temp.Stok*temp.Harga {
			(*A)[idx+1] = (*A)[idx]
			idx--
		}
		(*A)[idx+1] = temp
	}
}

// Procedure menghitung total velue barang di toko
func hitung_totStok(A tabBarang, n int) {
	total := 0
	for i := 0; i < n; i++ {
		total += A[i].Stok * A[i].Harga
	}
	fmt.Println("==============================")
	fmt.Println("Total nilai seluruh Stok Barang: Rp", total)
	fmt.Println("==============================")
}

// Procedur mengurutkan ascending berdsarkan kode, menggunakan Algoritma selectionsort.
func pengurutan_KodeMenaik(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		idx := pass - 1
		for i := pass; i < n; i++ {
			if (*A)[i].Kode < (*A)[idx].Kode { // pilih kode terkecil
				idx = i
			}
		}
		temp := (*A)[idx]
		(*A)[idx] = (*A)[pass-1]
		(*A)[pass-1] = temp
	}
}

// Procedur mengurutkan descending berdsarkan kode, menggunakan Algoritma insertionsort.
func pengurutan_KodeMenurun(A *tabBarang, n int) {
	for pass := 1; pass < n; pass++ {
		temp := (*A)[pass]
		idx := pass - 1
		for idx >= 0 && (*A)[idx].Kode < temp.Kode {
			(*A)[idx+1] = (*A)[idx]
			idx--
		}
	(*A)[idx+1] = temp
	}
}







// Procedure Menampilkan Statistik barang
func statistik_Barang(A tabBarang, n int) {
    if n == 0 {
        fmt.Println("⚠️  Belum ada data untuk dianalisis.")
        return
    }
    var i int
    var totalHarga, totalStok int
    var maxHarga, minHarga int
    maxHarga = A[0].Harga
    minHarga = A[0].Harga
    for i = 0; i < n; i++ {
        totalHarga += A[i].Harga
        totalStok  += A[i].Stok
        if A[i].Harga > maxHarga { maxHarga = A[i].Harga }
        if A[i].Harga < minHarga { minHarga = A[i].Harga }
    }
    rataHarga := float64(totalHarga) / float64(n)
    rataStok  := float64(totalStok)  / float64(n)
    totalNilai := 0
    for i = 0; i < n; i++ {
        totalNilai += A[i].Stok * A[i].Harga
    }
    fmt.Println(Purple + "========================================" + Reset)
    fmt.Println(Purple + "            STATISTIK BARANG            " + Reset)
    fmt.Println(Purple + "========================================" + Reset)
    fmt.Printf("Jumlah data         : %d\n", n)
    fmt.Printf("Total nilai stok    : Rp %d\n", totalNilai)
    fmt.Printf("Harga rata-rata     : %.2f\n", rataHarga)
    fmt.Printf("Harga maksimum      : %d\n", maxHarga)
    fmt.Printf("Harga minimum       : %d\n", minHarga)
    fmt.Println("----------------------------------------")
    fmt.Printf("Stok rata-rata      : %.2f\n", rataStok)
    fmt.Println(Purple + "========================================" + Reset)
}

// func clear cmd menggunakan os, os/exec
func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") // hanya Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}
