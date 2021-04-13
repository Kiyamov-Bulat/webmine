package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

const (
	ERROR   int = -1
	WARNING int = 0
	SUCCESS int = 1
)

func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Stou(sid string) (uint, error) {
	id, err := strconv.Atoi(sid)
	return uint(id), err
}

func MakeLine(arr []string) string {
	var result string
	for i, str := range arr {
		if i < len(arr)-1 {
			result += str + " "
		} else {
			result += str
		}
	}
	return result
}

func ArrContain(arr []string, value string) bool {
	for _, val := range arr {
		if val == value {
			return true
		}
	}
	return false
}

func GetType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetStructValue(i interface{}, fieldName string) (string, error) {
	structValue := reflect.ValueOf(i).Elem()
	structFieldValue := structValue.FieldByName(fieldName)
	if !structFieldValue.IsValid() || structFieldValue.Type().Name() == "string" {
		return "", errors.New("Field value is not valid")
	}
	return structFieldValue.String(), nil
}
