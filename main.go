package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// DotSeparatedStringList is a string list with a mapstructure decode hook
// that decodes a dot separated string list.
type DotSeparatedStringList []string

// DotSeparatedStringListHookFunc returns a DecodeHookFunc that converts
// strings to string slices, when the target type is DotSeparatedStringList.
func DotSeparatedStringListHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(DotSeparatedStringList{}) {
			return data, nil
		}

		return DotSeparatedStringList(strings.Split(data.(string), ".")), nil
	}
}

// SemicolonSeparatedStringList is a string list that implements
// encoding.TextUnmarshaler and decodes a semicolon separated string list.
type SemicolonSeparatedStringList []string

func (s *SemicolonSeparatedStringList) UnmarshalText(text []byte) error {
	*s = strings.Split(string(text), ";")

	return nil
}

// Config is a struct showcasing all examples.
type Config struct {
	Comma     []string
	Dot       DotSeparatedStringList
	Semicolon SemicolonSeparatedStringList
}

func main() {
	// Custom decode hook function
	{
		v := viper.New()
		v.Set("key", "foo.bar.baz.bat")

		var s DotSeparatedStringList

		v.UnmarshalKey("key", &s, viper.DecodeHook(DotSeparatedStringListHookFunc()))

		fmt.Printf("Dot separated list (DotSeparatedStringListHookFunc): %#v\n", s)
	}

	// TextUnmarshaller decode hook func
	{
		v := viper.New()
		v.Set("key", "foo;bar;baz;bat")

		var s SemicolonSeparatedStringList

		v.UnmarshalKey("key", &s, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))

		fmt.Printf("Semicolon separated list (TextUnmarshallerHookFunc): %#v\n", s)
	}

	// Builtin decode hook function
	{
		v := viper.New()
		v.Set("key", "foo,bar,baz,bat")

		var s []string

		v.UnmarshalKey("key", &s)

		fmt.Printf("Comma separated list (builtin decode hook function): %#v\n", s)
	}

	// All in one: Config struct
	{
		v := viper.New()
		v.Set("comma", "foo,bar,baz,bat")
		v.Set("dot", "foo.bar.baz.bat")
		v.Set("semicolon", "foo;bar;baz;bat")

		var config Config

		v.Unmarshal(&config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			DotSeparatedStringListHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(), // default hook
			mapstructure.StringToSliceHookFunc(","),     // default hook
		)))

		fmt.Printf("All in one config struct: %#v\n", config)
	}
}
