package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bonsai-oss/traefik-weight-control/internal/util"
)

type parameters struct {
	file  string
	debug bool
}

var params parameters

func init() {
	app := kingpin.New("aectl", "Traefik Weight Control")
	app.HelpFlag.Short('h')
	app.Flag("file", "Path to the Traefik configuration file").Short('f').Required().StringVar(&params.file)
	app.Flag("verbose", "Enable debug mode").Short('v').BoolVar(&params.debug)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if params.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	configFile := util.ConfigFile{FileName: params.file}
	traefikConfiguration, providerError := configFile.Read()
	if providerError != nil {
		log.Err(providerError).Msg("Failed to select configuration file traefikConfiguration")
		os.Exit(1)
	}

	if traefikConfiguration.HTTP == nil {
		log.Error().Msg("No HTTP configuration found")
		os.Exit(1)
	}

	log.Info().Msg("HTTP configuration found")

	for serviceName, service := range traefikConfiguration.HTTP.Services {
		if service.Weighted == nil {
			continue
		}
		log.Info().Str("service", serviceName).Msg("Service found")
	}

	configFile.FileName = params.file + ".new.yaml"
	err := configFile.Write(traefikConfiguration)
	if err != nil {
		log.Err(err).Msg("Failed to write configuration file")
	}
}
