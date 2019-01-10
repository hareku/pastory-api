package models

type Diary struct {
	ID        string `json:"id" dynamo:"ID"`
	UserID    string `json:"user_id" dynamo:"UserID"`
	Date      string `json:"date" dynamo:"Date"`
	Body      string `json:"body" dynamo:"Body"`
	UpdatedAt string `json:"updated_at" dynamo:"UpdatedAt"`
	CreatedAt string `json:"created_at" dynamo:"CreatedAt"`
}
