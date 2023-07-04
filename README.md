# unwindlogger

Reference implementation of a logger, that can unwind and
log previously skipped messages, after encountering an error


## Motivation

In prod, 99% of the time, we don't need logs.

So, we log at WARN or ERROR level.
To keep costs low, and drive our alerts/monitors etc.

However, when we _need_ logs, we want _all_ the logs.
Eg: what business processes fired, what happened in the
call before this error, etc.
Typically, these are logged at INFO level, and so we can't
see what happened.

So, wouldn't it be nice if, when an error was encountered, we
could see the INFO logs for the rest of the request?
Including from before we knew an error would occur?

That's what this logger intends to do

## How it works

See [example](./example_test.go) for example usage

Whenever a message with a context is logged, if it is below the
immediate log level, it is deferred.

If an error does not occur, the ignorable logs are not logged, and
dropped when the tracking stops
```go
logger := unwindlogger.
    NewLogger().
    WithImmediateLevel(unwindlogger.WARN).
    WithDeferredLevel(unwindlogger.INFO)

ctx := context.Background()
doSomething := func(ctx context.Context) {
    logger.WithContext(ctx).WithField("hello", "world").Info("info log!")
}

ctx = logger.StartTracking(ctx)
doSomething(ctx)
logger.WithContext(ctx).Warn("warn log!")
logger.EndTracking(ctx)
// Will only output: {"level":"WARN","msg":"warn log!","time":"..."}
// As the tracking was ended without an error

```

If an error does occur, the deferred logs are flushed, providing
more context
```go
logger := unwindlogger.
    NewLogger().
    WithOut(os.Stdout).
    WithImmediateLevel(unwindlogger.WARN).
    WithDeferredLevel(unwindlogger.INFO)

ctx := context.Background()
doSomething := func(ctx context.Context) {
    logger.WithContext(ctx).WithField("hello", "world").Info("info log!")
}

ctx = logger.StartTracking(ctx)
doSomething(ctx)
logger.WithContext(ctx).Warn("warn log!")
logger.EndTrackingWithError(ctx, fmt.Errorf("oh no"))
// Will output both: {"level":"INFO","msg":"info log!","time":"..."} and {"level":"WARN","msg":"warn log!","time":"..."}
// As the tracking was ended with an error
```

Also supports a mode called `fullDefer`, which defers all logs,
and logs those at immediate level, when tracking stops without an error,
or logs at defer level, when tracking stops with an error.

This `fullDefer` mode ensures that logs are always emitted in time order,
as deferring some logs and not others can create out of order logs.

## Performance

It ain't great, but it's a POC.

Easy improvements possible with better buffer management
and grouping writes to the out pipe.

```shell
$> go test -test.benchtime 10s -test.benchmem -test.bench .
goos: linux
goarch: amd64
pkg: github.com/Ben435/unwindlogger
cpu: AMD Ryzen 7 1700 Eight-Core Processor          
BenchmarkLogger_ImmediateLogging-16             17948122               674.4 ns/op          1218 B/op         22 allocs/op
BenchmarkLogger_DeferredLogging-16              16605958               729.1 ns/op          1226 B/op         23 allocs/op
BenchmarkLogger_MixedLoggingWithError-16         8963803              1344 ns/op            2260 B/op         43 allocs/op
BenchmarkLogger_MixedLoggingWithoutError-16     14674735               830.2 ns/op          1355 B/op         26 allocs/op
```
