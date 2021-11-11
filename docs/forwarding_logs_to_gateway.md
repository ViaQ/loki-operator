# Forwarding Logs to Loki Operator

This document will describe how to send application, infrastructure, and audit logs to the Lokistack Gateway as different tenants using Promtail. The gateway provides secure access to the distributor (and query-frontend) via consulting an OAuth/OIDC endpoint for the request subject.

__Please read the `hacking_loki_operator.md` document before proceeding with the following instructions.__

_Note: While this document will only give instructions for two methods of log forwarding into the gateway, the examples given in the Promtail section can be extrapolated to other log forwarders._

## Openshift Logging

Although there is a way to [forward logs to an external Loki instance](https://docs.openshift.com/container-platform/4.9/logging/cluster-logging-external.html#cluster-logging-collector-log-forward-loki_cluster-logging-external), [Openshift Logging](https://github.com/openshift/cluster-logging-operator) does not currently have support to send logs through the Lokistack Gateway.

Support will be added in the near future.

## Promtail

[Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) is an agent managed by Grafana which forwards logs to a Loki instance.

In order to generate an instance of Promtail with the necessary authorization (service account token) and authentication (rbac) resources for interacting with the gateway, perform the following steps:

1. Deploy the Loki Operator to the cluster with the gateway component
2. Execute the following commands in the terminal:

```console
kubectl -n openshift-logging create -f hack/promtail_client.yaml
```

## Troubleshooting

### Log Entries Out of Order

If the forwarder is configured to send too much data in a short span of time, Loki will back-pressure the forwarder and respond to the POST requests with `429` errors. In order to alleviate this, the ingestion rate (global or tenant) can be changed via configuration changes to `lokistack`:

```console
kubectl -n openshift-logging edit lokistack
```

```yaml
limits:
    tenants:
        4a5bb098-7caf-42ec-9b1a-8e1d979bfb95:
            IngestionLimits:
                IngestionRate: 10
                IngestionBurstSize: 20
```
