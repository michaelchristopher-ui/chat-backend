package conf

//Config stores the config values obtained from the yaml file
type Config struct {
	Server   Server   `yaml:"server"`
	Redis    Redis    `yaml:"redis"`
	Database Database `yaml:"database"`
}

//Redis defines the stored config values of the redis block within the yaml file
type Redis struct {
	Host        string `yaml:"host"`
	Password    string `yaml:"password"`
	UserTimeout int    `yaml:"user_timeout"`
}

//Server defines the stored config values of the server block within the yaml file
type Server struct {
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Port         string `yaml:"port"`
}

//Database defines the stored config values of the database block within the yaml file
type Database struct {
	Dsn string `yaml:"dsn"`
}
