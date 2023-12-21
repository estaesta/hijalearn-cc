package models

type ProgressBab struct {
	Subab   map[string]bool `json:"subab"`
	Selesai bool            `json:"selesai"`
}

type ProgressUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	// ProgresBab []ProgressBab `json:"bab"`
}

type UserProfile struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	ProfilePictureURL string `json:"profile_picture_url"`
}
