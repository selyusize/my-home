package git

import (
	"fmt"
	"sort"
	"sync"
)

// Option configures a Factory.
type Option func(*Factory)

// WithProvider registers a constructor for the given provider.
func WithProvider(provider Provider, ctor Constructor) Option {
	return func(f *Factory) {
		f.Register(provider, ctor)
	}
}

// Factory creates git.Client instances by provider name.
type Factory struct {
	mu           sync.RWMutex
	constructors map[Provider]Constructor
}

// NewFactory builds a factory and applies optional provider registrations.
func NewFactory(opts ...Option) *Factory {
	f := &Factory{constructors: make(map[Provider]Constructor)}
	for _, opt := range opts {
		if opt != nil {
			opt(f)
		}
	}
	return f
}

// Register adds or replaces a git client constructor.
func (f *Factory) Register(provider Provider, ctor Constructor) {
	if f == nil || ctor == nil || provider == "" {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.constructors[provider] = ctor
}

// New creates a client for the given provider.
func (f *Factory) New(provider Provider, cfg Config) (Client, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknownProvider, provider)
	}

	f.mu.RLock()
	ctor, ok := f.constructors[provider]
	f.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownProvider, provider)
	}

	return ctor(cfg)
}

// Has reports whether a provider is registered.
func (f *Factory) Has(provider Provider) bool {
	if f == nil {
		return false
	}

	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.constructors[provider]
	return ok
}

// Providers returns registered provider names in sorted order.
func (f *Factory) Providers() []Provider {
	if f == nil {
		return nil
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	out := make([]Provider, 0, len(f.constructors))
	for provider := range f.constructors {
		out = append(out, provider)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
