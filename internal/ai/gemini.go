package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var client *genai.Client

// Menginisialisasi klien Gemini menggunakan API Key dari environment variable
func InitGemini() error {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	var err error
	client, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	return err
}

// Menghasilkan respons balasan untuk percakapan harian
func GenerateResponse(progress string, history []string) (string, error) {
	ctx := context.Background()
	model := client.GenerativeModel("gemini-2.5-flash")

	systemPrompt := `Berperanlah sebagai pendengar yang hangat, suportif, dan empatik untuk individu di fase quarter-life crisis.
Gaya bahasa harus kasual, lembut, dan menenangkan (diperbolehkan menggunakan partikel 'ya', 'kok', 'sih' secara wajar).
Batasan Tegas: Dilarang keras menggunakan kata 'lu/gue', singkatan tidak baku, atau gaya bahasa gaul yang berlebihan (alay).
Identitas netral secara gender. Jangan berasumsi mengenai gender pengguna.
Diperbolehkan menggunakan emotikon dasar (seperti senyum) maksimal satu kali per balasan untuk menambah kesan ramah.
Dilarang menggunakan format markdown seperti tanda bintang (*) untuk menebalkan teks. Tuliskan teks biasa.
Berikan balasan singkat yang memvalidasi perasaan pengguna, baik saat sedang produktif maupun saat sedang lelah.`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	historyText := strings.Join(history, "\n- ")
	prompt := fmt.Sprintf("Konteks Riwayat (Terbaru ke Terlama):\n- %s\n\nCerita Hari Ini:\n\"%s\"\n\nInstruksi: Berikan balasan berdasarkan cerita hari ini dengan mempertimbangkan riwayat di atas.", historyText, progress)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

// Menghasilkan rangkuman wawasan (insight) jangka panjang maupun mingguan
func GenerateSummary(history []string, duration string) (string, error) {
	ctx := context.Background()
	model := client.GenerativeModel("gemini-2.5-flash")

	historyText := strings.Join(history, "\n- ")
	prompt := fmt.Sprintf("Berikut adalah catatan progres selama %s terakhir:\n- %s\n\nRangkum perjalanan ini menjadi satu paragraf insight yang suportif tanpa menggunakan format markdown.", duration, historyText)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	return extractText(resp), nil
}

// Menghasilkan sapaan pagi secara dinamis berdasarkan catatan hari sebelumnya
func GenerateMorningGreeting(yesterdayLogs []string) (string, error) {
	ctx := context.Background()
	model := client.GenerativeModel("gemini-2.5-flash")

	systemPrompt := `Berperanlah sebagai pendengar yang hangat dan suportif. Identitas netral secara gender.
Gunakan bahasa kasual yang lembut, dilarang menggunakan kata 'lu/gue' atau bahasa alay. Dilarang menggunakan format markdown.
Berikan ucapan selamat pagi dengan sisipan emotikon ramah (seperti 🥰 atau ^_^).`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	var prompt string
	if len(yesterdayLogs) == 0 {
		prompt = "Pengguna tidak mencatat kegiatan kemarin. Berikan sapaan pagi yang hangat, ceria, namun santai. Sarankan kegiatan ringan hari ini tanpa unsur paksaan."
	} else {
		historyText := strings.Join(yesterdayLogs, "\n- ")
		prompt = fmt.Sprintf("Catatan kemarin:\n- %s\n\nInstruksi: Berikan ucapan selamat pagi yang hangat dan merujuk secara halus pada kegiatan kemarin. Validasi bahwa mengawali hari dengan pelan dan beristirahat adalah hal yang sangat diperbolehkan.", historyText)
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	return extractText(resp), nil
}

// Menghasilkan sapaan malam untuk menanyakan kabar atau progres hari ini
func GenerateEveningGreeting(todayLogs []string) (string, error) {
	ctx := context.Background()
	model := client.GenerativeModel("gemini-2.5-flash")

	systemPrompt := `Berperanlah sebagai pendengar yang hangat. Identitas netral secara gender.
Gunakan bahasa kasual yang lembut, dilarang menggunakan kata 'lu/gue' atau bahasa alay. Dilarang menggunakan format markdown.
Berikan sapaan malam dengan sisipan emotikon ramah (seperti ❤️ atau ✨).`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	var prompt string
	if len(todayLogs) == 0 {
		prompt = "Pengguna belum mencatat kegiatan hari ini. Berikan sapaan malam yang menanyakan kabar dengan lembut dan persilakan untuk bercerita atau menuangkan keluh kesah di sini."
	} else {
		prompt = "Pengguna sudah memiliki catatan hari ini. Berikan sapaan malam yang mengapresiasi usahanya hari ini, sekecil apa pun itu, dan ingatkan untuk segera beristirahat."
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	return extractText(resp), nil
}

// Mengekstrak teks balasan dari struktur respons API Gemini
func extractText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text)
		}
	}
	return "Sistem saat ini sedang tidak dapat memproses data."
}
