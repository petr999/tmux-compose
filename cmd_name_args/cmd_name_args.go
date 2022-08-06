package cmd_name_args

import (
	_ "embed"
	"tmux_compose/types"
)

func Construct(os_struct types.CnaOsInterface) CmdNameArgs {
	cna := CmdNameArgs{}
	cna.New(os_struct)
	return cna
}

type CmdNameArgs struct {
}

func (cna CmdNameArgs) New(os_struct types.CnaOsInterface) {}
func (cna CmdNameArgs) Get(types.DcYmlValue) (types.CmdNameArgsValueType, error) {
	return types.CmdNameArgsValueType{}, nil
}

// type tmplType []struct {
// 	Cmd  string
// 	Args []string
// }

// type dcvBasedirType struct {
// 	Workdir         string
// 	DcServicesNames []string
// 	Basedir         string
// }

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

// 	baseDir := filepath.Base(dcConfig.Workdir)
// 	// XXX impossible?
// 	// if len(baseDir) == 0 {
// 	// 	return nilValue, fmt.Errorf("error finding base dir name '%v' for work dir: '%v'", baseDir, dcConfig.Workdir)
// 	// }
// 	if len(baseDir) >= len(dcConfig.Workdir) {
// 		return nilValue, fmt.Errorf("error finding base dir name '%v' same length for work dir: '%v'", baseDir, dcConfig.Workdir)
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
