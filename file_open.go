package xpad

import "os"

func openReadWriteOrReadOnly(path string) (*os.File, bool, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err == nil {
		return file, false, nil
	}
	file, err = os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, false, err
	}
	return file, true, nil
}
