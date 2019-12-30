package cmdutil

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/knieriem/cmdutil/confutil"
	"github.com/knieriem/g/ioutil/terminal"
	"github.com/knieriem/text/rc"
)

var NoExec bool

type CmdLine string

type RunOption func(c *exec.Cmd)

func (cl CmdLine) Start(tplData interface{}, options ...RunOption) (*exec.Cmd, error) {
	var err error
	args := rc.Tokenize(string(cl))
	for i := range args {
		if strings.Index(args[i], "{{") == -1 {
			continue
		}
		tplName := `post-process: "` + args[0] + `"[` + strconv.Itoa(i) + `]`
		args[i], err = confutil.StringTemplate(args[i]).Execute(tplName, tplData)
		if err != nil {
			return nil, err
		}
	}
	if NoExec {
		return nil, nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	for _, o := range options {
		o(cmd)
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (cl CmdLine) Run(tplData interface{}, options ...RunOption) error {
	cmd, err := cl.Start(tplData, options...)
	if err != nil {
		return err
	}
	if cmd != nil {
		return cmd.Wait()
	}
	return nil
}

func TerminalInput() (RunOption, io.Closer, error) {
	t, err := terminal.OpenInput()
	if err != nil {
		return nil, nil, err
	}
	a := func(cmd *exec.Cmd) {
		cmd.Stdin = t
	}
	return a, t, nil
}

func StderrOutput(cmd *exec.Cmd) {
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
}
