package domain

type User struct {
	SessionID string `json:"session_id"`
	Lied      int    `json:"lied"`
}

func (i *User) ScanValues() []interface{} {
	return []interface{}{
		&i.SessionID,
		&i.Lied,
	}
}

func (i *User) Values() []interface{} {
	return []interface{}{
		i.SessionID,
		i.Lied,
	}
}

func (i *User) IsValid() bool {
	if i.SessionID == "" {
		return false
	}
	return true
}
