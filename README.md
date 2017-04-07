# govuegui
A web GUI for Go

# Development settings

When developing on the html files or the vue app you need to set the unexported variable `useRice = false`.

That the changes takes effect to the package, after every change inside of the html folder you need to run:

```
rice embed-go
```

