# fluentd-metrics

Unfortunately fluentd only seems to offer up metrics based on its output buffer statuses.  So that's what we've got.

    usage: fluentd-metrics --statsd=STATSD --cluster=CLUSTER [<flags>]

    Flags:
      --help             Show context-sensitive help (also try --help-long and
                         --help-man).
      --statsd=STATSD    Host:Port of Datadog Statsd agent
      --cluster=CLUSTER  Name of kubernetes cluster
      --fluent=http://127.0.0.1:24220  
                         Fluentd HTTP API endpoint
      --interval=10s     Gap between metric probes
