package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	WebServerPort  string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort string `mapstructure:"GRPC_SERVER_PORT"`
	RcpServerPort  string `mapstructure:"RPC_SERVER_PORT"`
	MongUrl        string `mapstructure:"MONG_URL"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
