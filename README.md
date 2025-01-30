# golang-sample-with-pgx

  This repository is to demo how to use pgx

## what is pgx?

pgx is a golang library for handle database connection handle function to PostgreSQL 
it provide extensive function more than the origin standard database pg library such as connection pool, data serialize

## why is pgx? 

It could handle many basic stuff for handle database query and parameter sanitzation

## setup dependency

```shell
go get github.com/jackc/pgx/v5 
```

## prepare schema

```sql
CREATE TABLE IF NOT EXISTS authors (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS books (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  author_id INTEGER REFERENCES authors(id),
  published_year INTEGER,
  genre TEXT
);

CREATE TABLE IF NOT EXISTS members (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  join_date DATE NOT NULL DEFAULT CURRENT_DATE
);
```

## models 

```golang
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
```

## usage

### Connection Setup

```golang
  // Database connection with pool
	pool, err := pgxpool.New(context.Background(), config.AppConfig.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()
```

### Transaction(Handle multiple data modification in one transaction) 
```golang
	// Transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
		INSERT INTO authors (name, email)
		VALUES ($1, $2);
		`, "George R.R. Martin", "george.martin@sample.io")
	if err != nil {
		log.Fatalf("Error inserting author: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO books (title, author_id, published_year, genre)
		VALUES ($1, $2, $3, $4);
		`, "Harry Potter", 1, 1997, "Fantasy")
	if err != nil {
		log.Fatalf("Error inserting book: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO members (name, email)
		VALUES ($1, $2);
		`, "John Doe", "john.doe@sample.com")
	if err != nil {
		log.Fatalf("Error inserting members: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}
```

### Query Authors (Query multiple record)
```golang 
	// Query All authors from database
	rows, err := pool.Query(context.Background(),
		`SELECT id, name, email FROM authors;`,
	)
	if err != nil {
		log.Fatalf("Error querying authors: %v", err)
	}
	defer rows.Close()

	authors := make([]model.Author, 0, 100)
	for rows.Next() {
		var author model.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Email); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		authors = append(authors, author)
	}
	fmt.Println("Authors:", authors)
```

### Query A Book (Single Record)
```golang
	// Query single book from Database
	var book model.Book
	err = pool.QueryRow(context.Background(), `
	SELECT id, title, author_id, published_year, genre
	FROM books
	WHERE title=$1
	`, "Harry Potter").Scan(&book.ID, &book.Title, &book.AuthorID, &book.PublishedYear, &book.Genre)
	if err != nil {
		log.Fatalf("Error query book: %v", err)
	}
	fmt.Println("Book Detail:", book)
```

### Update Author
```golang
	// Update an author's name
	authorID := 1
	newName := "J.K. Rowling Updated"
	_, err = pool.Exec(context.Background(),
		`UPDATE authors SET name = $1
	 WHERE id = $2
	`, newName, authorID,
	)
	if err != nil {
		log.Fatalf("Error updating author: %v", err)
	}
	fmt.Println("Author updated successfully")
```

### DELETE Book 
```golang
	// DELETE a book by ID
	bookID := 1
	_, err = pool.Exec(context.Background(),
		`DELETE FROM books
	 WHERE id = $1;
	`, bookID,
	)
	if err != nil {
		log.Fatalf("Error deleting book: %v", err)
	}
	fmt.Println("Book delete successfully")
```