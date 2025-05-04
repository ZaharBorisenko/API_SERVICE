package models

import (
	"time"
)

type Author struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	DeathDate   *time.Time `json:"death_date,omitempty"`
	Nationality string     `json:"nationality,omitempty" gorm:"size:50"`
	Biography   string     `json:"biography,omitempty" gorm:"type:text"`
	Website     string     `json:"website,omitempty" gorm:"size:255"`
	CoverURL    string     `json:"cover_url" gorm:"type:text;default:'https://avatars.mds.yandex.net/i?id=e5be5c8d1fe86f031ac75d8cf920d292_l-5324012-images-thumbs&n=13.jpg'"`
	Books       []Book     `gorm:"foreignKey:AuthorID" json:"books"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
