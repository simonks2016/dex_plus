package kraken

import "log"

type Option func(public *Public)

func WithSymbols(symbols ...string) Option {
	return func(public *Public) {
		public.symbols = append(public.symbols, symbols...)
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(public *Public) {
		public.logger = logger
	}
}
