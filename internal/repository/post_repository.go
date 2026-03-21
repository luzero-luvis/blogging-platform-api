package repository

import (
	"database/sql"
	"encoding/json"

	"blogging-platform-api/internal/model"
)

type Postrepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *Postrepository {
	return &Postrepository{db: db}
}

func (r *Postrepository) Create(post *model.Post) error {
	tagJson, _ := json.Marshal(post.Tags)

	querry := `
	 INSERT INTO posts(title,content,category,tags)
	 VALUES ($1, $2, $3, $4)
	 RETURNING id, created_at ,updated_at
	`
	return r.db.QueryRow(querry, post.Title, post.Content, post.Category, tagJson).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}
