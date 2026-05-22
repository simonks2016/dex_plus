package internal

import "github.com/simonks2016/dex_plus/kraken/payload"

type Caller func(envelope *payload.KrakenEnvelope) error
