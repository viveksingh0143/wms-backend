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
		valueType := reflect.TypeOf(value)

		if dbTag == "" || IsZero(value) {
			continue
		}

		if strings.Contains(dbTag, ",") {
			dbTags := strings.Split(dbTag, ",")
			var orQueryStrings []string
			var orQueryValues []interface{}
			for _, tag := range dbTags {
				if whereTypeTag == "startswith" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s LIKE ?", tag))
					orQueryValues = append(orQueryValues, fmt.Sprintf("%v", value)+"%")
				} else if whereTypeTag == "like" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s LIKE ?", tag))
					orQueryValues = append(orQueryValues, "%"+fmt.Sprintf("%v", value)+"%")
				} else if whereTypeTag == "ne" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s != ?", tag))
					orQueryValues = append(orQueryValues, value)
				} else if whereTypeTag == "gte" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s >= ?", tag))
					orQueryValues = append(orQueryValues, value)
				} else if whereTypeTag == "lte" {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s <= ?", tag))
					orQueryValues = append(orQueryValues, value)
				} else if whereTypeTag == "in" {
					if valueType.Kind() == reflect.Slice {
						sliceValue := reflect.ValueOf(value)
						num := sliceValue.Len()
						args := make([]interface{}, num)
						for j := 0; j < num; j++ {
							args[j] = sliceValue.Index(j).Interface()
						}
						placeholders := strings.TrimRight(strings.Repeat("?,", len(args)), ",")

						orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s IN (%s)", tag, placeholders))
						orQueryValues = append(orQueryValues, args...)
					}
				} else {
					orQueryStrings = append(orQueryStrings, fmt.Sprintf("%s = ?", tag))
					orQueryValues = append(orQueryValues, value)
				}
			}
			query = query.Where(strings.Join(orQueryStrings, " OR "), orQueryValues...)
		} else {
			if whereTypeTag == "startswith" {
				query = query.Where(fmt.Sprintf("%s LIKE ?", dbTag), fmt.Sprintf("%v", value)+"%")
			} else if whereTypeTag == "like" {
				query = query.Where(fmt.Sprintf("%s LIKE ?", dbTag), "%"+fmt.Sprintf("%v", value)+"%")
			} else if whereTypeTag == "ne" {
				query = query.Where(fmt.Sprintf("%s != ?", dbTag), value)
			} else if whereTypeTag == "gte" {
				query = query.Where(fmt.Sprintf("%s >= ?", dbTag), value)
			} else if whereTypeTag == "lte" {
				query = query.Where(fmt.Sprintf("%s <= ?", dbTag), value)
			} else if whereTypeTag == "in" {
				if valueType.Kind() == reflect.Slice {
					sliceValue := reflect.ValueOf(value)
					num := sliceValue.Len()
					args := make([]interface{}, num)
					for j := 0; j < num; j++ {
						args[j] = sliceValue.Index(j).Interface()
					}

					placeholders := strings.TrimRight(strings.Repeat("?,", len(args)), ",")
					query = query.Where(fmt.Sprintf("%s IN (%s)", dbTag, placeholders), args...)
				}
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
		query = query.Limit(pagination.PageSize).Offset(pagination.Start)
	}
	return query
}

func IsZero(value interface{}) bool {
	v := reflect.ValueOf(value)
	return reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface())
}
