package v1

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secrity-svc/pkg/app"
	"github.com/secrity-svc/pkg/e"
	"github.com/secrity-svc/pkg/logging"
	"github.com/secrity-svc/pkg/util"
	"github.com/secrity-svc/service/content_service"
)

type ExperimentParams struct {
	Host     string `json:"host" binding:"required"`      //支持单ip，多ip（“,”分割），网段,目标ip: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12
	Ports    string `json:"ports" binding:"omitempty"`    //设置扫描的端口: 22 | 1-65535 | 22,80,3306 (default "21,22,80,81,135,139,443,445,1433,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017")
	Portadd  string `json:"portadd" binding:"omitempty"`  //新增需要扫描的端口,-pa 3389 (会在原有端口列表基础上,新增该端口)
	Pocname  string `json:"pocname" binding:"omitempty"`  //指定web poc的模糊名字, -pocname weblogic
	Pocpath  string `json:"pocpath" binding:"omitempty"`  //指定poc路径
	Path     string `json:"path" binding:"omitempty"`     //fcgi、smb romote file path
	User     string `json:"user" binding:"omitempty"`     //指定爆破时的用户名
	Userfile string `json:"userfile" binding:"omitempty"` //指定爆破时的用户名文件
	Password string `json:"password" binding:"omitempty"` //指定爆破时的密码
	Passfile string `json:"passfile" binding:"omitempty"` //指定爆破时的密码文件
	Url      string `json:"url" binding:"omitempty"`      //指定Url扫描
	Time     string `json:"time" binding:"omitempty"`     //端口扫描超时时间 (default 3)
	Threads  string `json:"threads" binding:"omitempty"`  //扫描线程 (default 600)
	Cookie   string `json:"cookie" binding:"omitempty"`   //设置cookie
	Proxy    string `json:"proxy" binding:"omitempty"`    //设置代理, -proxy http://127.0.0.1:8080
	Ping     string `json:"ping" binding:"omitempty"`     //使用ping代替icmp进行存活探测
	Nopoc    string `json:"nopoc" binding:"omitempty"`    //跳过web poc扫描
	SshKey   string `json:"sshkey" binding:"omitempty"`   //ssh连接时,指定ssh私钥
	NoPorts  string `json:"no_ports" binding:"omitempty"` //扫描时要跳过的端口,as: -pn 445
	Scantype string `json:"scantype" binding:"omitempty"` //置扫描模式: -m ssh (default "all")
	NoHosts  string `json:"no_hosts" binding:"omitempty"` //扫描时,要跳过的ip: -hn 192.168.1.1/24
	Command  string `json:"command" binding:"omitempty"`  //exec command (ssh|wmiexec)
}

func changeParamskey(key string) string {
	switch key {
	case "Host":
		return "h"
	case "Pocname":
		return "pocname"
	case "Ports":
		return "p"
	case "Portadd":
		return "pa"
	case "NoPorts":
		return "pn"
	case "Pocpath":
		return "pocpath"
	case "Path":
		return "path"
	case "User":
		return "user"
	case "Userfile":
		return "userf"
	case "Password":
		return "pwd"
	case "Passfile":
		return "pwdf"
	case "Url":
		return "u"
	case "Time":
		return "time"
	case "Cookie":
		return "cookie"
	case "Proxy":
		return "proxy"
	case "Ping":
		return "ping"
	case "Nopoc":
		return "nopoc"
	case "Scantype":
		return "m"
	case "NoHosts":
		return "hn"
	case "Command":
		return "c"
	case "SshKey":
		return "sshkey"
	default:
		return ""
	}
}

// @Summary Add Experiment content
// @Description 通过 JSON 创建一个新的安全实验
// @Tags 创建安全实验
// @Accept  json
// @Produce  json
// @Param params body ExperimentParams  true "host必填,其他选填"
// @Router /api/v1/content [post]
func AddExperiment(c *gin.Context) {
	var (
		appG             = app.Gin{C: c}
		experimentParams ExperimentParams
	)
	if err := c.ShouldBindJSON(&experimentParams); err != nil {
		// JSON解析失败，返回错误信息
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
		return
	}
	args := structToStr(experimentParams)
	logging.Info(args)
	uid, err := util.GenerateUid()
	if err != nil {
		logging.Error(err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, "uid 生成错误："+err.Error())
		return
	}
	go executeLongRunningScript(uid, args)
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"uid": uid})

}
func executeLongRunningScript(uid, args string) {
	osName := runtime.GOOS
	var cmd *exec.Cmd
	switch osName {
	case "windows":
		// CMD批处理脚本
		script := "   ./fscan/fscan.exe   "
		cmd = exec.Command("cmd.exe", "/C", script+" "+args)

	case "darwin": // macOS
		script := "   ./fscan/fscan_mac   "
		cmd = exec.Command("bash", "-c", script+" "+args)

	default: // 假设Linux和其他类Unix系统
		//cmd = exec.Command("bash", "-c", fmt.Sprintf("./%s %s", script, strings.Join(args, " ")))
		script := "   ./fscan/fscan   "
		cmd = exec.Command("bash", "-c", script+" "+args)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logging.Error(err.Error())
	}
	if err := cmd.Start(); err != nil {
		logging.Error(err.Error())
	}
	pid := cmd.Process.Pid
	//插入一条记录带uid，pid
	contentService := content_service.Content{
		Uid:     uid,
		Pid:     pid,
		Result:  "",
		StartAt: time.Now().Local(),
	}
	err = contentService.Add()
	if err != nil {
		logging.Error(err.Error())
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 使用定时器,每5秒记录一次结果记录一次结果到文件
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	reader := io.Reader(stdout)
	buf := make([]byte, 100)
	for {
		select {
		case <-ticker.C:
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				logging.Error(err.Error())
			}
			if n == 0 {
				break
			}
			output := string(buf[:n])

			//根据uid更新库
			contentService := content_service.Content{
				Uid:    uid,
				Result: output,
			}
			err = contentService.Edit(uid, output)
			if err != nil {
				logging.Error(err.Error())
			}
		case err := <-done:
			if err != nil {
				logging.Error(err.Error())
			} else {
				fmt.Println("Script executed successfully.")
				contentService := content_service.Content{
					IsEnd: 1,
				}
				err := contentService.EditIsEndStatus(uid)
				if err != nil {
					logging.Error(err.Error())
				}
			}
			return
		}
	}
}

// 处理参数
func structToStr(experimentParams ExperimentParams) string {
	// 获取User结构体的反射类型和值
	userType := reflect.TypeOf(experimentParams)
	userValue := reflect.ValueOf(experimentParams)
	// 创建一个缓冲区来构建字符串
	var buffer strings.Builder
	// 遍历结构体的每个字段
	for i := 0; i < userType.NumField(); i++ {
		value := userValue.Field(i).Interface()
		field := userType.Field(i)
		if value == "" || value == "string" || field.Name == "" {
			continue
		}
		newKey := changeParamskey(field.Name)
		// 将每个字段名加上前缀并转换为字符串
		strValue := fmt.Sprintf("-%s %v", newKey, value)
		buffer.WriteString(strValue)
		buffer.WriteString(" ")
	}
	// 转换缓冲区内容为字符串
	return buffer.String()
}

type ExperimentProcessId struct {
	Uid string `json:"uid" binding:"required"`
	Pid int    `json:"pid" binding:"required"`
}

// @Summary Destroy A Running Secrity Experiment
// @Description 通过 JSON 销毁一个正在执行的安全实验
// @Tags 销毁安全实验
// @Accept  json
// @Produce  json
// @Param params body ExperimentProcessId true "安全实验唯一标识pid和进程pid必填"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/pid [post]
func DestroyExperiment(c *gin.Context) {
	var (
		appG                = app.Gin{C: c}
		experimentProcessId ExperimentProcessId
	)
	if err := c.ShouldBindJSON(&experimentProcessId); err != nil {
		// JSON解析失败，返回错误信息
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
		return
	}
	pid := experimentProcessId.Pid
	uid := experimentProcessId.Uid
	err := terminateProcess(pid)
	if err != nil {
		logging.Error(err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
		return
	}
	contentService := content_service.Content{
		Uid:         uid,
		IsDestroyed: 1,
	}
	err = contentService.EditStatus(uid, 1)
	if err != nil {
		logging.Error(err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
	}
	appG.Response(http.StatusOK, e.SUCCESS, "destory experiment "+uid+"  sucess")
}

func terminateProcess(pid int) error {
	//syscall.Kill(pid, syscall.SIGTERM)在Go语言中是可以跨多个操作系统（主要是类Unix系统和Windows）尝试结束进程的
	err := syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to terminate process with PID %d: %w", pid, err)
	}
	return nil
}

//type ExperimentParams struct {
//	Host     string `json:"h" binding:"required"`        //支持单ip，多ip（“,”分割），网段,目标ip: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12
//	Port     string `json:"P" binding:"omitempty"`       //设置扫描的端口: 22 | 1-65535 | 22,80,3306 (default "21,22,80,81,135,139,443,445,1433,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017")
//	Portadd  string `json:"pa" binding:"omitempty"`      //新增需要扫描的端口,-pa 3389 (会在原有端口列表基础上,新增该端口)
//	Pocname  string `json:"pocname" binding:"omitempty"` //指定web poc的模糊名字, -pocname weblogic
//	Pocpath  string `json:"pocpath" binding:"omitempty"` //指定poc路径
//	Path     string `json:"path" binding:"omitempty"`    //fcgi、smb romote file path
//	User     string `json:"user" binding:"omitempty"`    //指定爆破时的用户名
//	Userfile string `json:"userf" binding:"omitempty"`   //指定爆破时的用户名文件
//	Password string `json:"pwd" binding:"omitempty"`     //指定爆破时的密码
//	Passfile string `json:"pwdf" binding:"omitempty"`    //指定爆破时的密码文件
//	Url      string `json:"url" binding:"omitempty"`     //指定Url扫描
//	Time     string `json:"time" binding:"omitempty"`    //端口扫描超时时间 (default 3)
//	Threads  string `json:"t" binding:"omitempty"`       //扫描线程 (default 600)
//	Cookie   string `json:"cookie" binding:"omitempty"`  //设置cookie
//	Proxy    string `json:"proxy" binding:"omitempty"`   //设置代理, -proxy http://127.0.0.1:8080
//	Ping     string `json:"ping" binding:"omitempty"`    //使用ping代替icmp进行存活探测
//	Nopoc    string `json:"nopoc" binding:"omitempty"`   //跳过web poc扫描
//	SshKey   string `json:"sshkey" binding:"omitempty"`  //ssh连接时,指定ssh私钥
//	NoPorts  string `json:"pn" binding:"omitempty"`      //扫描时要跳过的端口,as: -pn 445
//	Scantype string `json:"m" binding:"omitempty"`       //置扫描模式: -m ssh (default "all")
//	NoHosts  string `json:"hn" binding:"omitempty"`      //扫描时,要跳过的ip: -hn 192.168.1.1/24
//	Command  string `json:"c" binding:"omitempty"`       //exec command (ssh|wmiexec)
//}