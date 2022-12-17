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

var params parameters

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := kingpin.New("aectl", "Traefik Weight Control")
	app.HelpFlag.Short('h')

	//app.Flag("file", "Path to the Traefik configuration file").Short('f').Required().StringVar(&params.file)
	app.Flag("verbose", "Enable debug mode").Short('v').BoolVar(&params.debug)

	listCommand := &commands.ListCommand{}
	kpListCommand := app.Command("list", "List all services and servers").Action(listCommand.Execute)
	kpListCommand.Flag("file", "Path to the Traefik configuration file").Short('f').Required().StringVar(&listCommand.File)
	kpListCommand.Flag("format", "Output format").Short('o').Default("text").EnumVar(&listCommand.Format, "text", "json", "yaml", "yml")

	kingpin.MustParse(app.Parse(os.Args[1:]))
	if params.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
