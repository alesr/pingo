# PINGO

When you need to keep something alive with the *wrong* permissions.

![Pingo in action](https://media.giphy.com/media/snUaDwX2aPXV6RPjgw/giphy.gif)

*Use it with discretion.*

## Usage

```go
// Start pinger with default configuration
cancelPinger, err := pingo.StartPinger(context.TODO())
if err != nil {
    log.Fatal(err)
}
defer cancelPinger()

// Configure custom URL
cancelPinger, err := pingo.StartPinger(
    context.TODO(),
    pingo.WithURL("http://example.com"),
)
if err != nil {
    log.Fatal(err)
}
defer cancelPinger()

// Configure custom frequency
cancelPinger, _ = pingo.StartPinger(ctx,
    pingo.WithFrequency(5 * time.Second),
)

// Use custom HTTP client
customClient := &http.Client{
    Timeout: 1 * time.Second,
}
cancelPinger, err = pingo.StartPinger(ctx,
    pingo.WithHTTPClient(customClient),
)
```

The pinger will continue running until either the context is cancelled or the `cancelPinger` function is called.


