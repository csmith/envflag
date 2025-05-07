package envflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Option defines a configuration option for the envflag package.
type Option func(*config)

// WithFlagSet specifies which FlagSet envflag should operate on.
// If not specified, defaults to flag.CommandLine.
func WithFlagSet(flagSet *flag.FlagSet) Option {
	return func(c *config) {
		c.flagSet = flagSet
	}
}

// WithShowInUsage specifies whether the name of environment variables should be prepended to the usage text
// of each flag. If not specified, defaults to true.
func WithShowInUsage(showInUsage bool) Option {
	return func(c *config) {
		c.showInUsage = showInUsage
	}
}

// WithPrefix alters the mapping of flag names to environment variables so that they are prefixed with the given
// string. If not specified, no prefix is used.
func WithPrefix(prefix string) Option {
	return func(c *config) {
		originalMapper := c.mapper
		c.mapper = func(s string) string {
			return fmt.Sprintf("%s%s", prefix, originalMapper(s))
		}
	}
}

// WithArguments specifies the arguments that should be parsed into flags. If a flag is specified in both an
// environment variable and in the given arguments, the one in the arguments will be used. Defaults to the
// command-line arguments (os.Args[1:]).
func WithArguments(arguments []string) Option {
	return func(c *config) {
		c.arguments = arguments
	}
}

// config encapsulates the caller-configurable options for the envflag package.
type config struct {
	flagSet     *flag.FlagSet
	showInUsage bool
	mapper      func(string) string
	exit        func(int)
	arguments   []string
}

// configFor builds a config object for the envflag package using the given options. Options that aren't specified
// will fall back to their defaults.
func configFor(opts []Option) *config {
	c := &config{
		flagSet:     flag.CommandLine,
		showInUsage: true,
		mapper:      defaultMapper,
		exit:        os.Exit,
		arguments:   os.Args[1:],
	}

	for i := range opts {
		opts[i](c)
	}

	return c
}

// defaultMapper maps flag names to environment variable names by uppercasing them and replacing dashes with
// underscores.
func defaultMapper(flagName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(flagName), "-", "_"), ".", "_")
}
