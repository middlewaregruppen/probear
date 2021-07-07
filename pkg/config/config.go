package config

import "github.com/middlewaregruppen/probear/pkg/network"

type Config struct {
	Network network.Network `json:"network"`
}
