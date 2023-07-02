# Defer-style


The idea is still the same as ring buffer logging. We only want to enable all level if there is an error.

But instead of writing our own implementation, we just defer the log statement.

When an error is encounted, we just flip the switch to toggle the log level.
However, the order of the logs printed will be reversed (this may not be so important, because usually you can just reverse the sort order in your logs monitoring tools).


With error:
```bash
➜  examples git:(main) ✗ go run main.go
time=2023-07-03T00:10:17.085+08:00 level=ERROR msg="failed to call Bar" req_id=cigq2q7ltaq0d1n2nle0
time=2023-07-03T00:10:17.340+08:00 level=DEBUG msg="called Bar" req_id=cigq2q7ltaq0d1n2nle0 elapsed=342ns
time=2023-07-03T00:10:17.340+08:00 level=DEBUG msg="called Foo" req_id=cigq2q7ltaq0d1n2nle0 elapsed=458ns
```

Without error:

```bash
➜  examples git:(main) ✗ go run main.go
time=2023-07-03T00:10:19.632+08:00 level=INFO msg="called Bar" req_id=cigq2qvltaq0dcj5tat0
```
