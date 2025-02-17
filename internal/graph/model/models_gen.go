// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string  `json:"id"`
	PostID    string  `json:"postId"`
	ParentID  *string `json:"parentId,omitempty"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"createdAt"`
}

type Mutation struct {
}

type Post struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	AllowComments bool   `json:"allowComments"`
	CreatedAt     string `json:"createdAt"`
}

type PostWithComments struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	AllowComments bool       `json:"allowComments"`
	CreatedAt     string     `json:"createdAt"`
	StartPage     *int       `json:"start_page,omitempty"`
	PageSize      *int       `json:"page_size,omitempty"`
	Comments      []*Comment `json:"comments"`
}

type Query struct {
}

type Response struct {
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

type ResponseID struct {
	ID   string `json:"id"`
	Desc string `json:"desc"`
}

type Subscription struct {
}
