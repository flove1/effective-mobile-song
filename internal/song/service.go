package song

import (
	"context"
	"effective-mobile/go/config"
	"effective-mobile/go/internal/common"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type SongService struct {
	config *config.Config
	repo   *SongRepository
}

func NewSongService(cfg *config.Config, repo *SongRepository) *SongService {
	return &SongService{
		config: cfg,
		repo:   repo,
	}
}

func (s *SongService) CreateSong(ctx context.Context, song *SongModel) error {
	defaultDate, err := time.Parse(time.DateOnly, "2000-01-01")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("%s/info?group=%s&song=%s", s.config.SongDetailAPI, song.Group, song.Song)
	resp, err := http.Get(query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var details struct {
			ReleaseDate string `json:"releaseDate"`
			Text        string `json:"text"`
			Link        string `json:"link"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
			return err
		}

		release_date, err := time.Parse("2006-01-02", details.ReleaseDate)
		if err != nil {
			return err
		}

		song.ReleaseDate = release_date
		song.Text = strings.Split(details.Text, "\n\n")
		song.Link = details.Link
	default:
		if s.config.Mode != "development" {
			return ErrServiceUnavailable
		}

		log.Error("failed to get song details: ", resp.StatusCode)
	}

	if song.ReleaseDate.IsZero() {
		song.ReleaseDate = defaultDate
	}

	if err := s.repo.CreateSong(ctx, song); err != nil {
		return err
	}

	return nil
}

func (s *SongService) DeleteSong(ctx context.Context, songID int) error {
	return s.repo.DeleteSong(ctx, songID)
}

func (s *SongService) GetSongLyrics(ctx context.Context, songID int, page, limit int) ([]string, *common.PaginationMetadata, error) {
	return s.repo.GetSongLyrics(ctx, songID, page, limit)
}

func (s *SongService) GetSongs(ctx context.Context, filter SongFilter, page, limit int) ([]*SongModel, *common.PaginationMetadata, error) {
	return s.repo.GetSongs(ctx, filter, page, limit)
}

func (s *SongService) UpdateSong(ctx context.Context, dto UpdateSongDTO) error {
	return s.repo.UpdateSong(ctx, dto)
}
