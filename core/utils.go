package core

import (
	"fmt"
	"log"
	"os"
)

func openWriteFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0777)
}

func NewLogger(name string) (*log.Logger, error) {
	file, err := openWriteFile(fmt.Sprintf("log/%s.log", name))

	if err != nil {
		return nil, err
	}

	return log.New(file, "", 0), nil
}
