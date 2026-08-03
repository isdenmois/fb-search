package repositories

import (
	"context"
	"fmt"

	"fb-search/application/ports"
	"fb-search/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBooksRepository struct {
	pool *pgxpool.Pool
}

var (
	fields  string = "id, title, authors, series, serno, lang, size"
	ruWhere string = "to_tsvector('russian', search) @@ websearch_to_tsquery('russian', $1)"
	enWhere string = "to_tsvector('simple', search) @@ websearch_to_tsquery('simple', $1)"
	ruRank  string = "ts_rank(to_tsvector('russian', search), websearch_to_tsquery('russian', $1))"
	enRank  string = "ts_rank(to_tsvector('simple', search), websearch_to_tsquery('simple', $1))"
)

func searchQuery(cfg ports.SearchConfig) string {
	if cfg.Language == "russian" {
		return "SELECT " + fields + ", " + ruRank + " as rank FROM books WHERE " + ruWhere + " ORDER BY rank DESC LIMIT 100"
	}

	return "SELECT " + fields + ", " + enRank + " as rank FROM books WHERE " + enWhere + " ORDER BY rank DESC LIMIT 100"
}

var byIdQuery = "SELECT " + fields + ", 0.0 as rank FROM books WHERE id = $1 LIMIT 1"

func scanRow(rows pgx.Row, book *domain.Book) error {
	return rows.Scan(&book.Id, &book.Title, &book.Authors, &book.Series, &book.Serno, &book.Lang, &book.Size, &book.Rank)
}

func (self *PostgresBooksRepository) SearchBooks(ctx context.Context, q string, cfg ports.SearchConfig) ([]domain.Book, error) {
	query := searchQuery(cfg)

	rows, err := self.pool.Query(ctx, query, q)
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

func (self *PostgresBooksRepository) FindById(ctx context.Context, id string) (domain.Book, error) {
	row := self.pool.QueryRow(ctx, byIdQuery, id)

	var book domain.Book
	err := scanRow(row, &book)

	return book, err
}

func (self *PostgresBooksRepository) RebuildDb(ctx context.Context) error {
	if _, err := self.pool.Exec(ctx, "TRUNCATE TABLE books RESTART IDENTITY"); err != nil {
		return err
	}
	_, err := self.pool.Exec(ctx, "VACUUM")
	return err
}

func (self *PostgresBooksRepository) InsertBatch(ctx context.Context, rows ports.RowBatchSource) (uint64, error) {
	tableName := pgx.Identifier{"books"}
	columns := []string{"id", "title", "search", "authors", "series", "serno", "lang", "size"}

	res, err := self.pool.CopyFrom(ctx, tableName, columns, rows)

	return uint64(res), err
}

func NewPostgresBooksRepository(pool *pgxpool.Pool) *PostgresBooksRepository {
	return &PostgresBooksRepository{pool: pool}
}
