package config

type ServerConfig struct {
	RedisInfo     RedisConfig `mapstructure:"redis" json:"redis"`
	EthBlock      string      `mapstructure:"ethBlock" json:"ethBlock"`
	EthBlockList  string      `mapstructure:"ethBlockList" json:"ethBlockList"`
	EthUrl        string      `mapstructure:"ethUrl" json:"ethUrl"`
	WalletPools   string      `mapstructure:"WalletPools" json:"WalletPools"`
	MatchedBlocks string      `mapstructure:"MatchedBlocks" json:"MatchedBlocks"`
	OrderState    string      `mapstructure:"order_state" json:"order_state"`
	MysqlInfo     MysqlConfig `mapstructure:"mysql" json:"mysql"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Db       string `mapstructure:"db" json:"db"`
}
