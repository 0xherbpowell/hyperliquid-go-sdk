package types

import (
	"fmt"
	"strconv"
	"strings"

	"hyperliquid-go-sdk/pkg/errors"
)

// Cloid represents a client order ID (16 bytes hex string)
type Cloid struct {
	raw string
}

// NewCloidFromString creates a Cloid from a hex string
func NewCloidFromString(cloid string) (*Cloid, error) {
	c := &Cloid{raw: cloid}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// NewCloidFromInt creates a Cloid from an integer
func NewCloidFromInt(cloid int64) *Cloid {
	return &Cloid{raw: fmt.Sprintf("0x%032x", cloid)}
}

// validate checks if the cloid is a valid 16-byte hex string
func (c *Cloid) validate() error {
	if !strings.HasPrefix(c.raw, "0x") {
		return errors.ErrInvalidCloid
	}
	if len(c.raw[2:]) != 32 {
		return errors.ErrInvalidCloid
	}
	// Check if it's valid hex
	if _, err := strconv.ParseInt(c.raw[2:], 16, 64); err != nil {
		return errors.ErrInvalidCloid
	}
	return nil
}

// String returns the string representation of the cloid
func (c *Cloid) String() string {
	return c.raw
}

// Raw returns the raw hex string
func (c *Cloid) Raw() string {
	return c.raw
}

// MarshalJSON implements json.Marshaler
func (c *Cloid) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.raw + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (c *Cloid) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	cloid, err := NewCloidFromString(str)
	if err != nil {
		return err
	}
	c.raw = cloid.raw
	return nil
}
