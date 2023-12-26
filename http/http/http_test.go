package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

const (
	UserAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0`
)

func TestPUBGChatBot(t *testing.T) {
	var collector = colly.NewCollector()
	collector.UserAgent = UserAgent
	collector.OnRequest(func(r *colly.Request) {
		t.Logf("Visiting %s", r.URL.String())
	})
	collector.OnResponse(func(r *colly.Response) {
		t.Logf("Visited %s", r.Request.URL.String())
		t.Logf("Headers %s", r.Headers)
		t.Logf("Response %s", string(r.Body))

	})
	err := collector.Request("GET",
		"https://support.pubg.com/hc/zh-cn/articles/5048335129497",
		nil,
		nil,
		http.Header{
			"Accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"Acept-Encoding":     {"gzip, deflate, br"},
			"Acept-Language":     {"zh-CN,zh;q=0.9"},
			"Cahce-Control":      {"no-cache"},
			"Dnt":                {"1"},
			"User-Agent":         {UserAgent},
			"Sec-Ch-Ua":          {`"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`},
			"Sec-Ch-Ua-Mobile":   {"?0"},
			"Sec-Ch-Ua-Platform": {"Windows"},
			"Sec-Fetch-Dest":     {"document"},
			"Sec-Fetch-Mode":     {"navigate"},
			"Sec-Fetch-Site":     {"none"},
			"Sec-Fetch-User":     {"?1"},
		},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestHttpGet(t *testing.T) {
	var jar, _ = cookiejar.New(nil)
	var client = http.Client{
		Transport: http.DefaultTransport,
		Jar:       jar,
		Timeout:   0,
	}

	var req, err = http.NewRequest(http.MethodGet, "https://support.pubg.com/hc/zh-cn/articles/5048335129497", nil)
	assert.NoError(t, err)
	req.Header = http.Header{
		"Accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Acept-Encoding":     {"gzip, deflate, br"},
		"Acept-Language":     {"zh-CN,zh;q=0.9"},
		"Cahce-Control":      {"no-cache"},
		"Dnt":                {"1"},
		"User-Agent":         {UserAgent},
		"Sec-Ch-Ua":          {`"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {"Windows"},
		"Sec-Fetch-Dest":     {"document"},
		"Sec-Fetch-Mode":     {"navigate"},
		"Sec-Fetch-Site":     {"none"},
		"Sec-Fetch-User":     {"?1"},
	}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	t.Logf("Status %s", resp.Status)
	var data, _ = io.ReadAll(resp.Body)
	t.Logf("%s", string(data))
}
