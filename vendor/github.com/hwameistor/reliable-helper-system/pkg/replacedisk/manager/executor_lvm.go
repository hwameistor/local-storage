package manager

import (
	"github.com/hwameistor/reliable-helper-system/pkg/exechelper"
	"github.com/hwameistor/reliable-helper-system/pkg/exechelper/nsexecutor"
)

type lvmExecutor struct {
	cmdExec exechelper.Executor
}

var lvmExecutorInstance *lvmExecutor

func NewLVMExecutor() *lvmExecutor {
	if lvmExecutorInstance == nil {
		lvmExecutorInstance = &lvmExecutor{
			cmdExec: nsexecutor.New(),
		}
	}
	return lvmExecutorInstance
}

func (lvm *lvmExecutor) vgreduce(vgName, diskpath string, options []string) error {
	params := exechelper.ExecParams{
		CmdName: "vgreduce",
		CmdArgs: []string{vgName, diskpath},
	}
	params.CmdArgs = append(params.CmdArgs, options...)
	res := lvm.cmdExec.RunCommand(params)
	if res.ExitCode == 0 {
		return nil
	}
	return res.Error
}

func (lvm *lvmExecutor) pvremove(vgName, diskpath string, options []string) error {
	params := exechelper.ExecParams{
		CmdName: "pvmove",
		//CmdArgs: []string{diskpath, "--force", "--force"},
		CmdArgs: []string{diskpath},
	}
	params.CmdArgs = append(params.CmdArgs, options...)
	res := lvm.cmdExec.RunCommand(params)
	if res.ExitCode == 0 {
		return nil
	}
	return res.Error
}
