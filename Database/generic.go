package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func Upsert(guildId string, table string, idField string, idValue any, fieldsAndValues ...any) (int64, error) {
	if len(fieldsAndValues)%2 != 0 {
		err := errors.New("upsert failed, fields and values are an odd number")
		logger.Error(guildId, err)
		return 0, err
	}

	var (
		columns    []string
		values     []any
		updateList []string
	)

	for i := 0; i < len(fieldsAndValues); i += 2 {
		field, ok := fieldsAndValues[i].(string)
		if !ok {
			err := fmt.Errorf("upsert failed, field name at position %d is not a string", i)
			logger.Error(guildId, err)
			return 0, err
		}
		value := fieldsAndValues[i+1]

		columns = append(columns, field)
		values = append(values, value)
		updateList = append(updateList, fmt.Sprintf("%s = ?", field))
	}

	if helpers.IsZero(idValue) {
		placeholderList := strings.Repeat("?,", len(columns))
		placeholderList = placeholderList[:len(placeholderList)-1]

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), placeholderList)
		res, err := Db.ExecContext(context.Background(), query, values...)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	} else {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table, strings.Join(updateList, ", "), idField)
		values = append(values, idValue)
		_, err := Db.ExecContext(context.Background(), query, values...)
		if err != nil {
			return 0, err
		}
		return toInt64(idValue), nil
	}
}

// Helper to cast ID to int64
func toInt64(v any) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case uint:
		return int64(val)
	case uint64:
		return int64(val)
	default:
		return 0
	}
}
