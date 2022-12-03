[![test-and-linter](https://github.com/microup/vbalancer/actions/workflows/main.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/main.yml)

# What is VBalancer?

The VBalancer has realized for switch traffic between peers, and it uses Robin Round algorithm.
This is an implementation need to increase stability and downgrade high load to backend.

![Diagram](assets/vbalancer.png)

## Important: need set ENV to run

For normal run a VBalancer, it needs to set an environment OS "ProxyPort" and path to config file "ConfigFile".

## Settings

On file: config/config.yaml you can add or delete peer(s), and configure PROXY settings.

## Docker

### build

```bash
$docker build --tag vbalancer . -f Dockerfile
```

### run

```bash
$docker run --rm -p 8080:8080 vbalancer
```

## Features

- The proxy (VBalancer) has realized on TCP net.Listener
- All log write to 'csv' file
- Size log file can be changed in the config file
  
