# WPORG

An API wrapper made to simplify fetching data from the WordPress.org APIs in Go.

## Examples

Get Latest Revision
```go
rev, _ := api.GetRevision("plugins")

rev, _ := api.GetRevision("themes")
```

Get Directory List
```go
plugins, _ := api.GetList("plugins")

themes, _ := api.GetList("themes")
```

Get Directory Changelog
```go
list, _ := api.GetChangeLog("plugins", current, latest)

list, _ := api.GetChangeLog("themes", current, latest)
```

Get Info
```go
info, _ := api.GetInfo("plugins", "gutenberg")

info, _ := api.GetInfo("themes", "twentytwelve")
```