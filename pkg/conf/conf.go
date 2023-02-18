package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigStruct struct {
	PushGateway string   `yaml:"push_gateway"`
	DnsCheck    DnsCheck `yaml:"dns_check"`
	Instance    string   `yaml:"instance"`
	Ping        Ping     `yaml:"ping"`
	Job         string   `yaml:"job"`
	Cron        string   `yaml:"cron"`
}

type Ping struct {
	MachineListURL string   `yaml:"machine_list_url"`
	MachineList    []string `yaml:"machine_list"`
	Size           int      `yaml:"size"`
	Packages       int      `yaml:"packages"`
	Interval       int      `yaml:"interval"`
	Retry          int      `yaml:"retry"`
	Pool           int      `yaml:"pool"`
	TimeOut        int      `yaml:"timeout"`
}

type Domains struct {
	Resolve string `yaml:"resolve"`
	Type    string `yaml:"type"`
}

type DnsCheck struct {
	Domains []Domains `yaml:"domains"`
	Retry   int       `yaml:"retry"`
	Pool    int       `yaml:"pool"`
}

var Config ConfigStruct

func ReadConfigFile(path string) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Config file may not exists")
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		panic("Config file may not a yaml file, please check")
	}
}
