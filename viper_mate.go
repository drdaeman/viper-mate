package vipermate

import (
	"errors"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/gogap/logrus_mate"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/gogap/config"
	"fmt"
)

// provider is a Logrus Mate-compatible ConfigurationProvider
// that uses Viper, as well as the Configuration struct.
// It is intentionally not exported, as the path between
// the provider and Configuration is a hack. Just use NewMate.
// See some comments near the ParseConfig for more information.
type provider struct {
	V *viper.Viper
}

// Copied from go-akka/configuration sources (why they didn't export it?)
func splitDottedPathHonouringQuotes(path string) []string {
	tmp1 := strings.Split(path, "\"")
	var values []string
	for i := 0; i < len(tmp1); i++ {
		tmp2 := strings.Split(tmp1[i], ".")
		for j := 0; j < len(tmp2); j++ {
			if len(tmp2[j]) > 0 {
				values = append(values, tmp2[j])
			}
		}
	}
	return values
}

func (cfg *provider) getSub(path []string) *viper.Viper {
	v := cfg.V
	for _, s := range path {
		v = v.Sub(s)
		if v == nil {
			break
		}
	}
	return v
}

func (cfg *provider) getPath(path string) interface{} {
	p := splitDottedPathHonouringQuotes(path)
	if len(p) < 1 {
		return cfg.V
	}
	var v *viper.Viper
	if len(p) > 1 {
		v = cfg.getSub(p[:len(p)-1])
	} else {
		v = cfg.V
	}
	if v == nil {
		return nil
	}
	return v.Get(p[len(p)-1])
}

func (cfg *provider) GetBoolean(path string, defaultVal ...bool) bool {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case bool:
		return val
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return false
	}
}

func (cfg *provider) GetByteSize(path string) *big.Int {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case int64, int32, int16, int8, int:
		return big.NewInt(cast.ToInt64(val))
	default:
		return big.NewInt(0)
	}
}

func (cfg *provider) GetInt32(path string, defaultVal ...int32) int32 {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case int64, int32, int16, int8, int:
		return cast.ToInt32(val)
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return 0
	}
}

func (cfg *provider) GetInt64(path string, defaultVal ...int64) int64 {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case int64, int32, int16, int8, int:
		return cast.ToInt64(val)
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return 0
	}
}

func (cfg *provider) GetString(path string, defaultVal ...string) string {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case string:
		return val
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return ""
	}
}

func (cfg *provider) GetFloat32(path string, defaultVal ...float32) float32 {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case float64, float32:
		return cast.ToFloat32(val)
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return 0
	}
}

func (cfg *provider) GetFloat64(path string, defaultVal ...float64) float64 {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case float64, float32:
		return cast.ToFloat64(val)
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return 0
	}
}

func (cfg *provider) GetTimeDuration(path string,
	defaultVal ...time.Duration) time.Duration {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case time.Duration:
		return val
	default:
		if len(defaultVal) >= 1 {
			return defaultVal[0]
		}
		return 0
	}
}

func (cfg *provider) GetTimeDurationInfiniteNotAllowed(path string,
	defaultVal ...time.Duration) time.Duration {
	duration := cfg.GetTimeDuration(path, defaultVal...)
	if duration == time.Duration(-1) {
		panic("infinite time duration not allowed")
	}
	return duration
}

func (cfg *provider) GetStringList(path string) []string {
	v := cfg.getPath(path)
	switch val := v.(type) {
	case []string:
		return val
	}
	return nil
}

func (cfg *provider) IsEmpty() bool {
	return len(cfg.V.AllKeys()) == 0
}

func (cfg *provider) String() string {
	return fmt.Sprintf("%#v", cfg.V.AllSettings())
}

// Non-string list functions are not implemented in Viper
func (cfg *provider) GetBooleanList(path string) []bool {
	return []bool{}
}
func (cfg *provider) GetFloat32List(path string) []float32 {
	return []float32{}
}
func (cfg *provider) GetFloat64List(path string) []float64 {
	return []float64{}
}
func (cfg *provider) GetInt32List(path string) []int32 {
	return []int32{}
}
func (cfg *provider) GetInt64List(path string) []int64 {
	return []int64{}
}
func (cfg *provider) GetByteList(path string) []byte {
	return []byte{}
}

func (cfg *provider) GetConfig(path string) config.Configuration {
	sub := cfg.getSub(splitDottedPathHonouringQuotes(path))
	if sub != nil {
		return &provider{V: sub}
	}
	return nil
}

func (cfg *provider) WithFallback(fallback config.Configuration) config.Configuration {
	log.Panic("viperConfigProvider.WithFallback is not implemented")
	return nil
}

func (cfg *provider) HasPath(path string) bool {
	p := splitDottedPathHonouringQuotes(path)
	if len(p) < 1 {
		return false
	}
	var v *viper.Viper
	if len(p) > 1 {
		v = cfg.getSub(p[:len(p)-1])
	} else {
		v = cfg.V
	}
	if v == nil {
		return false
	}
	return v.IsSet(p[len(p)-1])
}

func (cfg *provider) Keys() []string {
	// XXX: This approach is ugly, but AllKeys is not what we want, too.
	// Let's hope passed configuration objects are small enough, so the
	// overhead of AllSettings is negligible.
	a := cfg.V.AllSettings()
	keys := make([]string, len(a))
	i := 0
	for k := range a {
		keys[i] = k
		i++
	}
	return keys
}

// Logrus Mate assumes ConfigurationProviders load configuration by themselves,
// either from string or file. With Viper, this assumption is not true.
// So, we implement a ConfigurationProvider stub that doesn't do anything but
// expects ParseString with a dummy value and returns a Configuration that
// uses pre-provided Viper instance.
// To not expose those dirty hacks to users we hide it in the NewViperMate

func (cfg *provider) ParseString(cfgStr string) config.Configuration {
	return cfg
}

func (cfg *provider) LoadConfig(filename string) config.Configuration {
	log.Panic("LoadConfig is not implemented and not supposed to be called")
	return nil
}

// NewMate returns Logrus Mate instance configured using Viper
func NewMate(cfg *viper.Viper) (*logrus_mate.LogrusMate, error) {
	if cfg == nil {
		return nil, errors.New("NewMate got a nil Viper reference")
	}
	return logrus_mate.NewLogrusMate(
		logrus_mate.ConfigString("/* viper */"), // Hack, see notes above
		logrus_mate.ConfigProvider(&provider{V: cfg}),
	)
}
