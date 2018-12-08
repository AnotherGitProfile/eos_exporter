# Eos Exporter for Prometheus
Prometheus exporter for [EOS](https://eos.io). Exposes metrics about accounts such as available CPU, NET, RAM, etc.

To run it:
```console
go mod download
go build
./eos_exporter [flags]
```
Help on flags:
```console
./eos_exporter --help
```

## Using Docker
```console
docker run -p 9386:9386 --volume=/path/to/config.yml:/etc/eos_exporter/config.yml anothergitprofile/eos_exporter
```