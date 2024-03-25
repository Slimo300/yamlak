package cmd

import (
	"io"
	"os"
)

var inPlaceFlag bool

func outputResult(file *os.File, r io.Reader, inPlaceFlag bool) error {

	if inPlaceFlag {
		if _, err := file.Seek(0, 0); err != nil {
			return err
		}
		size, err := io.Copy(file, r)
		if err != nil {
			return err
		}
		if err := file.Truncate(size); err != nil {
			return err
		}
	} else {
		if _, err := io.Copy(os.Stdout, r); err != nil {
			return err
		}
	}

	return nil

}
