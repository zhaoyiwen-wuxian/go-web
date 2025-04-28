package utils

import (
	"fmt"

	"github.com/panjf2000/ants/v2"
)

var Pool *ants.Pool

func InitWorkerPool() {
	var err error
	Pool, err = ants.NewPool(Conf.Pool.Size)
	if err != nil {
		panic("初始化协程池失败: " + err.Error())
	}
}

func SubmitTask(task func()) {
	err := Pool.Submit(task)
	if err != nil {
		fmt.Println("任务提交失败:", err)
	}
}
