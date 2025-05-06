package models

import (
	"github.com/google/uuid"
	"time"
)

type Author struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName   string     `gorm:"unique;not null;type:varchar(100);default:null" json:"first_name"`
	LastName    string     `gorm:"unique;not null;type:varchar(100);default:null" json:"last_name"`
	BirthDate   *time.Time `json:"birth_date"`
	DeathDate   *time.Time `json:"death_date"`
	Nationality string     `json:"nationality" gorm:"size:50"`
	Biography   string     `json:"biography" gorm:"type:text"`
	Website     string     `json:"website" gorm:"size:255"`
	CoverURL    string     `json:"cover_url" gorm:"type:text;default:'https://avatars.mds.yandex.net/i?id=e5be5c8d1fe86f031ac75d8cf920d292_l-5324012-images-thumbs&n=13.jpg'"`
	Books       []Book     `gorm:"foreignKey:AuthorID" json:"books"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
