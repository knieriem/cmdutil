package confutil

import (
	"os"

	"github.com/knieriem/cmdutil/pathutil"
)

type BaseInf struct {
	Exe pathutil.FilePath
}

func NewBaseInf() (*BaseInf, error) {
	b := new(BaseInf)
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	b.Exe = pathutil.FilePath(exe)
	return b, nil
}
