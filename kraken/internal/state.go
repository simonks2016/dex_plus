package internal

import (
	"strings"
	"sync"
)

type State int

const (
	Subscribing State = iota
	Subscribed
	Resubscribing
	SubscribeFailed
	Unsubscribed
)

type SubscribeChannelState struct {
	mu       sync.RWMutex
	stateMap map[string]State
}

func NewSubscribeChannelState() *SubscribeChannelState {
	return &SubscribeChannelState{
		stateMap: make(map[string]State),
	}
}

func subscribeKey(channel, symbol string) string {
	return strings.ToLower(channel) + "_" + strings.ToUpper(symbol)
}

func (s *SubscribeChannelState) Switch(channel, symbol string, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stateMap[subscribeKey(channel, symbol)] = state
}

func (s *SubscribeChannelState) Get(channel, symbol string) (State, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.stateMap[subscribeKey(channel, symbol)]
	return state, ok
}

func (s *SubscribeChannelState) IsSubscribed(channel, symbol string) bool {
	state, ok := s.Get(channel, symbol)
	return ok && state == Subscribed
}

func (s *SubscribeChannelState) IsSubscribing(channel, symbol string) bool {
	state, ok := s.Get(channel, symbol)
	return ok && state == Subscribing
}

func (s *SubscribeChannelState) Remove(channel, symbol string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.stateMap, subscribeKey(channel, symbol))
}
