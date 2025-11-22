package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type PostEntity struct {
	ID       string
	Platform string
	Caption  string
	Title    string
	Comments []string
}

func IsEligible(post PostEntity, brand, keyword string) bool {

	brand = strings.ToLower(brand)
	keyword = strings.ToLower(keyword)

	containsBoth := func(text string) bool {
		text = strings.ToLower(text)
		return strings.Contains(text, brand) && strings.Contains(text, keyword)
	}

	if containsBoth(post.Caption) {
		return true
	}

	if post.Platform == "tiktok" && containsBoth(post.Title) {
		return true
	}

	for _, comment := range post.Comments {
		if containsBoth(comment) {
			return true
		}
	}

	return false
}

func AnalyzeSentiment(text string) float64 {

	time.Sleep(50 * time.Millisecond)
	return 0.85
}

func ProcessPostWithRetry(ctx context.Context, post PostEntity, brand, keyword string, wg *sync.WaitGroup) {
	defer wg.Done()

	maxRetries := 3
	for i := 0; i <= maxRetries; i++ {

		if !IsEligible(post, brand, keyword) {

			fmt.Printf("[SKIP] Post %s: Kriteria tidak terpenuhi\n", post.ID)
			return
		}

		if rand.Float32() < 0.1 {
			fmt.Printf("[ERR] Post %s: Network error (Percobaan %d/%d)\n", post.ID, i+1, maxRetries+1)

			backoff := time.Duration(100*(1<<i)) * time.Millisecond
			jitter := time.Duration(rand.Intn(50)) * time.Millisecond
			sleepTime := backoff + jitter

			time.Sleep(sleepTime)
			continue
		}

		score := AnalyzeSentiment(post.Caption)

		fmt.Printf("[SUCCESS] Post %s disimpan. Sentimen: %.2f\n", post.ID, score)
		return
	}

	fmt.Printf("[DLQ] Post %s dipindahkan ke Dead Letter Queue setelah gagal berulang kali\n", post.ID)
}

func main() {

	brand := "Indomilk"
	keyword := "Susu"

	posts := []PostEntity{
		{ID: "1", Platform: "instagram", Caption: "Minum Indomilk rasa susu enak", Comments: []string{}},
		{ID: "2", Platform: "tiktok", Caption: "Joget asik", Title: "Indomilk Susu murni", Comments: []string{}},
		{ID: "3", Platform: "instagram", Caption: "Cuma foto", Comments: []string{"Wah indomilk susunya segar"}},
		{ID: "4", Platform: "instagram", Caption: "Ini Indomilk", Comments: []string{"Suka susu sapi"}},
	}

	var wg sync.WaitGroup
	fmt.Println("Memulai Pipeline...")

	for _, p := range posts {
		wg.Add(1)
		go ProcessPostWithRetry(context.Background(), p, brand, keyword, &wg)
	}

	wg.Wait()
	fmt.Println("Pipeline Selesai.")
}
