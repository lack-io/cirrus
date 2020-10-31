package net

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newClient() *HTTPClient {
	b := Builder()
	return b
}

func TestHTTPClient_Get(t *testing.T) {
	c := newClient()
	data, err := c.Get(context.Background(), "https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestHTTPClient_Post(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), JSON)
		assert.NotEqual(t, r.Header.Get("User-Agent"), "")
		assert.Contains(t, r.Header.Get("Authorization"), "token")

		assert.Contains(t, r.URL.RawQuery, "a=1")
		assert.Contains(t, r.URL.RawQuery, "b=2")

		data, _ := ioutil.ReadAll(r.Body)
		t.Log(string(data))
		r.Body.Close()
		rw.Write([]byte("ok"))
	})
	go http.ListenAndServe("127.0.0.1:30000", mux)

	c := newClient()
	c.SetContentType(JSON)
	c.SetToken("token")
	c.SetParams(map[string]string{"a": "1", "b": "2"})
	data, err := c.Post(context.Background(), "http://127.0.0.1:30000", strings.NewReader(`{"name":"xingyys"}`))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, []byte("ok"))
}

func TestHTTPClient_Post_Upload(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Content-Type"), MUL)
		assert.NotEqual(t, r.Header.Get("User-Agent"), "")
		assert.Contains(t, r.Header.Get("Authorization"), "token")

		err := r.ParseMultipartForm(4096)
		if err != nil {
			t.Fatal(err)
		}
		// 获取其他数据
		t.Log(r.MultipartForm.Value)
		// 获取上传的文件
		t.Log(r.MultipartForm.File)

		rw.Write([]byte("ok"))
	})
	go http.ListenAndServe("127.0.0.1:30000", mux)

	bb := &bytes.Buffer{}
	bw := multipart.NewWriter(bb)

	// 装载文件
	w, err := bw.CreateFormFile("uploadFile", "file1")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("123"))
	w, err = bw.CreateFormFile("uploadFile", "file2")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("456"))

	// 装载数据
	bw.WriteField("data1", "1")
	bw.WriteField("data2", "2")

	_ = bw.Close()

	c := newClient()
	c.SetContentType(bw.FormDataContentType())
	c.SetToken("token")
	data, err := c.Post(context.Background(), "http://127.0.0.1:30000", bb)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, []byte("ok"))
}
