# Elasticsearch Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [beta]    |
| Supported pipeline types | metrics   |
| Distributions            | [contrib] |

This receiver queries the Elasticsearch [node stats](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-nodes-stats.html), [cluster health](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html) and [index stats](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-stats.html) endpoints in order to scrape metrics from a running elasticsearch cluster.

## Prerequisites

This receiver supports Elasticsearch versions 7.9+

If Elasticsearch security features are enabled, you must have either the `monitor` or `manage` cluster privilege.
See the [Elasticsearch docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/authorization.html) for more information on authorization and [Security privileges](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-privileges.html).

## Configuration

The following settings are optional:
- `metrics` (default: see `DefaultMetricsSettings` [here](./internal/metadata/generated_metrics.go): Allows enabling and disabling specific metrics from being collected in this receiver.
- `nodes` (default: `["_all"]`): Allows specifying node filters that define which nodes are scraped for node-level and cluster-level metrics. See [the Elasticsearch documentation](https://www.elastic.co/guide/en/elasticsearch/reference/7.9/cluster.html#cluster-nodes) for allowed filters. If this option is left explicitly empty, then no node-level metrics will be scraped and cluster-level metrics will scrape only metrics related to cluster's health.
- `skip_cluster_metrics` (default: `false`): If true, cluster-level metrics will not be scraped.
- `indices` (default: `["_all"]`): Allows specifying index filters that define which indices are scraped for index-level metrics. See [the Elasticsearch documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-stats.html#index-stats-api-path-params) for allowed filters. If this option is left explicitly empty, then no index-level metrics will be scraped.
- `endpoint` (default = `http://localhost:9200`): The base URL of the Elasticsearch API for the cluster to monitor.
- `username` (no default): Specifies the username used to authenticate with Elasticsearch using basic auth. Must be specified if password is specified.
- `password` (no default): Specifies the password used to authenticate with Elasticsearch using basic auth. Must be specified if username is specified.
- `collection_interval` (default = `10s`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration). On larger clusters, the interval may need to be lengthened, as querying Elasticsearch for metrics will take longer on clusters with more nodes.

### Example Configuration

```yaml
receivers:
  elasticsearch:
    metrics:
      elasticsearch.node.fs.disk.available:
        enabled: false
    nodes: ["_local"]
    skip_cluster_metrics: true
    indices: [".geoip_databases"]
    endpoint: http://localhost:9200
    username: otel
    password: password
    collection_interval: 10s
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

The following metric are available with versions:
- `elasticsearch.indexing_pressure.memory.limit` >= [7.10](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/release-notes-7.10.0.html)
- `elasticsearch.node.shards.data_set.size` >= [7.13](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/release-notes-7.13.0.html)
- `elasticsearch.cluster.state_update.count` >= [7.16.0](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/release-notes-7.16.0.html)
- `elasticsearch.cluster.state_update.time` >= [7.16.0](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/release-notes-7.16.0.html)

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)

## Feature gate configurations

See the [Collector feature gates](https://github.com/open-telemetry/opentelemetry-collector/blob/main/featuregate/README.md#collector-feature-gates) for an overview of feature gates in the collector.

**BETA**: `receiver.elasticsearch.emitClusterHealthDetailedShardMetrics`

The feature gate `receiver.elasticsearch.emitClusterHealthDetailedShardMetrics` once enabled starts emitting the metric `elasticsearch.cluster.shards`
with two additional data points - one with `state` equal to `active_primary` and one with `state` equal to `unassigned_delayed`.

This is considered a breaking change for existing users of this receiver, and it is recommended to migrate to the new implementation when possible. Any new users planning to adopt this receiver should enable this feature gate to avoid having to migrate any visualisations or alerts.

This feature gate is enabled by default, and eventually the old implementation will be removed. It aims
to give users time to migrate to the new implementation. The target release for the old implementation to be removed
is 0.71.0.

**BETA**: `receiver.elasticsearch.emitAllIndexOperationMetrics`

The feature gate `receiver.elasticsearch.emitAllIndexOperationMetrics` once enabled starts emitting metrics `elasticsearch.index.operation.count`
and `elasticsearch.index.operation.time` with all possible data points - for every possible operation type and both shard aggregation types.

Because of the amount of added data points, this change might affect performance for existing users of this receiver.
It is recommended to migrate to the new implementation when possible.
Any new users planning to adopt this receiver should enable this feature gate to avoid risking unexpected slowdowns.

This feature gate is enabled by default, and eventually the old implementation will be removed. It aims
to give users time to migrate to the new implementation. The target release for the old implementation to be removed
is 0.71.0.

[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
