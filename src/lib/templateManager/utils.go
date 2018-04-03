package templateManager

import (
	"html/template"
	"path/filepath"

	"github.com/pchan37/studybuddy/src/utils"
)

func getLayoutFiles() (layoutFiles []string, err error) {
	layoutFiles, err = filepath.Glob(config.layoutPath + "*.tmpl")
	utils.LogFatalOnError(err)
	return
}

func getIncludeFiles() (includeFiles []string, err error) {
	includeFiles, err = filepath.Glob(config.includePath + "*.tmpl")
	utils.LogFatalOnError(err)
	return
}

func getMainTemplate() (mainTemplate *template.Template, err error) {
	mainTemplate, err = template.New("main").Parse(mainTmpl)
	utils.LogFatalOnError(err)
	return
}
