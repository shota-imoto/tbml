package args

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Args struct {
	YamlPath string
}

func Load() (Args, error) {
	args := os.Args
	if len(args) < 1 {
		return Args{}, fmt.Errorf("設定ファイルのパスを指定してください: %v", args)
	}

	return Args{YamlPath: args[1]}, nil

}

type YAMLExtentions []string

func NewYamlExtentions() YAMLExtentions {
	return YAMLExtentions{"yaml", "yml"}
}

func (exts YAMLExtentions) isYaml(ext string) bool {
	for _, e := range exts {
		if e == ext {
			return true
		}
	}
	return false
}

func (a Args) OutPath() string {
	filename := filepath.Base(a.YamlPath)
	split := strings.Split(filename, ".")
	ext := split[len(split)-1]

	exts := NewYamlExtentions()

	var fn_slice []string
	if exts.isYaml(ext) {
		fn_slice = split[0 : len(split)-1]
	} else {
		fn_slice = split
	}

	fn_slice = append(fn_slice, "svg")
	fn := strings.Join(fn_slice, ".")

	return filepath.Join(filepath.Dir(a.YamlPath), fn)
}
