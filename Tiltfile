load('ext://restart_process', 'docker_build_with_restart')

# Load centralized secrets (shared by all services)
k8s_yaml('./infra/development/k8s/secrets.yaml')

### API Gateway ###
# load config
k8s_yaml('./infra/development/k8s/api-gateway/config.yaml')

chat_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway/cmd/main.go'

local_resource(
  'api-gateway-compile',
  chat_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")

docker_build_with_restart(
  'horo/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway/deployment.yaml')
k8s_resource('api-gateway', port_forwards=8080, resource_deps=['api-gateway-compile'], labels="services")

### End of API Gateway ###

### Payment Service ###
k8s_yaml('./infra/development/k8s/payment-service/config.yaml')
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
# load config
k8s_yaml('./infra/development/k8s/chat-service/config.yaml')

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
# Load config
k8s_yaml('./infra/development/k8s/order-service/config.yaml')

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

### User Management Service ###

# Load config
k8s_yaml('./infra/development/k8s/user-management-service/config.yaml')

# Create firebase-key secret from local file
local_resource(
  'firebase-key-secret',
  'kubectl create secret generic firebase-key-secret --from-file=firebase-key.json=./firebase-key.json --dry-run=client -o yaml | kubectl apply -f -',
  labels=["secrets"]
)

user_management_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/user-management-service ./services/user-management-service/cmd/main.go'

local_resource(
  'user-management-service-compile',
  user_management_compile_cmd,
  deps=['./services/user-management-service', './shared'], 
  labels=["compiles"]
)

docker_build_with_restart(
  'horo/user-management-service',
  '.',
  entrypoint=['/app/build/user-management-service'],
  dockerfile='./infra/development/docker/user-management-service.Dockerfile',
  only=[
    './build/user-management-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/user-management-service/deployment.yaml')

k8s_resource(
  'user-management-service', 
  resource_deps=['user-management-service-compile', 'firebase-key-secret'], 
  labels=["services"]
)

### End of User Management Service ###

### Course Service ###

# load secrets
k8s_yaml('./infra/development/k8s/course-service/secrets.yaml')

course_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/course-service ./services/course-service/cmd/main.go'

local_resource(
  'course-service-compile',
  course_compile_cmd,
  deps=['./services/course-service', './shared'], labels="compiles")

docker_build_with_restart(
  'horo/course-service',
  '.',
  entrypoint=['/app/build/course-service'],
  dockerfile='./infra/development/docker/course-service.Dockerfile',
  only=[
    './build/course-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/course-service/deployment.yaml')
k8s_resource('course-service', resource_deps=['course-service-compile'], labels="services")

### End of Course Service ###