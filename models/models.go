package models

type ProgressSubab struct {
	Subab   int  `json:"subab"`
	Selesai bool `json:"selesai"`
}

type ProgressBab struct {
	Bab          int             `json:"bab"`
	ProgresSubab []ProgressSubab `json:"subab"`
	Selesai      bool            `json:"selesai"`
}

type ProgressUser struct {
	Id         string        `json:"id"`
	Username   string        `json:"username"`
	ProgresBab []ProgressBab `json:"bab"`
}
