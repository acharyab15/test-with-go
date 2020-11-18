# test-with-go

## What is a test?

> A test is a repeatable process that verifies whether or not something is working as intended.

## Why do tests matter?

1. Tests help find, fix, and prevent mistakes (bug, side effects, edge cases, etc)
2. Tests document expected behavior
3. Tests encourage us to write better code -- break down into smaller units
4. Tests can speed up development -- confidence, new developers, changes to code etc

## File naming conventions

Caveats to general file naming conventions:

- `export_test.go`: to access unexported variables in external tests
- `xxx_internal_test.go` for internal tests
- `example_xxx_test.go` for example in isolated files

## Variable naming conventions

Use `got` and `want` for variables that was received from a function call and that was expected, respectively.
