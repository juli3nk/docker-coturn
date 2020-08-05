package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/go-resty/resty"
	"github.com/juli3nk/go-utils/ip"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/go-playground/validator.v10"
)

type EnvOptions struct {
	ListeningPort        int      `split_words:"true" default:"3478" validate:"port"`
	MinPort              int      `split_words:"true" default:"49152" validate:"port"`
	MaxPort              int      `split_words:"true" default:"65535" validate:"port"`
	StaticAuthSecret     string   `split_words:"true"`
	StaticAuthSecretFile string   `split_words:"true" validate:"omitempty,file"`
	ServerName           string   `split_words:"true" default:"coturn"`
	Realm                string   `required:"true" validate:"fqdn"`
	UserQuota            int      `split_words:"true" default:"12"`
	TotalQuota           int      `split_words:"true" default:"1200"`
	DeniedPeerIP         []string `split_words:"true"`
	AllowedPeerIP        string   `split_words:"true" validate:"omitempty,ip"`
}

type Options struct {
	LocalIP  string
	RemoteIP string
	Config   *EnvOptions
}

func main() {
	envOpts := new(EnvOptions)

	if err := envconfig.Process("coturn", envOpts); err != nil {
		log.Fatal(err)
	}

	localIp := getLocalIP()
	remoteIp := getRemoteIP()

	if len(envOpts.StaticAuthSecretFile) > 0 {
		secretBytes, err := ioutil.ReadFile(envOpts.StaticAuthSecretFile)
		if err != nil {
			log.Fatal(err)
		}

		envOpts.StaticAuthSecret = string(secretBytes)
	}

	validate := validator.New()
	validate.RegisterValidation("port", validatePort)

	if err := validate.Struct(envOpts); err != nil {
		log.Fatal(err)
	}

	opts := Options{
		LocalIP:  localIp,
		RemoteIP: remoteIp,
		Config:   envOpts,
	}

	tBody := template.Must(template.New("turnserverConf").Parse(configTemplate))

	bufBody := new(bytes.Buffer)
	if err := tBody.Execute(bufBody, opts); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("/etc/coturn/turnserver.conf", bufBody.Bytes(), 0644);  err != nil {
		log.Fatal(err)
	}
}

func getLocalIP() string {
	infs := ip.New()

	if err := infs.Get(); err != nil {
		return ""
	}

	inf := infs.GetIntf("eth0")

	return inf.V4[0]
}

func getRemoteIP() string {
	client := resty.New()

	resp, err := client.R().
		Get("https://ifconfig.me/ip")
	if err != nil {
		return ""
	}

	return resp.String()
}
