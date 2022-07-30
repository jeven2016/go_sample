package common

type Config struct {
	ApiServerConfig ApiServerConfig `mapstructure:"api-server,omitempty"`
	MongoConfig     MongoConfig     `mapstructure:"mongodb,omitempty"`
	RedisConfig     RedisConfig     `mapstructure:"redis,omitempty"`
	AuthConfig      AuthConfig      `mapstructure:"auth,omitempty"`
}

type ApiServerConfig struct {
	LogLevel      string `mapstructure:"log_level,omitempty"`
	LogPath       string `mapstructure:"log_path,omitempty"`
	OutputConsole bool   `mapstructure:"output_console,omitempty"`
	Dev           bool   `mapstructure:"dev,omitempty"`
	ServiceName   string `mapstructure:"service_name,omitempty"`
	Port          int32  `mapstructure:"port,omitempty"`
	BindAddress   string `mapstructure:"bind_address,omitempty"`
}

type MongoConfig struct {
	Uri      string `bson:"uri"`
	Database string `bson:"database"`
}

type RedisConfig struct {
	Address      string `mapstructure:"address,omitempty"`
	Password     string `mapstructure:"password,omitempty"`
	DefaultDb    int    `mapstructure:"default_db,omitempty"`
	PoolSize     int    `mapstructure:"pool_size,omitempty"`
	PoolTimeout  int    `mapstructure:"pool_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout_seconds"`
	WriteTimeout int    `mapstructure:"write_timeout_seconds"`
}

type AuthConfig struct {
	EnableAuth    bool   `mapstructure:"enable_auth,omitempty"`
	KeycloakUrl   string `mapstructure:"keycloak_url,omitempty"`
	ClientId      string `mapstructure:"client_id,omitempty"`
	ClientSecret  string `mapstructure:"client_secret,omitempty"`
	KeycloakRealm string `mapstructure:"keycloak_realm,omitempty"`
}
