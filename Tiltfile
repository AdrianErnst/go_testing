docker_prune_settings( disable = False , max_age_mins = 360 , num_builds = 0 , interval_hrs = 1 , keep_recent = 2 ) 

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')
load('extensions/go-test-tiltfile', 'test_go')
load('extensions/k8s-yaml-object-selectors-tiltfile', 'k8s_yaml_object_selectors')

test_go("test-runner", ".", ".", recursive=True)

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/go_k8s_client ./cmd/web/'
if os.name == 'nt': compile_cmd = './scripts/build.bat'

local_resource(
    'go-build',
    compile_cmd,
    deps = ['cmd', "pkg"],
    ignore= [
        '**/test',
    ],
    resource_deps = []
)

docker_build_with_restart(
    'go_k8s_client',
    '.',
    entrypoint = ['/app/build/go_k8s_client'],
    build_args = {"PORT":"9292"},
    dockerfile = './deployments/Dockerfile_local',
    only=[
        './build',
    ],
    live_update = [
        sync('./build', '/app/build'),
    ],
)

helm_blob = helm('./deployments/helm', values='./deployments/helm/values.yaml')
k8s_yaml(helm_blob)

k8s_resource(
    '',
    'local-helm',
    objects = k8s_yaml_object_selectors(
        helm_blob,
        ignore={'go-k8s-client':bool},
        extra_resources=['go-k8s-client:serviceaccount']
    ),
)

k8s_resource(
    'go-k8s-client',
    port_forwards = 9292,
    resource_deps = ['go-build']
)