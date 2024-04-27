package core

import (
	"fmt"
	"log"
	"os"
)

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
}

func NewLogger(name string) (*log.Logger, error) {
	path := fmt.Sprintf("log/%s.log", name)
	file, err := openFile(path)

	if err != nil {
		return nil, err
	}

	return log.New(file, "", 0), nil
}
