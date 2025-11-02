# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

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