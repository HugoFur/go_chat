package main

import "os"

type Config struct {
	host     string
	port     string
	user     string
	password string
	database string
}

func NewConfig() *Config {
	return &Config{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		user:     os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		database: os.Getenv("DB_DATABASE"),
	}
}

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetUser() string {
	return c.user
}

func (c *Config) GetPassword() string {
	return c.password
}

func (c *Config) GetDatabase() string {
	return c.database
}
