package config

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	Server    ServerConfig
	BBS       BBSConfig
	Paths     PathsConfig
}

type ServerConfig struct {
	Port    int
	Address string
}

type BBSConfig struct {
	BBSName   string
	SysopName string
}

type PathsConfig struct {
	AnsiPath     string
	AsciiPath    string
	DoorPath     string
	MenuPath     string
	MessagePath  string
	DBPath       string
	ConfigsPath  string
}

func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port:    cfg.Section("server").Key("port").MustInt(2323),
			Address: cfg.Section("server").Key("address").MustString("0.0.0.0"),
		},
		BBS: BBSConfig{
			BBSName:   cfg.Section("bbs").Key("bbsname").MustString("My BBS"),
			SysopName: cfg.Section("bbs").Key("sysopname").MustString("Sysop"),
		},
		Paths: PathsConfig{
			AnsiPath:    cfg.Section("paths").Key("ansipath").MustString("ansi/"),
			AsciiPath:   cfg.Section("paths").Key("asciipath").MustString("ascii/"),
			DoorPath:    cfg.Section("paths").Key("doorpath").MustString("doors/"),
			MenuPath:    cfg.Section("paths").Key("menupath").MustString("menus/"),
			MessagePath: cfg.Section("paths").Key("messagepath").MustString("messages/"),
			DBPath:      cfg.Section("paths").Key("dbpath").MustString("bbs.db"),
			ConfigsPath: cfg.Section("paths").Key("configspath").MustString("configs/"),
		},
	}

	return config, nil
}
