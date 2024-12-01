package song

import (
	"effective-mobile/go/config"
	"effective-mobile/go/internal/common"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	cfg     *config.Config
	service *SongService
}

func NewSongHandler(cfg *config.Config, service *SongService) *SongHandler {
	return &SongHandler{
		cfg:     cfg,
		service: service,
	}
}

// swagger:route POST /songs Songs CreateSong
// Create a new song by providing the group and song name
//
// responses:
//
//	200: Response
//	400: ErrorResponse
//	401: ErrorResponse
//	500: ErrorResponse
func (h *SongHandler) CreateSong(ctx *gin.Context) {
	// swagger:parameters CreateSong
	type Request struct {
		// in: body
		Body struct {
			// Name of the song
			// required: true
			// example: Angel
			Song string `json:"song" binding:"required"`
			// Group of the song
			// required: true
			// example: Massive Attack
			Group string `json:"group" binding:"required"`
		}
	}

	var req Request
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid request", err))
		return
	}

	song := &SongModel{
		Group: req.Body.Group,
		Song:  req.Body.Song,
		Text:  make([]string, 0),
	}

	if err := h.service.CreateSong(ctx, song); err != nil {
		switch err {
		case ErrServiceUnavailable:
			ctx.JSON(http.StatusServiceUnavailable, common.FormatErrorResponse("service unavailable, try again later", err))
		default:
			log.Error("failed to create song: ", err)
			ctx.JSON(http.StatusInternalServerError, common.FormatErrorResponse("failed to create song", err))
		}

		return
	}

	ctx.JSON(http.StatusCreated, common.Response{
		Message: "song successfully created",
	})
}

// swagger:route DELETE /songs/:id Songs DeleteSong
// Delete a song by providing the song ID
//
// responses:
//
//	200: Response
//	400: ErrorResponse
//	401: ErrorResponse
//	500: ErrorResponse
func (h *SongHandler) DeleteSong(ctx *gin.Context) {
	// swagger:parameters DeleteSong
	type requestDefinition struct {
		// in: path
		// required: true
		ID int `uri:"id" json:"id" binding:"required"`
	}

	var req requestDefinition
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid song id", err))
		return
	}

	if err := h.service.DeleteSong(ctx, req.ID); err != nil {
		log.Error("failed to delete song: ", err)
		ctx.JSON(http.StatusInternalServerError, common.FormatErrorResponse("failed to delete song", err))
		return
	}

	ctx.JSON(http.StatusOK, common.Response{Message: "song successfully deleted"})
}

// swagger:route GET /songs/:id/lyrics Songs GetSongLyrics
// Get lyrics for a song
//
// responses:
//
//	200: LyricsResponse
//	400: ErrorResponse
//	401: ErrorResponse
//	500: ErrorResponse
func (h *SongHandler) GetSongLyrics(ctx *gin.Context) {
	// swagger:parameters GetSongLyrics
	type requestDescription struct {
		// ID of the song
		// in: path
		// required: true
		ID int `uri:"id" binding:"required" json:"id"`
		// Page number
		// in: query
		// required: false
		// default: 1
		Page int `form:"page,default=1" json:"page" binding:"min=1"`
		// Number of couplets per page
		// in: query
		// required: false
		// default: 1
		Limit int `form:"limit,default=1" json:"limit" binding:"min=1,max=10"`
	}

	var req requestDescription
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid song id", err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid query", err))
		return
	}

	couplets, metadata, err := h.service.GetSongLyrics(ctx, req.ID, req.Page, req.Limit)
	if err != nil {
		log.Error("failed to get song's lyrics: ", err)
		ctx.JSON(http.StatusInternalServerError, common.FormatErrorResponse("failed to get song", err))
		return
	}

	// swagger:response LyricsResponse
	type responseDescription struct {
		// in: body
		Body struct {
			PaginationMetadata common.PaginationMetadata `json:"metadata"`
			Message            string                    `json:"message"`
			Body               []string                  `json:"body"`
		}
	}

	ctx.JSON(http.StatusOK, common.PaginationResponse[string](responseDescription{
		Body: struct {
			PaginationMetadata common.PaginationMetadata `json:"metadata"`
			Message            string                    `json:"message"`
			Body               []string                  `json:"body"`
		}{
			Message:            "songs successfully retrieved",
			PaginationMetadata: *metadata,
			Body:               couplets,
		},
	}.Body))
}

// swagger:route GET /songs Songs GetSongs
// Get list of songs with optional filters
//
// responses:
//
//	200: SongsResponse
//	400: ErrorResponse
//	401: ErrorResponse
//	500: ErrorResponse
func (h *SongHandler) GetSongs(ctx *gin.Context) {
	// swagger:parameters GetSongs
	type requestDescription struct {
		// Page number
		// in: query
		// required: false
		// default: 1
		Page int `form:"page,default=1" json:"page" binding:"min=1"`
		// Number of songs per page
		// in: query
		// required: false
		// default: 10
		Limit int `form:"limit,default=10" json:"limit" binding:"min=1,max=10"`
	}

	var req requestDescription
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid query", err))
		return
	}

	// swagger:parameters GetSongs
	type filterDescription struct {
		// Name of the song
		// in: query
		// example: Angel
		// required: false
		Song *string `form:"song" json:"song"`
		// Group of the song
		// in: query
		// example: Massive Attack
		// required: false
		Group *string `form:"group" json:"group"`
		// Release date of the song
		// in: query
		// example: 2021-01-01
		// required: false
		ReleaseDate *time.Time `form:"release_date" time_format:"2006-01-02" json:"release_date"`
		// Lyrics of the song
		// in: query
		// example: Blah-blah-blah
		// required: false
		Text *string `form:"text" json:"text"`
		// Link to the song
		// in: query
		// example: https://example.com
		// required: false
		Link *string `form:"link" json:"link"`
	}

	var filter filterDescription
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid query", err))
		return
	}

	songs, metadata, err := h.service.GetSongs(ctx, SongFilter(filter), req.Page, req.Limit)
	if err != nil {
		log.Error("failed to get songs: ", err)
		ctx.JSON(http.StatusInternalServerError, common.FormatErrorResponse("failed to get songs", err))
		return
	}

	songsDTO := make([]SongDTO, 0, len(songs))
	for _, song := range songs {
		songsDTO = append(songsDTO, SongDTO{
			ID:          song.ID,
			Song:        song.Song,
			Group:       song.Group,
			ReleaseDate: DateOnly(song.ReleaseDate),
			Text:        song.Text,
			Link:        song.Link,
		})
	}

	// swagger:response SongsResponse
	type responseDescription struct {
		// in: body
		Body struct {
			PaginationMetadata common.PaginationMetadata `json:"metadata"`
			Message            string                    `json:"message"`
			Body               []SongDTO                 `json:"body"`
		}
	}

	ctx.JSON(http.StatusOK, common.PaginationResponse[SongDTO](responseDescription{
		Body: struct {
			PaginationMetadata common.PaginationMetadata `json:"metadata"`
			Message            string                    `json:"message"`
			Body               []SongDTO                 `json:"body"`
		}{
			Message:            "songs successfully retrieved",
			PaginationMetadata: *metadata,
			Body:               songsDTO,
		},
	}.Body))
}

// swagger:route PATCH /songs/:id Songs UpdateSong
// Update a song by providing the song ID and new data
//
// responses:
//
//	200: Response
//	400: ErrorResponse
//	401: ErrorResponse
//	500: ErrorResponse

func (h *SongHandler) UpdateSong(ctx *gin.Context) {
	// swagger:parameters UpdateSong
	type requestDescription struct {
		// Song ID
		// in: path
		// required: true
		SongID int `uri:"id" json:"id"`
		// in: body
		Body struct {
			// Name of the song
			// example: Angel
			// required: false
			Song *string `json:"song"`
			// Group of the song
			// example: Massive Attack
			// required: false
			Group *string `json:"group"`
			// Release date of the song
			// example: 2021-01-01
			// required: false
			ReleaseDate *time.Time `json:"release_date" time_format:"2006-01-02"`
			// Lyrics of the song
			// example: Blah-blah-blah
			// required: false
			Text *[]string `json:"text"`
			// Link to the song
			// example: https://example.com
			// required: false
			Link *string `json:"link"`
		}
	}

	var req requestDescription
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid song id", err))
		return
	}

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		ctx.JSON(http.StatusBadRequest, common.FormatErrorResponse("invalid request", err))
		return
	}

	if err := h.service.UpdateSong(ctx, UpdateSongDTO{
		SongID:      req.SongID,
		Group:       req.Body.Group,
		Song:        req.Body.Song,
		ReleaseDate: req.Body.ReleaseDate,
		Link:        req.Body.Link,
		Text:        req.Body.Text,
	}); err != nil {
		log.Error("failed to update song: ", err)
		ctx.JSON(http.StatusInternalServerError, common.FormatErrorResponse("failed to update song", err))
		return
	}

	ctx.JSON(http.StatusOK, common.Response{Message: "song successfully updated"})
}
