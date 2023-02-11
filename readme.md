![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/microup/vbalancer) ![GitHub repo file count](https://img.shields.io/github/directory-file-count/microup/vbalancer) [![test-and-linter](https://github.com/microup/vbalancer/actions/workflows/main.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/main.yml) 
[![Release](https://github.com/microup/vbalancer/actions/workflows/release.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/release.yml) [![Docker CI](https://github.com/microup/vbalancer/actions/workflows/docker-ci.yml/badge.svg)](https://github.com/microup/vbalancer/actions/workflows/docker-ci.yml)

# What is VBalancer?

VBalancer is a project written on Go lang that provides a simple and efficient way to balance traffic between identical servers(peers). It is designed to be lightweight and fast. Vbalancer supports TCP protocol. With its easy-to-use configuration options, it is an ideal solution for any organization looking to improve their network application performance and reliability. VBalancer is an open source project released under the AGPL-3.0 license.

![Image description](https://raw.githubusercontent.com/microup/vbalancer/33860ed4ee5dcfce78f9cb760208346eac426739/assets/vbalancer.png)

## Getting Started:

This solution it allows users to easily configure and manage their load balancers with minimal effort.

1. Clone the repository: Use the following command to clone the repository from GitHub: ```git clone https://github.com/microup/vbalancer.git```
2. Install the dependencies required for the project: ```go get github.com/microup/vbalancer```
3. Build the project: Navigate to the project directory and run the following command to build the project: ```go build```
4. Run the project with the following command: ```./vbalancer -config <path_to_config_file>```
5. Configure your load balancers using the configuration file (config.yaml):
    - Set up your load balancers by specifying their IP addresses, ports, and other settings in the config file
    - Add or remove services from your load balancers by adding or removing them from the config file
7. Enjoy!

## Docker

Using VBalancer, you can quickly set up a load balancing solution for your Docker containers. It allows you to easily configure the rules for routing traffic between containers, as well as set up health checks to ensure that the containers are running properly. 

### You can build manually

#### build

```bash
$docker build --tag vbalancer . -f Dockerfile
```

#### run

```bash
$docker run --restart=always -p 8080:8080 vbalancer
```
### You can take a ready image.

#### login to Github 

```bash
$docker login docker.pkg.github.com -u USERNAME -p ${{ secrets.TOKEN }}
```

#### pull the image from GitHub Packages

```bash
$docker pull docker.pkg.github.com/microup/vbalancer/vbalancer:latest
```
#### then start it

before using the ready-made image, you need to enter the container and modify the configuration file by adding or removing the required hosts

```bash
$docker run -d -e PORT=8080 -p 8080:8080 docker.pkg.github.com/microup/vbalancer/vbalancer:latest
```
Congratulations! Now you can use the high-load balancer by sending network requests to it.

#### to stop it, get the container id:

```bash
$ docker container ls

CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                    NAMES
271b4abb963s        vbalancer            ./vbalancer"       7 minutes ago       Up 7 minutes        0.0.0.0:8080->8080/tcp   beeraf_skills
```
#### and stop it

```bash
$ docker container stop 271b4abb963s
```

## Features

- The proxy (VBalancer) has realized on TCP net.Listener
- A limit on connections and a deadline time can be set
- All log write to 'CSV' file
- Size log file can be changed in the config file

### License

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free
Software Foundation, version 3.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
See the GNU Affero General Public License for more details.

####

Copyright (C) 2022-2023 https://microup.ru
