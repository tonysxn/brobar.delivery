package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/tonysanin/brobar/web-service/internal/models"
)

type SettingRepository struct {
	db *sqlx.DB
}

func NewSettingRepository(db *sqlx.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting
	query := `SELECT key, type, value FROM settings`
	err := r.db.Select(&settings, query)
	return settings, err
}

func (r *SettingRepository) GetByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	query := `SELECT key, type, value FROM settings WHERE key = $1`
	err := r.db.Get(&setting, query, key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &setting, err
}

func (r *SettingRepository) Update(setting *models.Setting) error {
	query := `
		INSERT INTO settings (key, type, value) 
		VALUES ($1, $2, $3)
		ON CONFLICT (key) 
		DO UPDATE SET type = EXCLUDED.type, value = EXCLUDED.value
	`
	_, err := r.db.Exec(query, setting.Key, setting.Type, setting.Value)
	return err
}
