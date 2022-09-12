package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ZipString(filename string, text []byte) {
	dirPath := filepath.Dir(filename)
	if _, e := os.Stat(dirPath); e != nil {
		if e := os.MkdirAll(dirPath, 0777); e != nil {
			fmt.Println(e)
		}
	}
	var content bytes.Buffer
	b := []byte(text)
	w := zlib.NewWriter(&content)
	w.Write(b)
	w.Close()
	if e := ioutil.WriteFile(filename, content.Bytes(), 0777); e != nil {
		fmt.Println(e)
	}
}

func UnzipString(filename string) []byte {
	if _, e := os.Stat(filename); e != nil {
		return []byte{}
	}
	s, _ := ioutil.ReadFile(filename)
	var out bytes.Buffer
	r, e := zlib.NewReader(bytes.NewBuffer(s))
	if e != nil {
		return []byte{}
	}
	if _, e := io.Copy(&out, r); e != nil {
		return []byte{}
	}
	return out.Bytes()
}
