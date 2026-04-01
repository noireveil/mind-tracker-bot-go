package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Membuat koneksi pool ke Supabase menggunakan URL dari variabel lingkungan
func ConnectDB() error {
	dbURL := os.Getenv("SUPABASE_DB_URL")
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	return err
}

func SaveUser(userID int64, username string) error {
	query := `INSERT INTO users (telegram_id, username) VALUES ($1, $2) ON CONFLICT (telegram_id) DO NOTHING`
	_, err := DB.Exec(context.Background(), query, userID, username)
	return err
}

func SaveProgress(userID int64, text string) error {
	query := `INSERT INTO progress_logs (user_id, log_text) VALUES ($1, $2)`
	_, err := DB.Exec(context.Background(), query, userID, text)
	return err
}

func GetTodayLogs(userID int64) ([]string, error) {
	query := `
		SELECT log_text
		FROM progress_logs
		WHERE user_id = $1
		  AND DATE(created_at AT TIME ZONE 'Asia/Jakarta') = DATE(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
		ORDER BY created_at ASC
	`
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var log string
		if err := rows.Scan(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// Mengambil log pada H-1 untuk keperluan konteks pesan sapaan pagi hari
func GetYesterdayLogs(userID int64) ([]string, error) {
	query := `
		SELECT log_text
		FROM progress_logs
		WHERE user_id = $1
		  AND DATE(created_at AT TIME ZONE 'Asia/Jakarta') = DATE(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta') - INTERVAL '1 day'
		ORDER BY created_at ASC
	`
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var log string
		if err := rows.Scan(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func GetRecentLogs(userID int64, limit int) ([]string, error) {
	query := `SELECT log_text FROM progress_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := DB.Query(context.Background(), query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var log string
		if err := rows.Scan(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func GetAllUsers() ([]int64, error) {
	query := `SELECT telegram_id FROM users`
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		users = append(users, id)
	}
	return users, nil
}

func DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE telegram_id = $1`
	_, err := DB.Exec(context.Background(), query, userID)
	return err
}

// Menghitung ukuran database untuk fitur peringatan batas penyimpanan gratis
func GetDatabaseSizeMB() (float64, error) {
	var sizeMB float64
	query := "SELECT pg_database_size(current_database()) / 1024.0 / 1024.0"
	err := DB.QueryRow(context.Background(), query).Scan(&sizeMB)
	return sizeMB, err
}

// Menghapus log lama untuk efisiensi ruang penyimpanan
func PurgeOldLogs(userID int64, months int) error {
	query := `
		DELETE FROM progress_logs
		WHERE user_id = $1
		AND created_at < NOW() - INTERVAL '1 month' * $2
	`
	_, err := DB.Exec(context.Background(), query, months)
	return err
}

func SaveSummary(userID int64, text string, sType string) error {
	query := `INSERT INTO summaries (user_id, summary_text, summary_type) VALUES ($1, $2, $3)`
	_, err := DB.Exec(context.Background(), query, userID, text, sType)
	return err
}

func GetPersistentSummaries(userID int64) ([]string, error) {
	query := `SELECT summary_text FROM summaries WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []string
	for rows.Next() {
		var log string
		if err := rows.Scan(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
