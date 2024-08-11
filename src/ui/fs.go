package ui

// 将 embed fs 解压到临时目录下

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type MyReader struct {
	src []byte
	pos int
}

func (r *MyReader) Read(dst []byte) (n int, err error) {
	n = copy(dst, r.src[r.pos:])
	r.pos += n
	if r.pos == len(r.src) {
		return n, io.EOF
	}
	return
}

func NewMyReader(b []byte) *MyReader { return &MyReader{b, 0} }

func expandEmbed(eFS embed.FS) (string, error) {
	// expand embedded dir into temp fs
	dir, err := os.MkdirTemp("", "DK4HDTrainer")
	if err != nil {
		return "", err
	}

	fmt.Println("expanding to temp dir:", dir)

	err = fs.WalkDir(eFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileName := filepath.Join(dir, path)
		if d.IsDir() {
			fmt.Println("dir", fileName)
			os.MkdirAll(fileName, os.ModePerm)
		} else {
			fmt.Println("file", fileName)
			destination, err := os.Create(fileName)
			if err != nil {
				return err
			}
			defer destination.Close()
			file, err := eFS.ReadFile(path)
			nBytes, err := io.Copy(destination, NewMyReader(file))
			_ = nBytes
			return err
		}

		return nil
	})

	return dir, err
}