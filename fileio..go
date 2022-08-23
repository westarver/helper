package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	apath "github.com/rhysd/abspath"
)

//─────────────┤ GetLogFile ├─────────────

func GetLogFile(file string) *os.File {
	if file == "std" {
		return os.Stderr
	}
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if f == nil {
		fmt.Println("Failure to open Log file " + file + ". Will default to Stderr.")
		f = os.Stderr
	}
	return f
}

//─────────────┤ OpenFile ├─────────────

func OpenFileRead(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.Open(path)
	}
	return nil, err
}

//─────────────┤ OpenAppend ├─────────────

func OpenAppend(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	}
	return nil, err
}

//─────────────┤ OpenTrunc ├─────────────

func OpenTrunc(file string) (*os.File, error) {
	path, err := GetPath(file)

	if err == nil {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	}
	return nil, err
}

//─────────────┤ CopyFileStr ├─────────────

func CopyFileStr(dst, src string) error {
	if DirExists(dst) {
		_, base := filepath.Split(src)
		dst = filepath.Join(dst, base)
	}
	d, err := OpenTrunc(dst)
	if err != nil {
		return err
	}
	s, err := OpenFileRead(src)
	if err != nil {
		return err
	}
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}

//─────────────┤ GetPath  ├─────────────

func GetPath(name string) (string, error) {
	a, err := apath.ExpandFrom(name)
	return a.String(), err
}

//─────────────┤ DirExists ├─────────────

func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

//─────────────┤ ValidatePath ├─────────────

func ValidatePath(path string) (string, error) {
	a, err := apath.ExpandFrom(path)
	if err != nil {
		return "", err
	}
	path = a.String()
	d, f := filepath.Split(path)
	err = os.MkdirAll(d, 0777)
	if err != nil {
		return path, err
	}
	return d + f, nil
}
