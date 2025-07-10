# Allure

An [Allure](https://allurereport.org/) plugin for `testo`.

Take a look at [the example](./examples/simple).

## Steps and sub-tests

Allure plugin provides step abstraction.

Both, sub-tests and steps are shown in allure report as steps under parent test.
However, `allure.Step` propagates fatal errors to the parent.
Fatal errors are triggered by the `t.FailNow()` function, commonly called from `t.Fatal`.

Take a look at the example:

```go
func (Suite) TestStep(t T) {
    // trigger fatal error
    allure.Step(t, "first", func(t T) { t.FailNow() })

    // ❌ this code won't be executed
    t.Log("Hi")
}

func (Suite) TestRun(t T) {
    // trigger fatal error
    testo.Run(t, "first", func(t T) { t.FailNow() })

    // ✅ this code will be executed
    t.Log("Hi")
}
```
