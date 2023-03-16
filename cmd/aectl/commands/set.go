package commands

import (
	"errors"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/rs/zerolog/log"

	"github.com/bonsai-oss/traefik-weight-control/internal/util"
)

type SetCommand struct {
	File       string
	DryRunMode bool
	GlobalMode bool
	Service    string
	Server     string
	Weight     int
}

func (setCommand *SetCommand) ValidateParameters(ctx *kingpin.ParseContext) error {
	if setCommand.Weight < 0 {
		return errors.New("weight must be a positive integer")
	}
	if setCommand.GlobalMode && setCommand.Service != "" {
		return errors.New("cannot set both global and service")
	}
	if !setCommand.GlobalMode && setCommand.Service == "" {
		return errors.New("either global or service must be set")
	}
	return nil
}

func (setCommand *SetCommand) Execute(ctx *kingpin.ParseContext) error {
	configFile := util.ConfigFile{FileName: setCommand.File}
	traefikConfiguration, providerError := configFile.Read()
	if providerError != nil {
		log.Err(providerError).Msg("Failed to read configuration file")
		os.Exit(1)
	}

	if traefikConfiguration.HTTP == nil {
		log.Error().Msg("No HTTP configuration found")
		os.Exit(1)
	}

	set := false
	setServerWeight := func(serviceName string) {
		service, ok := traefikConfiguration.HTTP.Services[serviceName]
		if !ok {
			log.Error().Str("service", serviceName).Msgf("Service not found")
			os.Exit(1)
		}
		if service.Weighted == nil {
			log.Debug().Str("service", serviceName).Msgf("Service is not weighted")
			return
		}
		if service.Weighted.Services == nil {
			log.Error().Str("service", serviceName).Msgf("Service has no servers")
			return
		}
		for serverIndex := range traefikConfiguration.HTTP.Services[serviceName].Weighted.Services {
			if traefikConfiguration.HTTP.Services[serviceName].Weighted.Services[serverIndex].Name != setCommand.Server {
				continue
			}
			old := traefikConfiguration.HTTP.Services[serviceName].Weighted.Services[serverIndex].Weight
			traefikConfiguration.HTTP.Services[serviceName].Weighted.Services[serverIndex].Weight = &setCommand.Weight
			set = true

			log.Info().Str("service", serviceName).Str("server", setCommand.Server).Int("old", *old).Int("new", setCommand.Weight).Msg("Server weight updated")
		}
	}
	switch setCommand.GlobalMode {
	case true:
		for serviceName := range traefikConfiguration.HTTP.Services {
			setServerWeight(serviceName)
		}
	default:
		setServerWeight(setCommand.Service)
	}

	if !set {
		log.Error().Str("service", setCommand.Service).Str("server", setCommand.Server).Msgf("Server not found in service")
		os.Exit(1)
	}

	if setCommand.DryRunMode {
		log.Info().Msg("Dry run, skipping configuration file update")
		return nil
	}
	configFile.Write(traefikConfiguration)

	return nil
}
