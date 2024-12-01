package song

import (
	"encoding/json"
	"time"
)

type UpdateSongDTO struct {
	SongID      int        `uri:"id"`
	Group       *string    `json:"group"`
	Song        *string    `json:"song"`
	ReleaseDate *time.Time `json:"release_date" time_format:"2006-01-02"`
	Text        *[]string  `json:"text"`
	Link        *string    `json:"link"`
}

// swagger:model SongDTO
type SongDTO struct {
	ID    int    `json:"id"`
	Song  string `json:"song"`
	Group string `json:"group"`
	// example: 2021-01-01
	ReleaseDate DateOnly `json:"release_date" time_format:"2006-01-02"`
	Text        []string `json:"text"`
	Link        string   `json:"link"`
}

// swagger:type DateOnly
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	b := make([]byte, 0, len(time.DateOnly)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.DateOnly)
	b = append(b, '"')

	return b, nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	*d = DateOnly(t)
	return nil
}
