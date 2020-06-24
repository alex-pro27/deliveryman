package common

import (
	"git.samberi.com/dois/delivery_api/proto/gen"
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
)

type PaginateData struct {
	Paginate *gen.PaginateMessage `json:"paginate"`
	Data     interface{}          `json:"result"`
}

func Paginate(
	model interface{},
	queryset *gorm.DB,
	page int,
	limit int,
	preloading []string,
	serializeFunc func(interface{}) interface{},
) PaginateData {
	t := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		obj := t.Elem()
		if obj.Kind() == reflect.Slice {
			if page == 0 {
				page = 1
			}
			count := 0
			queryset.Model(model).Count(&count)
			start := page*limit - limit
			for _, preload := range preloading {
				queryset = queryset.Preload(preload)
			}
			queryset.Offset(start).Limit(limit).Find(model)

			var result []interface{}

			for i := 0; i < obj.Len(); i++ {
				result = append(result, serializeFunc(obj.Index(i).Interface()))
			}

			var length int
			if length = limit; limit != len(result) {
				length = len(result)
			}
			return PaginateData{
				Paginate: &gen.PaginateMessage{
					CurrentPage: int32(page),
					Count:       int32(count),
					CountPage:   int32(math.Ceil(float64(count) / float64(limit))),
					Length:      int32(length),
				},
				Data: result,
			}
		} else {
			panic("model not Slice")
		}
	} else {
		panic("model not Ptr")
	}
}
