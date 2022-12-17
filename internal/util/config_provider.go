package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
)

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

type ConfigFile struct {
	FileName string
}

func (c *ConfigFile) SelectProvider(fileHandle io.ReadWriter) (Decoder, Encoder, error) {
	switch path.Ext(c.FileName) {
	case ".yaml", ".yml":
		return yaml.NewDecoder(fileHandle), yaml.NewEncoder(fileHandle), nil
	case ".json":
		return json.NewDecoder(fileHandle), json.NewEncoder(fileHandle), nil
	}

	return nil, nil, fmt.Errorf("unsupported file extension")
}

func (c *ConfigFile) Read() (*dynamic.Configuration, error) {
	fh, osOpenError := os.OpenFile(c.FileName, os.O_RDONLY, 0644)
	if osOpenError != nil {
		return nil, osOpenError
	}
	defer func() { fh.Close() }()

	decoder, _, providerError := c.SelectProvider(fh)
	if providerError != nil {
		return nil, providerError
	}

	log.Debug().
		Str("file", c.FileName).
		Str("provider", fmt.Sprintf("%T", decoder)).
		Msg("Reading configuration file")

	config := dynamic.Configuration{}
	decodeError := decoder.Decode(&config)
	return &config, decodeError
}

func (c *ConfigFile) Write(config *dynamic.Configuration) error {
	fh, osOpenError := os.OpenFile(c.FileName, os.O_WRONLY|os.O_CREATE, 0644)
	if osOpenError != nil {
		return osOpenError
	}
	defer func() { fh.Close() }()

	_, encoder, providerError := c.SelectProvider(fh)
	if providerError != nil {
		return providerError
	}

	log.Debug().
		Str("file", c.FileName).
		Str("provider", fmt.Sprintf("%T", encoder)).
		Msg("Writing configuration file")

	encodeError := encoder.Encode(config)
	return encodeError
}
