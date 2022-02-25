```bash
export DOCKER_BUILD_NAME="gibbs/godanuk"
export DOCKER_BUILD_TAG="0.5"

# Build
docker build -t ${DOCKER_BUILD_NAME}:${DOCKER_BUILD_TAG} -t ${DOCKER_BUILD_NAME}:latest .

# Run
docker run -d -p 8084:8084 --name=godanuk --net=appdangibbsuk_sail ${DOCKER_BUILD_NAME}:${DOCKER_BUILD_TAG}

# Request example
curl -s -X POST http://localhost:8084/tools/pwgen \
   -H 'Content-Type: application/json' \
   -d '{"num-passwords":1, "length": 32}' | jq
```
