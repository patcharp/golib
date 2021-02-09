package database

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Model struct
type Model struct {
	Seq       int64      `json:"seq" gorm:"primary_key;auto_increment:false;" sql:"index"`
	Uid       uuid.UUID  `json:"uid" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// BeforeCreate hook table
func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	if f, ok := scope.FieldByName("uid"); ok {
		if f.IsBlank {
			if err := scope.SetColumn("uid", uuid.NewV4()); err != nil {
				return err
			}
		}
	}
	return scope.SetColumn("seq", time.Now().UnixNano())
}

// Preload
type Preload struct {
	Column     string        `json:"column"`
	Conditions []interface{} `json:"conditions"`
}

// Where
type Where struct {
	Condition string        `json:"condition"`
	Arguments []interface{} `json:"args"`
}
