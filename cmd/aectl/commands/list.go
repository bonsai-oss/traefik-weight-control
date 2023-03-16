package commands

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rs/zerolog/log"

	"github.com/bonsai-oss/traefik-weight-control/internal/util"
)

type outputFormat struct {
	Service string   `json:"service" yaml:"service"`
	Servers []Server `json:"servers" yaml:"servers"`
}

type Server struct {
	Name   string `json:"name" yaml:"name"`
	Weight int    `json:"weight" yaml:"weight"`
}

type ListCommand struct {
	File   string
	Format string
}

func (listCommand *ListCommand) Execute(ctx *kingpin.ParseContext) error {
	configFile := util.ConfigFile{FileName: listCommand.File}
	traefikConfiguration, providerError := configFile.Read()
	if providerError != nil {
		log.Err(providerError).Msg("Failed to read configuration file")
		os.Exit(1)
	}

	if traefikConfiguration.HTTP == nil {
		log.Error().Msg("No HTTP configuration found")
		os.Exit(1)
	}

	var output []outputFormat

	for serviceName, service := range traefikConfiguration.HTTP.Services {
		if service.Weighted == nil {
			continue
		}
		var servers []Server
		for _, server := range service.Weighted.Services {
			servers = append(servers, Server{Name: server.Name, Weight: *server.Weight})
		}
		output = append(output, outputFormat{Service: serviceName, Servers: servers})
	}

	switch listCommand.Format {
	case "text":
		for _, o := range output {
			fmt.Println(o.Service)
			for _, s := range o.Servers {
				fmt.Printf("  %s: %d\n", s.Name, s.Weight)
			}
		}
	default:
		// reuse the configuration_provider code to get the output encoder for the selected format
		_, encoder, _ := (&util.ConfigFile{FileName: "test." + listCommand.Format}).SelectProvider(os.Stdout)
		encoder.Encode(output)
	}

	return nil
}
