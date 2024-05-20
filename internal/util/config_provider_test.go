package util

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
)

func TestSelectProvider(t *testing.T) {
	t.Run("SuccessYaml", func(t *testing.T) {
		// Create a test file handle
		fileHandle := bytes.NewBuffer([]byte{})

		// Create a ConfigFile with a .yaml file extension
		configFile := &ConfigFile{FileName: "test.yaml"}

		// Call the selectProvider function
		decoder, encoder, err := configFile.SelectProvider(fileHandle)

		// Assert that no error is returned
		assert.Nil(t, err)

		// Assert that the returned decoder and encoder are of the correct types
		assert.IsType(t, &yaml.Decoder{}, decoder)
		assert.IsType(t, &yaml.Encoder{}, encoder)
	})

	t.Run("SuccessJson", func(t *testing.T) {
		// Create a test file handle
		fileHandle := bytes.NewBuffer([]byte{})

		// Create a ConfigFile with a .json file extension
		configFile := &ConfigFile{FileName: "test.json"}

		// Call the selectProvider function
		decoder, encoder, err := configFile.SelectProvider(fileHandle)

		// Assert that no error is returned
		assert.Nil(t, err)

		// Assert that the returned decoder and encoder are of the correct types
		assert.IsType(t, &json.Decoder{}, decoder)
		assert.IsType(t, &json.Encoder{}, encoder)
	})

	t.Run("ErrorUnsupportedExtension", func(t *testing.T) {
		// Create a test file handle
		fileHandle := bytes.NewBuffer([]byte{})

		// Create a ConfigFile with an unsupported file extension
		configFile := &ConfigFile{FileName: "test.csv"}

		// Call the selectProvider function
		decoder, encoder, err := configFile.SelectProvider(fileHandle)

		// Assert that an error is returned
		assert.NotNil(t, err)
		assert.Equal(t, "unsupported file extension", err.Error())

		// Assert that the returned decoder and encoder are nil
		assert.Nil(t, decoder)
		assert.Nil(t, encoder)
	})
}

func TestConfigFile_Read(t *testing.T) {
	// Test data
	for _, testCase := range []struct {
		name     string
		filename string
		expected *dynamic.Configuration
		err      string
	}{
		{
			name:     "SuccessYaml",
			filename: "/tmp/test.yaml",
			expected: &dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"foo": {
							Rule: "bar",
						},
					},
				},
			},
			err: "",
		},
		{
			name:     "SuccessJson",
			filename: "/tmp/test.json",
			expected: &dynamic.Configuration{
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"foo": {
							Rule: "bar",
						},
					},
				},
			},
			err: "",
		},
		{
			name:     "ErrorUnsupportedExtension",
			filename: "/tmp/test.csv",
			expected: nil,
			err:      "unsupported file extension",
		},
	} {

		t.Run(testCase.name, func(t *testing.T) {
			// Create a test ConfigFile
			configFile := &ConfigFile{FileName: testCase.filename}

			// Create a test file handle
			err := configFile.Write(testCase.expected)
			if testCase.err == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, testCase.err, err.Error())
			}

			// Call the Read method
			config, err := configFile.Read()

			// Assert the error is as expected
			if testCase.err == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, testCase.err, err.Error())
			}

			// Assert that the returned config is as expected
			assert.Equal(t, testCase.expected, config)
		})
	}
}
