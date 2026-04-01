package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/username/mind-tracker-bot-go/internal/ai"
	"github.com/username/mind-tracker-bot-go/internal/database"
	"github.com/username/mind-tracker-bot-go/internal/handlers"
)

func StartCronJobs() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal("Gagal memuat zona waktu:", err)
	}

	s := gocron.NewScheduler(loc)

	// Jadwal Sapaan Pagi
	s.Every(1).Day().At("05:00").Do(func() {
		users, _ := database.GetAllUsers()
		for _, userID := range users {
			yesterdayLogs, _ := database.GetYesterdayLogs(userID)
			greeting, err := ai.GenerateMorningGreeting(yesterdayLogs)
			if err != nil {
				greeting = "Pagi!!! Semangat yaa buat hari ini, semoga ada hal kecil yang bikin kamu senyum hari ini 🥰"
			}
			msg := tgbotapi.NewMessage(userID, greeting)
			handlers.Bot.Send(msg)
		}
	})

	// Jadwal Pemeriksaan Kapasitas Database
	s.Every(1).Day().At("09:00").Do(func() {
		users, _ := database.GetAllUsers()
		for _, userID := range users {
			handlers.CheckStorageWarning(userID)
		}
	})

	// Jadwal Sapaan Malam
	s.Every(1).Day().At("20:00").Do(func() {
		users, _ := database.GetAllUsers()
		for _, userID := range users {
			todayLogs, _ := database.GetTodayLogs(userID)
			greeting, err := ai.GenerateEveningGreeting(todayLogs)
			if err != nil {
				greeting = "Haloo, udah jam 8 malam nih. Kalau ada cerita atau unek-unek hari ini, ceritain di sini aja yaa ❤️"
			}
			msg := tgbotapi.NewMessage(userID, greeting)
			handlers.Bot.Send(msg)
		}
	})

	s.StartAsync()
}
