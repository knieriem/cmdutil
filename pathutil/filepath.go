package pathutil

import (
	"os"
	"path/filepath"
)

type FilePath string

func (f FilePath) Dir() string {
	return filepath.Dir(string(f))
}

func (f FilePath) Base() string {
	return filepath.Base(string(f))
}

func (f FilePath) Ext() string {
	return filepath.Ext(string(f))
}

func (f FilePath) Abs() string {
	a, err := filepath.Abs(string(f))
	if err != nil {
		return string(f)
	}
	return a
}

func (f FilePath) RelIfShorter() string {
	return filepathRelIfShorter(string(f))
}

func filepathRelIfShorter(fileName string) string {
	p, err := filepath.Abs(fileName)
	if err != nil {
		p = fileName
	}

	wd, err := os.Getwd()
	if err == nil {
		rel, err := filepath.Rel(wd, p)
		if err == nil && len(rel) < len(fileName) {
			fileName = rel
		}
	}
	return fileName
}

func (f FilePath) Stem() string {
	b := f.Base()
	return b[:len(b)-len(f.Ext())]
}

func (f FilePath) String() string {
	return string(f)
}

func (f FilePath) RelTo(basepath string) string {
	rel, err := filepath.Rel(basepath, f.String())
	if err != nil {
		return string(f)
	}
	return rel
}
