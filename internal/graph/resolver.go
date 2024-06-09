package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"sync"

	"github.com/randnull/posts-service/internal/graph/model"
)


type PostsRepository interface {
	GetAll() ([]*model.Post, error)
	Create(title string, content string, allowComments bool) (string, error)
	ChangeCommentVisible(id string, is_allow bool) error
	AddComment(postID string, parentID *string, content string) (*model.Comment, error)
	GetPost(postID string, startPage *int, pageSize *int) (*model.PostWithComments, error)
}

type Resolver struct {
	Repo             PostsRepository

	commentObservers map[string][]chan *model.Comment
	mu               sync.Mutex
}

func NewResolver(Repo PostsRepository) *Resolver{
	return &Resolver{
		Repo: Repo,
		commentObservers: make(map[string][]chan *model.Comment),
	}
}