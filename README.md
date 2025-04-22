# s3-aggregated-metrics-collector

`s3-aggregated-metrics-collector` is a tool to gather aggregated object metrics from S3 buckets and export them to OpenTelemetry (OTel) collectors.

## Features

- Aggregates object sizes and number of objects metrics from S3 buckets.
- Exports metrics to OpenTelemetry collectors for observability and monitoring.
- Supports running as a one-time collection or as a service.

## Usage

Run the collector with the following options:
-h, --help help for s3-aggregated-metrics-collector -t, --toggle Help message for toggle


For more details, see the [documentation](docs/s3-aggregated-metrics-collector.md).

---

## Exported Metrics

The following metrics are exported by the collector:

| Metric Name                   | Description                              | Labels                     |
|-------------------------------|------------------------------------------|----------------------------|
| `s3.aggregated.storage.size`  | Aggregated object size in bytes          | `bucket`, `prefix`         |
| `s3.aggregated.objects.total` | Total number of objects in the prefix    | `bucket`, `prefix`         |

### Metric Details

- **`s3.aggregated.storage.size`**: Tracks the total size of all objects in a specific S3 prefix, measured in bytes. Useful for monitoring storage usage and costs.
- **`s3.aggregated.objects.total`**: Tracks the total number of objects in a specific S3 prefix. Helps monitor bucket usage and lifecycle policy effects.
