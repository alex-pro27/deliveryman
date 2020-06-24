package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var ConfPath = os.Getenv("DELIVERY_API_CONF")

type TypeConfig struct {
	Database struct {
		User     string `yaml:"USER"`
		Database string `yaml:"DATABASE"`
		Port     string `yaml:"PORT"`
		Password string `yaml:"PASSWORD"`
		Host     string `yaml:"HOST"`
	}

	System struct {
		AppName     string `yaml:"APP_NAME"`
		SecretKey   string `yaml:"SECRET_KEY"`
		Debug       bool   `yaml:"DEBUG"`
		GRPCServer  string `yaml:"GRPC_SERVER"`
		HTTPServer  string `yaml:"HTTP_SERVER"`
		LogPath     string `yaml:"LOG_PATH"`
		ServerUrl   string `yaml:"SERVER_URL"`
		MemProfiler string `yaml:"MEM_PROFILER"`
	}

	Session struct {
		Key    string `yaml:"KEY"`
		MaxAge int    `yaml:"MAX_AGE"`
	}

	Admins []struct {
		Name  string `yaml:"NAME"`
		Email string `yaml:"EMAIL"`
	}

	Email struct {
		Name     string `yaml:"NAME"`
		From     string `yaml:"FROM"`
		Host     string `yaml:"HOST"`
		Port     int    `yaml:"PORT"`
		User     string `yaml:"USER"`
		Password string `yaml:"PASSWORD"`
	}

	Static struct {
		StaticRoot string `yaml:"STATIC_ROOT"`
		MediaRoot  string `yaml:"MEDIA_ROOT"`
	}

	Firebase struct {
		CertPath string `yaml:"CERT_PATH"`
	}

	SMS struct {
		URL      string `yaml:"URL"`
		Login    string `yaml:"LOGIN"`
		Password string `yaml:"PASSWORD"`
		Source   string `yaml:"SOURCE"`
	}
	PaymentService struct {
		Server  string `yaml:"SERVER"`
		Methods struct {
			ConfirmOrder string `yaml:"CONFIRM_ORDER"`
			CreateOrder  string `yaml:"CHECK_STATUS"`
			CheckStatus  string `yaml:"CHECK_STATUS"`
			RefundPay    string `yaml:"REFUND_PAY"`
		}
		Login    string `yaml:"LOGIN"`
		Password string `yaml:"PASSWORD"`
	}
}

var Config *TypeConfig

func Init() {
	data, _ := ioutil.ReadFile(ConfPath)
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatal(err)
	} else {
		for _, p := range []string{
			Config.Static.MediaRoot,
			Config.Static.StaticRoot,
		} {
			err := os.MkdirAll(p, 0644)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
