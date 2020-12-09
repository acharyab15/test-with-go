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


## 04 - Failing Tests

The `testing.T` type has a `Log` and `Logf` method. Work similar to `Print` and `Printf` in the `fmt` package.

There are two ways to signal that a test has failed:
- Fail = fail, but keep running
- FailNow = fail now and stop test

Most times, people call
- Error: Log + Fail
- Errorf: Logf + Fail
- Fatal: Log + FailNow
- Fatalf: Logf + FailNow

Which do you use?
- If you can let a test keep running, use Error/Errorf
- If a test is completely over and running further won't help at all, use Fatal/Fatalf

## 07 - Running tests in parallel

Sometimes running many tests in parallel can provide a ton of value.

Example use cases:
1. Simulating a real-world scenario
    - A web app with many users
2. Verify that a type is truly threadsafe
    - Verify that in-memory cache can handle multiple concurrent web requests using it

Parallelism could mean more work:
- Tests can't use as many hard-coded values; eg unique email constraints
- Tests might try to use shared resources incorrectly; eg image manipulation on the same image or sharing a DB that doesn't support multiple concurrent connecions

## 08 - Race conditions


## 09 - Comparing objects for equality

### Golden files

If we need to compare big files (large csv, image, etc) then trying to recreate the `want` variable can be hard.

Common solution: store a "golden file" - a file representing the desired test output - in our actual test source directory
and to just compare to it directly.

## 11 - Controlling which tests are run

- `go test -v -run TestSomething`
- `go test ./...`
- `for pkg in \*\; do go test "./$pkg"; done`
- `// +build integration (build tags)`

## 13 - External and Internal Testing

`xxx_internal_test.go` => internal tests
`xxx_test.go` => external tests
`export_test.go` => export unexported stuff for testing

## 14 - Types of tests

### Unit tests: 

Testing very small things, like a function. Usually in isolation

Example:

```go
// This is the unit - a function
func Magic(a, b int) int {
    return (a+b) * (a+b)
}

// This is the unit test
func TestMagic(t *testing.T) {
    got := Magic(1,2)
    if got != 9 {
        t.Errorf("Magic() = %v, want %v", got, 9)
    }
}
```

Very common; and require very little setup.

### Integration tests:

Testing 2+ systems together

Example:
```go
type UserStore struct {
    db *sql.DB
}
func (us *UserStore) Create(user *User) error {
    // .. this uses the us.db (the sql database) to 
    // create a new user entry from the user object passed in.
}

// Integration tests might use a REAL database, meaning it is 
// testing the integration of our UserStore with a real SQL
// DB and not some mocked out DB.
func TestUserStore_Create(t *testing.T) {
    // ..
}
```

Unit tests, especially ones with  mocks, only test that other systems works as we expect it to work.
Integration tests will verify that our expectations of how the system should work are correct.

### End-to-end tests

Testing the entire application, or most of it. 
Usually in a way similar to how end users would use the app
There can be a fuzzy line between integration and E2E tests. Typically involves using the entire system in a way similar to how end users would use it.

Pros:
Great for simulating real user scenarios.
Great for catching bugs - touches a ton of code.
Could involve *multiple systems* 

Con:
Not great at pointing at WHY bugs occured or how to fix them quickly/clearly.

## 16 - Dependency Injection


