# ev

A helper package to use
[github.com/tidwall/evio](https://github.com/tidwall/evio) more easily

### Dependencies

* [github.com/tidwall/evio](https://github.com/tidwall/evio)

### How To Use

See
[http example](https://github.com/JesusIslam/ev/tree/master/examples/http/http.go)

### TODO

* Improve performance (Even with just wake up and immediately replying, it could
  only reach 17k rps on my machine, (just like Node.JS 8.9.2; while with the
  http example, it averaged at 13k rps. I think this could be hard to do, the
  command I use is `go-wrk -M POST -H 'Content-Type: application/json' -body
  '{"data":"ping"}' -c 10 -d 10 -T 500 http://localhost:5000`)

* Make this pluggable with Gin (I still don't know how to make
  http.ResponseWriter to feed Gin engine to)
