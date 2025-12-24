package repositories

import (
	"database/sql"
	"fmt"
)

// LocaleRepository handles locale data queries
type LocaleRepository struct {
	db *sql.DB
}

// NewLocaleRepository creates a new locale repository
func NewLocaleRepository(db *sql.DB) *LocaleRepository {
	return &LocaleRepository{db: db}
}

// InsertLocale inserts a locale string
func (r *LocaleRepository) InsertLocale(key, language, text string) error {
	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO atlasloot_locale (locale_key, language, text)
		VALUES (?, ?, ?)
	`, key, language, text)
	return err
}

// GetLocale retrieves a localized string
func (r *LocaleRepository) GetLocale(key, language string) (string, error) {
	var text string
	err := r.db.QueryRow(`
		SELECT text FROM atlasloot_locale WHERE locale_key = ? AND language = ?
	`, key, language).Scan(&text)
	if err != nil {
		// Fallback to English
		err = r.db.QueryRow(`
			SELECT text FROM atlasloot_locale WHERE locale_key = ? AND language = 'en'
		`, key).Scan(&text)
	}
	return text, err
}

// GetAllLocalesForLanguage gets all locale strings for a language
func (r *LocaleRepository) GetAllLocalesForLanguage(language string) (map[string]string, error) {
	rows, err := r.db.Query(`SELECT locale_key, text FROM atlasloot_locale WHERE language = ?`, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, text string
		if err := rows.Scan(&key, &text); err != nil {
			return nil, err
		}
		result[key] = text
	}
	return result, nil
}

// ClearLocaleData removes all locale data
func (r *LocaleRepository) ClearLocaleData() error {
	_, err := r.db.Exec("DELETE FROM atlasloot_locale")
	return err
}

// GetLocaleStats returns statistics about locale data
func (r *LocaleRepository) GetLocaleStats() (map[string]int, error) {
	stats := make(map[string]int)
	rows, err := r.db.Query(`SELECT language, COUNT(*) FROM atlasloot_locale GROUP BY language`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var lang string
		var count int
		if err := rows.Scan(&lang, &count); err != nil {
			return nil, err
		}
		stats[fmt.Sprintf("Language_%s", lang)] = count
	}
	var total int
	r.db.QueryRow("SELECT COUNT(*) FROM atlasloot_locale").Scan(&total)
	stats["Total"] = total
	return stats, nil
}
