package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Result menampung hasil dari satu kali request
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

func main() {
	// 1. Menerima Input dari Terminal (Flags)
	// Contoh: ./go-blast -url=http://localhost:8080 -n=100 -c=10
	targetUrl := flag.String("url", "", "URL target untuk di-test")
	totalRequests := flag.Int("n", 100, "Total jumlah request")
	concurrency := flag.Int("c", 10, "Jumlah worker (concurrent users)")
	flag.Parse()

	// Validasi input
	if *targetUrl == "" {
		fmt.Println("❌ Error: Harap masukkan URL target. Contoh: -url=http://google.com")
		return
	}

	fmt.Printf("\n🚀 GO-BLAST STARTED!\n")
	fmt.Printf("Target: %s\n", *targetUrl)
	fmt.Printf("Requests: %d | Concurrency: %d\n\n", *totalRequests, *concurrency)

	// 2. Setup Channel dan WaitGroup
	// jobs: saluran untuk mengirim tugas ke worker
	// results: saluran untuk menerima laporan dari worker
	jobs := make(chan struct{}, *totalRequests)
	results := make(chan Result, *totalRequests)
	var wg sync.WaitGroup

	// Catat waktu mulai
	startTime := time.Now()

	// 3. Spawn Workers (Membuat Pasukan Goroutines)
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(*targetUrl, jobs, results, &wg)
	}

	// 4. Kirim Jobs (Mengisi antrian tugas)
	for i := 0; i < *totalRequests; i++ {
		jobs <- struct{}{} // Kirim sinyal kosong sebagai tanda "kerjakan 1 tugas"
	}
	close(jobs) // Tutup antrian jobs karena semua tugas sudah didaftarkan

	// 5. Tunggu semua worker selesai
	wg.Wait()
	close(results) // Tutup antrian results

	// 6. Hitung Statistik (Report)
	printReport(results, time.Since(startTime))
}

// Fungsi Worker: Ini yang dijalankan oleh banyak Goroutines sekaligus
func worker(url string, jobs <-chan struct{}, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	// Client HTTP yang reusable
	client := &http.Client{Timeout: 10 * time.Second}

	// Selama masih ada job di channel, kerjakan!
	for range jobs {
		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)

		res := Result{
			Duration: duration,
			Error:    err,
		}

		if err == nil {
			res.StatusCode = resp.StatusCode
			resp.Body.Close() // Jangan lupa tutup body biar gak memory leak
		}

		results <- res
	}
}

// Fungsi untuk mencetak laporan akhir
func printReport(results chan Result, totalTime time.Duration) {
	var (
		successCount int
		failCount    int
		totalDuration time.Duration
	)

	// Baca semua hasil dari channel results
	for res := range results {
		if res.Error != nil || res.StatusCode >= 400 {
			failCount++
		} else {
			successCount++
		}
		totalDuration += res.Duration
	}

	totalRequests := successCount + failCount
	avgDuration := totalDuration / time.Duration(totalRequests)
	rps := float64(totalRequests) / totalTime.Seconds()

	fmt.Println("📊 --- FINAL REPORT ---")
	fmt.Printf("Total Time:     %v\n", totalTime)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success:        %d ✅\n", successCount)
	fmt.Printf("Failed:         %d ❌\n", failCount)
	fmt.Printf("Avg Latency:    %v\n", avgDuration)
	fmt.Printf("Throughput:     %.2f Req/Sec\n", rps)
	fmt.Println("-----------------------")
}