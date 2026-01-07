package option

import (
	"log"
	"strings"

	"github.com/panjf2000/ants/v2"
)

type Option interface {
	Name() string
	Value() any
}

func GetOption(name string, opts ...Option) any {

	if len(opts) == 0 {
		return nil
	}

	for _, opt := range opts {
		if strings.ToLower(opt.Name()) == strings.ToLower(name) {
			return opt.Value()
		}
	}
	return nil
}

func WithReadBufferSize(size int64) Option  { return NewCustomOption("read_buffer_size", size) }
func WithWriteBufferSize(size int64) Option { return NewCustomOption("write_buffer_size", size) }
func WithURL(url string) Option             { return NewCustomOption("url", url) }
func WithThreadPool(pool *ants.Pool) Option { return NewCustomOption("thread_pool", pool) }
func WithForbidIpv6(ok bool) Option         { return NewCustomOption("forbid_ipv6", ok) }
func WithLogger(log *log.Logger) Option {
	return NewCustomOption("logger", log)

}
func WithSandBoxEnvironment() Option    { return NewCustomOption("is_sandbox_environment", true) }
func WithProductionEnvironment() Option { return NewCustomOption("is_sandbox_environment", false) }

type CustomOption struct {
	name  string
	value any
}

func (c *CustomOption) Name() string { return c.name }
func (c *CustomOption) Value() any   { return c.value }

func NewCustomOption(name string, val any) *CustomOption {
	return &CustomOption{
		name:  name,
		value: val,
	}
}
