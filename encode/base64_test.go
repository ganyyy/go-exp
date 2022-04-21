package main

import (
	"encoding/base64"
	"testing"
)

func TestBase64(t *testing.T) {
	var src = "AAAAA"
	var encode = base64.StdEncoding.EncodeToString([]byte(src))

	t.Logf("%+v", encode)

	var dst, _ = base64.URLEncoding.DecodeString(encode)

	t.Logf("%v,encode:%s, dst %s", src, encode, dst)
}
