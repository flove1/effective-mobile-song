package http

import (
	"github.com/gin-gonic/gin"
)

func newRouter(handlers Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/healthcheck", healthcheck)

	r.GET("/songs", handlers.SongHandler.GetSongs)
	r.GET("/songs/:id/lyrics", handlers.SongHandler.GetSongLyrics)
	r.POST("/songs", handlers.SongHandler.CreateSong)
	r.DELETE("/songs/:id", handlers.SongHandler.DeleteSong)
	r.PATCH("/songs/:id", handlers.SongHandler.UpdateSong)

	return r
}

func healthcheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
