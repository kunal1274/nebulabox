# NebulaBox Performance Testing and Benchmarking

This document describes the performance testing and benchmarking tools available in NebulaBox.

## Overview

NebulaBox includes comprehensive performance testing tools:
- **Benchmark tests**: Measure performance of individual operations
- **Load tests**: Test system behavior under concurrent load
- **Performance metrics**: Real-time and historical performance data collection

## Running Benchmarks

### All Benchmarks
```bash
make benchmark-all
```

### API Benchmarks Only
```bash
make benchmark-api
```

### Containerd Benchmarks Only
```bash
make benchmark-containerd
```

### Specific Package
```bash
go test -bench=. -benchmem -benchtime=3s ./internal/api/...
```

## Load Testing

Run load tests to verify system behavior under concurrent requests:

```bash
make load-test
```

This runs tests that simulate concurrent load on:
- Container listing endpoint
- Network listing endpoint  
- Authentication endpoint

## Performance Metrics

### Real-time Metrics
The API server tracks performance metrics for all endpoints:

- **Request rate**: Requests per second (1m and 5m windows)
- **Latency**: P95 latency in milliseconds
- **Error rate**: Percentage of failed requests

Access metrics via:
- `GET /api/perf/metrics` - Current snapshot
- `GET /api/perf/stream` - Real-time SSE stream
- `GET /api/perf/endpoints` - Per-endpoint detailed metrics

### Per-Endpoint Metrics

The system tracks detailed metrics for each API endpoint:
- Total request count
- Average, min, and max latency
- Error count and error rate
- Last request timestamp

## Benchmark Tests

### API Endpoints
- `BenchmarkAPI_ListContainers` - Container listing performance
- `BenchmarkAPI_ListNetworks` - Network listing performance
- `BenchmarkAPI_ListTeams` - Team listing performance
- `BenchmarkAPI_ListTenants` - Tenant listing performance
- `BenchmarkAPI_CreateNetwork` - Network creation performance
- `BenchmarkAPI_AuthLogin` - Authentication performance

### Containerd Operations
- `BenchmarkListContainers` - Container listing
- `BenchmarkPullImage` - Image pulling
- `BenchmarkCreateContainer` - Container creation

## Load Test Results

Load tests report:
- Total requests processed
- Successful vs failed requests
- Average, min, and max response times
- Requests per second
- Error rate percentage

Example output:
```
📊 Load Test Results - List Containers
=====================================
Total Requests:    100
Successful:        100
Failed:            0
Total Duration:     2.5s
Avg Response Time:  25ms
Min Response Time:  20ms
Max Response Time:  45ms
Requests/sec:      40.00
Error Rate:        0.00%
=====================================
```

## Integration with Monitoring

Performance metrics are integrated with the monitoring dashboard. View real-time performance data in the Performance page of the web dashboard.

## Best Practices

1. **Run benchmarks before deployments** - Establish baseline performance
2. **Monitor error rates** - Keep error rates below 1% under normal load
3. **Track latency trends** - Watch for gradual latency increases
4. **Load test new features** - Verify new endpoints handle concurrent requests
5. **Compare before/after** - Use benchmarks to verify performance improvements

## Continuous Performance Monitoring

The performance middleware automatically tracks all API requests. Metrics are available via:
- REST API endpoints
- Server-Sent Events (SSE) streams
- Web dashboard visualization

