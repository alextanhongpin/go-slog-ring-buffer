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


## Thoughts

- storing in buffer can take up a lot of space, should probably limit the amount of records the buffer can hold
- pass in correlation id to tie logs together
- buffering logs before logging them can cause logs to be out-of-order. Store the original timestamp in another field, e.g. `event_time` to distinguish between the time the event happened, and the time of logging.
- TBH, this can be accomplished by wrapping errors
