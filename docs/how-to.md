# How-to

How to use the features of `testo`.

## How to write parallel tests

You can use your regular `t.Parallel` method to mark test as parallel.

```go
func (Suite) TestFoo(t *testo.T) {
    t.Parallel()

    // your test here
}
```

You can expect all `AfterEach` and `AfterAll` hooks to execute at the end of each test properly.

The only limitation here is that top-level sub-tests can't be parallel.
Consider this example:

```go
func (Suite) TestFoo(t *testo.T) {
    // this is ok and will work as expected
    t.Parallel()

    testo.Run(t, "top-level subtest", func(t *testo.T) {
        // this is not supported.
        // you can call Parallel here, but it will become
        // a no-op with a warning in logs.
        t.Parallel()

        testo.Run(t, "nested subtest", func(t *testo.T) {
            // this is ok and will work as expected
            t.Parallel()
        })
    })
}
```
