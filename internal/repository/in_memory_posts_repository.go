package repository


import (
	"errors"
	"sync"
	"time"
	"log"

	"github.com/google/uuid"

	"github.com/randnull/posts-service/internal/graph/model"
	"github.com/randnull/posts-service/internal/config"
)


type InMemoryRepository struct {
	posts        map[string]*model.Post
	comments     map[string]*model.Comment
	post_mu      sync.Mutex
	comment_mu 	 sync.Mutex

	max_com_len  int
}


func NewInMemoryRepository(cfg *config.Config) *InMemoryRepository {
	log.Print("In-memory storage is ready")

	return &InMemoryRepository{
		posts:       make(map[string]*model.Post),
		comments:    make(map[string]*model.Comment),
		max_com_len: cfg.MaxComLen,
	}
}


func (repo *InMemoryRepository) Create(title string, content string, allowComments bool) (string, error) {
	repo.post_mu.Lock()
    defer repo.post_mu.Unlock()

    postID := uuid.New().String()

    post := &model.Post{
        ID:            postID,
        Title:         title,
        Content:       content,
        AllowComments: allowComments,
		CreatedAt:     time.Now().String(),
    }

    repo.posts[postID] = post

    return postID, nil
}


func (repo *InMemoryRepository) GetAll() ([]*model.Post, error) {
    repo.post_mu.Lock()
    defer repo.post_mu.Unlock()

    posts := make([]*model.Post, 0, len(repo.posts))

    for _, post := range repo.posts {
        posts = append(posts, post)
    }

    return posts, nil
}


func (repo *InMemoryRepository) ChangeCommentVisible(postID string, is_allow bool) error {
    repo.post_mu.Lock()
    defer repo.post_mu.Unlock()

    post, exists := repo.posts[postID]

    if !exists {
        return errors.New("post not found")
    }

	post.AllowComments = is_allow

    return nil
}


func (repo* InMemoryRepository) AddComment(postID string, parentID *string, content string) (*model.Comment, error) {
	repo.post_mu.Lock()
	post, exists := repo.posts[postID]
	repo.post_mu.Unlock()

	if !exists {
		return nil, errors.New("post not found")
	}

	if !post.AllowComments {
		return nil, errors.New("comment not avalible")
	}

	if parentID != nil {
		_, exists := repo.comments[*parentID]

		if !exists {
			return nil, errors.New("parent comment not found")
		}
	}

	repo.comment_mu.Lock()
	defer repo.comment_mu.Unlock()

	commentID := uuid.New().String()

	comment := &model.Comment{
		ID:        commentID,
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now().String(),
	}

	repo.comments[commentID] = comment

	return comment, nil
}


func (repo* InMemoryRepository) GetPost(postID string, startPage *int, pageSize *int) (*model.PostWithComments, error) {
	repo.post_mu.Lock()
	post, exists := repo.posts[postID]
	repo.post_mu.Unlock()

	if !exists {
		return nil, errors.New("post not found")
	}

	repo.comment_mu.Lock()
	defer repo.comment_mu.Unlock()

	comments := make([]*model.Comment, 0)

	for _, comment := range repo.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	if startPage != nil && pageSize != nil {
		if *startPage > 0 && *pageSize > 0 {
			size_comment := len(comments)
			start := (*startPage - 1) * *pageSize
			end := *startPage * *pageSize

			if start > size_comment {
				return nil, nil
			}

			if end > size_comment {
				end = size_comment
			}

			comments = comments[start:end]
		} else {
			return nil, errors.New("start_page or page_size <= 0")
		}
	}

	postWithComments := &model.PostWithComments{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt,
		Comments:      comments,
	}

	return postWithComments, nil
}
