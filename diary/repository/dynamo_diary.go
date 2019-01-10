package repository

import (
	"github.com/guregu/dynamo"
	"github.com/hareku/pastory-api/diary"
	"github.com/hareku/pastory-api/models"
	"github.com/pkg/errors"
)

type dynamoDiaryRepository struct {
	DB *dynamo.DB
}

// NewDynamoDiaryRepository will create an object that represent the diary.Repository interface
func NewDynamoDiaryRepository(DB *dynamo.DB) diary.Repository {
	return &dynamoDiaryRepository{
		DB: DB,
	}
}

func (d *dynamoDiaryRepository) FetchMany(userID string) ([]*models.Diary, error) {
	diaries := []*models.Diary{}

	table := d.DB.Table("Diaries")
	err := table.Get("UserID", userID).Index("UserID-Date-index").Order(true).All(&diaries)

	if err != nil {
		return diaries, errors.Wrap(err, "Failed to fetch user diaries from DynamoDB")
	}

	return diaries, nil
}

func (d *dynamoDiaryRepository) FetchOne(diaryID string, userID string) (*models.Diary, error) {
	diary := new(models.Diary)

	table := d.DB.Table("Diaries")
	err := table.Get("ID", diaryID).Range("UserID", dynamo.Equal, userID).One(diary)

	if err != nil {
		return diary, errors.Wrap(err, "Failed to fetch user diary by id from DynamoDB")
	}

	return diary, nil
}

func (d *dynamoDiaryRepository) Create(diary *models.Diary) error {
	table := d.DB.Table("Diaries")
	err := table.Put(diary).Run()

	if err != nil {
		return errors.Wrap(err, "Failed to put diary to DynamoDB")
	}

	return nil
}

func (d *dynamoDiaryRepository) Update(diary *models.Diary) error {
	table := d.DB.Table("Diaries")
	err := table.Update("ID", diary.ID).Range("UserID", diary.UserID).Set("Body", diary.Body).Set("UpdatedAt", diary.UpdatedAt).Value(diary)

	if err != nil {
		return errors.Wrap(err, "Failed to update diary in DynamoDB")
	}

	return nil
}

func (d *dynamoDiaryRepository) Delete(diaryID string, userID string) error {
	table := d.DB.Table("Diaries")
	err := table.Delete("ID", diaryID).Range("UserID", userID).Run()

	if err != nil {
		return errors.Wrap(err, "Failed to delete diary from DynamoDB")
	}

	return nil
}
