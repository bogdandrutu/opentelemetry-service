10 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 967438              1386 ns/op             519 B/op          2 allocs/op
BenchmarkMapV1_NoAlloc-16        5117604               228 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                6403716               188 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 5.530s
```

15 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 728233              1720 ns/op             954 B/op          2 allocs/op
BenchmarkMapV1_NoAlloc-16        4881985               244 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                4621915               263 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 6.598s
```

25 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 603705              1976 ns/op            1068 B/op          2 allocs/op
BenchmarkMapV1_NoAlloc-16        5167207               232 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                3985296               292 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 7.432s
```

50 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 366093              3248 ns/op            2116 B/op          3 allocs/op
BenchmarkMapV1_NoAlloc-16        5335215               248 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                3831584               310 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 6.557s
```

100 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 183198              6067 ns/op            4237 B/op          3 allocs/op
BenchmarkMapV1_NoAlloc-16        4661228               238 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                3122168               374 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 5.646s
```

200 Attributes:
```
go test -benchmem -run=^$ github.com/open-telemetry/opentelemetry-collector/internal/data -bench "^(BenchmarkMapV.*)$"
goos: darwin
goarch: amd64
pkg: github.com/open-telemetry/opentelemetry-collector/internal/data
BenchmarkMapV1-16                 100332             11404 ns/op            8301 B/op          3 allocs/op
BenchmarkMapV1_NoAlloc-16        5007284               252 ns/op               0 B/op          0 allocs/op
BenchmarkMapV2-16                2448896               487 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/open-telemetry/opentelemetry-collector/internal/data 5.750s
```