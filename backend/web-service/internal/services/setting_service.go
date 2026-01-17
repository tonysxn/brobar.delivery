package services

import (
	"github.com/tonysanin/brobar/web-service/internal/models"
	"github.com/tonysanin/brobar/web-service/internal/repositories"
)

type SettingService struct {
	repo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetAllSettings() ([]models.Setting, error) {
	return s.repo.GetAll()
}

func (s *SettingService) GetSetting(key string) (*models.Setting, error) {
	return s.repo.GetByKey(key)
}

func (s *SettingService) UpdateSetting(key string, value string, typeStr string) error {
	setting := &models.Setting{
		Key:   key,
		Value: value,
		Type:  typeStr,
	}
	return s.repo.Update(setting)
}
