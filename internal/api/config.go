package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort     string `envconfig:"APP_PORT" default:":8080"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"user-service"`
	InstanceID  string `envconfig:"INSTANCE_ID" default:""`
	AppHostName string `envconfig:"APP_HOST_NAME" default:"localhost:8080"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
