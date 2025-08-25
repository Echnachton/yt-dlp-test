package db

import (
	"database/sql"
	"sync"

	"github.com/Echnachton/yt-dlp-test/internal/logger"
)

var (
	globalDb *sql.DB
	once sync.Once
)

func Init() {
	var err error
	once.Do(func() {
		logger.Println("Initializing database...")

		globalDb, err = sql.Open("sqlite3", "../../yt-dlp.db")
		if err != nil {
			logger.Println("Error opening database:", err)
			return
		}

	globalDb.Exec("CREATE TABLE IF NOT EXISTS videos (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, internal_video_id TEXT, owner_id TEXT, )")
	})
}

func GetDB() *sql.DB {
	if globalDb == nil {
		Init()
	}
	return globalDb
}