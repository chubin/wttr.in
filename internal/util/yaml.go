package util

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// YamlUnmarshalStrict unmarshals YAML data with an error when unknown fields are present.
func YamlUnmarshalStrict(in []byte, out interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(in))
	dec.KnownFields(true)
	return dec.Decode(out)
}
