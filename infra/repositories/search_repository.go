package repositories

import (
	"context"
	"fb-search/domain"

	"github.com/jackc/pgx/v5"
	"github.com/typesense/typesense-go/v4/typesense"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

type SearchRepository struct {
	ts *typesense.Client
}

func (self SearchRepository) Clean() {
	self.ts.Collection("books").Documents().Delete(context.Background(), &api.DeleteDocumentsParams{})
}

func (self SearchRepository) InsertBatch(rows pgx.CopyFromSource) (int, error) {
	var docs []interface{}

	for {
		if !rows.Next() {
			break
		}
		row, err := rows.Values()
		if err != nil {
			continue
		}
		// columns := []string{"title", "search", "authors", "series", "serno", "file", "path", "lang", "size"}
		docs = append(docs, Book{
			ID:      row[6].(string),
			Idx:     row[1].(string),
			Title:   row[0].(string),
			Authors: row[2].(string),
			Series:  row[3].(string),
			File:    row[5].(string),
		})
	}
	params := &api.ImportDocumentsParams{
		Action:    (*api.IndexAction)(pointer.String("create")),
		BatchSize: pointer.Int(100),
	}

	_, err := self.ts.Collection("books").Documents().Import(context.Background(), docs, params)

	return len(docs), err
}

func (self SearchRepository) Search(q string) ([]domain.Book, error) {
	searchParameters := &api.SearchCollectionParams{
		Q:       pointer.String(q),
		QueryBy: pointer.String("idx"),
		Limit:   pointer.Int(100),
	}

	res, err := self.ts.Collection("books").Documents().Search(context.Background(), searchParameters)

	if err != nil {
		return nil, err
	}

	var books []domain.Book

	for _, hit := range *res.Hits {
		doc := *hit.Document
		authours := doc["authors"].(string)
		series := doc["series"].(string)

		books = append(books, domain.Book{
			Id:      doc["id"].(string),
			Title:   doc["title"].(string),
			Authors: &authours,
			Series:  &series,
		})
	}

	return books, nil
}

func NewSearchRepository(ts *typesense.Client) *SearchRepository {
	schema := &api.CollectionSchema{
		Name: "books",
		Fields: []api.Field{
			{
				Name: "idx",
				Type: "string",
			},
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "authors",
				Type: "string",
			},
			{
				Name: "series",
				Type: "string",
			},
			{
				Name: "file",
				Type: "string",
			},
		},
	}

	ts.Collections().Create(context.Background(), schema)

	return &SearchRepository{ts: ts}
}

type Book struct {
	ID      string `json:"id"`
	Idx     string `json:"idx"`
	Title   string `json:"title"`
	Authors string `json:"authors"`
	Series  string `json:"series"`
	File    string `json:"file"`
}
