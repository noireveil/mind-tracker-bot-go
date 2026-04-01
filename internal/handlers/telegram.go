package handlers

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/username/mind-tracker-bot-go/internal/ai"
	"github.com/username/mind-tracker-bot-go/internal/database"
)

var Bot *tgbotapi.BotAPI

func InitBot() error {
	var err error
	Bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return err
	}
	log.Printf("Otorisasi berhasil pada akun %s", Bot.Self.UserName)
	return nil
}

func StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go processMessage(update.Message)
		} else if update.CallbackQuery != nil {
			go processCallbackQuery(update.CallbackQuery)
		}
	}
}

func processMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	username := message.From.UserName
	text := message.Text

	database.SaveUser(userID, username)

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			welcomeText := `Ini adalah bot untuk mencatat keseharian, jurnal harian, dan ruang untuk menyimpan perkembangan personal. Sistem ini dirancang untuk mendampingi dalam mengelola produktivitas dan menjadi ruang aman tanpa penghakiman.

Contoh pencatatan:
- "Hari ini berhasil menyelesaikan revisi dokumen, tapi merasa cukup penat."
- "Tidak melakukan apa-apa hari ini, hanya berbaring dan memulihkan energi."

Daftar Perintah:
/start - Menampilkan informasi sistem ini
/mulai_percakapan - Memulai obrolan hari ini
/hari_ini - Melihat catatan progres pada hari ini
/rekap_mingguan - Melihat rangkuman progres 1 minggu terakhir
/rekap_bulanan - Melihat insight perkembangan dari data historis
/stop - Menghentikan layanan dan menghapus semua data secara permanen`
			msg := tgbotapi.NewMessage(userID, welcomeText)
			Bot.Send(msg)

		case "mulai_percakapan":
			msg := tgbotapi.NewMessage(userID, "Haii, kamu hari ini ngapain ajaa?~ Ayoo ceritakan padaku ^_^")
			Bot.Send(msg)

		case "hari_ini":
			logs, err := database.GetTodayLogs(userID)
			if err != nil || len(logs) == 0 {
				Bot.Send(tgbotapi.NewMessage(userID, "Belum terdapat catatan yang dimasukkan hari ini."))
				return
			}
			reply := "Catatan hari ini:\n"
			for i, logText := range logs {
				reply += fmt.Sprintf("%d. %s\n", i+1, logText)
			}
			Bot.Send(tgbotapi.NewMessage(userID, reply))

		case "rekap_mingguan":
			handleRekap(userID, 7, "1 minggu")

		case "rekap_bulanan":
			handleRekapPersistent(userID)

		case "stop":
			warningText := "Apakah yakin ingin menghentikan layanan bot? Seluruh data riwayat akan dihapus secara permanen dan tidak dapat dipulihkan.\n\nKirim perintah /yakin_stop untuk melanjutkan proses penghapusan."
			Bot.Send(tgbotapi.NewMessage(userID, warningText))

		case "yakin_stop":
			err := database.DeleteUser(userID)
			if err != nil {
				Bot.Send(tgbotapi.NewMessage(userID, "Terjadi kesalahan sistem saat mencoba menghapus data."))
				return
			}
			Bot.Send(tgbotapi.NewMessage(userID, "Data berhasil dihapus sepenuhnya. Bot telah dihentikan."))
		}
		return
	}

	// Alur penyimpanan dan respons AI untuk pesan reguler
	err := database.SaveProgress(userID, text)
	if err != nil {
		log.Println("Database error:", err)
		return
	}

	history, _ := database.GetRecentLogs(userID, 3)
	replyText, err := ai.GenerateResponse(text, history)
	if err != nil {
		replyText = "Sistem sedang mengalami gangguan sementara, mohon tunggu beberapa saat."
		log.Println("AI Error:", err)
	}

	Bot.Send(tgbotapi.NewMessage(userID, replyText))
}

func processCallbackQuery(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	data := callback.Data

	callbackCfg := tgbotapi.NewCallback(callback.ID, "")
	Bot.Request(callbackCfg)

	var responseText string
	switch data {
	case "purge_1":
		err := database.PurgeOldLogs(userID, 1)
		if err == nil {
			responseText = "Log pesan mentah yang lebih dari 1 bulan berhasil dihapus. Data insight perkembangan tetap dipertahankan."
		} else {
			responseText = "Terjadi kegagalan saat mencoba menghapus log."
		}
	case "purge_3":
		err := database.PurgeOldLogs(userID, 3)
		if err == nil {
			responseText = "Log pesan mentah yang lebih dari 3 bulan berhasil dihapus."
		} else {
			responseText = "Terjadi kegagalan saat mencoba menghapus log."
		}
	case "upgrade_info":
		responseText = "Untuk mengelola kapasitas, silakan akses pengaturan penyimpanan pada dashboard Supabase."
	}

	Bot.Send(tgbotapi.NewMessage(userID, responseText))
}

func handleRekap(userID int64, days int, durationName string) {
	history, _ := database.GetRecentLogs(userID, days*3)
	if len(history) == 0 {
		Bot.Send(tgbotapi.NewMessage(userID, "Data historis belum mencukupi untuk menyusun rangkuman."))
		return
	}

	replyText, err := ai.GenerateSummary(history, durationName)
	if err != nil {
		Bot.Send(tgbotapi.NewMessage(userID, "Gagal menyusun rangkuman pada saat ini."))
		return
	}

	database.SaveSummary(userID, replyText, "weekly")
	Bot.Send(tgbotapi.NewMessage(userID, replyText))
}

func handleRekapPersistent(userID int64) {
	summaries, _ := database.GetPersistentSummaries(userID)
	if len(summaries) == 0 {
		Bot.Send(tgbotapi.NewMessage(userID, "Belum terdapat riwayat rangkuman perkembangan jangka panjang."))
		return
	}

	replyText, err := ai.GenerateSummary(summaries, "jangka panjang")
	if err != nil {
		Bot.Send(tgbotapi.NewMessage(userID, "Gagal menyusun insight jangka panjang."))
		return
	}
	Bot.Send(tgbotapi.NewMessage(userID, replyText))
}

func CheckStorageWarning(userID int64) {
	size, err := database.GetDatabaseSizeMB()
	if err == nil && size > 480 {
		msgText := fmt.Sprintf("Peringatan Sistem: Kapasitas penyimpanan mencapai %.2f MB. Jika batas terpenuhi, pencatatan log baru akan terhenti.", size)
		msg := tgbotapi.NewMessage(userID, msgText)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Hapus Log > 1 Bulan", "purge_1"),
				tgbotapi.NewInlineKeyboardButtonData("Hapus Log > 3 Bulan", "purge_3"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Abaikan Peringatan", "upgrade_info"),
			),
		)
		msg.ReplyMarkup = keyboard
		Bot.Send(msg)
	}
}
