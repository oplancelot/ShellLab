package database

import "fmt"

// InitAtlasLootLocaleSchema creates the locale table for multi-language support
func (s *SQLiteDB) InitAtlasLootLocaleSchema() error {
	schema := `
	-- AtlasLoot Locale table for multi-language support
	CREATE TABLE IF NOT EXISTS atlasloot_locale (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		locale_key TEXT NOT NULL,
		language TEXT NOT NULL,  -- 'en', 'cn', 'de', 'fr', 'es', etc.
		text TEXT NOT NULL,
		UNIQUE(locale_key, language)
	);

	CREATE INDEX IF NOT EXISTS idx_atlasloot_locale_key ON atlasloot_locale(locale_key);
	CREATE INDEX IF NOT EXISTS idx_atlasloot_locale_lang ON atlasloot_locale(language);
	`

	_, err := s.db.Exec(schema)
	return err
}

// LocaleRepository handles locale data queries
type LocaleRepository struct {
	db *SQLiteDB
}

// NewLocaleRepository creates a new locale repository
func NewLocaleRepository(db *SQLiteDB) *LocaleRepository {
	return &LocaleRepository{db: db}
}

// InsertLocale inserts a locale string
func (r *LocaleRepository) InsertLocale(key, language, text string) error {
	_, err := r.db.DB().Exec(`
		INSERT OR REPLACE INTO atlasloot_locale (locale_key, language, text)
		VALUES (?, ?, ?)
	`, key, language, text)
	return err
}

// GetLocale retrieves a localized string
func (r *LocaleRepository) GetLocale(key, language string) (string, error) {
	var text string
	err := r.db.DB().QueryRow(`
		SELECT text FROM atlasloot_locale 
		WHERE locale_key = ? AND language = ?
	`, key, language).Scan(&text)

	if err != nil {
		// Fallback to English if not found
		err = r.db.DB().QueryRow(`
			SELECT text FROM atlasloot_locale 
			WHERE locale_key = ? AND language = 'en'
		`, key).Scan(&text)
	}

	return text, err
}

// GetAllLocalesForLanguage gets all locale strings for a language
func (r *LocaleRepository) GetAllLocalesForLanguage(language string) (map[string]string, error) {
	rows, err := r.db.DB().Query(`
		SELECT locale_key, text FROM atlasloot_locale
		WHERE language = ?
	`, language)
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
	_, err := r.db.DB().Exec("DELETE FROM atlasloot_locale")
	return err
}

// GetLocaleStats returns statistics about locale data
func (r *LocaleRepository) GetLocaleStats() (map[string]int, error) {
	stats := make(map[string]int)

	// Get count per language
	rows, err := r.db.DB().Query(`
		SELECT language, COUNT(*) 
		FROM atlasloot_locale 
		GROUP BY language
	`)
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

	// Get total count
	var total int
	err = r.db.DB().QueryRow("SELECT COUNT(*) FROM atlasloot_locale").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["Total"] = total

	return stats, nil
}
