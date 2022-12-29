### Local Mode

Traefik also offers a developer mode that can be used for temporary testing of plugins not hosted on GitHub.
To use a plugin in local mode, the Traefik static configuration must define the module name (as is usual for Go packages) and a path to a [Go workspace](https://golang.org/doc/gopath_code.html#Workspaces), which can be the local GOPATH or any directory.

The plugins must be placed in `./plugins-local` directory,
which should be in the working directory of the process running the Traefik binary.
The source code of the plugin should be organized as follows:

```
 ├── docker-compose.yml
 └── plugins-local
    └── src
        └── github.com
            └── ghnexpress
                └── traefik-ratelimit
                    ├── main.go
                    ├── vendor
                    ├── go.mod
                    └── ...

```

```yaml
# docker-compose.yml
version: "3.6"

services:
  memcached:
    image: launcher.gcr.io/google/memcached1
    container_name: some-memcached
    ports:
      - "11211:11211"
    networks:
      - traefik-network
  traefik:
    image: traefik:v2.9.6
    container_name: traefik
    depends_on:
      - memcached
    command:
      # - --log.level=DEBUG
      - --log.level=INFO
      - --api
      - --api.dashboard
      - --api.insecure=true
      - --providers.docker=true
      - --entrypoints.web.address=:80
      - --experimental.localPlugins.ratelimit.moduleName=github.com/ghnexpress/traefik-ratelimit
    ports:
      - "80:80"
      - "8080:8080"
    networks:
      - traefik-network
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./plugins-local/src/github.com/ghnexpress/traefik-ratelimit:/plugins-local/src/github.com/ghnexpress/traefik-ratelimit
    labels:
      - traefik.http.middlewares.rate-limit.plugin.ratelimit.memcachedConfig.address=some-memcached:11211
      - traefik.http.middlewares.rate-limit.plugin.ratelimit.windowTime=100
      - traefik.http.middlewares.rate-limit.plugin.ratelimit.maxRequestInWindow=10
  whoami:
    image: traefik/whoami
    container_name: simple-service
    depends_on:
      - traefik
    networks:
      - traefik-network
    labels:
      - traefik.enable=true
      - traefik.http.routers.whoami.rule=Host(`localhost`)
      - traefik.http.routers.whoami.entrypoints=web
      - traefik.http.routers.whoami.middlewares=rate-limit
networks:
  traefik-network:
    driver: bridge
```