package w2g_database_package

import (
	"github.com/jinzhu/gorm"
	guuid "github.com/satori/go.uuid"
	"time"
	"unicode"
)

type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time `json:"-"`
	DeletedBy string
}

func (base *Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", GetHash())
}

func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

func ConvertStructField2DatabaseField(s string) string {

	var field string

	for ind, ch := range s {
		if unicode.IsUpper(ch) {
			if ind == 0 {
				field += string(unicode.ToLower(ch))
			} else {
				field = field + "_" + string(unicode.ToLower(ch))
			}
		} else {
			field += string(ch)
		}
	}

	return field
}

func FindInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
