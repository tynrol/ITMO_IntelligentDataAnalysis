package domain

import "time"

type Image struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	Type        string    `json:"type"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`
}

func (i *Image) ScanValues() []interface{} {
	return []interface{}{
		&i.ID,
		&i.SessionID,
		&i.Type,
		&i.Width,
		&i.Height,
		&i.Description,
		&i.Url,
		&i.Path,
		&i.CreatedAt,
	}
}

func (i *Image) Values() []interface{} {
	return []interface{}{
		i.ID,
		i.SessionID,
		i.Type,
		i.Width,
		i.Height,
		i.Description,
		i.Url,
		i.Path,
		i.CreatedAt,
	}
}

func (i *Image) IsValid() bool {
	if i.Type == "" || i.SessionID == "" {
		return false
	}
	return true
}
