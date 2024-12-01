package song

import "time"

type SongModel struct {
	ID          int       `db:"id"`
	Song        string    `db:"song"`
	Group       string    `db:"group"`
	ReleaseDate time.Time `db:"release_date"`
	Text        []string  `db:"text"`
	Link        string    `db:"link"`
}

type SongFilter struct {
	Song        *string
	Group       *string
	ReleaseDate *time.Time
	Text        *string
	Link        *string
}
