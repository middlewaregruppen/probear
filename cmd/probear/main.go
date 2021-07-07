package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ghodss/yaml"
	"github.com/middlewaregruppen/probear/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	conf, err := ioutil.ReadFile("/config/probear.yaml")
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}

	out, err := yaml.YAMLToJSON(conf)
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}
	log.Printf("%s", out)

	cnf := &config.Config{}
	err = json.Unmarshal(out, cnf)
	if err != nil {
		log.Fatalf("error loading config %s", err)
	}
	log.Printf("%+v", cnf)

	cnf.Network.Probe()

	b, err := json.Marshal(cnf)
	if err != nil {
		log.Fatalf("error probing %s", err)
	}
	log.Printf("%s", b)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	for {
		time.Sleep(10 * time.Second)

		cnf.Network.Probe()
		b, err = json.Marshal(cnf)
		if err != nil {
			log.Fatalf("error probing %s", err)
		}
		log.Printf("%s", b)
	}

}
