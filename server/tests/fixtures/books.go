package fixtures

import (
	"context"

	"fb-search/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestBook struct {
	Id      string
	Title   string
	Search  string
	Authors *string
	Series  *string
	Serno   *string
	Lang    *string
	Size    *uint
}

func Ptr(s string) *string {
	return &s
}

func PtrUint(u uint) *uint {
	return &u
}

var (
	RussianBooks = []TestBook{
		{
			Id:      "fb2-1.zip/1.fb2",
			Title:   "Война и мир",
			Search:  "война и мир толстой лев",
			Authors: Ptr("Толстой Лев"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("ru"),
			Size:    PtrUint(123456),
		},
		{
			Id:      "fb2-1.zip/2.fb2",
			Title:   "Преступление и наказание",
			Search:  "преступление и наказание достоевский федор",
			Authors: Ptr("Достоевский Федор"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("ru"),
			Size:    PtrUint(98765),
		},
		{
			Id:      "fb2-1.zip/3.fb2",
			Title:   "Мастер и Маргарита",
			Search:  "мастер и маргарита булгаков михаил",
			Authors: Ptr("Булгаков Михаил"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("ru"),
			Size:    PtrUint(87654),
		},
	}

	EnglishBooks = []TestBook{
		{
			Id:      "fb2-2.zip/4.fb2",
			Title:   "The Great Gatsby",
			Search:  "the great gatsby fitzgerald f scott",
			Authors: Ptr("F. Scott Fitzgerald"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("en"),
			Size:    PtrUint(54321),
		},
		{
			Id:      "fb2-2.zip/5.fb2",
			Title:   "1984",
			Search:  "1984 orwell george",
			Authors: Ptr("George Orwell"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("en"),
			Size:    PtrUint(45678),
		},
		{
			Id:      "fb2-2.zip/6.fb2",
			Title:   "To Kill a Mockingbird",
			Search:  "to kill a mockingbird lee harper",
			Authors: Ptr("Harper Lee"),
			Series:  nil,
			Serno:   nil,
			Lang:    Ptr("en"),
			Size:    PtrUint(67890),
		},
	}

	SeriesBooks = []TestBook{
		{
			Id:      "fb2-3.zip/7.fb2",
			Title:   "Harry Potter and the Philosopher's Stone",
			Search:  "harry potter philosopher stone rowling j k",
			Authors: Ptr("J.K. Rowling"),
			Series:  Ptr("Harry Potter"),
			Serno:   Ptr("1"),
			Lang:    Ptr("en"),
			Size:    PtrUint(112233),
		},
		{
			Id:      "fb2-3.zip/8.fb2",
			Title:   "Harry Potter and the Chamber of Secrets",
			Search:  "harry potter chamber secrets rowling j k",
			Authors: Ptr("J.K. Rowling"),
			Series:  Ptr("Harry Potter"),
			Serno:   Ptr("2"),
			Lang:    Ptr("en"),
			Size:    PtrUint(223344),
		},
	}

	AllBooks = append(append(RussianBooks, EnglishBooks...), SeriesBooks...)
)

func InsertBooks(ctx context.Context, pool *pgxpool.Pool, books []TestBook) error {
	for _, b := range books {
		_, err := pool.Exec(ctx,
			"INSERT INTO books (id, title, search, authors, series, serno, lang, size) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			b.Id, b.Title, b.Search, b.Authors, b.Series, b.Serno, b.Lang, b.Size,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func ToDomainBooks(books []TestBook) []domain.Book {
	result := make([]domain.Book, len(books))
	for i, b := range books {
		result[i] = domain.Book{
			Id:      b.Id,
			Title:   b.Title,
			Authors: b.Authors,
			Series:  b.Series,
			Serno:   b.Serno,
			Lang:    b.Lang,
			Size:    b.Size,
		}
	}
	return result
}
