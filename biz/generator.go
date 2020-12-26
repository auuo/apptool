package biz

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
)

var tmpl *template.Template

type project struct {
	idlParser idlParser
	idlPath   string
	dir       string
	modName   string
	isNew     bool
}

type genTmpl struct {
	tmplName string
	file     string
	param    interface{}
}

type idlParser interface {
	genRouterCode() (string, error)
	genModelFiles() []genTmpl
}

func init() {
	// load template
	tmpl = template.New("")
	template.Must(tmpl.New("confTmpl").Parse(confTmpl))
	template.Must(tmpl.New("yamlTmpl").Parse(yamlTmpl))
	template.Must(tmpl.New("baseModelTmpl").Parse(baseModelTmpl))
	template.Must(tmpl.New("userModelTmpl").Parse(userModelTmpl))
	template.Must(tmpl.New("mainTmpl").Parse(mainTmpl))
	template.Must(tmpl.New("pingTmpl").Parse(pingTmpl))
	template.Must(tmpl.New("appResponseTmpl").Parse(appResponseTmpl))
	template.Must(tmpl.New("appWrapperTmpl").Parse(appWrapperTmpl))
	template.Must(tmpl.New("errTmpl").Parse(errTmpl))
	template.Must(tmpl.New("errCommonTmpl").Parse(errCommonTmpl))
	template.Must(tmpl.New("genLogIdTmpl").Parse(genLogIdTmpl))
	template.Must(tmpl.New("logTmpl").Parse(logTmpl))
	template.Must(tmpl.New("logIdTmpl").Parse(logIdTmpl))
	template.Must(tmpl.New("modTmpl").Parse(modTmpl))
	template.Must(tmpl.New("daoTmpl").Parse(daoTmpl))
	template.Must(tmpl.New("gitIgnoreTmpl").Parse(gitIgnoreTmpl))
}

func Generator(cmd, dir, modName, idlPath string) error {
	isNew := cmd == "new"
	if isNew {
		if err := makeProjectDirs(dir); err != nil {
			return err
		}
	}
	p := project{
		dir:     dir,
		modName: modName,
		isNew:   isNew,
	}
	if idlPath == "" {
		p.idlParser = newEmptyParser(dir)
	} else {
		abs, err := filepath.Abs(idlPath)
		if err != nil {
			return err
		}
		p.idlPath = abs
	}
	return p.generator()
}

func (p *project) generator() error {
	modParam := map[string]string{
		"modName": p.modName,
	}
	var files []genTmpl
	if p.isNew {
		files = append(files,
			genTmpl{tmplName: "gitIgnoreTmpl", file: filepath.Join(p.dir, ".gitignore")},
			genTmpl{tmplName: "modTmpl", file: filepath.Join(p.dir, "go.mod"), param: modParam},
			genTmpl{tmplName: "confTmpl", file: filepath.Join(p.dir, "conf/conf.go")},
			genTmpl{tmplName: "yamlTmpl", file: filepath.Join(p.dir, "conf/conf.yaml")},
			genTmpl{tmplName: "mainTmpl", file: filepath.Join(p.dir, "main.go"), param: modParam},
			genTmpl{tmplName: "pingTmpl", file: filepath.Join(p.dir, "handler/ping.go")},
			genTmpl{tmplName: "appWrapperTmpl", file: filepath.Join(p.dir, "pkg/app/wrapper.go")},
			genTmpl{tmplName: "appResponseTmpl", file: filepath.Join(p.dir, "pkg/app/response.go"), param: modParam},
			genTmpl{tmplName: "errTmpl", file: filepath.Join(p.dir, "pkg/e/err.go")},
			genTmpl{tmplName: "errCommonTmpl", file: filepath.Join(p.dir, "pkg/e/common.go")},
			genTmpl{tmplName: "logTmpl", file: filepath.Join(p.dir, "pkg/logs/log.go")},
			genTmpl{tmplName: "logIdTmpl", file: filepath.Join(p.dir, "pkg/logs/logid.go")},
			genTmpl{tmplName: "genLogIdTmpl", file: filepath.Join(p.dir, "middleware/log.go"), param: modParam},
			genTmpl{tmplName: "daoTmpl", file: filepath.Join(p.dir, "dao/dao.go"), param: modParam},
		)
	}

	files = append(files, p.idlParser.genModelFiles()...)

	var buf bytes.Buffer
	for _, f := range files {
		if err := tmpl.ExecuteTemplate(&buf, f.tmplName, f.param); err != nil {
			return err
		}
		if err := ioutil.WriteFile(f.file, buf.Bytes(), 0644); err != nil {
			return err
		}
		buf.Reset()
	}
	return nil
}

func makeProjectDirs(dir string) error {
	return MakeDirs(
		filepath.Join(dir, "conf"),
		filepath.Join(dir, "dao"),
		filepath.Join(dir, "handler"),
		filepath.Join(dir, "middleware"),
		filepath.Join(dir, "model"),
		filepath.Join(dir, "pkg/app"),
		filepath.Join(dir, "pkg/e"),
		filepath.Join(dir, "pkg/logs"),
		filepath.Join(dir, "service"),
	)
}
