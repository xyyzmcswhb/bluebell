package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community-id":1,
		"title":"test",
		"content":"just a test",
	}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	//方法一：判断响应的内容是否按预期返回了需要登录的错误
	assert.Contains(t, w.Body.String(), "需要登录")
	// assert.Equal("")

	//方法2：将响应内容反序列化到Responsedata,判断字段是否与预期一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil { //反序列化
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedAuth)
}
