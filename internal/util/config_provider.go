package util

import (
	"encoding/json"
	"fmt"
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
	FileName   string
	fileHandle *os.File
}

func (c *ConfigFile) selectProvider(strict bool) (Decoder, Encoder, error) {
	fh, osOpenError := os.OpenFile(c.FileName, os.O_RDWR|os.O_CREATE, 0644)

	if osOpenError != nil {
		return nil, nil, osOpenError
	}
	c.fileHandle = fh

	switch path.Ext(c.FileName) {
	case ".yaml", ".yml":
		dec := yaml.NewDecoder(c.fileHandle)
		if strict {
			dec.KnownFields(true)
		}
		return dec, yaml.NewEncoder(c.fileHandle), nil
	case ".json":
		dec := json.NewDecoder(c.fileHandle)
		if strict {
			dec.DisallowUnknownFields()
		}
		return dec, json.NewEncoder(c.fileHandle), nil
	}

	return nil, nil, fmt.Errorf("unsupported file extension")
}

func (c *ConfigFile) Read() (*dynamic.Configuration, error) {
	decoder, _, providerError := c.selectProvider(false)
	if providerError != nil {
		return nil, providerError
	}
	defer func() {
		c.fileHandle.Close()
	}()

	log.Debug().
		Str("file", c.FileName).
		Str("provider", fmt.Sprintf("%T", decoder)).
		Msg("Reading configuration file")

	config := dynamic.Configuration{}
	decodeError := decoder.Decode(&config)
	return &config, decodeError
}

func (c *ConfigFile) Write(config *dynamic.Configuration) error {
	_, encoder, providerError := c.selectProvider(false)
	if providerError != nil {
		return providerError
	}
	defer func() {
		c.fileHandle.Close()
	}()

	log.Debug().
		Str("file", c.FileName).
		Str("provider", fmt.Sprintf("%T", encoder)).
		Msg("Writing configuration file")

	encodeError := encoder.Encode(config)
	return encodeError
}
