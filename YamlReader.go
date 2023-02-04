package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"regexp"
)

type server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type rateLimit struct {
	IpRequestPerSec   int `yaml:"ip_request_per_sec"`
	OtpSmsIntervalSec int `yaml:"otp_sms_interval_sec"`
}
type conf struct {
	Server    server    `yaml:"server"`
	RateLimit rateLimit `yaml:"rate_limit"`
}

func checkValidation(c *conf) (*conf, error) {
	if c.Server == (server{}) {
		return c, errors.New("server is Nil")
	}

	matchedPort, err := regexp.MatchString("^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$", c.Server.Port)
	if !matchedPort || err != nil {
		return c, errors.New("port is not valid")
	}

	if c.RateLimit == (rateLimit{}) {
		return c, nil
	}

	if c.RateLimit.IpRequestPerSec == 0 || (c.RateLimit.IpRequestPerSec < 60 && c.RateLimit.IpRequestPerSec > 1000) {
		return c, errors.New("IpRequestPerSec is not valid")
	}

	if c.RateLimit.OtpSmsIntervalSec != 0 && c.RateLimit.OtpSmsIntervalSec < 60 && c.RateLimit.OtpSmsIntervalSec > 300 {
		return c, errors.New("OtpSmsIntervalSec is not valid")
	}
	return c, nil
}

func setDefaults(c *conf) *conf {

	if len(c.Server.Host) <= 0 {
		c.Server.Host = "localhost"
	}

	return c
}
func (c *conf) getConf() *conf {

	yamlFile, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c

}

func main() {
	var c conf
	c.getConf()
	_, err := checkValidation(&c)

	if err != nil {
		fmt.Println(err)
		return
	}
	setDefaults(&c)

	fmt.Println(c)
}
