package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"veg-store-backend/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
This file handles loading configuration settings from YAML/YML files using Viper.
Logic:
- Determine the mode (dev, prod, etc.) from the MODE environment variable (default to "dev").
- Load configuration from ./config/config.{mode}.yaml or .yml, with a fallback to config.yaml.
- Support environment variable expansion in the format ${VAR} or ${VAR:default}.
- Log loaded configuration values (masking sensitive info like passwords).
- Unmarshal the configuration into a config struct for easy access.

Example YAML structure:
server:
  port: "${SERVER_PORT:8080}"
  api_prefix: "/restful"
  api_version: "v1"

database:
  host: "${DB_HOST:localhost}"
  port: 5432

==> config Struct after unmarshalling:
Server:
  Port: "8080"
  ApiPrefix: "/restful"
  ApiVersion: "v1"

Database:
  Host: "localhost"
  Port: 5432
*/

// Config - Mapping configuration with yaml structure
type Config struct {
	Mode string

	Server struct {
		Port       string `mapstructure:"port"`
		ApiPrefix  string `mapstructure:"api_prefix"`
		ApiVersion string `mapstructure:"api_version"`
	} `mapstructure:"server"`

	JWT struct {
		ExpectedIssuer    string   `mapstructure:"expected_issuer"`
		ExpectedAudiences []string `mapstructure:"expected_audiences"`
		AccessDuration    string   `mapstructure:"access_duration"`
		RefreshDuration   string   `mapstructure:"refresh_duration"`
		PrivateKeyPath    string   `mapstructure:"private_key_path"`
		PublicKeyPath     string   `mapstructure:"public_key_path"`
	} `mapstructure:"jwt"`

	Cors struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	} `mapstructure:"cors"`

	Swagger struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"swagger"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
}

// Load LoadConfig loads configuration from ./config/config.{mode}.yaml or .yml
func Load() *Config {
	Logger.Info(fmt.Sprintf("Load configs for '%s' mode.", Configs.Mode))

	// Load .env file
	_ = godotenv.Load()

	readConfigWithFallback()

	// Expand ${VAR[:default]} syntax using env values
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		expanded := expandEnvWithDefault(val)
		if val != expanded {
			//component.Logger.Info()("Expanding variable for '%s': '%s' → '%s'\n", key, val, expanded)
			viper.Set(key, expanded)
		}
	}

	// Log all loaded values (masking passwords)
	if Configs.Mode != "prod" && Configs.Mode != "production" {
		logAppConfig()
	}

	// Unmarshal (is equivalent to Decode) into config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		Logger.Fatal("failed to unmarshal config", zap.Error(err))
	}

	return &config
}

// --- Helper: fallback config reader ---
func readConfigWithFallback() {
	// setup Viper with fallback order: config.{mode}.yaml → config.{mode}.yml → config.yaml
	viper.SetConfigName(fmt.Sprintf("config.%s", Configs.Mode))
	viper.SetConfigType("yaml")

	// Set config path to .../.../config
	configPath := util.GetConfigPathFromGoMod("config")
	Logger.Info(fmt.Sprintf("config path: %s", configPath))
	viper.AddConfigPath(configPath)

	// Try .yaml first
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	// Try .yml
	ymlPath := strings.TrimSuffix(viper.ConfigFileUsed(), filepath.Ext(viper.ConfigFileUsed())) + ".yml"
	viper.SetConfigFile(ymlPath)
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	// Try config.yaml (default)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err == nil {
		return
	}

	Logger.Fatal("no valid config file found")
}

// --- Helper: supports ${VAR} and ${VAR:default} ---
func expandEnvWithDefault(input string) string {
	// Regex to match ${VAR} or ${VAR:default}
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(?::([^}]*))?}`)

	return re.ReplaceAllStringFunc(input, func(s string) string {
		// Extract variable name and default value
		matches := re.FindStringSubmatch(s)

		// Ensure we have at least the variable name
		if len(matches) < 2 {
			return s
		}

		// matches[1] is the variable name, matches[2] is the default value (if any)
		key := matches[1]
		def := ""

		// If default value is provided
		if len(matches) == 3 {
			def = matches[2]
		}

		if val, ok := os.LookupEnv(key); ok && val != "" {
			return val
		}
		return def
	})
}

func logAppConfig() {
	configFile := viper.ConfigFileUsed()
	var fields []zap.Field

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		if strings.Contains(strings.ToLower(key), "password") {
			val = "********"
		}
		fields = append(fields, zap.Any(key, val))
	}

	Logger.Info("Application configuration loaded",
		zap.String("config_file", configFile),
		zap.Any("configs", fields),
	)
}
