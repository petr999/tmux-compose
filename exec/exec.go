package exec

import (
	osExec "os/exec"
	"regexp"
	"strings"
	"tmux_compose/types"
)

func Construct(osStruct types.ExecOsInterface, config types.ConfigInterface) *Exec {
	exec := &Exec{}
	stdHandles := osStruct.GetStdHandles()
	exec.New(osStruct, stdHandles, config)
	return exec
}

type dryRunCmd struct{}

func (obj *dryRunCmd) Run() error { return nil }

type Exec struct {
	Cmd        *types.CmdType
	Selector   any
	osStruct   types.ExecOsInterface
	stdHandles types.StdHandlesType
	config     types.ConfigInterface
}

func (exec *Exec) New(osStruct types.ExecOsInterface, stdHandles types.StdHandlesType, config types.ConfigInterface) {
	exec.Selector, exec.osStruct, exec.stdHandles, exec.config = false, osStruct, stdHandles, config
}

func (exec *Exec) execCommand(cna types.CmdNameArgsValueType) *osExec.Cmd {
	obj := osExec.Command(cna.Cmd, cna.Args...)
	obj.Stdout = exec.stdHandles.Stdout
	obj.Stderr = exec.stdHandles.Stderr
	obj.Stdin = exec.stdHandles.Stdin

	return obj
}

func (exec *Exec) dryRun(cna types.CmdNameArgsValueType) types.CmdInterface {
	return &dryRunCmd{}
}

func cmdEscape(str string) string {
	str = string(regexp.MustCompile(`^;$`).ReplaceAll([]byte(str), []byte(`\;`)))
	if strings.ContainsAny(str, " \t\n'") {
		return `'` + strings.ReplaceAll(str, `'`, `'\''`) + `'`
	}
	return str
}

func (exec *Exec) GetSelector() (selector any) {
	if len(exec.config.GetDryRun()) > 0 { // dry run
		return true
	}
	return selector
}

func (exec *Exec) GetCommand(cna types.CmdNameArgsValueType) *types.CmdType {
	var obj types.CmdInterface
	selector := exec.GetSelector()
	var runFunc func() error
	if dryRun, ok := selector.(bool); ok {
		if dryRun {
			obj = exec.dryRun(cna)
			runFunc = func() error {
				slice := []string{"#!/usr/bin/env bash\n\ncd " + cmdEscape(cna.Workdir) + "\n\n" + cmdEscape(cna.Cmd)}
				args := make([]string, len(cna.Args))
				for i, arg := range cna.Args {
					args[i] = cmdEscape(arg)
				}
				slice = append(slice, args...)

				if _, err := exec.osStruct.GetStdHandles().Stdout.Write([]byte(strings.Join(slice, " ") + "\n")); err != nil {
					return err
				}
				return nil
			}
		} else { // obj.Run()
			obj = exec.execCommand(cna)
			runFunc = func() error {
				if err := exec.osStruct.Chdir(cna.Workdir); err != nil {
					return err
				}
				return obj.Run()
			}
		}
	} else {
		if objWithRun, ok := selector.(interface{ Run() error }); ok {
			obj = objWithRun
		}
		runFunc = func() error { return obj.Run() }
	}
	exec.Cmd = &types.CmdType{Obj: obj}
	runFuncNew := func() error {
		exec.Cmd.RunData.WasCalledTimes++
		exec.Cmd.RunData.Cnas = append(exec.Cmd.RunData.Cnas, cna)
		return runFunc()
	}
	exec.Cmd.Run = runFuncNew

	exec.Cmd.Stdhandles = exec.stdHandles
	return exec.Cmd
}
