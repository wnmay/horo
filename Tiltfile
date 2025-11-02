load('ext://restart_process', 'docker_build_with_restart')

k8s_yaml('./infra/development/k8s/payment-service/secrets.yaml')

payment_compile_cmd = (
    'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 '
    'go build -o build/payment-service ./services/payment-service/cmd/main.go'
)

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
# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

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
