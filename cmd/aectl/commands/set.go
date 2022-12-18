package commands

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bonsai-oss/traefik-weight-control/internal/util"
)

type SetCommand struct {
	File    string
	DryRun  bool
	Service string
	Server  string
	Weight  int
}

func (lsc *SetCommand) Execute(ctx *kingpin.ParseContext) error {
	configFile := util.ConfigFile{FileName: lsc.File}
	traefikConfiguration, providerError := configFile.Read()
	if providerError != nil {
		log.Err(providerError).Msg("Failed to read configuration file")
		os.Exit(1)
	}

	if traefikConfiguration.HTTP == nil {
		log.Error().Msg("No HTTP configuration found")
		os.Exit(1)
	}

	service, ok := traefikConfiguration.HTTP.Services[lsc.Service]
	if !ok {
		log.Error().Str("service", lsc.Service).Msgf("Service not found")
		os.Exit(1)
	}
	if service.Weighted == nil {
		log.Error().Str("service", lsc.Service).Msgf("Service is not weighted")
		os.Exit(1)
	}
	if service.Weighted.Services == nil {
		log.Error().Str("service", lsc.Service).Msgf("Service has no servers")
		os.Exit(1)
	}

	set := false
	for serverIndex := range traefikConfiguration.HTTP.Services[lsc.Service].Weighted.Services {
		if traefikConfiguration.HTTP.Services[lsc.Service].Weighted.Services[serverIndex].Name != lsc.Server {
			continue
		}
		old := traefikConfiguration.HTTP.Services[lsc.Service].Weighted.Services[serverIndex].Weight
		traefikConfiguration.HTTP.Services[lsc.Service].Weighted.Services[serverIndex].Weight = &lsc.Weight
		set = true

		log.Info().Str("service", lsc.Service).Str("server", lsc.Server).Int("old", *old).Int("new", lsc.Weight).Msg("Server weight updated")
	}
	if !set {
		log.Error().Str("service", lsc.Service).Str("server", lsc.Server).Msgf("Server not found in service")
		os.Exit(1)
	}

	if lsc.DryRun {
		log.Info().Msg("Dry run, skipping configuration file update")
		return nil
	}
	configFile.Write(traefikConfiguration)

	return nil
}
