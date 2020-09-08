package w2g_database_package

import (
	"fmt"
	"github.com/jinzhu/gorm"
	guuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
func SetOrder(r *http.Request, db *gorm.DB) *gorm.DB {

	order := r.URL.Query().Get("order-by")

	params := strings.Split(order, "|")
	if len(params) != 2 {
		return db
	}

	field := ConvertStructField2DatabaseField(params[0])
	direction := params[1]
	DB := db.Order(field + " " + direction)

	return DB

}

func SetFilters(r *http.Request, db *gorm.DB) *gorm.DB {

	var listToIgnore = []string{"page", "per-page", "order-by"}

	DB := db
	// TODO we need to get a way how to know the field type in DB.
	if r.URL.RawQuery != "" {
		fmt.Println("RawQuery:", r.URL.RawQuery)
		m, err := url.ParseQuery(r.URL.RawQuery)
		if err == nil {
			for k, v := range m {
				_, result := FindInSlice(listToIgnore, k)
				if result {
					continue
				}

				field := ConvertStructField2DatabaseField(k)
				// Here we have to use additional char, it help us to define a type of field.
				if strings.Index(v[0], "|") != -1 {
					tmp := strings.Split(v[0], "|")
					operation := tmp[0]
					condition := tmp[1]
					fieldType := tmp[2]

					switch fieldType {

					case "b":
						c, err := strconv.Atoi(condition)
						if err != nil {
							DB = DB.Where(field+" "+operation+" ?", c)
						}

					case "s":
						DB = DB.Where(field+" "+operation+" ?", condition)

					default:

					}

					fmt.Printf("Key: %q Values: %q\n", k, v)
				}
			}
		}
	}

	return DB

}

func SetPagePerPageValues(r *http.Request, db *gorm.DB) *gorm.DB {

	var page int
	var perPage int

	pageStr := r.URL.Query().Get("page")

	if pageStr == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageStr)
	}

	perPageStr := r.URL.Query().Get("per-page")
	if perPageStr == "" {
		perPage = 15
	} else {
		perPage, _ = strconv.Atoi(perPageStr)
		if perPage > 1000 {
			perPage = 1000
		}
	}

	DB := db.Offset((page - 1) * perPage).Limit(perPage)
	return DB

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
