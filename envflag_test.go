package envflag

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapsEnvironmentVariables(t *testing.T) {
	var (
		set        = flag.NewFlagSet("test", flag.ExitOnError)
		stringFlag = set.String("my-string", "foo", "Some string flag")
		boolFlag   = set.Bool("my-bool", true, "Some bool flag")
		unsetFlag  = set.String("other-string", "baz", "A third flag")
	)

	os.Setenv("MY_STRING", "bar")
	os.Setenv("MY_BOOL", "false")

	Parse(WithFlagSet(set), WithArguments([]string{}))

	assert.Equal(t, "bar", *stringFlag)
	assert.Equal(t, false, *boolFlag)
	assert.Equal(t, "baz", *unsetFlag)
}

func TestParsesAndPrefersArguments(t *testing.T) {
	var (
		set        = flag.NewFlagSet("test", flag.ExitOnError)
		stringFlag = set.String("my-string", "foo", "Some string flag")
		boolFlag   = set.Bool("my-bool", true, "Some bool flag")
		unsetFlag  = set.String("other-string", "baz", "A third flag")
	)

	os.Setenv("MY_STRING", "bar")
	os.Setenv("MY_BOOL", "false")

	Parse(WithFlagSet(set), WithArguments([]string{"-my-string", "qux", "-other-string=foo"}))

	assert.Equal(t, "qux", *stringFlag)
	assert.Equal(t, false, *boolFlag)
	assert.Equal(t, "foo", *unsetFlag)
}

func TestWithPrefix(t *testing.T) {
	var (
		set        = flag.NewFlagSet("test", flag.ExitOnError)
		stringFlag = set.String("my-string", "foo", "Some string flag")
		boolFlag   = set.Bool("my-bool", true, "Some bool flag")
		unsetFlag  = set.String("other-string", "baz", "A third flag")
	)

	os.Setenv("MY_STRING", "bar2")
	os.Setenv("MY_BOOL", "true")
	os.Setenv("PREFIX_MY_STRING", "bar")
	os.Setenv("PREFIX_MY_BOOL", "false")

	Parse(WithFlagSet(set), WithArguments([]string{}), WithPrefix("PREFIX_"))

	assert.Equal(t, "bar", *stringFlag)
	assert.Equal(t, false, *boolFlag)
	assert.Equal(t, "baz", *unsetFlag)
}

func TestExitsStatusZeroForHelp(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
	)

	exitCode := -1
	exit := func(code int) {
		if exitCode == -1 {
			exitCode = code
		}
	}

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{"-help"}), withExit(exit))

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, buf.String(), "Usage of test:")
}

func TestExitsStatusTwoForBadFlags(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
	)

	exitCode := -1
	exit := func(code int) {
		if exitCode == -1 {
			exitCode = code
		}
	}

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{"-bazinga"}), withExit(exit))

	assert.Equal(t, 2, exitCode)
	assert.Contains(t, buf.String(), "flag provided but not defined")
	assert.Contains(t, buf.String(), "Usage of test:")
}

func TestExitsStatusTwoForBadEnvVars(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
		_   = set.Int("my-int", 0, "Some int flag")
	)

	exitCode := -1
	exit := func(code int) {
		if exitCode == -1 {
			exitCode = code
		}
	}

	os.Setenv("MY_INT", "seven")

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{}), withExit(exit))

	assert.Equal(t, 2, exitCode)
	assert.Contains(t, buf.String(), "Unable to set flag my-int value to seven from environment:")
	assert.Contains(t, buf.String(), "Usage of test:")
}

func TestShowsEnvNamesInUsage(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
		_ = set.String("test-flag-plz-ignore", "foo", "Some usage")
	)

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{"-help"}), withExit(func(code int) {}))

	assert.Contains(t, buf.String(), `[TEST_FLAG_PLZ_IGNORE] Some usage (default "foo")`)
}

func TestShowsPrefixedEnvNamesInUsage(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
		_ = set.String("test-flag-plz-ignore", "foo", "Some usage")
	)

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{"-help"}), WithPrefix("BOB_"), withExit(func(code int) {}))

	assert.Contains(t, buf.String(), `[BOB_TEST_FLAG_PLZ_IGNORE] Some usage (default "foo")`)
}

func TestDoesntUpdateUsageIfShowInUsageDisabled(t *testing.T) {
	var (
		set = flag.NewFlagSet("test", flag.ContinueOnError)
		_ = set.String("test-flag-plz-ignore", "foo", "Some usage")
	)

	buf := &bytes.Buffer{}
	set.SetOutput(buf)
	Parse(WithFlagSet(set), WithArguments([]string{"-help"}), WithShowInUsage(false), withExit(func(code int) {}))

	assert.NotContains(t, buf.String(), `TEST_FLAG_PLZ_IGNORE`)
}

func withExit(exit func(int)) Option {
	return func(c *config) {
		c.exit = exit
	}
}
