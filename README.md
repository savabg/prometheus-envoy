# Enphase Envoy Golang Client

This is a prometheus collector for pulling metrics from an Envoy Enphase unit. The collector
utilizes the local interface exposed by the device rather than the Enlighten API, and is set to work with V7. Enphase units
are embedded devices, so the collector is implemented as a proxy collector similar to the
snmp_exporter tool.

<https://enphase.com/en-us/support/what-envoy>

## Example

Set environment variable ENPHASE_TOKEN to your JWT token obtained from https://entrez.enphaseenergy.com


```bash
ENPHASE_TOKEN=

```

```yml
  - job_name: 'prometheus-envoy'
    params:
      target: ['<envoy-ip>']
    static_configs:
      - targets: ['127.0.0.1:2112']
        replacement: 127.0.0.1:2112  # The prometheus-smarthome's real hostname:port.
```

## Building and running

```sh
cd cmd/prometheus-envoy
go build
./prometheus-envoy -port 2112 -listen 0.0.0.0
```

## License

This library is provided under the [MIT License](LICENSE.md)
