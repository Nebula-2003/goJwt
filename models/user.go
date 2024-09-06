package models

import (
	"time"

	"github.com/Nebula-2003/goJwt/config"
)

type User struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt *time.Time  `json:"deletedAt,omitempty"`
	Email     string      `gorm:"unique" json:"email"`
	Name      string      `json:"name"`
	Password  string      `json:"-"`
	Role      config.Role `json:"role" gorm:"type:varchar(20)"`
}
