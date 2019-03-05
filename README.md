# Eos Exporter for Prometheus [![Build Status](https://travis-ci.org/AnotherGitProfile/eos_exporter.svg?branch=master)](https://travis-ci.org/AnotherGitProfile/eos_exporter)
Prometheus exporter for [EOS](https://eos.io). Exposes metrics about accounts such as available CPU, NET, RAM, etc.
## Exported metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| eos_account_balance | Currency balance for given account | account, token |
| eos_account_cpu_max | Maximum amount of CPU that can be used by given account | account |
| eos_account_cpu_used | Current value of used CPU for given account | account |
| eos_account_net_max | Maximum amount of NET than can be used by given account | account |
| eos_account_net_used | Current value of used NET for given account | account |
| eos_account_ram_quota | Amount of available ram for given account | account |
| eos_account_ram_used | Total amount of used ram for given account | account |

## Flags
```bash
$ ./eos_exporter --help
Usage of ./eos_exporter:
  -config.file string
    	path to configuration file (default "config.yml")
  -h	show this help
  -port uint
    	port to listen (default 9386)
```

## Builing and running
```bash
go mod download
go build
./eos_exporter [flags]
```

## Using Docker
```bash
docker run -p 9386:9386 --volume=/path/to/config.yml:/etc/eos_exporter/config.yml anothergitprofile/eos_exporter
```
