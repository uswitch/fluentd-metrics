# fluentd-metrics

    usage: fluentd-metrics --statsd=STATSD --cluster=CLUSTER [<flags>]

    Flags:
      --help             Show context-sensitive help (also try --help-long and
                         --help-man).
      --statsd=STATSD    Host:Port of Datadog Statsd agent
      --cluster=CLUSTER  Name of kubernetes cluster
      --fluent=http://127.0.0.1:24220  
                         Fluentd HTTP API endpoint
      --interval=10s     Gap between metric probes
