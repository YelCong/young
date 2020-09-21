package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"os/exec"
	"time"
)

func main() {
	//tLookPath()
	//tCmd()
	//tCmdCombinedOutput()
	//tCron()
	tRunCronJob()
}

func tLookPath() {
	//str, err := exec.LookPath("./run.sh")
	//str, err := exec.LookPath("/Users/yecihong/go/src/young/run.sh")
	str, err := exec.LookPath("go") //str->/usr/local/Cellar/go@1.14/1.14.8/libexec/bin/go
	if err != nil {
		fmt.Printf("%T\n", err)  //*exec.Error
		fmt.Println("err:", err) //err:exec: "../run.sh": stat ../run.sh: no such file or directory
		return
	}
	fmt.Println("str:", str)
}

func tCmd() {
	cmd := exec.Command("/bin/bash", "-c", "echo hello;sleep 1;ls -ll;")
	fmt.Printf("%T,%v", cmd, cmd) //*exec.Cmd,/bin/bash -c echo hello;sleep 1;ls -ll;%
	err := cmd.Run()
	if err != nil {
		fmt.Println("err", err)
		return
	}
}

func tCmdCombinedOutput() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)
	cmd = exec.Command("/bin/bash", "-c", "echo hello;sleep 10;ls -ll;")
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("output:", string(output))
}

func tCron() {
	var (
		expr *cronexpr.Expression
		//err      error
		now      time.Time
		nextTime time.Time
	)

	//if expr, err = cronexpr.Parse("* * * * *"); err != nil {
	//	fmt.Println("err", err) //err syntax error in day-of-week field: '*qqq'
	//}

	expr = cronexpr.MustParse("*/2 * * * * * *") //相信cron表达式一定是正确的，失败立马退出

	//计算下一次执行时间
	now = time.Now()
	nextTime = expr.Next(now)
	fmt.Println("time ", now)
	fmt.Println("next time ", nextTime)

	//等待定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了", nextTime)
	})
	//Q1:没有sleep直接就执行结束了。。。这又不像个协程
	//Q2:只能执行一次
	//time.Sleep(time.Second * 12)
}
