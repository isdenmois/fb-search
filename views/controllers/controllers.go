package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Bind(*gin.Engine) error
}

// func CreateControllers(r *gin.Engine, pool *pgxpool.Pool) {
// 	pingController(r)
// 	parserControler(r, pool)
// 	bookControler(r, pool)
// }

// progress := parser.ParserProgress{}

// // Define a simple GET endpoint

// r.GET("/api/search", func(c *gin.Context) {
// 	q := c.Query("q")

// 	if len(q) < 2 {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Q should be at least 2 symbols"})
// 		return
// 	}

// 	books, err := searchBooksRu(pool, c.Query("q"))

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, books)

// })
// //   const file = await findFileById(id)

// //   if (file?.file && file.path) {
// //     const fileId = parse(file.path).name
// //     const stream = await getCover(file.file, fileId)

// //     if (stream) {
// //       set.headers['content-type'] = 'image/jpeg'

// //       if (stream.length) {
// //         set.headers['content-length'] = String(stream.length)
// //       }

// //       return new Response(stream!).bytes()
// //     }
// //   }

// //   set.status = 404
// //   return { ok: false }

// r.GET("/cover/:id", func(ctx *gin.Context) {
// 	start := time.Now()
// 	file, err := findFileById(pool, ctx.Param("id"))

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if file.Path == nil || len(*file.Path) == 0 {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no path"})
// 		return
// 	}

// 	fmt.Println("find took ", time.Since(start))
// 	start = time.Now()

// 	img, err := getCover(*file.File, *file.Path)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Println("xl to img took ", time.Since(start))
// 	start = time.Now()

// 	jpg, err := imageToJpg(img)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	imageDataBytes := jpg.Bytes()

// 	fmt.Println("img to jpg took ", time.Since(start))

// 	ctx.Header("Content-Size", string(len(imageDataBytes)))
// 	ctx.Data(http.StatusOK, "image/jpeg", imageDataBytes)
// })

// r.GET("/api/parser", func(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"files": progress.files,
// 		"books": progress.books,
// 		"time":  progress.time,
// 	})
// })

// r.POST("/api/parser", func(c *gin.Context) {
// 	progress = parser.ParserProgress{
// 		files: 0,
// 		books: 0,
// 		time:  0,
// 	}

// 	db.RebuildDb(pool, &progress)

// 	c.JSON(http.StatusOK, gin.H{
// 		"files": progress.files,
// 		"books": progress.books,
// 		"time":  progress.time,
// 	})
// })

// r.GET("/dl/:id", func(ctx *gin.Context) {
// 	start := time.Now()

// 	file, err := findFileById(pool, ctx.Param("id"))

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if file.Path == nil || len(*file.Path) == 0 {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no path"})
// 		return
// 	}

// 	fmt.Println("find took ", time.Since(start))
// 	start = time.Now()

// 	fb2, err := getFb2(*file.File, *file.Path)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// filename = [
// 	//   file.authors ? slugify(file.authors) : null,
// 	//   file.title ? slugify(file.title) : null,
// 	//   file.path,
// 	// ]
// 	//   .filter(Boolean)
// 	//   .join('.')

// 	ctx.Header("content-disposition", "attachment; filename=\"test.fb2\"")
// 	ctx.Header("content-type", "text/fb2+xml")
// 	ctx.String(http.StatusOK, fb2)

// 	fmt.Println("took ", time.Since(start))
// })
