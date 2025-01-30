package model

type Author struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Book struct {
	ID            int32  `json:"id"`
	Title         string `json:"title"`
	AuthorID      int32  `json:"author_id"`
	PublishedYear int32  `json:"published_year"`
	Genre         string `json:"genre"`
}

type Member struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinDate string `json:"join_date"`
}
