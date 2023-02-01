package discordfx

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// Guild represents a Guild ID or Guild URL in the config.
type Guild string

func (g *Guild) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// If it's numeric, just return it.
	if _, err := strconv.Atoi(s); err == nil {
		return nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	urlParts := strings.Split(u.Path, "/")
	*g = Guild(urlParts[len(urlParts)-1])

	// The result must be numeric.
	if _, err := strconv.Atoi(s); err != nil {
		return err
	}

	return nil
}

// Channel represents a Channel ID or Channel URL in the config.
type Channel string

func (c *Channel) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// If it's numeric, just return it.
	if _, err := strconv.Atoi(s); err == nil {
		return nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	urlParts := strings.Split(u.Path, "/")
	*c = Channel(urlParts[len(urlParts)-1])

	// The result must be numeric.
	if _, err := strconv.Atoi(s); err != nil {
		return err
	}

	return nil
}
