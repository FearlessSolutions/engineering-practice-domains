## Containerization Tool Ecosystem

- Kubernetes - container orchestration system that manages containerized
  applications across multiple hosts
  - Self-Healing: Kubernetes continuously monitors the state of pods and
    nodes and can restart or replace containers that fail, migrate
    containers when nodes die, and kill containers that don't respond to
    user-defined health checks.
  - Helm charts package all the necessary components needed to run an
    application or service on Kubernetes.
- Docker Swarm - turns a pool of Docker hosts into a single, virtual
  Docker host
- Podman - run containers without daemon, without root access.
  - Manage groups of containers together in pods like kubernetes
  - `podman-compose` compatible with docker-compose
- [Colima](https://github.com/abiosoft/colima) - Docker Desktop
  replacement. Container runtimes on macOS (and Linux) with minimal
  setup. 16k+ stars on github
- Amazon Elastic Container service, and Azure Container Instances
- BuildKit
  - Parallelize run steps that do no depend on eachother
  - Use secrets during the build process without leaving traces of those
    secrets in the final image

## General

- Containerization stands in contrast to traditional architectures:
  - Monolithic - easiest to do at first, difficult to scale
  - Service-oriented - smaller, loosely-coupled services developed
    deployed and scaled independently
  - Microservices architectures - SOA + well defined APIs
- Less secure than virtual machines as containers share the host
  system’s kernel, whereas VMs provide hardware-level virtualization.
  - If the kernel is compromised, so are the containers.
- Containers don't include OSes and are lightwieght.
- Containers build for one CPU architecture can be run on another using
  QEMU (open-source machine emulator and virtualizer). Beware emulation
  overhead, i.e., translating instructions from one architecture to
  another at runtime.

### Images/Containers lifecycle

- Dockerfile instructions are only executed during image build process.
  Results in a Docker image which is a static snapshot of the
  environment and code as defined in Dockerfile
- When you docker run, docker creates a container based on that static
  image.
- Containers do no autoupdate - you must rebuild the image and recreate
  the containers from the updated image
- Run `docker build `<dockerfile> to rebuild image or
  `docker compose build` in a directory with a docker-compose.yml to
  rebuild all the images/containers in your docker compose project.

### Terminology

- Daemon/Engine/Server - a process running on a host os/system
- Image - READ-ONLY: image does not have state and never changes.
- Container- A "stateful" instance of a docker image
- Layer - A modification to the image, represented by an instruction in
  the Dockerfile
  - Caching- When an image is updated or rebuilt, only layers that
    change need to be updated, and unchanged layers are cached locally.
    This is part of why Docker images are so fast and lightweight.
- Dockerfile - All the commands you would normally execute manually in
  order to build a Docker image
- Docker Hub - Docker certified images/template hosting, with SECVULN
  scanning
  - Official images vs open source images
  - also Github Container Repository (ghcr.io), Red Hat Quay (quay.io),
    AWS Elastic Container Registry (ECR), or Google Container Registry
    (GCR)
- Compose - Defining and running multi-container Docker applications
- Docker Desktop - includes engine, cli client, compose, content trust,
  kubernetes, credential helper

### Alpine Linux

- Alpine Linux's combination of a small footprint, security features,
  and efficiency has made it a popular choice for Docker containers
- Lightweight, security-oriented Linux distribution based on musl libc
  (a lightweight standard library for C language) and BusyBox
- Alpine Linux does not derive from one of the major families like
  Debian or Red Hat but is independently developed, but rather created
  by Natanael Copa in 2005 as a fork of the LEAF Project (Linux Embedded
  Appliance Framework)
- Minimal size 5 MB base image!!!
- You probably shouldn't use Alpine for Python projects, instead use the
  slim Docker image versions that are still based on Debian, but are
  smaller. - source:
  [tiangolo](https://github.com/tiangolo/uvicorn-gunicorn-fastapi-docker)
  - Alpine is more useful for other languages where you build a static
    binary in one Docker image stage (using multi-stage Docker building)
    and then copy it to a simple Alpine image, and then just execute
    that binary. For example, using Go. But for Python, as Alpine
    doesn't use the standard tooling used for building Python
    extensions, when installing packages, in many cases Python (pip)
    won't find a precompiled installable package (a "wheel") for Alpine.

### Bitnami

Bitnami containers, including the Python container, are designed to be
lightweight and secure, focusing on running applications. As such, they
often do not include unnecessary packages, tools like SSH, or even a
full shell environment that you might be accustomed to in standard Linux
distributions. This design choice helps minimize the attack surface and
the size of the container image.

### Low-level implementation of containerization

- Namespaces - within Linux kernel, isolates and virtualizes the system
  resources, restricts processes within the namespace to only interact
  with other processes in the same namespace
- Control groups - limit and isolate resource usage (CPU, memory, I/O,
  network
- Union file system - allowing files and directories of separate file
  systems, known as layers, to be transparently overlaid to form a
  single coherent file system.

## Docker CLI Commands

### Version

- `docker version` - will give you client AND server version info

### Info

- `docker info` - more verbose docker server information

### List installed Images

- `docker images -a` - list the installed images, -a includes the exited
  ones

### List containers

- `docker ps -a` - running process image, age, status, name, and port
  information
- `-a` flag lists exited containers as well
- `--size` flag lists container size, i.e., disk space
  - Actual size: The size utilized by the writable container layer
    consists of the container’s filesystem and any modifications made
  - Virtual Size: combined size of all the container's layers, which
    include any shares, base pictures, and additional layers
- `docker ps -a --format '{{.ID}}' | xargs docker rm`

### top

- `docker stats` - Resource utilization like top broken down by
  container

### Logs

- `docker logs [OPTIONS] CONTAINER` - --details, -f, -tail n,
  --timestamps, --until

### Exec

- `docker exec -it HEXADECIMAL bash` - ssh into container, just do this
  to begin with and keep a tmux open with it.

### Pull

- `docker pull` - pull from dockerhub or wherever
- `docker pull nvcr.io/nvidia/digits:18.06`
- `docker pull nvidia/caffe`

### Remove Container

- `docker rm Container_ID_or_Name` - remove container
- `docker rm $(docker container ls -q -a)`

### Remove Image

- `docker rmi $(docker images -a -q)` - remove all images

### History

- `docker history --no-trunc IMAGE_NAME` - Show all layers of an image
  in reverse chronological order (the last layer being the one on top).

### Run

- For interactive container session, usage is
  `docker run -it `<image>` `<command>
  - \<-it: This is a combination of -i and -t. -i keeps STDIN open even
    if not attached, and -t allocates a pseudo-TTY, making the terminal
    interactive.
  - command could be sh or bash
- Examples:
  - `docker run --gpus all -it --name tensorflow_1910_py3 -d --net=host -v /home/colettace/projects:/projects --ipc=host --ulimit memlock=-1 nvcr.io/nvidia/tensorflow:19.10-py3`
  - `docker run --runtime=nvidia --name digits -d -p 5000:5000 -v /jobs:/jobs -v /datasets:/datasets -it nvidia/digits`
  - `docker run -P` - uppercase P publishes all ports that are exposed
    (?)
  - `docker run --entrypoint /bin/sh -it myimage:latest -c "cat /app/config/app.conf"`
    Override entrypoint for debugging purposes

### Login

- Log in with your Docker ID or email address to push and pull images
  from Docker Hub

### Search

- `docker search --filter=is-official=true nginx` - show official nginx
  image

### Manifest

- `docker manifest inspect `<image_name> - see all available manifests
  and their architectures.
- A manifest list allows you to use one name to refer to the same image
  built for multiple architectures

### Build

- <cdoe>docker build --tag myapp:latest .</code>
- `--squash` Squash layers to save space?

#### When you need secrets to build

- Do not use ARG to pass secrets, which may leave sensitive information
  in image layers
- `docker build --secret id=mysecret,src=/path/to/secret/file.txt -t myapp:latest .` -
  build with a secret
  - `RUN --mount=type=secret,id=mysecret cat /run/secrets/mysecret` -
    Use the secret without exposing it in the image
- Use when
  - When your build requires pulling dependencies from private
    repositories.
  - When you need to use private keys for SSH access or to decrypt files
    during the build
  - When your build process involves accessing APIs that require
    authentication

#### .dockerignore

- tells Docker which files and directories to ignore when building an
  image

<!-- -->

    .git
    .gitignore
    Dockerfile*
    *.md
    node_modules
    temp/

### Prune

- prune images, containers, volumes, networks, and build cache

### Environment Vars

    # Enable Docker BuildKit
    export DOCKER_BUILDKIT=1

### cp

- `docker cp /host/directory/path container_name:/container/directory/path`
- `-a` archive
- `-L` follow links

## Dockerfile Syntax

### FROM

- Initialize new build stage, specifies base image
- Mandatory first non-comment instruction
- could also be `FROM scratch`

### LABEL

- Adds metadata to an image, like the maintainer, version, or any other
  key-value pairs you wish to include.
- Replaces deprecated MAINTAINER statement

### RUN

- Executes any commands on top of the current image layer and commits
  the results
- Install software packages, build code, or change file system contents

### CMD

- defaults for executing a container
- There can be only one of these in a Dockerfile

### COPY

- copy from source working directory in current build context (where you
  run docker build) or build container, to destination filesystem

### ADD

- like COPY but with url and tarball auto-extraction support
- Avoid since downloading from URL could introduce unexpected behavior
- Use RUN wget or RUN curl to provide more control over the download
  process (e.g., checksum validation)

### ENV

- Set environment variables
- Can be overridden at runtime using -e

### EXPOSE

- Exposes a port to be accessed internally by other containers within
  app.
- Does not actually publish the port to the outside, or make a service
  available outside Docker environment.
- Not strictly necessary since two containers on the same network within
  Docker Compose will be able to listen to any port regardless
- Helps to be explicit, helps document intentions, helps some
  orchestration tools

### ENTRYPOINT

- Configures a container that will run as an executable. Command line
  arguments to docker run <image> will be appended.

### USER

- specify username/UID, otherwise it's root. Will be the user of any
  command run by RUN, CMD and ENTRYPOINT instructions

### VOLUME

- Persist data outside containers on host OS, because data within
  containers is ephemeral.
- Creates a mount point with the specified name and marks it as holding
  externally mounted volumes from the native host or other containers.
- Critical for database
- Containers themselves are stateless; it's the volumes that hold the
  state

### WORKDIR

- Sets the target directory inside the container
- Change the current working directory for the duration of the
  Dockerfile
- Sets the working directory for any RUN, CMD, ENTRYPOINT, COPY, and ADD
  instruction
- Creates the dir in the container if it doesn't exist

### ARG

- Defines a variable that users can pass at build-time to the builder
  with the docker build command using the
  `--build-arg `<varname>`=`<value> flag.
- Does not add a layer to the image.
- Variables can then be used using the traditional BASH syntax
  `${variable}`
- Validate like so:
  - `RUN if [ "$config" != "expected_value" ]; then echo "Invalid config value" && exit 1; fi`
  - `RUN `$$"${variable}" == "true"$$` && ...`
- ARG values aren't available in the final image like ENV values are
- Security: ARG values can still be seen in the intermediate layers of
  the build process and potentially in the Dockerfile itself

### STOPSIGNAL

- How you indicate the stop signal that your application or main process
  is configured to handle for a graceful shutdown.
- Default is SIGTERM followed by SIGKILL when you hit Control-C again or
  after 10 second grace period passes.

### Others

- ONBUILD

<!-- -->

- HEALTHCHECK
- SHELL

## Dockerfile workflow

- Try to order instructions to optimize build cache usability

### Multistage builds

- reduces final image size as well as reduces attack surface by
  excluding unnecessary tools
- Use multiple `FROM` statements in Dockerfile, give earlier ones
  aliases using the `AS` keyword
- `COPY --from=builder /app/myapp /app/myapp`

## Docker Compose

- Installation - Included with Docker desktop for Windows and MacOS,
  otherwise download binary for linux based on CPU architecture
- Can define volumes, configs, secrets at top-level, otherwise define at
  service level
- Service names as hostnames - Containers within the same network can
  reach and communicate with each other using the service names as
  hostnames.
- Merge multiple Compose files
  - e.g., production.yml for production-appropriate configuration
- Profiles: A service with no profile specified means it's always
  active, otherwise must specify when running up command
- Can interpolate variables into values used in compose files at runtime
  using bash-like syntax
  - `$VARIABLE`
  - `${VARIABLE}`
  - `${VARIABLE-default}` - "-" dash delimiter specifies default value
    only if `VARIABLE` is unset
  - `${VARIABLE:-default}` - ":-" colon-dash delimiter specifies default
    value if `VARIABLE` is unset OR empty
  - `${VARIABLE?err_msg}` - "?" error out with err_msg if `VARIABLE` is
    unset (":?" if unset or empty)
  - Interpolation can be nested: `${VARIABLE:-${FOO:-default}}`
  - `$$` - a literal dollarsign
  - Can't interpolate YAML keys, only values

### Compose file syntax

- [Compose file syntax
  specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md)

#### services

- `services` - key is the service name, values is the service definition
- `platform: linux/amd64` - os/arch/variant - specify if the image is
  specific to a platform
- `image` - no build step necessary if you use an image
- `build`
  - Either path to build context, or break out into details
  - `context` - path is relative to current compose file, else defaults
    to .
  - `dockerfile: path/to.Dockerfile` - look for a specific named
    dockerfile, otherwise just look for a file named `Dockerfile`
  - If list both image and build, will try to pull image first before
    building if that fails
- `ports: [HOST IP:]host_port:container_port` - ports can also be ranges
- `depends_on` - expresses startup & shutdown dependencies between
  services
- `environment` - set env vars inside container
  - Use "true", "false", "yes", "no", any boolean values should be
    quoted, otherwise will be converted by YAML parser
  - can use YAML map syntax `ENV_VAR: value` or array syntax with one
    dash per line `- ENV_VAR=value`
- `command: `<linux shell command> - overrides the default command
  declared by container image within dockerfile `CMD`
- `volumes`
  - directories on host machine's filesystem mounted into container
  - `HOST_PATH or top-level volume name:CONTAINER_PATH:ACCESS_MODE` -
    access modes are rw, ro, z, Z SELinux shared and unshared with other
    containers respectively.
  - Creates the host path if not there by default
- `profiles` - only run these services when I explicitly ask for them in
  my docker compose up command
- `container_name` - specify a custom container name, rather the default
  generated name <project_name>`_`<service_name>`_`<instance_number>

#### others

- `networks` - services communicate with each other through networks.
  Establish an IP route between containers to connect services
  - specify a `default` network
- `volumes` - How services store and share persistent data. Could be
  high-level filesystem mount, bypasses union filesystem
  - If declared at top-level, volumes can be reused across multiple
    services
- `configs` - files mounted into the container
- `secrets` - flavor of configuration data
- Project - project name set with top-level `name` attribute
  - Can then expose `COMPOSE_PROJECT_NAME` in
    `services/`<SERVICE NAME>`/environment` yaml section which will then
    get pushed to env vars within container

## Docker Compose CLI Commands

### Up

#### --profile

- `docker compose up --profile `<profile name>` -f docker-compose.custom.yml`
- Note --profile != -p - this will fail silently!!

#### Overrides

- `docker compose -f docker-compose.yml -f production.yml up -d` -
  overrides are done in order, i.e., LIFO

#### detached mode

- `docker compose up -d` - send process into the background
- `docker compose stop` - stops the containers
- `docker compose down` stops the containers and also removes them along
  with any networks that were created
  - can help remove any orphans containers from previous runs that don't
    stop when told to for some reason

### Exec

- `docker compose exec stac-browser /bin/sh` - exec into the raw image

### Build

- `docker compose build --no-cache` - force rebuild from scratch for all
  images
- `docker compose build `<service_name> - only rebuild the image that
  corresponds to `service_name` in your docker-compose.yml.

### Verbose

- `docker compose --verbose up -f docker-compose.custom.yml`

### logs

- `docker compose logs`

### Config

- `docker compose config` - outputs the final configuration derived from
  your docker-compose.yml and .env files, after all variable
  substitutions

### Env Vars

- `DOCKER_CLI_EXPERIMENTAL=enabled`
- `COMPOSE_DOCKER_CLI_BUILD=1` - make Docker Compose use the Docker CLI
  directly for image building, which might provide more detailed output

## Install/Enable Docker Server (Daemon)

- `sudo systemctl start docker`
  - `sudo systemctl enable docker` - for CentOS7
- `journalctl -u docker` - daemon log for CentOS7

## Terraform using Docker

    provider "docker" {
      # Provider configuration
    }

    resource "docker_image" "nginx" {
      name         = "nginx:latest"
      keep_locally = false
    }

    resource "docker_container" "nginx" {
      image = docker_image.nginx.latest
      name  = "nginx-example"
      ports {
        internal = 80
        external = 8080
      }
    }

## NVIDIA Images

- [ngc.nvidia.com](https://ngc.nvidia.com) - NGC registry page
- On website registry page, generate NGC API key then store it somewhere
  safe
  - I put it in a file in the home directory
  - not a big deal if you lose it, just generate a new one, old key
    become invalid
- then, on server: `docker login nvcr.io` - user \$oauthtoken, password
  use NGC API key

## Reference

- [Docker Glossary](https://docs.docker.com/glossary/)
- [Get
  Started](https://github.com/docker/docker.github.io/tree/master/get-started) -
  End-to-end tutorial
- [How to create-react-app with
  docker](https://www.peterbe.com/plog/how-to-create-react-app-with-docker)
- [react-base](https://hub.docker.com/r/bayesimpact/react-base/) Base
  docker image for using React on top of npm
- [Another](https://hub.docker.com/r/ontouchstart/react/)
- [WhatWhyHow Introduction to
  Docker](https://kulkarniamit.github.io/whatwhyhow/one-zero-one/introduction-to-docker.html)
- [13 Docker Tricks You Didn’t
  Know](https://overcast.blog/13-docker-tricks-you-didnt-know-47775a4f678f)
