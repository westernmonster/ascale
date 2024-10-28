package def

import (
	"errors"
	"reflect"
	"strings"

	"github.com/codemodus/kace"
)

type Int32StatusItem struct {
	Key   int32  `json:"value"`
	Value string `json:"label"`
}

// Status
// All anonymous in this file used as global constants, DO NOT change it's value in runtime.
// currently only allow int32 and string field type

func GetStatusName(status, statusGlobalValue interface{}) (name string) {
	statusType := reflect.TypeOf(status)
	switch statusType.Kind() {
	case reflect.Int32:
		break
	case reflect.String:
		break
	default:
		panic(errors.New("invalid status"))
	}

	t := reflect.TypeOf(statusGlobalValue)
	v := reflect.ValueOf(statusGlobalValue)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)
		if fieldVal.Kind() != statusType.Kind() {
			panic(errors.New("Invalid status"))
		}

		switch fieldVal.Kind() {
		case reflect.Int32:
			statusVal := status.(int32)
			if int64(statusVal) == fieldVal.Int() {
				name = strings.ReplaceAll(kace.KebabUpper(field.Name), "-", " ")
				label := field.Tag.Get("label")
				if label != "" {
					name = label
				}
			}
			break
		case reflect.String:
			statusVal := status.(string)
			if statusVal == fieldVal.String() {
				name = strings.ReplaceAll(kace.KebabUpper(field.Name), "-", " ")
				label := field.Tag.Get("label")
				if label != "" {
					name = label
				}
			}
			break
		default:
			panic(errors.New("Invalid Kind"))
		}
	}

	return
}

func GetInt32StatusItems(statusGlobalValue interface{}) (items []*Int32StatusItem) {
	items = make([]*Int32StatusItem, 0)
	t := reflect.TypeOf(statusGlobalValue)
	v := reflect.ValueOf(statusGlobalValue)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		switch fieldVal.Kind() {
		case reflect.Int32:
			name := strings.ReplaceAll(kace.KebabUpper(field.Name), "_", " ")
			label := field.Tag.Get("label")
			if label != "" {
				name = label
			}
			items = append(items, &Int32StatusItem{
				Key:   int32(fieldVal.Int()),
				Value: name,
			})
			break
		// case reflect.String:
		// 	statusVal := status.(string)
		// 	if statusVal == fieldVal.String() {
		// 		return field.Name
		// 	}
		// 	break
		default:
			panic(errors.New("Invalid Kind"))
		}
	}
	return
}

var ClientID = struct {
	Portal string
}{
	Portal: "7ee9c4f86007ba41bc79bbfab1cd8a68",
}

var ValcodeType = struct {
	Login               int32
	Register            int32
	ResetPassword       int32
	StudentVerification int32
}{
	Login:               int32(1),
	Register:            int32(2),
	ResetPassword:       int32(3),
	StudentVerification: int32(4),
}

var UserStatus = struct {
	Active     int32
	Inactive   int32
	Discharged int32
}{
	Active:     int32(1),
	Inactive:   int32(2),
	Discharged: int32(3),
}

var AppHeader = struct {
	Version  string
	Platform string
}{
	Version:  "app-version",
	Platform: "app-platform",
}

var Origin = struct {
	Web string
	App string
}{
	Web: "web",
	App: "app",
}

var UserRole = struct {
	User  string
	Admin string
}{
	User:  "user",
	Admin: "admin",
}

var MysqlErrorCode = struct {
	DeadLock uint16
}{
	DeadLock: 1213,
}
