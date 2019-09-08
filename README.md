## ProxyBlock

ProxyBlock is a simple add-in for network proxies to keep a sliding-window count of hits and recommend a deny via a simple `ShouldBlock()` method in the case of the limit being exceeded.

###  Sliding Window

``` previousCount * ratio of (how many seconds into current minute/60 seconds) + currentCount ```

### http example

```go run http/server.go```

Returns `http.StatusTooManyRequests` which is status code `429` when the count is exceeded.