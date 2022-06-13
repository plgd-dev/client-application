package main

import (
	"fmt"
	"io"
	"os"

	"github.com/plgd-dev/client-application/pkg/tar"
	"github.com/plgd-dev/hub/v2/pkg/log"
)

var UISeparator = "-----------------------------------------------------"

func setOffsetToTar(f *os.File) (bool, error) {
	pos := 0
	uiSeparator := "\n" + UISeparator + "\n"
	buffer := make([]byte, 64*1024)
	for {
		buffer = buffer[:cap(buffer)]
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		data := buffer[:n]
		for idx, b := range data {
			if pos+1 == len(uiSeparator) {
				_, err := f.Seek(-int64(len(data)-idx-1), io.SeekCurrent)
				if err != nil {
					return false, err
				}
				return true, nil
			}
			if uiSeparator[pos] == b {
				pos++
			} else {
				pos = 0
			}
		}
	}
	return false, nil
}

func extractUI(directory string) (errRet error) {
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}
	f, err := os.OpenFile(ex, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("cannot open executable: %w", err)
	}
	defer func() {
		err1 := f.Close()
		if errRet == nil {
			errRet = err1
		}
	}()

	ok, err := setOffsetToTar(f)
	if err != nil {
		return fmt.Errorf("cannot find tar offset: %w", err)
	}
	if !ok {
		log.Warn("cannot find tar offset for extract UI")
		return nil
	}
	if err = tar.Untar(directory, f); err != nil {
		return fmt.Errorf("cannot untar files to directory %v: %w", directory, err)
	}
	return nil
}
