# gRPC demo in Go

This is a simple example of using gRPC with Go.  It demonstrates a simple RPC call and bidirectional streaming.  To run it you will need two open terminals and the protocol buffers compiler with Go generator installed (follow installation instructions here: [https://github.com/golang/protobuf/](https://github.com/golang/protobuf/)).  All commands are run from the root of this repository.

Terminal 1:
```
$ go generate github.com/alittlebrighter/grpc-demo
$ go run server/main.go
```

Terminal 2:
```
$ go run client.go
```
