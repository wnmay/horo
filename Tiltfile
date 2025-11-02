load('ext://restart_process', 'docker_build_with_restart')

### Start Payment Service ###

payment_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/payment-service ./services/payment-service/cmd/main.go'

local_resource(
    'payment-service-compile',
    payment_compile_cmd,
    deps=[
        './services/payment-service',
        './shared',
    ],
    labels="compiles",
)

docker_build_with_restart(
    'horo/payment-service',
    '.',
    entrypoint=['/app/build/payment-service'],
    dockerfile='./infra/development/docker/payment-service.Dockerfile',
    only=[
        './build/payment-service',
        './shared',
    ],
    live_update=[
        sync('./build', '/app/build'),
        sync('./shared', '/app/shared'),
    ],
)

k8s_yaml('./infra/development/k8s/payment-service/deployment.yaml')

k8s_resource(
    'payment-service',
    resource_deps=['payment-service-compile'],
    labels="services",
)
### End Payment Service ###


### Chat Service ###
# load secrets
k8s_yaml('./infra/development/k8s/chat-service/secrets.yaml')

chat_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/chat-service ./services/chat-service/cmd/main.go'

local_resource(
  'chat-service-compile',
  chat_compile_cmd,
  deps=['./services/chat-service', './shared'], labels="compiles")

docker_build_with_restart(
  'horo/chat-service',
  '.',
  entrypoint=['/app/build/chat-service'],
  dockerfile='./infra/development/docker/chat-service.Dockerfile',
  only=[
    './build/chat-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/chat-service/deployment.yaml')
k8s_resource('chat-service', resource_deps=['chat-service-compile'], labels="services")

### End of Chat Service ###

### Order Service ###
# Load secrets
k8s_yaml('./infra/development/k8s/order-service/secrets.yaml')

order_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/order-service ./services/order-service/cmd/main.go'

local_resource(
  'order-service-compile',
  order_compile_cmd,
  deps=['./services/order-service', './shared'], labels="compiles")

docker_build_with_restart(
  'horo/order-service',
  '.',
  entrypoint=['/app/build/order-service'],
  dockerfile='./infra/development/docker/order-service.Dockerfile',
  only=[
    './build/order-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/order-service/deployment.yaml')
k8s_resource('order-service', resource_deps=['order-service-compile'], labels="services")

### End of Order Service ###
