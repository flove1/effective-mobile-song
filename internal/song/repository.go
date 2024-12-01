package song

import (
	"context"
	"effective-mobile/go/config"
	"effective-mobile/go/internal/common"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"

	log "github.com/sirupsen/logrus"
)

type SongRepository struct {
	config *config.Config
	db     *pgxpool.Pool
}

const songsTable = "songs"

func NewSongRepository(cfg *config.Config, db *pgxpool.Pool) *SongRepository {
	return &SongRepository{
		config: cfg,
		db:     db,
	}
}

func (r *SongRepository) CreateSong(ctx context.Context, song *SongModel) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (song, "group", release_date, "text", link) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, songsTable)
	err := r.db.QueryRow(ctx, query, song.Song, song.Group, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		return err
	}

	log.Debug("song created with ID: ", song.ID)
	return nil
}

func (r *SongRepository) DeleteSong(ctx context.Context, songID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, songsTable)
	_, err := r.db.Exec(context.Background(), query, songID)
	if err != nil {
		return err
	}

	log.Debug("song deleted with ID: ", songID)
	return nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, dto UpdateSongDTO) error {
	query := fmt.Sprintf(`
        UPDATE %s SET 
            song = COALESCE($1, song), 
            "group" = COALESCE($2, "group"), 
            release_date = COALESCE($3, release_date), 
            "text" = COALESCE($4, "text"), 
            link = COALESCE($5, link)
        WHERE id = $6
    `, songsTable)

	_, err := r.db.Exec(context.Background(), query,
		dto.Song,
		dto.Group,
		dto.ReleaseDate,
		dto.Text,
		dto.Link,
		dto.SongID,
	)

	if err != nil {
		return err
	}

	log.Debug("song updated with ID: ", dto.SongID)
	return nil
}

func (r *SongRepository) GetSongs(ctx context.Context, filter SongFilter, page, limit int) ([]*SongModel, *common.PaginationMetadata, error) {
	page = max(1, page)
	limit = min(10, max(1, limit))

	if filter.Text != nil {
		var text = "%" + strings.Trim(*filter.Text, " %") + "%"
		filter.Text = &text
	}

	totalQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT DISTINCT ON (id)
				id, 
				song, 
				"group", 
				release_date, 
				"text", 
				link,
				couplet
			FROM %s s
			JOIN LATERAL unnest(coalesce(nullif(s.text,'{}'),array[null::text])) AS couplet ON true
			WHERE 
				($1::text IS NULL OR LOWER(song) LIKE $1) AND
				($2::text IS NULL OR LOWER("group") LIKE $2) AND
				($3::date IS NULL OR release_date = $3) AND
				($4::text IS NULL OR LOWER(couplet) LIKE $4) AND
				($5::text IS NULL OR LOWER(link) LIKE $5)
		) as filtered
	`, songsTable)

	query := fmt.Sprintf(`
		SELECT DISTINCT ON (id)
			id, 
			song, 
			"group", 
			release_date, 
			"text", 
			link,
			couplet
		FROM %s s
		JOIN LATERAL unnest(coalesce(nullif(s.text,'{}'),array[null::text])) AS couplet ON true
		WHERE 
			($1::text IS NULL OR LOWER(song) LIKE $1) AND
			($2::text IS NULL OR LOWER("group") LIKE $2) AND
			($3::date IS NULL OR release_date = $3) AND
			($4::text IS NULL OR LOWER(couplet) LIKE $4) AND
			($5::text IS NULL OR LOWER(link) LIKE $5)
        LIMIT $6 OFFSET $7
    `, songsTable)

	var totalCount int
	var songs []*SongModel

	err := r.db.QueryRow(ctx, totalQuery,
		filter.Song,
		filter.Group,
		filter.ReleaseDate,
		filter.Text,
		filter.Link,
	).Scan(&totalCount)

	if err != nil {
		return nil, nil, err
	}

	metadata := common.CalculateMetadata(totalCount, page, limit)

	rows, err := r.db.Query(context.Background(), query,
		filter.Song,
		filter.Group,
		filter.ReleaseDate,
		filter.Text,
		filter.Link,
		limit, max(0, page-1)*limit,
	)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song SongModel
		err = rows.Scan(
			&song.ID,
			&song.Song,
			&song.Group,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			nil,
		)

		if err != nil {
			return nil, nil, err
		}

		songs = append(songs, &song)
	}

	return songs, &metadata, nil
}

func (r *SongRepository) GetSongLyrics(ctx context.Context, songID int, page, limit int) ([]string, *common.PaginationMetadata, error) {
	page = max(1, page)
	limit = min(10, max(1, limit))

	totalQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT 
				unnest(text)
			FROM %s 
			WHERE id = $1
		) as couplets
	`, songsTable)

	query := fmt.Sprintf(`
		SELECT 
			unnest(text)
		FROM %s 
		WHERE id = $1
        LIMIT $2 OFFSET $3
	`, songsTable)

	var totalCount int
	var couplets []string

	if err := r.db.QueryRow(ctx, totalQuery, songID).Scan(&totalCount); err != nil {
		return nil, nil, err
	}

	metadata := common.CalculateMetadata(totalCount, page, limit)

	rows, err := r.db.Query(context.Background(), query, songID, limit, max(0, page-1)*limit)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var couplet string
		err = rows.Scan(&couplet)
		if err != nil {
			return nil, nil, err
		}

		couplets = append(couplets, couplet)
	}

	return couplets, &metadata, nil

}
