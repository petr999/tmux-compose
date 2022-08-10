package cmd_name_args

import (
	_ "embed"
	"fmt"
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
	// dcvBasedir
	_, err = cna.getDcvBasedir(dcYmlValue)
	if err != nil {
		return
	}

	// tmplStr :=
	_, err = cna.getTmplStr()
	if err != nil {
		return
	}

	return types.CmdNameArgsValueType{Workdir: dcYmlValue.Workdir}, nil
}

// type tmplType []struct {
// 	Cmd  string
// 	Args []string
// }

type dcvBasedirType struct {
	Workdir         string
	DcServicesNames []string
	Basedir         string
}

func (cna *CmdNameArgs) getDcvBasedir(dcYmlValue types.DcYmlValue) (dcvBasedir dcvBasedirType, err error) {
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
		fName = filepath.Join(fName, `tmux-compose-template.json`)
	}
	fileInfo, err = cna.osStruct.Stat(fName)
	if err != nil {
		return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
	}
	if !fileInfo.IsFile() {
		return tmplStr, fmt.Errorf(`error reading file '%v': not a file`, fName)
	}

	var buf []byte
	buf, err = cna.osStruct.ReadFile(fName)
	if err == nil {
		tmplStr = string(buf)
	} else {
		if fName == fNameByConf {
			return tmplStr, fmt.Errorf(`error reading file '%v': %w`, fName, err)
		}
		tmplStr = string(cnaDefaultTemplateContents)
	}
	return tmplStr, nil
}

// //go:embed templates/bash-new-window.gson
// var tmplJson []byte

// var nilValue = types.CmdNameArgsType{Workdir: ``, Name: ``, Args: nil}

// func getTmplJson() []byte {
// 	return tmplJson
// }

// func tmplExecute(tmplJson string, dcvBasedir dcvBasedirType) (tmplJsonNew string, err error) {
// 	tmplObj := template.New(`tmux_compose`).Funcs(template.FuncMap{
// 		"IsNotLast": func(i int, length int) bool {
// 			return i < length-1
// 		}})
// 	nameTmplObj, err := tmplObj.Parse(tmplJson)
// 	if err != nil {
// 		return ``, fmt.Errorf("error reading name template: '%v'\n\t%w", tmplJson, err)
// 	}
// 	var nameBuf bytes.Buffer
// 	err = nameTmplObj.Execute(&nameBuf, dcvBasedir)

// 	if err != nil {
// 		return ``, fmt.Errorf("error executing name template: '%v' on '%v':\n\t%w", tmplJson, dcvBasedir, err)
// 	}

// 	tmplJsonNew = nameBuf.String()

// 	return
// }

// func CmdNameArgs(dcConfigReader dc_config.ReaderInterface, osStruct OsStructCmdNameArgs) (types.CmdNameArgsType, error) {
// 	dcConfig, err := dcConfigReader.Read()
// 	if err != nil {
// 		return nilValue, fmt.Errorf("error reading config:\n\t%w", err)
// 	}

// 	dcvBasedir := dcvBasedirType{dcConfig.Workdir, dcConfig.DcServicesNames, baseDir}

// 	tmplJsonStr := string(getTmplJson())

// 	tmplJsonNew, err := tmplExecute(tmplJsonStr, dcvBasedir)
// 	if err != nil {
// 		return nilValue, fmt.Errorf("error applying template '%s':\n\t%w", tmplJsonStr, err)
// 	}

// 	var tmpl tmplType
// 	err = json.Unmarshal([]byte(tmplJsonNew), &tmpl)
// 	if err != nil {
// 		return nilValue, fmt.Errorf("error parsing json '%v':\n\t%w", `templates/bash-new-window.json`, err)
// 	}

// 	cmdName, args := tmpl[0].Cmd, tmpl[0].Args

// 	// ID=0; try_next=1; trap 'echo "trap pid: ${PID}"; kill -INT $PID; try_next="";' SIGINT; while [ 'x1' == "x${try_next}" ]; do sleep 1 & PID=$!; echo "pid: ${PID}" ; wait $PID; done
// 	return types.CmdNameArgsType{Workdir: dcConfig.Workdir, Name: cmdName, Args: args}, nil
// }
