package common

import (
	"encoding/json"
	"reflect"
)

func TaskName(task ITask) string {
	return reflect.TypeOf(task).String()
}

func TaskToDetail(task ITask) string {
	var bs, _ = json.Marshal(task)
	return string(bs)
}
