package main

import (
	"testing"

	"github.com/astaxie/beego/httplib"
)

func TestPostForm(t *testing.T) {
	var request = httplib.Post("http://localhost:8899/postform")
	request.
		Param("age", "99").
		Param("name", "lalala").
		Param("address", "123131")

	var resp, err = request.Bytes()
	t.Logf("%v, %v", string(resp), err)
}
