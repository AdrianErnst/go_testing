docker_prune_settings( disable = False , max_age_mins = 360 , num_builds = 0 , interval_hrs = 1 , keep_recent = 2 ) 

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')
load('extensions/go-test-tiltfile', 'test_go')
load('extensions/k8s-yaml-object-selectors-tiltfile', 'k8s_yaml_object_selectors')
load('extensions/secrets-tiltfile', 'deploy_secrets')

# create/delete cluster and registry
ctlptl_filepath = 'ctlptl-kind.yaml'
registry_name = 'ctlptl-registry'
if config.tilt_subcommand == 'up':
    ctlptl_cmd = 'ctlptl apply -f {file}'.format(file=ctlptl_filepath)
    local_resource(
        "ctl_kind_cluster_registry",
        ctlptl_cmd,
        trigger_mode=TRIGGER_MODE_MANUAL
    )
if config.tilt_subcommand == 'down':
    local('ctlptl delete -f {} | true'.format(ctlptl_filepath), quiet=True, echo_off=True)
    local('ctlptl delete registry {}'.format(registry_name), quiet=True, echo_off=True)

# build app
compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/go_k8s_client ./cmd/web/'
if os.name == 'nt': compile_cmd = './scripts/build.bat'

local_resource(
    'go-build',
    compile_cmd,
    deps = ['cmd', "pkg"],
    ignore= [
        '**/test',
        '**/tests',
    ],
    resource_deps = ['ctl_kind_cluster_registry']
)

# test app
test_go(
    "test-runner",
    ".",
    ".",
    recursive=True,
    extra_args=['-v', '-coverpkg=./...'],
    resource_deps=['go-build'],
    ignore=['**/*_integration_test.go']
)

# integration tests in cluster setup
test_go(
    "incluster-integration-test-runner",
    ".",
    ".",
    tags=['integration', 'incluster'],
    recursive=True,
    extra_args=['-v', '-coverpkg=./...'],
    trigger_mode=TRIGGER_MODE_MANUAL
)

# integration tests out of cluster setup
test_go(
    "outcluster-integration-test-runner",
    ".",
    ".",
    tags=['integration', 'outcluster'],
    recursive=True,
    extra_args=['-v', '-coverpkg=./...'],
    trigger_mode=TRIGGER_MODE_MANUAL
)

# deploy app with live update and restart
docker_build_with_restart(
    'go_k8s_client',
    '.',
    entrypoint = ['/app/build/go_k8s_client'],
    build_args = {"PORT":"9292"},
    dockerfile = './deployments/local.Dockerfile',
    only=[
        './build',
    ],
    live_update = [
        sync('./build', '/app/build'),
    ],
)


# deploy local secrets into cluster
local_secrets = {}
renovate_secret = []
deploy_secrets("secrets", local_secrets)

# manage deployments into tilt resource
helm_blob = helm('./deployments/helm', values='./deployments/helm/values.yaml')
k8s_yaml(helm_blob)
k8s_resource(
    '',
    'local-helm',
    objects = k8s_yaml_object_selectors(
        helm_blob,
        ignore={'go-k8s-client':bool},
        extra_resources=[
            'go-k8s-client:serviceaccount',
        ]
    ),
)

k8s_resource(
    'go-k8s-client',
    port_forwards = 9292,
    resource_deps = ['go-build']
)
