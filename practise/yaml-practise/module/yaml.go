package model

type Yaml struct {
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
	Cache struct {
		Enable bool     `yaml:"enable"`
		List   []string `yaml:"list,flow"`
	}
	Student Student `yaml:"student"`
	Name    string  `yaml:"name"`
}

type Student map[string]Item

type Item struct {
	School string `yaml:"school" json:"school"`
	Age    int    `yaml:"age" json:"age,omitempty"`
}
