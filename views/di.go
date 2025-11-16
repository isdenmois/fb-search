package views

import (
	"fb-search/app"
	"fb-search/infra/db"
	"fb-search/infra/repositories"
	"fb-search/views/controllers"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di/v2"
	"github.com/typesense/typesense-go/v4/typesense"
)

var DbDef = &di.Def{
	Name:  "db",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		pool, err := db.Connect()

		return pool, err
	},
	Close: func(obj interface{}) error {
		obj.(*pgxpool.Pool).Close()
		return nil
	},
}

var TsDef = &di.Def{
	Name:  "typesense",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		client := typesense.NewClient(
			typesense.WithServer(os.Getenv("TS_URL")),
			typesense.WithAPIKey(os.Getenv("TS_KEY")),
		)

		return client, nil
	},
}

var BooksRepositoryDef = &di.Def{
	Name:  "booksRepository",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		pool := ctn.Get(DbDef).(*pgxpool.Pool)
		booksRepo := repositories.NewBooksRepository(pool)

		return booksRepo, nil
	},
}

var SearchRepositoryDef = &di.Def{
	Name:  "searchRepository",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		ts := ctn.Get(TsDef).(*typesense.Client)

		return repositories.NewSearchRepository(ts), nil
	},
}

var InpParserDef = &di.Def{
	Name:  "inpParser",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		booksRepo := ctn.Get(BooksRepositoryDef).(*repositories.BooksRepository)
		searchRepo := ctn.Get(SearchRepositoryDef).(*repositories.SearchRepository)

		return app.NewInpParserCase(booksRepo, searchRepo), nil
	},
}

var ControllersDef = &di.Def{
	Name:  "controllers",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		inpParser := ctn.Get(InpParserDef).(*app.InpParserCase)
		booksRepo := ctn.Get(BooksRepositoryDef).(*repositories.BooksRepository)
		searchRepo := ctn.Get(SearchRepositoryDef).(*repositories.SearchRepository)

		ping := &controllers.PingController{}
		parser := controllers.NewParserController(inpParser)
		books := controllers.NewBookController(booksRepo, searchRepo)

		return &[]controllers.Controller{ping, parser, books}, nil
	},
}

var HttpServerDef = &di.Def{
	Name:  "http-server",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		pool := ctn.Get(DbDef).(*pgxpool.Pool)
		ctrls := ctn.Get(ControllersDef).(*[]controllers.Controller)
		server := NewHttpServer(pool, ctrls)

		return server, nil
	},
}

func CreateDi() (di.Container, error) {
	builder, _ := di.NewEnhancedBuilder()
	builder.Add(DbDef)
	builder.Add(TsDef)
	builder.Add(BooksRepositoryDef)
	builder.Add(SearchRepositoryDef)
	builder.Add(InpParserDef)
	builder.Add(ControllersDef)
	builder.Add(HttpServerDef)

	return builder.Build()
}
