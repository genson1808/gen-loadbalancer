package mailer

import (
	"github.com/genson1808/balancer/config"
	"log"
)

var Cfg *config.Config
var AppMailer Mailer

func init() {
	var err error
	Cfg, err = config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = Cfg.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	AppMailer = New(Cfg.SmtpHost, Cfg.SmtpPort, Cfg.Username, Cfg.Password, Cfg.Sender)

}
