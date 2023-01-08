[![test-and-linter](https://github.com/microup/vbalancer/actions/workflows/main.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/main.yml) [![Release](https://github.com/microup/vbalancer/actions/workflows/release.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/release.yml)

# What is VBalancer?

VBalancer TCP/IP socket reverse proxy for highload switch traffic between identical peers, and it uses Robin Round algorithm.
This is an implementation need to increase stability and downgrade high load to back-end.

![Diagram](assets/vbalancer.png)

## Settings

There is file config.yaml, where you can add or delete peer(s), and configure PROXY settings.

## Docker

### build

```bash
$docker build --tag vbalancer . -f Dockerfile
```

### run

```bash
$docker run --restart=always -p 8080:8080 vbalancer
```

## Features

- The proxy (VBalancer) has realized on TCP net.Listener
- All log write to 'CSV' file
- Size log file can be changed in the config file

## License

Copyright (C) 2022-2023 https://microup.ru

### vblancer

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free
Software Foundation, version 3.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
See the GNU Affero General Public License for more details.
