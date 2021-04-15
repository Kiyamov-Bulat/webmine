package utils

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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
		return "", errors.New("field value is not valid")
	}
	return structFieldValue.String(), nil
}

func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	prefix := strings.TrimPrefix(pathToZip, filepath.Base(pathToZip))
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, prefix)
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	err = destinationFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	log.Println(strconv.Quote(filepath.Base(filePath)), filePath)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

func GetResultOfExecCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	res := string(stdout)
	return res[:len(res)-1], nil
}
