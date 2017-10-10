package fileupload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestProcessFileUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ProcessFileUpload(r)
		if err != nil {
			t.Fatal("error:", err)
		}
	}))

	defer ts.Close()

	// write temp file for upload
	tmpdir := os.TempDir()
	fname := path.Join(tmpdir, "testupload.txt")
	err := ioutil.WriteFile(fname, []byte("hello"), 0600)
	if err != nil {
		t.Fatal("error:", err)
	}

	fncheck := func() error {
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}

		if string(b) != "hello" {
			return fmt.Errorf("should be hello")
		}

		log.Println(string(b))
		return nil
	}

	err = fncheck()
	if err != nil {
		t.Fatal("error:", err)
	}

	// setup upload file
	buf := &bytes.Buffer{}
	bw := multipart.NewWriter(buf)
	fw, err := bw.CreateFormFile("uploadfile", fname)
	if err != nil {
		t.Fatal("error:", err)
	}

	fh, err := os.Open(fname)
	if err != nil {
		t.Fatal("error:", err)
	}

	_, err = io.Copy(fw, fh)
	if err != nil {
		t.Fatal("error:", err)
	}

	// also set the upload path
	loc, err := bw.CreateFormField("uploadpath")
	if err != nil {
		t.Fatal("error:", err)
	}

	// use tmp folder as upload location
	_, err = loc.Write([]byte(""))
	if err != nil {
		t.Fatal("error:", err)
	}

	ctype := bw.FormDataContentType()
	bw.Close()
	_, err = http.Post(ts.URL, ctype, buf)
	if err != nil {
		t.Fatal("error:", err)
	}

	err = fncheck()
	if err != nil {
		t.Fatal("error:", err)
	}
}
