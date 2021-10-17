package envflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Parse parses flag values from the environment and the arguments list. Must be called after all flags are defined
// and before flags are accessed by the program.
//
// The parsing behaviour can be customised by passing one or more options, but will use sensible defaults if no
// options are passed.
//
// If an error occurs trying to set a flag, error and usage information will be printed to the flag.FlagSet's Output
// and os.Exit(2) will be called. This is the same behaviour as flag.ExitOnError, which is the default behaviour for
// the flag.CommandLine set.
//
// NOTE: This method will call flag.Parse(). Callers do not need to call it directly.
func Parse(opts ...Option) {
	c := configFor(opts)

	if c.showInUsage {
		updateUsage(c)
	}

	setValuesFromEnv(c)
	if err := c.flagSet.Parse(c.arguments); err != nil {
		if err == flag.ErrHelp {
			c.exit(0)
		}
		c.exit(2)
	}
}

// updateUsage updates the help message for each flag to include its corresponding environment variable.
func updateUsage(c *config) {
	c.flagSet.VisitAll(func(f *flag.Flag) {
		prefix := fmt.Sprintf("[%s] ", c.mapper(f.Name))
		if !strings.HasPrefix(f.Usage, prefix) {
			f.Usage = fmt.Sprintf("%s%s", prefix, f.Usage)
		}
	})
}

// setValuesFromEnv iterates over each flag and sets it with the corresponding environment variable, if it exists.
// If the flag cannot be set (e.g. the string value cannot be converted to the appropriate data type) then a usage
// message is displayed and the application will exit. This is the same behaviour as a flag.FlagSet with error handling
// set to flag.ExitOnError.
func setValuesFromEnv(c *config) {
	c.flagSet.VisitAll(func(f *flag.Flag) {
		varName := c.mapper(f.Name)
		if value, ok := os.LookupEnv(varName); ok {
			if err := c.flagSet.Set(f.Name, value); err != nil {
				_, _ = fmt.Fprintf(c.flagSet.Output(), "Unable to set flag %s value to %s from environment: %v\n", f.Name, value, err)
				c.flagSet.Usage()
				c.exit(2)
			}
		}
	})
}
