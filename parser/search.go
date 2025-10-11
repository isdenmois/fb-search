package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var fields string = "id, title, authors, series, serno, file, path, lang, size"
var ruWhere string = "to_tsvector('russian', search) @@ websearch_to_tsquery('russian', $1)"
var enWhere string = "to_tsvector('simple', search) @@ websearch_to_tsquery('simple', $1)"

func searchQuery(q string) string {
	if ContainsCyrillic(q) {
		return "SELECT " + fields + " FROM books WHERE " + ruWhere + " LIMIT 100"
	}

	return "SELECT " + fields + " FROM books WHERE " + enWhere + " LIMIT 100"
}

var byIdQuery = "SELECT " + fields + " FROM books WHERE id = $1 LIMIT 1"

func scanRow(rows pgx.Row, book *Book) error {
	return rows.Scan(&book.Id, &book.Title, &book.Authors, &book.Series, &book.Serno, &book.File, &book.Path, &book.Lang, &book.Size)
}

func searchBooksRu(pool *pgxpool.Pool, q string) ([]Book, error) {
	query := searchQuery(q)

	rows, err := pool.Query(context.Background(), query, q)
	if err != nil {
		fmt.Println("Query error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	var books []Book = []Book{}

	for rows.Next() {
		var book Book
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

func findFileById(pool *pgxpool.Pool, id string) (Book, error) {
	row := pool.QueryRow(context.Background(), byIdQuery, id)

	var book Book
	err := scanRow(row, &book)

	return book, err
}
