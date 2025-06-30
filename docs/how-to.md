# How-to

Learn how to use the features of `testo`.

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

## How to inherit `T`

When writing tests, you may want to add some plugins:

```go
package common

type T struct {
    *testo.T

    *myplugin.MyPlugin
    *otherplugin.OtherPlugin
}
```

It may useful to define some base `T` once and inherit it in when you need to extend it.

You can do you it like that:

```go
type T struct {
    *common.T

    *extraplugin.ExtraPlugin
}
```

Testo understands this pattern and handles it as you would expect - plugins
of parent `T` are registered along with plugins of the inherited `T`.
