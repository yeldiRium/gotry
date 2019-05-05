# GoTry

A small and highly flexible Go library for non-blockingly retrying potentially
failing operations and preserving their return values once they succeed.

Heavily inspired by [avast/retry-go](https://github.com/avast/retry-go) and
[giantswarm/retry-go](https://github.com/giantswarm/retry-go) and partially
based on the latter.

I didn't want to use either of the two because neither
handles return values of the retried operations, which is crucial for things
like connecting to a database which may be offline for a short period of time.

Also I like my libraries tested.

## Usage

```go
package something

import "github.com/yeldiRium/gotry"

func main() {
    resultChannel, err := gotry.Try(func() (*ConnectionHandle, error) {
        return connectToSomeDatabaseWhichMightFailButOtherwiseReturnsAHandle()
    })

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
