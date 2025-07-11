<h1>
    <img width="50" alt="Allure logo" src="https://github.com/user-attachments/assets/1b83a23a-3045-4d30-9336-97541b28923e"/>
    Allure
</h1>

An [Allure report](https://allurereport.org/) plugin for `testo`.

Take a look at [the example](./examples/simple) as a starting point.

> Allure Report is a popular open source tool for visualizing the results of a test run.
> It can be added to your testing workflow with little to zero configuration.
> It produces reports that can be opened anywhere and can be read by anyone, no deep technical knowledge required.

## Asserts

[Testify]-based asserts are available with `allure.Assert` and `allure.Require` functions.

Each assertion call is reflected in the allure report as steps with parameters.

For example, the following code:

```go
allure.Require(t).Equal(4, 2+2)
allure.Assert(t).True(false)
```

Is converted to the following steps:

```txt
require: equal
    expected: 4
    actual:   4

assert: true
    value: false
```

## Steps and sub-tests

Allure plugin provides step abstraction.

Both, sub-tests and steps are shown in allure report as steps under parent test.
However, `allure.Step` propagates fatal errors to the parent.
Fatal errors are triggered by the `t.FailNow()` function, commonly called from `t.Fatal`.

Example:

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

[Testify]: https://github.com/stretchr/testify
