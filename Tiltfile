analytics_settings(enable=False)
load('ext://color', 'color')

# Remember to add `local.tiltfile` to .gitignore
if os.path.exists('Tiltfile.local'):
  load_dynamic('Tiltfile.local')

if config.tilt_subcommand == 'down':
  print('Goodbye, closing down')

## Check Min Version
version_settings(check_updates=True, constraint='>=0.30.0')
os.putenv("MAGEFILE_HASHFAST", '1')



docker_prune_settings( disable = False , max_age_mins = 360 , num_builds = 10  , keep_recent = 2 )
# Some common test configurations to easily get up and running with Tilt for Tests
# (see docs: https://docs.tilt.dev/tests_in_tilt.html)

def test_go(name, package, deps, timeout='', tags=None, mod='', recursive=False, ignore=None,  extra_args=None, **kwargs):
    if recursive and not package.endswith('...'):
        package = package.rstrip('/')
        package = package + '/...'

    timeout_str = ''
    if timeout:  # expects a go-parsable dur
        timeout_str = '-timeout {}'.format(timeout)
    tags_str = ''
    if tags:
        tags_str = '-tags {}'.format(','.join(tags))

    mod_str = ''
    if mod:
        mod_str = '-mod {}'.format(mod)

    extra_args_str = ''
    if extra_args:
        extra_args_str = ' '.join(extra_args)

    if not ignore:
        ignore=[]

    # TODO: ability to pass multiple packages (as in `go test`)

    cmd = 'go test {mod_str} {tags_str} {timeout_str} {extra_args_str} {package}'.format(
          mod_str=mod_str, tags_str=tags_str, timeout_str=timeout_str,
          extra_args_str=extra_args_str, package=package)

    local_resource(name, cmd, deps=deps, ignore=ignore,labels=['test'], allow_parallel=True, **kwargs)
## Tool Requirements

# block Tiltfile execution if missing required tool (e.g. Helm)

def require_tool(tool):
    tool = shlex.quote(tool)
    local(
        command='command -v {tool} >/dev/null 2>&1 && echo "✅ {tool} available" || {{ echo >&2 "❗ {tool} is required but was not found in PATH"; exit; }}'.format( # Exit 1 to stop loading tilt
            tool=tool
        ),
        # `cmd_bat`, when present, is used instead of `cmd` on Windows.
        command_bat=[
            "powershell.exe",
            "-Noninteractive",
            "-Command",
            '& {{if (!(Get-Command {tool} -ErrorAction SilentlyContinue)) {{ Write-Error "{tool} is required but was not found in PATH"; exit 1 }}}}'.format(
                tool=tool
            ),
        ],
        echo_off=True,
    )
    return



require_tool("helm")
require_tool("go")
require_tool("docker")
require_tool("kubectl")
require_tool("mage")
require_tool("direnv")


allow_k8s_contexts('dsvtest') # ensure dsvtest is scoped so we don't load against anything other local test environment
def current_namespace():
  namespace=str(local("kubectl config view --minify --output 'jsonpath={..namespace}'",echo_off=True,quiet=True))
  if namespace != "dsv":
    print(color.red("""WARNING: You are not in the dsv namespace."""))
  return namespace
print("""
-----------------------------------------------------------------
✨ Hello Tilt!
-----------------------------------------------------------------

""".strip() + color.green("\n\nCurrent K8 Namespace: ") + color.green(current_namespace())
)


load('ext://uibutton', 'cmd_button', 'location', 'text_input')
cmd_button(
  name='mage-init-button',
  argv=['mage', 'init'],
  text='init',
  location=location.NAV,
  icon_name='system_update'
)
cmd_button(
  name='mage-doctor-button',
  argv=['mage', 'go:doctor'],
  text='go:doctor',
  location=location.NAV,
  icon_name='medication'
)

local_resource(
  "mage:refreshtasks",
  cmd="zsh -l -c \"mage -f -l\"",
  trigger_mode=TRIGGER_MODE_AUTO,
  auto_init=True,
  deps=['magefiles/*.go'],
  labels=["startup"]
)
local_resource(
  "direnv:allow",
  cmd="zsh -l -c \"direnv allow\"",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=['.envrc'],
  auto_init=True,
  labels=["startup"]
)
local_resource(
  "mage:init",
  cmd="mage init",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=[],
  auto_init=False,
  labels=["setup"]
)
test_go(
  "go-test",
  trigger_mode=TRIGGER_MODE_MANUAL,
  auto_init=False,
  package="./...",
  deps=[]
)

local_resource(
  "job:init",
  cmd="mage job:init",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=['.envrc'],
  auto_init=True,
  labels=["setup"],
)
local_resource(
  "job:redeploy",
  cmd="mage job:redeploy",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=['.cache/'],
  resource_deps=['minikube:init'],
  auto_init=True,
  labels=["deploy"],
)
local_resource(
  "k8s:logs",
  serve_cmd="mage k8s:logs",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=['.cache/'],
  auto_init=True,
  labels=["logs"],
)
local_resource(
  "optional:stern-without-mage",
  serve_cmd="stern --kubeconfig .cache/config --namespace dsv  --timestamps . ",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=[],
  auto_init=False,
  labels=["logs"],
)
local_resource(
  "minikube:init",
  cmd="mage minikube:init",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=[],
  auto_init=False,
  labels=["setup"],
)
local_resource(
  "minikube:destroy",
  cmd="mage minikube:destroy",
  trigger_mode=TRIGGER_MODE_MANUAL,
  deps=[],
  auto_init=False,
  labels=["setup"],
)

# k8s_resource('injector', resource_deps='minikube:init', pod_readiness='ignore')