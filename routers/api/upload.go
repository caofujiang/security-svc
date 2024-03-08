package api

import (
	_ "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/secrity-svc/pkg/app"
	"github.com/secrity-svc/pkg/e"
	"github.com/secrity-svc/pkg/logging"
	"github.com/secrity-svc/pkg/upload"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})
}

// 方案1  运行完结果存储
func Test(c *gin.Context) {
	appG := app.Gin{C: c}
	params := c.PostForm("params")
	fmt.Println("params=", params)

	var cmd *exec.Cmd

	cmd = exec.CommandContext(c, "/bin/sh", "-c", "./fscan/fscan_mac"+" "+params)
	//output, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("error=", err.Error())
	//}
	//outMsg := string(output)
	//fmt.Println("outinfo:=", outMsg)

	cmd.Start()

	//if runtime.GOOS == "linux" {
	//	cmd = exec.Command("ls")
	//}

	//log.Debugf(ctx, "execScript Command out: %s", outIsPython2)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"message": "Script execution started asynchronously"})
}

// 方案2  运行5秒存储一次
func AsyncExecuteShell(c *gin.Context) {
	appG := app.Gin{C: c}
	params := c.PostForm("params")
	fmt.Println("params=", params)
	outputChan := make(chan string)
	go executeLongRunningScript(c, outputChan)
	go recordOutputToFile(outputChan)
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"message": "Script execution started asynchronously"})
}

func executeLongRunningScript(c *gin.Context, outputChan chan string) {
	params := "-h 192.168.123.93 -pocname weblogic"
	//params := "-h 192.168.123.93 -p 139"
	//params := "-h 192.168.123.93 -c whoami;id"
	cmd := exec.Command("sh", "-c", "   ./fscan/fscan_mac   "+params)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Process.Pid)
	// 每隔一段时间记录一次结果到文件
	ticker := time.NewTicker(5 * time.Second) // 每10秒记录一次结果
	defer ticker.Stop()

	reader := io.Reader(stdout)
	buf := make([]byte, 100)
	for {
		select {
		case <-ticker.C:
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
			output := string(buf[:n])
			outputChan <- output
		}
	}
	close(outputChan)
}

// 改为追加方式写库
func recordOutputToFile(outputChan chan string) {
	file, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for output := range outputChan {
		_, err := file.WriteString(output)
		//写库

		if err != nil {
			log.Fatal(err)
		}
	}
}

//
//	// 使用context来控制超时或者取消请求
//	//ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) // 设置超时时间
//	//defer cancel()
//
//	// 从channel中等待获取结果
//	//select {
//	//case result := <-outCh:
//	//	c.JSON(http.StatusOK, gin.H{"output": string(result)})
//	//case <-ctx.Done():
//	//	c.JSON(http.StatusRequestTimeout, gin.H{"error": "Script execution timed out"})
//	//}
//
//	//result := <-outCh
//	//id, err := GenerateUid()
//	//if err != nil {
//	//
//	//}
//	//
//	//appG.Response(http.StatusOK, e.SUCCESS, id)
//}

//方案3  sse方式
//func AsyncExecuteShell(c *gin.Context) {
//
//	//params := c.PostForm("params")
//
//	params := "-h 192.168.123.93 -pocname weblogic -o cfj.txt"
//	fmt.Println("params=", params)
//	//	// 定义用来存放输出结果的channel
//
//	// 从请求参数中获取脚本名称和参数
//	//scriptName := c.Query("script_name")
//	//scriptArgs := c.QueryArray("script_args")
//
//	// 创建参数数组用于执行命令
//	//args := append([]string{scriptName}, scriptArgs...)
//
//	// 创建一个上下文和取消函数，以便在必要时取消脚本执行
//	ctx, cancel := context.WithCancel(c.Request.Context())
//
//	// 创建一个goroutine来异步执行shell脚本
//	outCh := make(chan []byte)
//	errCh := make(chan error)
//
//	go func() {
//		defer close(outCh)
//		defer close(errCh)
//		cmd := exec.Command("sh", "-c", "   ./fscan/fscan_mac   "+params) // 替换为你的shell脚本路径和参数
//
//		//cmd := exec.CommandContext(ctx, "sh", "-c", strings.Join(args, " "))
//		var stdoutBuf, stderrBuf bytes.Buffer
//		cmd.Stdout = &stdoutBuf
//		cmd.Stderr = &stderrBuf
//
//		// 异步执行命令
//		err := cmd.Run()
//		if err != nil {
//			errCh <- err
//			return
//		}
//
//		// 将执行结果定期发送至channel
//		ticker := time.NewTicker(time.Second) // 设置每隔一秒发送一次输出（可按需调整）
//		for {
//			select {
//			case <-ticker.C:
//				outCh <- stdoutBuf.Bytes()
//				//写库
//
//			case <-ctx.Done(): // 当context被取消时停止发送
//				ticker.Stop()
//				return
//			}
//		}
//	}()
//
//	// 发送脚本执行结果
//	c.Stream(func(w io.Writer) bool {
//		select {
//		case data := <-outCh:
//			// 创建并发送SSE事件
//			event := sse.Event{
//				Data: string(data),
//			}
//			sse.Encode(w, event)
//			return true
//		case err := <-errCh:
//			// 发送错误并结束Stream
//			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return false
//		}
//	})
//
//	// 当c.Stream返回时，意味着客户端已经断开连接，取消脚本执行上下文
//	defer cancel()
//}
