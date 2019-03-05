package sandbox

import (
	"errors"
)

type Lang struct {
	Name     string
	Compiler string
	FileName string
}

var Langs = make(map[string]Lang)

func init() {
	Langs["nodejs"] = Lang{
		Name:     "NodeJS",
		Compiler: "node ",
		FileName: "index.js",
	}

	Langs["python"] = Lang{
		Name:     "Python",
		Compiler: "python",
		FileName: "main.py",
	}

	Langs["ruby"] = Lang{
		Name:     "Ruby",
		Compiler: "ruby",
		FileName: "main.rb",
	}
}

func GetLang(name string) (Lang, error) {
	if srcLang, ok := Langs[name]; ok {
		return srcLang, nil
	} else {
		return srcLang, errors.New("Can't find language " + name)
	}
}
