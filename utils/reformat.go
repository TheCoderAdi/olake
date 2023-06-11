package utils

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/piyushsingariya/syndicate/types"
)

var DateTimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.000000",
}

func getFirstNotNullType(datatypes []types.DataType) types.DataType {
	for _, datatype := range datatypes {
		if datatype != types.Null {
			return datatype
		}
	}

	return types.Null
}

func ReformatValueOnDataTypes(datatypes []types.DataType, format string, v any) (any, error) {
	return ReformatValue(getFirstNotNullType(datatypes), format, v)
}

func ReformatValue(dataType types.DataType, format string, v any) (any, error) {
	switch dataType {
	case types.Null:
		return nil, nil
	case types.Boolean:
		// reformat boolean
		booleanValue, ok := v.(bool)
		if ok {
			return booleanValue, nil
		}
		return v, fmt.Errorf("found to be boolean, but value is not boolean : %v", v)
	case types.Integer:
		return ReformatInt64(v)
	case types.String:
		if format == "date" || format == "date-time" {
			return ReformatDate(v)
		}

		return fmt.Sprintf("%v", v), nil
	case types.Number:
		return ReformatFloat64(v)
	case types.Array:
		if value, isArray := v.([]any); isArray {
			return value, nil
		}

		// make it an array
		return []any{v}, nil
	default:
		return v, nil
	}
}

// reformat date
func ReformatDate(v interface{}) (time.Time, error) {
	parsed, err := func() (time.Time, error) {
		switch v := v.(type) {
		// we assume int64 is in seconds and don't currently scale to the precision
		case int64:
			return time.Unix(v, 0), nil
		case *int64:
			switch {
			case v != nil:
				return time.Unix(*v, 0), nil
			default:
				return time.Time{}, nil
			}
		case time.Time:
			return v, nil
		case *time.Time:
			switch {
			case v != nil:
				return *v, nil
			default:
				return time.Time{}, nil
			}
		case sql.NullTime:
			switch v.Valid {
			case true:
				return v.Time, nil
			default:
				return time.Time{}, nil
			}
		case *sql.NullTime:
			switch v.Valid {
			case true:
				return v.Time, nil
			default:
				return time.Time{}, nil
			}
		case nil:
			return time.Time{}, nil
		case string:
			return parseCHDateTime(v)
		case *string:
			if v == nil || *v == "" {
				return time.Time{}, nil
			} else {
				return parseCHDateTime(*v)
			}
		}
		return time.Time{}, nil
	}()
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func parseCHDateTime(value string) (time.Time, error) {
	var tv time.Time
	var err error
	for _, layout := range DateTimeFormats {
		tv, err = time.Parse(layout, value)
		if err == nil {
			return time.Date(
				tv.Year(), tv.Month(), tv.Day(), tv.Hour(), tv.Minute(), tv.Second(), tv.Nanosecond(), tv.Location(),
			), nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse datetime from available formats: %v : %s", DateTimeFormats, err)
}

func ReformatInt64(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	}

	return int64(0), fmt.Errorf("failed to change %v (type:%T) to int64", v, v)
}

func ReformatFloat64(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return float64(0), fmt.Errorf("failed to change string %v to float64: %w", v, err)
		}
		return f, nil
	}

	return float64(0), fmt.Errorf("failed to change %v (type:%T) to float64", v, v)
}
