package domain

type User struct {
	SessionID string `json:"session_id"`
	IsLying   bool   `json:"is_lying"`
}

func (i *User) ScanValues() []interface{} {
	return []interface{}{
		&i.SessionID,
		&i.IsLying,
	}
}

func (i *User) Values() []interface{} {
	return []interface{}{
		i.SessionID,
		i.IsLying,
	}
}
