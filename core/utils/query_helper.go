package utils

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	commonModels "star-wms/core/common/requests"
	"strings"
)

func BuildQuery(query *gorm.DB, filter interface{}) *gorm.DB {
	v := reflect.ValueOf(filter)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		whereTypeTag := field.Tag.Get("whereType")
		value := v.Field(i).Interface()

		if dbTag == "" || IsZero(value) {
			continue
		}

		if strings.Contains(dbTag, ",") {
			dbTags := strings.Split(dbTag, ",")
			var orQueryStrings []string
			var orQueryValues []interface{}
			for _, tag := range dbTags {
				if whereTypeTag == "like" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s LIKE ?", tag))
					orQueryValues = append(orQueryValues, "%"+fmt.Sprintf("%v", value)+"%")
				} else {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s = ?", tag))
					orQueryValues = append(orQueryValues, value)
				}
			}
			query = query.Or(strings.Join(orQueryStrings, " OR "), orQueryValues...)
		} else {
			if whereTypeTag == "like" {
				query = query.Where(fmt.Sprintf("%s LIKE ?", dbTag), "%"+fmt.Sprintf("%v", value)+"%")
			} else {
				query = query.Where(fmt.Sprintf("%s = ?", dbTag), value)
			}
		}
	}
	return query
}

func ApplySorting(query *gorm.DB, sorting commonModels.Sorting) *gorm.DB {
	if sorting.OrderBy != "" {
		direction := "ASC"
		if sorting.Desc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sorting.OrderBy, direction))
	}
	return query
}

func ApplyPagination(query *gorm.DB, pagination commonModels.Pagination) *gorm.DB {
	if pagination.PageSize > 0 {
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Limit(pagination.PageSize).Offset(offset)
	}
	return query
}

func IsZero(value interface{}) bool {
	v := reflect.ValueOf(value)
	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}
