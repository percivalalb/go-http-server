# Example HTTP Application

This is a simple [Golang][1] HTTP server, that prints "Hello, world\n" on only the base path (`/`) to [GET requests][2].

It demonstrates how to correctly, listen for interrupt signals and shutdown gracefully handling outstanding requests.

Run the application with `go run ./` and Ctrl-C to send an interrupt signal.

```console
$ go run ./
time=2024-02-11T14:12:41.474Z level=INFO msg=Starting
time=2024-02-11T14:12:41.474Z level=INFO msg="Listening for interrupt signals"
time=2024-02-11T14:12:41.474Z level=INFO msg="Setting up mux"
time=2024-02-11T14:12:41.476Z level=INFO msg="Listening on" addr=[::]:8080
time=2024-02-11T14:12:41.476Z level=INFO msg=Serving time=1.642225ms
time=2024-02-11T14:12:42.778Z level=INFO msg="serving request" source=127.0.0.1:40404
^Ctime=2024-02-11T14:12:45.035Z level=INFO msg="Initiating shutdown" timeout=10s
time=2024-02-11T14:12:45.036Z level=INFO msg=Shutdown time=982.991µ
```

[1]: https://go.dev/
[2]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/GET
