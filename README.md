# gofmtc
Custom Go code formatting

```
go build gofmtc.go
./gofmtc gofmtc.go
```

## Examples
|Before   	|After   	|
|---	|---	|
|`vlog.Info().Stack().Msg("hello world")`   	|`vlog.Info().Stack().Msg("Hello world")`   	|
|`fmt.Println(fmt.Errorf("Hello world:%w", errors.New("FF")))`   	|`fmt.Println(fmt.Errorf("Hello world: %w", errors.New("FF")))`   	|
|`errors.New("FF")`   	|`errors.New("fF")`   	|
