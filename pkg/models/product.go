package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"time"
	"unicode"
)

//nolint:gochecknoinits
func init() {
	govalidator.CustomTypeTagMap.Set("imgurl", func(i interface{}, o interface{}) bool {
		if i == nil {
			return true
		}

		imgUrl, ok := i.(string)
		if !ok {
			return false
		}

		if imgUrl == "" || strings.HasSuffix(imgUrl, ".png") || strings.HasSuffix(imgUrl, ".jpeg") ||
			strings.HasSuffix(imgUrl, ".jpg") {
			return true
		}

		return false
	})
}

type Product struct {
	ID          uint64    `json:"id"              valid:"required"`
	SalerID     uint64    `json:"saler_id"        valid:"required"`
	Title       string    `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Description string    `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	ImageUrl    string    `json:"image_url"       valid:"optional, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Price       uint64    `json:"price"           valid:"required"`
	CreatedAt   time.Time `json:"created_at"      valid:"required"`
}

type ProductWithIsMy struct {
	ID          uint64    `json:"id"              valid:"required"`
	SalerID     uint64    `json:"saler_id"        valid:"required"`
	Title       string    `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Description string    `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"` //nolint:nolintlint
	ImageUrl    string    `json:"image_url"       valid:"optional, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`   //nolint:nolintlint
	Price       uint64    `json:"price"           valid:"required"`
	IsMy        bool      `json:"is_my"           valid:"required"`
	CreatedAt   time.Time `json:"created_at"      valid:"required"`
}

type PreProduct struct {
	SalerID     uint64 `json:"saler_id"        valid:"required"`
	Title       string `json:"title"           valid:"required, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"`         //nolint:nolintlint
	Description string `json:"description"     valid:"required, length(1|4000)~Описание должно быть длинной от 1 до 4000 симвволов"`       //nolint:nolintlint
	ImageUrl    string `json:"image_url"       valid:"imgurl, optional, length(1|256)~Заголовок должен быть длинной от 1 до 256 символов"` //nolint:nolintlint
	Price       uint64 `json:"price"           valid:"required"`
}

func (p *PreProduct) Trim() {
	p.Title = strings.TrimFunc(p.Title, unicode.IsSpace)
	p.Description = strings.TrimFunc(p.Description, unicode.IsSpace)
}

func (p *Product) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.Title = sanitizer.Sanitize(p.Title)
	p.Description = sanitizer.Sanitize(p.Description)
}

func (p *ProductWithIsMy) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	p.Title = sanitizer.Sanitize(p.Title)
	p.Description = sanitizer.Sanitize(p.Description)
}
