package settings

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	Cig *ini.File

	ServerAddr string

	RedisAddr     string
	RedisPassword string
)

func Load() {
	configFile := resolveConfigFile()
	var err error
	Cig, err = ini.Load(configFile)
	if err != nil {
		log.Fatalf("Fail to load file %v: %v", configFile, err)
	}
	LoadServer()
	LoadDB()
}

func resolveConfigFile() string {
	configFileName := "settings.ini"
	if exePath, err := os.Executable(); err == nil {
		configFile := filepath.Join(filepath.Dir(exePath), configFileName)
		if _, err := os.Stat(configFile); err == nil {
			return configFile
		}
	}

	if configFile, err := filepath.Abs(configFileName); err == nil {
		if _, err := os.Stat(configFile); err == nil {
			return configFile
		}
	}

	log.Fatalf("Fail to locate %v in executable directory or current directory", configFileName)
	return ""
}

func LoadServer() {
	sectionName := "server"
	sec, err := Cig.GetSection(sectionName)
	if err != nil {
		log.Fatalf("Fail to load section %v: %v", sectionName, err)
	}
	ServerAddr = sec.Key("server_addr").MustString(":8000")
}

func LoadDB() {
	sectionName := "db"
	sec, err := Cig.GetSection(sectionName)
	if err != nil {
		log.Fatalf("Fail to load section %v: %v", sectionName, err)
	}
	RedisAddr = sec.Key("redis_addr").MustString("localhost:6379")
	RedisPassword = sec.Key("redis_password").MustString("")
}
