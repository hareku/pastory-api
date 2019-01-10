package diary

import (
	"github.com/hareku/pastory-api/models"
)

// Repository represent the diary's repository contract
type Repository interface {
	FetchMany(userID string) ([]*models.Diary, error)
	FetchOne(diaryID string, userID string) (*models.Diary, error)
	Create(diary *models.Diary) error
	Update(diary *models.Diary) error
	Delete(diaryID string, userID string) error
}
