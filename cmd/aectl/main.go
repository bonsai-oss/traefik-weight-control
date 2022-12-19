package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bonsai-oss/traefik-weight-control/cmd/aectl/commands"
)

type parameters struct {
	file  string
	debug bool
}

// version is the version of the binary. It is set at build time using the -ldflags -X option.
var version = "dev"

func main() {
	// Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	var params parameters

	app := kingpin.New("aectl", "Traefik Weight Control")
	app.HelpFlag.Short('h')
	app.Version(version)

	listCommand := &commands.ListCommand{}
	kpListCommand := app.Command("list", "List all services and servers").Action(listCommand.Execute)
	kpListCommand.Flag("format", "Output format").Short('o').Envar("AECTL_OUTPUT_FORMAT").Default("text").EnumVar(&listCommand.Format, "text", "json", "yaml", "yml")

	setCommand := &commands.SetCommand{}
	kpSetCommand := app.Command("set", "Set the weight of a server").PreAction(setCommand.ValidateParameters).Action(setCommand.Execute)
	kpSetCommand.Flag("dry-run", "Dry run").Short('d').BoolVar(&setCommand.DryRunMode)
	kpSetCommand.Flag("global", "Set the weight for all services").Short('g').BoolVar(&setCommand.GlobalMode)
	kpSetCommand.Flag("service", "Service name").Short('s').StringVar(&setCommand.Service)
	kpSetCommand.Flag("server", "Server name").Short('n').Required().StringVar(&setCommand.Server)
	kpSetCommand.Flag("weight", "Server weight").Short('w').Required().IntVar(&setCommand.Weight)

	app.Flag("verbose", "Enable debug mode").Short('v').BoolVar(&params.debug)
	app.Flag("file", "Path to the Traefik configuration file").Short('f').Envar("AECTL_FILE").Required().StringVar(&params.file)
	app.PreAction(func(ctx *kingpin.ParseContext) error {
		// apply global `file` option
		listCommand.File = params.file
		setCommand.File = params.file

		if params.debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
