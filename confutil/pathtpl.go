package confutil

import (
	"bytes"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type StringTemplate string

func (pt StringTemplate) parse(name string) (t *template.Template, err error) {
	t = template.New(name)
	t.Funcs(template.FuncMap{
		"fromSlash": func(path string) string {
			path = strings.Replace(path, `\`, "/", -1)
			return filepath.FromSlash(path)
		},
		"toSlash": func(path string) string {
			return filepath.ToSlash(path)
		},
		"fileurl": func(path string) string {
			psl := filepath.ToSlash(path)
			u := &url.URL{Scheme: "file", Path: psl}
			e := u.EscapedPath()
			if filepath.IsAbs(path) {
				if !strings.HasPrefix(psl, "/") {
					return "file:///" + e
				}
				return "file://" + e
			}
			return e
		},
		"stripExt": func(path string) string {
			return path[:len(path)-len(filepath.Ext(path))]
		},
		"trimRE": func(expr, s string) (string, error) {
			re, err := regexp.Compile(expr)
			if err != nil {
				return "", err
			}
			i := re.FindStringIndex(s)
			if i == nil {
				return s, nil
			}
			return s[:i[0]] + s[i[1]:], nil
		},
	})
	t, err = t.Parse(string(pt))
	return
}

func (pt StringTemplate) Postprocess() (err error) {
	_, err = pt.parse("")
	return
}

func (pt StringTemplate) Execute(name string, data interface{}) (string, error) {
	if pt == "" {
		return "", nil
	}
	t, err := pt.parse(name)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

type PathTemplate string

func (pt PathTemplate) Postprocess() (err error) {
	_, err = StringTemplate(pt).parse("")
	return
}

func (pt PathTemplate) Execute(name string, data interface{}) (string, error) {
	result, err := StringTemplate(pt).Execute(name, data)
	if err != nil {
		return "", err
	}
	result = strings.Replace(result, `\`, "/", -1)
	return filepath.FromSlash(result), nil
}
