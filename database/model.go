package database

import (
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
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

// KSModel is uid alternative
type KSModel struct {
	Uid       ksuid.KSUID    `json:"uid" gorm:"primary_key;type:varbinary(27);index"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (ks *KSModel) BeforeCreate(tx *gorm.DB) error {
	if ksuid.Compare(ksuid.Nil, ks.Uid) == 0 {
		ks.Uid = ksuid.New()
	}
	return nil
}
