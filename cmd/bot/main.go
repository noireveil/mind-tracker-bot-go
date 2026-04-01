package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/username/mind-tracker-bot-go/internal/ai"
	"github.com/username/mind-tracker-bot-go/internal/database"
	"github.com/username/mind-tracker-bot-go/internal/handlers"
	"github.com/username/mind-tracker-bot-go/internal/scheduler"
)

func main() {
	_ = godotenv.Load()

	if err := database.ConnectDB(); err != nil {
		log.Fatal("Koneksi database gagal:", err)
	}

	if err := ai.InitGemini(); err != nil {
		log.Fatal("Inisialisasi AI gagal:", err)
	}

	if err := handlers.InitBot(); err != nil {
		log.Fatal("Inisialisasi bot gagal:", err)
	}

	scheduler.StartCronJobs()
	go handlers.StartListening()

	// Server dummy untuk memenuhi syarat health check pada layanan cloud (misal: Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Bot beroperasi secara normal.")
	})

	log.Printf("Mendengarkan port %s untuk Health Check...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("HTTP Server gagal:", err)
	}
}
