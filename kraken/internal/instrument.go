package internal

import (
	"sync"

	"github.com/simonks2016/dex_plus/kraken/payload"
)

type InstrumentService struct {
	tradingPair map[string]payload.Pair
	assets      map[string]payload.Asset
	pairMu      sync.Mutex
	assetMu     sync.Mutex
}

func NewInstrumentService() *InstrumentService {
	return &InstrumentService{
		tradingPair: make(map[string]payload.Pair),
		assets:      make(map[string]payload.Asset),
	}
}

func (s *InstrumentService) AddTradingPairs(pairs ...payload.Pair) {
	s.pairMu.Lock()
	defer s.pairMu.Unlock()

	for _, pair := range pairs {
		s.tradingPair[pair.Symbol] = pair
	}
}

func (s *InstrumentService) AddAsset(assets ...payload.Asset) {
	s.assetMu.Lock()
	defer s.assetMu.Unlock()

	for _, asset := range assets {
		s.assets[asset.Id] = asset
	}
}

func (s *InstrumentService) GetTradingPair(symbol string) (payload.Pair, bool) {
	s.pairMu.Lock()
	defer s.pairMu.Unlock()

	pair, ok := s.tradingPair[symbol]
	return pair, ok
}

func (s *InstrumentService) GetAsset(symbol string) (payload.Asset, bool) {
	s.assetMu.Lock()
	defer s.assetMu.Unlock()

	asset, ok := s.assets[symbol]
	return asset, ok
}
