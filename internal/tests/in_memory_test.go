package tests

import (
	"testing"

	"github.com/randnull/posts-service/internal/config"
	"github.com/randnull/posts-service/internal/repository"
	"github.com/stretchr/testify/assert"
)


func TestCreatePost(t *testing.T) {
	cfg, _ := config.NewConfig()

	repo := repository.NewInMemoryRepository(cfg)

	id, err := repo.Create("some_title", "content", true)

	assert.NoError(t, err)

	post, err := repo.GetPost(id, nil, nil)

	assert.NoError(t, err)

	assert.Equal(t, "some_title", post.Title)
	assert.Equal(t, "content", post.Content)
	assert.Equal(t, true, post.AllowComments)
}


func TestGetAllPost(t *testing.T) {
	cfg, _ := config.NewConfig()

	repo := repository.NewInMemoryRepository(cfg)

	repo.Create("some_title", "content", true)
	repo.Create("some_title", "content", true)

	posts, err := repo.GetAll()

	assert.NoError(t, err)

	assert.Equal(t, len(posts), 2)
}


func TestAddComment(t *testing.T) {
	cfg, _ := config.NewConfig()

	repo := repository.NewInMemoryRepository(cfg)

	id, _ := repo.Create("some_title", "content", true)

	_, err := repo.AddComment(id, nil, "comment")

	assert.NoError(t, err)

	post, _ := repo.GetPost(id, nil, nil)

	assert.Equal(t, len(post.Comments), 1)
	assert.Equal(t, post.Comments[0].Content, "comment")
}
