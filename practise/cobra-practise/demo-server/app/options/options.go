package options

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/cobra"

	"go-learning/practise/cobra-practise/demo-server/app/config"
	"go-learning/practise/cobra-practise/demo-server/dao"
)

// Options has all the params needed to run a Autoscaler
type Options struct {
	ComponentConfig *config.Config

	// ConfigFile is the location of the autoscaler's configuration file.
	ConfigFile string

	Master string

	DBFactory dao.ShareDBFactory

	DB *gorm.DB
}

func NewOptions() (*Options, error) {
	return &Options{
		Master: "demo-master",
	}, nil
}

// BindFlags binds the demo Configuration struct fields
func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "configfile", "", "ConfigFile is the location of the demo's configuration file.")
}

const (
	defaultConfigFile = "democonfig.yaml"
)

func (o *Options) registerDatabase() error {
	sqlConfig := o.ComponentConfig.Mysql
	dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Host,
		sqlConfig.Port,
		sqlConfig.Name)
	DB, err := gorm.Open("mysql", dbConnection)
	if err != nil {
		return err
	}

	// set the connect pools
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	o.DB = DB

	// 注册或者通过 factory 的方式参数调用
	dao.Register(o.DB)
	return nil
}

func (o *Options) Complete() error {
	configFile := defaultConfigFile
	if len(o.ConfigFile) != 0 {
		configFile = o.ConfigFile
	}

	cfg, err := loadConfigFromFile(configFile)
	if err != nil {
		return err
	}
	o.ComponentConfig = cfg

	// 注册数据库链接池
	if err = o.registerDatabase(); err != nil {
		return err
	}

	// 构造数据库工厂函数
	o.DBFactory = dao.NewDBFactory(o.DB)

	// TODO 其他化客户端初始化

	return nil
}
