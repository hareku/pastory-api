package usecase

import (
	"time"

	"github.com/hareku/pastory-api/diary"
	"github.com/hareku/pastory-api/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type diaryUsecase struct {
	diaryRepo diary.Repository
}

// NewDiaryUsecase will create new an diaryUsecase object representation of diary.Usecase interface
func NewDiaryUsecase(d diary.Repository) diary.Usecase {
	return &diaryUsecase{
		diaryRepo: d,
	}
}

func (d *diaryUsecase) FetchMany(userID string) ([]*models.Diary, error) {
	return d.diaryRepo.FetchMany(userID)
}

func (d *diaryUsecase) Create(input *diary.CreateInput) (*models.Diary, error) {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	diary := &models.Diary{
		UserID:    input.UserID,
		Body:      input.Body,
		Date:      input.Date,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	// Set ID
	u2, err := uuid.NewV4()
	if err != nil {
		return diary, errors.Wrap(err, "Failed to make new UUID v4")
	}
	diary.ID = u2.String()

	err = d.diaryRepo.Create(diary)

	if err != nil {
		return diary, errors.Wrap(err, "Failed to create diary to repository")
	}

	return diary, nil
}

func (d *diaryUsecase) Update(input *diary.UpdateInput) (*models.Diary, error) {
	diary, err := d.diaryRepo.FetchOne(input.DiaryID, input.UserID)
	if err != nil {
		return diary, errors.Wrap(err, "Failed to fetch diary")
	}

	diary.Body = input.Body
	diary.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = d.diaryRepo.Update(diary)
	if err != nil {
		return diary, errors.Wrap(err, "Failed to update diary")
	}

	return diary, nil
}

func (d *diaryUsecase) Delete(input *diary.DeleteInput) error {
	err := d.diaryRepo.Delete(input.ID, input.UserID)

	if err != nil {
		return errors.Wrap(err, "Failed to delete diary from repository")
	}

	return nil
}
