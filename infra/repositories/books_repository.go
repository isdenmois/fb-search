package repositories

import (
	"context"
	"fb-search/domain"
	"fb-search/shared/utils"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BooksRepository struct {
	pool *pgxpool.Pool
}

var fields string = "id, title, authors, series, serno, file, path, lang, size"
var ruWhere string = "to_tsvector('russian', search) @@ websearch_to_tsquery('russian', $1)"
var enWhere string = "to_tsvector('simple', search) @@ websearch_to_tsquery('simple', $1)"

func searchQuery(q string) string {
	if utils.ContainsCyrillic(q) {
		return "SELECT " + fields + " FROM books WHERE " + ruWhere + " LIMIT 100"
	}

	return "SELECT " + fields + " FROM books WHERE " + enWhere + " LIMIT 100"
}

var byIdQuery = "SELECT " + fields + " FROM books WHERE id = $1 LIMIT 1"

func scanRow(rows pgx.Row, book *domain.Book) error {
	return rows.Scan(&book.Id, &book.Title, &book.Authors, &book.Series, &book.Serno, &book.File, &book.Path, &book.Lang, &book.Size)
}

func (self BooksRepository) SearchBooks(q string) ([]domain.Book, error) {
	query := searchQuery(q)

	rows, err := self.pool.Query(context.Background(), query, q)

	if err != nil {
		fmt.Println("Query error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	var books []domain.Book = []domain.Book{}

	for rows.Next() {
		var book domain.Book
		if err := scanRow(rows, &book); err != nil {
			fmt.Println("Scan error: ", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Rows error: ", err.Error())
		return nil, err
	}

	return books, nil
}

func (self BooksRepository) FindFileById(id string) (domain.Book, error) {
	row := self.pool.QueryRow(context.Background(), byIdQuery, id)

	var book domain.Book
	err := scanRow(row, &book)

	return book, err
}

func (self BooksRepository) RebuildDb() {
	self.pool.Exec(context.Background(), "TRUNCATE TABLE books RESTART IDENTITY")
	self.pool.Exec(context.Background(), "VACUUM")
}

func (self BooksRepository) InsertBatch(rows pgx.CopyFromSource) (uint64, error) {
	tableName := pgx.Identifier{"books"}
	columns := []string{"title", "search", "authors", "series", "serno", "file", "path", "lang", "size"}

	res, err := self.pool.CopyFrom(context.Background(), tableName, columns, rows)

	return uint64(res), err
}

func NewBooksRepository(pool *pgxpool.Pool) *BooksRepository {
	return &BooksRepository{pool: pool}
}
