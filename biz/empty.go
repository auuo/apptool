package biz

import "path/filepath"

type emptyParser struct {
	dir string
}

func newEmptyParser(dir string) *emptyParser {
	return &emptyParser{
		dir: dir,
	}
}

func (p *emptyParser) genRouterCode() (string, error) {
	panic("implement me")
}

func (p *emptyParser) genModelFiles() []genTmpl {
	return []genTmpl{
		{tmplName: "baseModelTmpl", file: filepath.Join(p.dir, "model/model.go"), param: nil},
		{tmplName: "userModelTmpl", file: filepath.Join(p.dir, "model/user.go"), param: nil},
	}
}
