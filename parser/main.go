package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func rebuildDb(pool *pgxpool.Pool, progress *ParserProgress) {
	pool.Exec(context.Background(), "TRUNCATE TABLE books RESTART IDENTITY")
	pool.Exec(context.Background(), "VACUUM")

	parseInpx("./files/fb2.Flibusta.Net.7z.inpx", pool, progress)
}

func main() {
	godotenv.Load()
	ctx := context.Background()
	databaseUrl := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Ping the database to verify connection
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to the database!")

	r := gin.Default()

	progress := ParserProgress{}

	// Define a simple GET endpoint
	r.GET("/api/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/search", func(c *gin.Context) {
		q := c.Query("q")

		if len(q) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Q should be at least 2 symbols"})
			return
		}

		books, err := searchBooksRu(pool, c.Query("q"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, books)

	})
	//   const file = await findFileById(id)

	//   if (file?.file && file.path) {
	//     const fileId = parse(file.path).name
	//     const stream = await getCover(file.file, fileId)

	//     if (stream) {
	//       set.headers['content-type'] = 'image/jpeg'

	//       if (stream.length) {
	//         set.headers['content-length'] = String(stream.length)
	//       }

	//       return new Response(stream!).bytes()
	//     }
	//   }

	//   set.status = 404
	//   return { ok: false }

	r.GET("/cover/:id", func(ctx *gin.Context) {
		start := time.Now()
		file, err := findFileById(pool, ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if file.Path == nil || len(*file.Path) == 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no path"})
			return
		}

		fmt.Println("find took ", time.Since(start))
		start = time.Now()

		img, err := getCover(*file.File, *file.Path)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("xl to img took ", time.Since(start))
		start = time.Now()

		jpg, err := imageToJpg(img)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageDataBytes := jpg.Bytes()

		fmt.Println("img to jpg took ", time.Since(start))

		ctx.Header("Content-Size", string(len(imageDataBytes)))
		ctx.Data(http.StatusOK, "image/jpeg", imageDataBytes)
	})

	r.GET("/api/parser", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"files": progress.files,
			"books": progress.books,
			"time":  progress.time,
		})
	})

	r.POST("/api/parser", func(c *gin.Context) {
		progress = ParserProgress{
			files: 0,
			books: 0,
			time:  0,
		}

		rebuildDb(pool, &progress)

		c.JSON(http.StatusOK, gin.H{
			"files": progress.files,
			"books": progress.books,
			"time":  progress.time,
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
