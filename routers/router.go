package routers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/secrity-svc/docs"
	"github.com/secrity-svc/routers/api/v1"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	//r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	//r.POST("/upload", api.UploadImage)
	//r.POST("/auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	apiv1.Use()
	{
		//创建安全演练执行结果
		apiv1.POST("/content", v1.AddExperiment)
		//销毁正在执行的安全演练
		apiv1.POST("/pid", v1.DestroyExperiment)
	}
	return r
}
