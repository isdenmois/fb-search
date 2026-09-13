package http

import (
	"log"
	"os"

	"fb-search/application/usecases"
	"fb-search/delivery/http/controllers"
	"fb-search/infrastructure/db"
	"fb-search/infrastructure/repositories"
	"fb-search/infrastructure/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di/v2"
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

var BooksRepositoryDef = &di.Def{
	Name:  "booksRepository",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		pool := ctn.Get(DbDef).(*pgxpool.Pool)
		booksRepo := repositories.NewPostgresBooksRepository(pool)

		return booksRepo, nil
	},
}

var BookFileStorageDef = &di.Def{
	Name:  "bookFileStorage",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		return storage.NewZipBookFileStorage("files"), nil
	},
}

var SearchBooksCaseDef = &di.Def{
	Name:  "searchBooksCase",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		booksRepo := ctn.Get(BooksRepositoryDef).(*repositories.PostgresBooksRepository)

		return usecases.NewSearchBooksCase(booksRepo), nil
	},
}

var DownloadBookCaseDef = &di.Def{
	Name:  "downloadBookCase",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		booksRepo := ctn.Get(BooksRepositoryDef).(*repositories.PostgresBooksRepository)
		fileStorage := ctn.Get(BookFileStorageDef).(*storage.ZipBookFileStorage)

		return usecases.NewDownloadBookCase(booksRepo, fileStorage), nil
	},
}

var InpParserDef = &di.Def{
	Name:  "inpParser",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		booksRepo := ctn.Get(BooksRepositoryDef).(*repositories.PostgresBooksRepository)

		return usecases.NewInpParserCase(booksRepo, "files/flibusta_fb2_local.inpx"), nil
	},
}

var ControllersDef = &di.Def{
	Name:  "controllers",
	Scope: di.App,
	Build: func(ctn di.Container) (interface{}, error) {
		inpParser := ctn.Get(InpParserDef).(*usecases.InpParserCase)
		searchBooksCase := ctn.Get(SearchBooksCaseDef).(*usecases.SearchBooksCase)
		downloadBookCase := ctn.Get(DownloadBookCaseDef).(*usecases.DownloadBookCase)

		ping := &controllers.PingController{}
		apiKey := os.Getenv("ADMIN_API_KEY")
		if apiKey == "" {
			log.Println("ADMIN_API_KEY is not set; POST /api/parse/rebuild is disabled (requests will get 401)")
		}
		parser := controllers.NewParserController(inpParser, apiKey)
		books := controllers.NewBookController(searchBooksCase, downloadBookCase)

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
	builder.Add(BooksRepositoryDef)
	builder.Add(BookFileStorageDef)
	builder.Add(SearchBooksCaseDef)
	builder.Add(DownloadBookCaseDef)
	builder.Add(InpParserDef)
	builder.Add(ControllersDef)
	builder.Add(HttpServerDef)

	return builder.Build()
}
