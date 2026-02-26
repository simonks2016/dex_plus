package internal

type Caller func(envelope *KrakenEnvelope) error
