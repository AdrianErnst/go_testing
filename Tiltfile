docker_prune_settings( disable = False , max_age_mins = 360 , num_builds = 0 , interval_hrs = 1 , keep_recent = 2 ) 

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

# Records the current time, then kicks off a server update.
# Normally, you would let Tilt do deploys automatically, but this
# shows you how to set up a custom workflow that measures it.
local_resource(
    'deploy',
    'python3 ./scripts/record-start-time.py',
)

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/go_k8s_client ./'
if os.name == 'nt': compile_cmd = 'build.bat'

local_resource(
    'go_k8s_client',
    compile_cmd,
    deps=['./main.go', './start.go'],
    resource_deps = ['deploy']
)

docker_build_with_restart(
    'go_k8s_client',
    '.',
    entrypoint=['/app/build/go_k8s_client'],
    build_args={"PORT":"9292"},
    dockerfile='./Dockerfile_local',
    only=[
        './build',
        './web',
    ],
    live_update=[
        sync('./build', '/app/build'),
    ],
)
k8s_yaml(helm('./helm'))
k8s_resource('go-k8s-client',
    port_forwards=9292,
    resource_deps=['deploy', 'go_k8s_client']
)