package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	flag "github.com/spf13/pflag"

	"github.com/ghodss/yaml"
	"github.com/middlewaregruppen/probear/pkg/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var flagConfig string

func init() {
	flag.StringVar(&flagConfig, "config", "config.yaml", "path to config file")
}

func main() {

	flag.Parse()

	conf, err := ioutil.ReadFile(flagConfig)
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}

	out, err := yaml.YAMLToJSON(conf)
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}
	log.Printf("%s", out)

	probes := &probe.Probes{}
	err = json.Unmarshal(out, probes)
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}
	log.Printf("%+v", probes)

	probes.Start()

	k8sprobes := probe.K8SProbes{}
	k8sprobes.Scan(10)

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":2112", nil)

}
