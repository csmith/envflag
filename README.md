# Yet another `envflag` library

[![PkgGoDev](https://pkg.go.dev/badge/github.com/csmith/envflag)](https://pkg.go.dev/github.com/csmith/envflag)

This is a simple library that allows that enables reading the value of flags from environment variables.

## Why?

The Go flags package is well-thought-out, battle-tested, and extensively used, but modern
[deployment practices](https://12factor.net/config) tend to favour storing configuration in environment variables rather
than passing them in on the command-line.

`envflag` and its predecessors bridge this gap by mapping environment variables onto flags.

## How?

Import this package, and use `envflag.Parse` instead of `flag.Parse`:

```go
package main

import (
	"flag"

	"github.com/csmith/envflag"
)

var (
	myFlag = flag.String("my-flag", "woohoo", "Something or other")
)

func main() {
	envflag.Parse()

	println(*myFlag)
}
```

## What?

In its default configuration, `envflag` will:

* Map all flag names to appropriate environment variable names (`my-flag` -> `MY_FLAG`)
* Update the usage information of flags to include the environment variable name
* Show usage and exit with an appropriate status code on error
* Use the default set of flags provided by `flag.CommandLine`
* Call `flag.Parse` to also parse any command-line flags

You can customise the behaviour of `envflag` by passing in options:

### Prefix all environment variables

```go
envflag.Parse(envflag.WithPrefix("MYAPP_"))
```

Changes the mapping of environment variables to always include the given prefix
(e.g. `my-flag` -> `MYAPP_MY_FLAG`).

### Use a different flag set

```go
envflag.Parse(envflag.WithFlagSet(someFlagSet))
```

Use the given `flag.FlagSet` instead of `flag.CommandLine`.

### Don't show environment variable names in usage

```go
envflag.Parse(envflag.WithShowInUsage(false))
```

Suppresses the default behaviour of updating flag usage to include the environment variable names.

### Use a different set of arguments

```go
envflag.Parse(envflag.WithArguments([]string{"-option1", "-option2"}))
```

When parsing the command-line arguments, use the given slice instead of `os.Args[1:]`.

## Licence/credits/contributions etc

Released under the MIT licence. See LICENCE for full details.

Heavily inspired by [kouhin/envflag](https://github.com/kouhin/envflag) which does
mostly the same job but can't easily add prefixes and is a bit harder to configure.

Contributions are welcome! Please feel free to open issues or send pull requests.
