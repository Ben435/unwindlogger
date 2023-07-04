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


