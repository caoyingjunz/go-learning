package config

type Config struct {
	Default DefaultConfig
	Mysql   MysqlConfig
}

type DefaultConfig struct {
	Name string `yaml:"name"`
}

type MysqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}
