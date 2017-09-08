package main

type Server struct {
	Port int `yaml:"port"`
}

type Config struct {
	Server Server `yaml:"server"`
}
