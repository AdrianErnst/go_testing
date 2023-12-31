docker_prune_settings( disable = False , max_age_mins = 360 , num_builds = 0 , interval_hrs = 1 , keep_recent = 2 ) 

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')
load('extensions/go-test-tiltfile', 'test_go')
load('extensions/k8s-yaml-object-selectors-tiltfile', 'k8s_yaml_object_selectors')
load('extensions/secrets-tiltfile', 'deploy_secrets')

# create/delete cluster and registry
setup_grp = "Setup"
ctlptl_filepath = 'ctlptl-kind.yaml'
registry_name = 'ctlptl-registry'
if config.tilt_subcommand == 'up':
    ctlptl_cmd = 'ctlptl apply -f {file}'.format(file=ctlptl_filepath)
    local_resource(
        "ctl_kind_cluster_registry",
        ctlptl_cmd,
        trigger_mode=TRIGGER_MODE_MANUAL,
        labels=[setup_grp]
    )
if config.tilt_subcommand == 'down':
    local('ctlptl delete -f {} | true'.format(ctlptl_filepath), quiet=True, echo_off=True)
    local('ctlptl delete registry {}'.format(registry_name), quiet=True, echo_off=True)

# generate app code
openapi_yaml = 'openapi.yaml'
openapi_generate_cmd = 'docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i ./local/{} \
    -g go-gin-server \
    -o /local/cmd/web/openapi-gin \
    --additional-properties=packageName=openapi'.format(openapi_yaml, PWD = 'PWD')

openapi_generate_docu_cmd = 'docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i ./local/{} \
    -g html2 \
    -o /local/docs/generated/'.format(openapi_yaml, PWD = 'PWD')

generation_grp = "Generate"
openapi_generate_name = 'generate_openapi'
openapi_generate_docu_name = openapi_generate_name + '_docu'

local_resource(
    openapi_generate_name,
    openapi_generate_cmd,
    deps=[openapi_yaml],
    labels=[generation_grp]
)

local_resource(
    openapi_generate_docu_name,
    openapi_generate_docu_cmd,
    deps=[openapi_yaml],
    labels=[generation_grp]
)

# build app
compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/go_k8s_client ./cmd/web/'
if os.name == 'nt': compile_cmd = './scripts/build.bat'

build_grp = "Build"
local_resource(
    'go-build',
    compile_cmd,
    deps = ['cmd', "pkg"],
    ignore= [
        '**/test',
        '**/tests',
    ],
    resource_deps = [
        'ctl_kind_cluster_registry',
        openapi_generate_name,
        openapi_generate_docu_name],
    labels=[build_grp]
)

# test app
test_grp = "Tests"
test_go(
    "test-runner",
    ".",
    ".",
    recursive=True,
    extra_args=['-v', '-coverpkg=./...'],
    resource_deps=['go-build'],
    ignore=['**/*_integration_test.go', 'deployments'],
    labels=[test_grp]
)

# integration tests out of cluster setup
test_go(
    "outcluster-integration-test-runner",
    ".",
    ".",
    tags=['integration'],
    recursive=True,
    extra_args=['-v', '-coverpkg=./...', '-external'],
    trigger_mode=TRIGGER_MODE_MANUAL,
    labels=[test_grp]
)

def inClusterIntegrationsTests():
    namespace = 'client-test'
    image = 'go-k8s-client-testing'
    pod_name = 'k8s-client-test'
    helm_test_blob = helm('./deployments/helm-test', namespace=namespace)
    res_name = "incluster-test"
    docker_build(image, '.', dockerfile='./deployments/test.Dockerfile')
    k8s_yaml(helm_test_blob)
    k8s_resource(pod_name, res_name, auto_init=False,
                 trigger_mode=TRIGGER_MODE_MANUAL,
                 objects=k8s_yaml_object_selectors(
                    helm_test_blob,
                    ignore={'{}'.format(pod_name):bool},
                    extra_resources=[
                        "{}:serviceaccount".format(pod_name),
                    ]
                ),
                labels=[test_grp]
    )
    images_cmd = "docker images '*/{}' -a -q".format(image)
    delete_testing_images_cmd = "docker rmi $({})".format(images_cmd)
    cmd = "if [ -n $({}) ]; then {}; fi".format(images_cmd, delete_testing_images_cmd)
    local_resource('delete testing images', cmd,
     resource_deps=[res_name], labels=[test_grp],
     trigger_mode=TRIGGER_MODE_MANUAL)
inClusterIntegrationsTests()

# deploy app with live update and restart
deployments_grp = 'Deployments'
docker_build_with_restart(
    'go_k8s_client',
    '.',
    entrypoint = ['/app/build/go_k8s_client'],
    build_args = {"PORT":"9292"},
    dockerfile = './deployments/local.Dockerfile',
    only=[
        './build',
        './docs'
    ],
    live_update = [
        sync('./build', '/app/build'),
        sync('./docs', '/app/docs'),
    ],
)


# deploy local secrets into cluster
local_secrets = {}
renovate_secret = []
deploy_secrets("secrets", local_secrets)

# manage deployments into tilt resource
helm_blob = helm('./deployments/helm', namespace='client')
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
    labels=[deployments_grp]
)

k8s_resource(
    'go-k8s-client',
    port_forwards = 9292,
    resource_deps = ['go-build'],
    labels=[deployments_grp]
)
