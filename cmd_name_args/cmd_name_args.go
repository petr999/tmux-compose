package cmd_name_args

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
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
	GetCnaTemplateFname func() string
	osStruct            types.CnaOsInterface
}

func (cna *CmdNameArgs) New(osStruct types.CnaOsInterface, config types.ConfigInterface) {
	cna.osStruct = osStruct
	cna.GetCnaTemplateFname = config.GetCnaTemplateFname
}
func (cna *CmdNameArgs) Get(dcYmlValue types.DcYmlValue) (cnaValue types.CmdNameArgsValueType, err error) {
	var dcvBasedir dcvBasedirType
	dcvBasedir, err = cna.getDcvBasedir(dcYmlValue)
	if err != nil {
		return
	}

	var tmplStr string
	tmplStr, err = cna.getTmplStr()
	if err != nil {
		return
	}

	var tmplJsonNew string
	tmplJsonNew, err = tmplExecute(tmplStr, dcvBasedir)
	if err != nil {
		err = fmt.Errorf(`error executing template '%s' on with vars from: '%s': %w`, tmplStr, dcvBasedir, err)
		return
	}

	var cnaValues []types.CmdNameArgsValueType
	cnaValues, err = tmplUnserialize(tmplJsonNew)
	if err != nil {
		err = fmt.Errorf(`error unserializing template '%s': %w`, tmplJsonNew, err)
		return
	}
	if len(cnaValues) != 1 {
		err = fmt.Errorf(`error unserializing template '%s': amount of commands '%v' is not '1' from '%v'`, tmplJsonNew, len(cnaValues), cnaValues)
		return
	}

	cnaValue = cnaValues[0]
	cnaValue.Workdir = dcYmlValue.Workdir

	return cnaValue, nil
}

func tmplUnserialize(tmplJson string) (cnaValues []types.CmdNameArgsValueType, err error) {
	err = json.Unmarshal([]byte(tmplJson), &cnaValues)
	return
}

type dcvBasedirType struct {
	Workdir         string
	DcServicesNames []string
	Basedir         string
}

func (cna *CmdNameArgs) getDcvBasedir(dcYmlValue types.DcYmlValue) (dcvBasedir dcvBasedirType, err error) {
	dcvBasedir.Workdir = dcYmlValue.Workdir
	dcvBasedir.DcServicesNames = dcYmlValue.DcServicesNames
	baseDir := filepath.Base(dcYmlValue.Workdir)

	if len(baseDir) == 0 { // XXX impossible?
		return dcvBasedir, fmt.Errorf("error finding base dir name '%v' for work dir: '%v'", baseDir, dcYmlValue.Workdir)
	}
	if len(baseDir) >= len(dcYmlValue.Workdir) {
		return dcvBasedir, fmt.Errorf("error finding base dir name '%v' same length for work dir: '%v'", baseDir, dcYmlValue.Workdir)
	}

	dcvBasedir.Basedir = baseDir
	return
}

func (cna *CmdNameArgs) getTmplStr() (tmplStr string, err error) {
	fNameByConf := cna.GetCnaTemplateFname()
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

func tmplExecute(tmplJson string, dcvBasedir dcvBasedirType) (tmplJsonNew string, err error) {
	tmplObj := template.New(`tmux_compose`).Funcs(template.FuncMap{
		"IsNotLast": func(i int, length int) bool {
			return i < length-1
		}})
	nameTmplObj, err := tmplObj.Parse(tmplJson)
	if err != nil {
		return ``, fmt.Errorf("error reading name template: '%v': %w", tmplJson, err)
	}

	var nameBuf bytes.Buffer
	err = nameTmplObj.Execute(&nameBuf, dcvBasedir)
	if err != nil {
		return ``, fmt.Errorf("error executing name template: '%v' on '%v': %w", tmplJson, dcvBasedir, err)
	}

	tmplJsonNew = nameBuf.String()

	return
}
