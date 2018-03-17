package sysio

import (
	"io"
	"io/ioutil"
	"os"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform"

	"github.com/xi2/xz"
	"strings"
	"bytes"
)

type FileReaderFunc func(path string) (io.ReadCloser, error)

var NewFileReader FileReaderFunc = func(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func ReadFile(path string) ([]byte, error) {
	reader, err := NewFileReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if (strings.HasSuffix(path, ".xz")) {
		dataBytes, err := buf.ReadAllToBytes(reader)
		if err != nil {
			return nil, err
		}
		r, err := xz.NewReader(bytes.NewReader(dataBytes), 0)
		if err != nil {
			return nil, err
		}
		return buf.ReadAllToBytes(r)
	}
	return buf.ReadAllToBytes(reader)
}

func ReadAsset(file string) ([]byte, error) {
	return ReadFile(platform.GetAssetLocation(file))
}

func CopyFile(dst string, src string) error {
	bytes, err := ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, bytes, 0644)
}
