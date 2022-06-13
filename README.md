## Build

```bash
source .env
docker build -t ${DOCKER_BUILD_NAME}:${DOCKER_BUILD_TAG} -t ${DOCKER_BUILD_NAME}:latest .
```

## Run

```bash
source .env
docker run -d -p 8084:8084 --name=godanuk --net=appdangibbsuk_sail ${DOCKER_BUILD_NAME}:${DOCKER_BUILD_TAG}
```

## Request example

```bash
curl -s -X POST http://localhost:8084/tools/pwgen \
   -H 'Content-Type: application/json' \
   -d '{"num-passwords":1, "length": 32}' | jq
```
