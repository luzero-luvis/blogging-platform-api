package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

func (r *Postrepository) GetAll() ([]model.Post, error) {
	rows, err := r.db.Query("SELECT  id, title, content, category, tags, created_at,updated_at FROM posts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var posts []model.Post

	for rows.Next() {
		var post model.Post
		var tagJson []byte

		rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &tagJson, &post.CreatedAt, &post.UpdatedAt)
		json.Unmarshal(tagJson, &post.Tags)

		posts = append(posts, post)

	}

	return posts, nil
}

func (r *Postrepository) GetById(id int) (*model.Post, error) {
	var post model.Post
	var tagJson []byte

	err := r.db.QueryRow("SELECT  id, title, content, category, tags, created_at,updated_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &tagJson, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(tagJson, &post.Tags)

	return &post, nil
}

func (r *Postrepository) Put(post *model.Post) error {
	tagJson, _ := json.Marshal(post.Tags)
	query := `
  UPDATE posts 
  SET title = $1,content = $2,category = $3,tags = $4,updated_at = NOW()
	WHERE id = $5
	RETURNING updated_at,created_at
	`
	err := r.db.QueryRow(query, post.Title, post.Content, post.Category, tagJson, post.ID).
		Scan(&post.UpdatedAt, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}

	return err
}
