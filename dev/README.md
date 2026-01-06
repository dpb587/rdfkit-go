# Developer Notes

Configure a [Go Workspace](https://go.dev/ref/mod#workspaces) to ensure all modules are kept in sync. This repository is versioned in its entirety with separate modules being used to avoid polluting the dependency tree. Nested modules should still include a `replace` directive and not rely on `go.work`.

```
go work init
find . -name go.mod -exec dirname {} \; | xargs go work use
```
