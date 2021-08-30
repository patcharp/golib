package database

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

// Model struct
type Model struct {
	Seq       int64          `json:"seq" gorm:"primary_key;auto_increment:false;"`
	Uid       uuid.UUID      `json:"uid" gorm:"primary_key;index"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// BeforeCreate hook table
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if uuid.Equal(m.Uid, uuid.Nil) {
		m.Uid = uuid.NewV4()
	}
	m.Seq = time.Now().UnixNano()
	return nil
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
