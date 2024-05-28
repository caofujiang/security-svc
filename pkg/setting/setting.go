package setting

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type App struct {
	JwtSecret string `yaml:"jwtSecret"`
	PageSize  int    `yaml:"pageSize"`
	PrefixUrl string `yaml:"prefixUrl"`

	RuntimeRootPath string `yaml:"runtimeRootPath"`

	ImageSavePath string `yaml:"imageSavePath"`
	ImageMaxSize  int    `yaml:"imageMaxSize"`

	ExportSavePath string `yaml:"exportSavePath"`
	QrCodeSavePath string `yaml:"qrCodeSavePath"`
	FontSavePath   string `yaml:"fontSavePath"`

	LogSavePath string `yaml:"logSavePath"`
	LogSaveName string `yaml:"logSaveName"`
	LogFileExt  string `yaml:"logFileExt"`
	TimeFormat  string `yaml:"timeFormat"`
}

var AppSetting = &App{}

type Server struct {
	RunMode  string `yaml:"runMode"`
	HttpPort int    `yaml:"httpPort"`
}

var ServerSetting = &Server{}

type Database struct {
	Type        string `yaml:"type"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	DbName      string `yaml:"dbName"`
	TablePrefix string `yaml:"tablePrefix"`
}

var DatabaseSetting = &Database{}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	App      App      `yaml:"app"`
}

//var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	//var configData map[string]interface{}
	var err error
	yamlFile, err := os.ReadFile("conf/app.yml")
	//cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf file': %v", err)
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("setting.Setup  Unmarshal yaml fail：", err)
		return
	}

	DatabaseSetting = &config.Database
	ServerSetting = &config.Server
	AppSetting = &config.App
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
}
