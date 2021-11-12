package api

import (
	"github.com/gin-gin-gonic/gin",
	db "github.com/EemeliSaari/turso/db/sqlc"
)

type Server struct {
	queries db.Queries
	router *gin.Engine
}

func NewServer(queries db.Queries) (*Server, error) {
	server := &Server{
		queries: queries,
		router: server.initRouter()
	}
	return server, nil
}

func initRouter() (*gin.Engine) {
	router := gin.Default()

	router.POST("/articles", server.createArticle)
	router.GET("/articles/:id", server.getArticle)
	router.GET("/articles", server.listArticles)

	return &router
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
