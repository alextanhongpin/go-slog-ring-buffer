# Ring Buffer Logging


When there's no error, only logs at `INFO` level will be logged.

```bash
➜  go-slog-ring-buffer go run main.go | jq
{
  "time": "2023-07-02T15:10:38.369088+08:00",
  "level": "INFO",
  "source": {
    "function": "main.Bar",
    "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/go-slog-ring-buffer/main.go",
    "line": 40
  },
  "msg": "bar",
  "req_id": "cigi5rnltaqfe68kt3sg",
  "event_time": "2023-07-02T15:10:38.369076+08:00"
}
```



However, when there is an error (for that request), then the log level will be updated to log all levels.


```bash
➜  go-slog-ring-buffer go run main.go | jq
{
  "time": "2023-07-02T15:10:58.166532+08:00",
  "level": "DEBUG",
  "source": {
    "function": "main.Foo",
    "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/go-slog-ring-buffer/main.go",
    "line": 33
  },
  "msg": "foo",
  "req_id": "cigi60nltaqff0faqi00",
  "event_time": "2023-07-02T15:10:58.166491+08:00"
}
{
  "time": "2023-07-02T15:10:58.166744+08:00",
  "level": "ERROR",
  "source": {
    "function": "main.Bar",
    "file": "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/go-slog-ring-buffer/main.go",
    "line": 42
  },
  "msg": "bar",
  "req_id": "cigi60nltaqff0faqi00",
  "event_time": "2023-07-02T15:10:58.166519+08:00"
}
```

## Improvement

- condense the `source` representation into one line
- show paths relative to source project directory



Without error:

```bash
➜  go-slog-ring-buffer git:(main) ✗ go run main.go | jq
{
  "time": "2023-07-02T21:07:36.062283+08:00",
  "level": "INFO",
  "src": "bar/bar.go:12 bar.Bar",
  "msg": "calling Bar",
  "event_time": "2023-07-02T21:07:36.062236+08:00",
  "req_id": "cignd67ltaq6un5v0vtg"
}
{
  "time": "2023-07-02T21:07:36.062616+08:00",
  "level": "INFO",
  "src": "bar/bar.go:15 bar.Bar",
  "msg": "Bar called",
  "event_time": "2023-07-02T21:07:36.062267+08:00",
  "req_id": "cignd67ltaq6un5v0vtg"
}
```

With error:

```bash
➜  go-slog-ring-buffer git:(main) ✗ go run main.go | jq
{
  "time": "2023-07-02T21:09:05.40141+08:00",
  "level": "DEBUG",
  "src": "main.go:28 main.Foo",
  "msg": "calling Foo",
  "user": {
    "name": "John"
  },
  "event_time": "2023-07-02T21:09:05.401359+08:00",
  "req_id": "cigndsfltaq74eikhs90"
}
{
  "time": "2023-07-02T21:09:05.4017+08:00",
  "level": "DEBUG",
  "src": "main.go:29 main.Foo",
  "msg": "Foo called",
  "event_time": "2023-07-02T21:09:05.401367+08:00",
  "req_id": "cigndsfltaq74eikhs90"
}
{
  "time": "2023-07-02T21:09:05.401708+08:00",
  "level": "INFO",
  "src": "bar/bar.go:12 bar.Bar",
  "msg": "calling Bar",
  "event_time": "2023-07-02T21:09:05.401372+08:00",
  "req_id": "cigndsfltaq74eikhs90"
}
{
  "time": "2023-07-02T21:09:05.401715+08:00",
  "level": "DEBUG",
  "src": "bar/bar.go:13 bar.Bar",
  "msg": "SELECT 1 + $1",
  "args": {
    "$1": 42
  },
  "event_time": "2023-07-02T21:09:05.401375+08:00",
  "req_id": "cigndsfltaq74eikhs90"
}
{
  "time": "2023-07-02T21:09:05.401724+08:00",
  "level": "ERROR",
  "src": "bar/bar.go:17 bar.Bar",
  "msg": "failed to call Bar",
  "event_time": "2023-07-02T21:09:05.401401+08:00",
  "req_id": "cigndsfltaq74eikhs90"
}
```


## Thoughts

- storing in buffer can take up a lot of space, should probably limit the amount of records the buffer can hold
- pass in correlation id to tie logs together
- buffering logs before logging them can cause logs to be out-of-order. Store the original timestamp in another field, e.g. `event_time` to distinguish between the time the event happened, and the time of logging.
- TBH, this can be accomplished by wrapping errors
- what to log?
  - the start of a step, with the input
  - the end of the step, with the output and whether it is successful or failed

## References

- https://yiblet.com/posts/ring-buffer-logging
- https://tersesystems.com/blog/2019/07/28/triggering-diagnostic-logging-on-exception/
- https://www.komu.engineer/blogs/09/log-without-losing-context?utm_source=pocket_saves
