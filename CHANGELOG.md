# Changelog

## 2.0.0 - 2025-05-07

### Breaking changes

- `.` in flag names will now also be mapped to `_` in environment vars.
  e.g. `--foo.bar` would map to `FOO.BAR` in v1, `FOO_BAR` in v2.

## 1.0.0 - 2021-10-17

_Initial release._