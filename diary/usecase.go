package diary

import (
	"time"

	"github.com/hareku/pastory-api/models"
	validator "gopkg.in/go-playground/validator.v9"
)

type CreateInput struct {
	UserID string
	Date   string `json:"date" validate:"required,date"`
	Body   string `json:"body" validate:"required,max=255"`
}

type DeleteInput struct {
	ID     string
	UserID string
}

type UpdateInput struct {
	DiaryID string
	UserID  string
	Body    string `json:"body" validate:"required,max=255"`
}

// Usecase represent the diary's usecases
type Usecase interface {
	FetchMany(userID string) ([]*models.Diary, error)
	Create(input *CreateInput) (*models.Diary, error)
	Update(input *UpdateInput) (*models.Diary, error)
	Delete(input *DeleteInput) error
}

func DateValidation(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
