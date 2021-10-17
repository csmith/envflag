/*
Package envflag facilitates setting standard Go flags using environment variables.

Basic usage

Simply replace the typical call to flag.Parse() with envflag.Parse(). Flag names are mapped to environment variables
by converting them to uppercase, and replacing dashes with underscores (e.g. `my-flag` => `MY_FLAG`). Command-line
arguments will take precedence over environment variables where both are specified.

Advanced usage

You can customise the behaviour of envflag by passing in options to the Parse() method. For example, to add a prefix
to all environment variables:

    envflag.Parse(envflag.WithPrefix("MYAPP_"))

This will map a flag named `my-flag` to the environment variable `MYAPP_MY_FLAG`.
*/
package envflag
