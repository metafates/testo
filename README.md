# Testo

<img src="https://github.com/user-attachments/assets/66844de4-4b13-428a-b924-1f26718cee41" align="right" width="250" alt="testo mascot">

Testo is a modular testing framework for Go built on top of `testing.T`.
It is focused on suite based tests and has an extensive plugin system.

> Work in progress.

## Get started

Your project must use go 1.22 or newer.

No stable version _yet_ - the following command installs the version from the latest commit:

```bash
go get github.com/metafates/testo@main
```

## Next steps

- Take [a guided tour of `testo`](https://github.com/metafates/testo/tree/main/docs/tutorial.md) by making a simple plugins and running the tests using various features.
- Learn [how to use various `testo` features](https://github.com/metafates/testo/tree/main/docs/how-to.md).
- Read a [brief description and technical overview](https://github.com/metafates/testo/tree/main/docs/technical-overview.md) of `testo`.
- View [API documentation](https://pkg.go.dev/github.com/metafates/testo).
- See [Allure plugin](./pkg/plugins/allure) as an example.

## To do

- [ ] Mock generation.
- [ ] Stabilize API.
- [ ] Interface generator CLI. Similar to [ifacemaker] but simplified for project needs. The goal is to simplify plugin development.
- [ ] Move Allure plugin into separate repository

[ifacemaker]: https://github.com/vburenin/ifacemaker
