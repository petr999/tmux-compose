package cmd_name_args

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"tmux_compose/types"
)

//go:embed templates/bash-new-window.gson
var cnaDefaultTemplateContents []byte

func Construct(osStruct types.CnaOsInterface, config types.ConfigInterface) *CmdNameArgs {
	cna := &CmdNameArgs{}
	cna.New(osStruct, config)
	return cna
}

type CmdNameArgs struct {
	config   types.ConfigInterface
	osStruct types.CnaOsInterface
}

func (cna *CmdNameArgs) New(osStruct types.CnaOsInterface, config types.ConfigInterface) {
	cna.osStruct = osStruct
	cna.config = config
}
func (cna *CmdNameArgs) Get(dcYmlValue types.DcYmlValue) (cnaValue types.CmdNameArgsValueType, err error) {
	var tmplVars tmplVarsType
	tmplVars, err = cna.getTmplVars(dcYmlValue)
	if err != nil {
		return
	}

	var tmplStr string
	tmplStr, err = cna.getTmplStr()
	if err != nil {
		return
	}

	var tmplJsonNew string
	tmplJsonNew, err = cna.tmplExecute(tmplStr, &tmplVars)
	if err != nil {
		err = fmt.Errorf(`error executing template '%s' on with vars from: '%s': %w`, tmplStr, tmplVars, err)
		return
	}

	var cnaValues []types.CmdNameArgsValueType
	cnaValues, err = cna.tmplUnserialize(tmplJsonNew)
	if err != nil {
		err = fmt.Errorf(`error unserializing template '%s': %w`, tmplJsonNew, err)
		return
	}
	if len(cnaValues) != 1 {
		err = fmt.Errorf(`error unserializing template '%s': amount of commands '%v' is not '1' from '%v'`, tmplJsonNew, len(cnaValues), cnaValues)
		return
	}

	cnaValue = cnaValues[0]
	cna.getShebangToCnaValues(&cnaValue)
	newArgs := make([]string, len(cnaValue.Args))
	for i, arg := range cnaValue.Args {
		newArg := strings.ReplaceAll(arg, `{{.Shebang}}`, cnaValue.Shebang)
		// newArg = strings.ReplaceAll(newArg, `{{$shebang}}`, tmplVars.Shebang)
		newArgs[i] = newArg
	}
	cnaValue.Args = newArgs
	cnaValue.Workdir = tmplVars.Workdir
	if len(cnaValue.Shebang) == 0 {
		cnaValue.Shebang = tmplVars.Shebang
	}
	return cnaValue, nil
}

func (cna *CmdNameArgs) getShebangToCnaValues(cnaValue *types.CmdNameArgsValueType) {
	shebangFromOs := cna.config.GetShell()
	if len(shebangFromOs) > 0 {
		cnaValue.Shebang = shebangFromOs
	} else {
		if len(cnaValue.Shebang) == 0 {
			cnaValue.Shebang = `bash`
		}
	}
}

func (cna *CmdNameArgs) tmplUnserialize(tmplJson string) (cnaValues []types.CmdNameArgsValueType, err error) {
	err = json.Unmarshal([]byte(tmplJson), &cnaValues)
	return
}

type tmplVarsType struct {
	Workdir         string
	DcServicesNames []string
	Basedir         string
	Shebang         string
}

func (cna *CmdNameArgs) getTmplVars(dcYmlValue types.DcYmlValue) (tmplVars tmplVarsType, err error) {
	tmplVars.Workdir = dcYmlValue.Workdir
	tmplVars.DcServicesNames = dcYmlValue.DcServicesNames
	baseDir := filepath.Base(dcYmlValue.Workdir)

	if len(baseDir) == 0 { // XXX impossible?
		return tmplVars, fmt.Errorf("error finding base dir name '%v' for work dir: '%v'", baseDir, dcYmlValue.Workdir)
	}
	if len(baseDir) >= len(dcYmlValue.Workdir) {
		return tmplVars, fmt.Errorf("error finding base dir name '%v' same length for work dir: '%v'", baseDir, dcYmlValue.Workdir)
	}

	tmplVars.Basedir = baseDir
	return
}

func (cna *CmdNameArgs) getTmplStr() (tmplStr string, err error) {
	fNameByConf := cna.config.GetCnaTemplateFname()
	var fName string
	if len(fNameByConf) == 0 {
		fName = `.`
	} else {
		fName = fNameByConf
	}
	fName, err = filepath.Abs(fName)
	if err != nil {
		return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
	}

	var fileInfo types.FileInfoStruct
	fileInfo, err = cna.osStruct.Stat(fName)
	if err != nil {
		return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
	}
	if fileInfo.IsDir() {
		fName = filepath.Join(fName, `tmux-compose-template.gson`)
	}
	fileInfo, err = cna.osStruct.Stat(fName)
	if err != nil {
		if len(fNameByConf) > 0 {
			return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
		} else {
			return string(cnaDefaultTemplateContents), nil
		}
	}
	if !fileInfo.IsFile() {
		if len(fNameByConf) > 0 {
			return tmplStr, fmt.Errorf(`error reading file '%v': not a file`, fName)
		} else {
			return string(cnaDefaultTemplateContents), nil
		}
	}

	var buf []byte
	buf, err = cna.osStruct.ReadFile(fName)
	if err == nil {
		tmplStr = string(buf)
	} else {
		if len(fNameByConf) > 0 {
			return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
		}
		tmplStr = string(cnaDefaultTemplateContents)
	}
	return tmplStr, nil
}

func (cna *CmdNameArgs) tmplExecute(tmplJson string, tmplVars *tmplVarsType) (tmplJsonNew string, err error) {
	tmplObj := template.New(`tmux_compose`).Funcs(template.FuncMap{
		"IsNotLast": func(i int, length int) bool {
			return i < length-1
		}})
	nameTmplObj, err := tmplObj.Parse(tmplJson)
	if err != nil {
		return ``, fmt.Errorf("error reading name template: '%v': %w", tmplJson, err)
	}

	var nameBuf bytes.Buffer
	tmplVars.Shebang = `{{.Shebang}}`
	err = nameTmplObj.Execute(&nameBuf, tmplVars)
	if err != nil {
		return ``, fmt.Errorf("error executing name template: '%v' on '%v': %w", tmplJson, tmplVars, err)
	}

	tmplJsonNew = nameBuf.String()

	return
}
