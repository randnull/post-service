package repository

import (
	"context"
	"fmt"
	"errors"
	"time"

	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/randnull/posts-service/internal/config"
	"github.com/randnull/posts-service/internal/graph/model"
)

type PostgresRepository struct {
	db *sqlx.DB
	max_com_len int
}


func NewRepository(cfg *config.Config) *PostgresRepository {
	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.UserDB, cfg.PasswordDB, cfg.HostDB, cfg.PortDB, cfg.NameDB)

	db, err := sqlx.Open("postgres", link)

	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Database is ready")

	return &PostgresRepository{
		db: db,
		max_com_len: cfg.MaxComLen,
	}
}


func (repo* PostgresRepository) Create(title string, content string, allowComments bool) (string, error) {
	query :=`INSERT INTO posts
				(title, content, allow_comments, create_datetime) 
			VALUES 
				($1, $2, $3, $4)
			RETURNING id`
	
	var postID string

	err := repo.db.QueryRow(query, title, content, allowComments, time.Now()).Scan(&postID)

	if err != nil {
		return "", err
	}

	return postID, nil
}


func (repo* PostgresRepository) GetAll() ([]*model.Post, error) {
	query := `SELECT id, title, content, allow_comments, create_datetime FROM posts`

	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post

	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt)
		
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}


func (repo *PostgresRepository) ChangeCommentVisible(postID string, isAllow bool) error {
	query := `UPDATE posts SET allow_comments = $1 WHERE id = $2`
	_, err := repo.db.Exec(query, isAllow, postID)

	if err != nil {
		return err
	}

	return nil
}


func (repo *PostgresRepository) AddComment(postID string, parentID *string, content string) (*model.Comment, error) {
	var allowComments bool

	query := `SELECT allow_comments FROM posts WHERE id = $1`

	err := repo.db.QueryRow(query, postID).Scan(&allowComments)

	if err != nil {
		return nil, errors.New("post not found")
	}

	if !allowComments {
		return nil, errors.New("comments not available")
	}

	if (len([]rune(content))) > repo.max_com_len {
		return nil, errors.New(fmt.Sprintf("comment len more than %v", repo.max_com_len))
	}

	if parentID != nil {
		var exists int

		query := `SELECT COUNT(*) FROM comments WHERE id = $1`

		err := repo.db.QueryRow(query, *parentID).Scan(&exists)

		if err != nil || (exists == 0) {
			return nil, errors.New("parent comment not found")
		}
	}

	query = `INSERT INTO comments
				(post_id, parent_id, content, create_datetime) 
			VALUES 
				($1, $2, $3, $4)
			RETURNING id`
	
	var commentID string

	err = repo.db.QueryRow(query, postID, parentID, content, time.Now()).Scan(&commentID)

	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		ID:        commentID,
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now().String(),
	}

	return comment, nil
}


func (repo *PostgresRepository) GetPost(postID string, startPage *int, pageSize *int) (*model.PostWithComments, error) {
	query := `SELECT id, title, content, allow_comments FROM posts WHERE id = $1`
	
	var post model.Post

	err := repo.db.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments)

	if err != nil {
		return nil, errors.New("post not found")
	}

	query = `SELECT id, post_id, parent_id, content, create_datetime FROM comments WHERE post_id = $1`

	rows, err := repo.db.Query(query, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt) 

		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
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
		Comments:      comments,
	}

	return postWithComments, nil
}
