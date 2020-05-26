package worker

import (
	"context"
	"github.com/Kscorpion/common"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExcuteInfo) {
	go func() {
		var (
			cmd    *exec.Cmd
			err    error
			output []byte
			result *common.JobExecuteResult
		)
		//任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			OutPut:      make([]byte, 0),
		}
		//记录任务开始时间
		result.StartTime = time.Now()

		//执行shell
		//linux
		//cmd = exec.CommandContext(context.TODO(),"/bin/bash","-c",info.Job.Command)
		//window
		cmd = exec.CommandContext(context.TODO(), "E:\\Git\\bin\\bash.exe", "-c", info.Job.Command)
		//执行并捕获输出
		output, err = cmd.CombinedOutput()

		//记录任务结束时间
		result.EndTime = time.Now()
		result.OutPut = output
		result.Err = err
		//任务执行完成，把结果返回给scheduler,scheduler从excutingTable中删除掉执行的记录
		G_scheduler.PushJobResult(result)

	}()
}

//初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
