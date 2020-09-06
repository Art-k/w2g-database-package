package w2g_database_package

import (
	"github.com/jinzhu/gorm"
	guuid "github.com/satori/go.uuid"
	"time"
)

type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time `json:"-"`
	// DeletedBy string
}

func (base *Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GetHash())
}

func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}
