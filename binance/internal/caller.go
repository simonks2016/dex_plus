package internal

import "github.com/goccy/go-json"

type Caller func(string,json.RawMessage) error
