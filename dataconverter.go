package dataconverter

import (
	"encoding/json"
	"errors"
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

// ErrUnableToDecode is returned when unable to decode.
var (
	// ErrUnableToEncode is returned when unable to encode.
	ErrUnableToEncode = errors.New("unable to encode")
	// ErrUnableToDecode is returned when unable to decode.
	ErrUnableToDecode = errors.New("unable to decode")
)

const (
	// MetadataEncoding is "encoding"
	MetadataEncoding = "encoding"
	// MetadataEncodingJSON is "json/plain"
	MetadataEncodingJSON = "json/plain"
)

// JSONPayloadConverter converts to/from JSON.
type JSONPayloadConverter struct {
}

// NewJSONPayloadConverter creates a new instance of JSONPayloadConverter.
func NewJSONPayloadConverter() *JSONPayloadConverter {
	return &JSONPayloadConverter{}
}

// ToPayload converts a single value to a payload.
func (c *JSONPayloadConverter) ToPayload(value interface{}) (*commonpb.Payload, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnableToEncode, err)
	}
	return newPayload(data, c), nil
}

// FromPayload converts a single payload to a value.
func (c *JSONPayloadConverter) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	err := json.Unmarshal(payload.GetData(), valuePtr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnableToDecode, err)
	}
	return nil
}

// ToString converts a payload object into a human-readable string.
func (c *JSONPayloadConverter) ToString(payload *commonpb.Payload) string {
	return string(payload.GetData())
}

// Encoding returns MetadataEncodingJSON.
func (c *JSONPayloadConverter) Encoding() string {
	return MetadataEncodingJSON
}

func newPayload(data []byte, c converter.PayloadConverter) *commonpb.Payload {
	return &commonpb.Payload{
		Metadata: map[string][]byte{
			MetadataEncoding: []byte(c.Encoding()),
		},
		Data: data,
	}
}
