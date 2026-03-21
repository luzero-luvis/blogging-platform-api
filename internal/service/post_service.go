package service

import (
	"errors"

	"blogging-platform-api/internal/model"
	"blogging-platform-api/internal/repository"
)

type PostService struct {
	repo *repository.Postrepository
}

func NewPostService(repo *repository.Postrepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(req *model.CreatePostRequest) (*model.Post, error) {
	if req.Title == "" {
		return nil, errors.New("title required")
	}
	if req.Content == "" {
		return nil, errors.New("Content required")
	}
	if req.Category == "" {
		return nil, errors.New("Category required")
	}
	if req.Tags == nil {
		req.Tags = []string{}
	}

	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
	}

	err := s.repo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
