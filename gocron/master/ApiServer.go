package master

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
	"young/gocron/common"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer //单例对象
)

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
		staticDir http.Dir	// 静态文件根目录
		staticHandler http.Handler	// 静态文件的HTTP回调
	)

	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handlerJobSave)
	mux.HandleFunc("/job/list", handlerJobList)
	mux.HandleFunc("/job/delete", handlerJobDelete)
	mux.HandleFunc("/job/kill", handleJobKill)

	// 静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	fmt.Println(staticDir)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))	//   ./webroot/index.html

	//启动端口监听

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		fmt.Println("服务启动失败")
		return err
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	//赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	return
}

//保存任务接口
//POST job={"name":"job1","command":"echo hello;","cronExpr":"5/* * * * * *"}
func handlerJobSave(w http.ResponseWriter, r *http.Request) () {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	postJob = r.PostForm.Get("job")

	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//success
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		_, _ = w.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = w.Write(bytes)
	}
}

//删除任务接口
//POST  /job/delete name=job1
func handlerJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	if err = req.ParseForm(); err != nil {
		return
	}

	name = req.PostForm.Get("name")

	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		_, _ = resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = resp.Write(bytes)
	}

}

func handlerJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		jobList []*common.Job
		bytes   []byte
	)

	if jobList, err = G_jobMgr.ListJob(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		_, _ = resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = resp.Write(bytes)
	}
}

//杀死任务
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要杀死的任务名
	name = req.PostForm.Get("name")

	// 杀死任务
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		_, _ = resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = resp.Write(bytes)
	}
}
