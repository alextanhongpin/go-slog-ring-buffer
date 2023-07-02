# Defer-style


The idea is still the same as ring buffer logging. We only want to enable all level if there is an error.

But instead of writing our own implementation, we just defer the log statement.

When an error is encounted, we just flip the switch to toggle the log level.
However, the order of the logs printed will be reversed (this may not be so important, because usually you can just reverse the sort order in your logs monitoring tools).


```bash
➜  examples git:(main) ✗ go run main.go
time=2023-07-02T21:30:13.622+08:00 level=ERROR msg="failed to call Bar"
time=2023-07-02T21:30:13.872+08:00 level=DEBUG msg="called Bar" elapsed=232ns
time=2023-07-02T21:30:13.872+08:00 level=DEBUG msg="called Foo" elapsed=281ns
```
