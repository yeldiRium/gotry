# GoTry

[![GoDoc](https://godoc.org/github.com/yeldiRium/gotry?status.svg)](https://godoc.org/github.com/yeldiRium/gotry)
[![codecov](https://codecov.io/gh/yeldiRium/gotry/branch/master/graph/badge.svg)](https://codecov.io/gh/yeldiRium/gotry)
[![Go Report Card](https://goreportcard.com/badge/github.com/yeldiRium/gotry)](https://goreportcard.com/report/github.com/yeldiRium/gotry)
[![GitHub license](https://img.shields.io/github/license/yeldiRium/gotry.svg)](https://github.com/yeldiRium/gotry/blob/master/LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/yeldiRium/gotry.svg)](https://github.com/yeldiRium/gotry/releases)

A small and highly flexible Go library for non-blockingly retrying potentially
failing operations and preserving their return values once they succeed.

Heavily inspired by [avast/retry-go](https://github.com/avast/retry-go) and
[giantswarm/retry-go](https://github.com/giantswarm/retry-go) and partially
based on the latter.

I didn't want to use either of the two because neither
handles return values of the retried operations, which is crucial for things
like connecting to a database that may be offline for a short period of time.

## Usage

You decide whether you want to run Try as a goroutine or not.

```go
package something

import "github.com/yeldiRium/gotry"

func main() {
    resultChannel := make(chan *RetryResult)
    go gotry.Try(
        func() (*ConnectionHandle, error) {
            return connectToSomeDatabaseWhichMightFailButOtherwiseReturnsAHandle()
        },
        resultChannel,
    )

    // do some other things

    for {
        select:
        case res <-resultChannel:
            if res.StopReason != nil {
                // Retrying was stopped because of something.
                switch res.StopReason.(type) {
                    // If you want to find out why retrying stopped failed.
                    case gotry.TooManyRetriesError:
                        //...
                    case gotry.TimeoutError:
                        //...
                    default:
                        // Should not happen.
                }
                // The last error the operation returned
                err := res.LastError
            }
            value := res.Value.(*ConnectionHandle)
            // work with it!
        default:
            doOtherStuffInTheMeanTime()
    }
}
```

Take a look at the [available options](./options.go) for more.

## Known Issues / Future Goals

* Currently operations that are long and blocking block the Try function and delay aborting due to timeouts.

  I.e. if I call Try(op) with op being a complicated computation that takes a while but I want to set a timeout for it,
  the ErrTimeout can at the earliest be returned after one full execution of f.

  There is already a commented out test for this in [gotry_test.go](./gotry_test.go).

  This will probably require a rewrite of the `Try` logic so that `f` is run in another goroutine and raced against the
  timeout.
* Rewriting this using [Contracts](https://go.googlesource.com/proposal/+/master/design/go2draft-contracts.md) would make it even better, since we'd have type safety at compile time.
