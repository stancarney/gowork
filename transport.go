package gowork

import (
	"bytes"
	"os"
)

type Transport interface {
	Connect(args ...interface{}) error
	Upload(path string, content *bytes.Buffer) error
	Download(path string, content *bytes.Buffer) error
	List(path string) ([]os.FileInfo, error)
	Raw(args ...interface{}) ([]interface{}, error)
	Close() error
}
