docker_compose('integration_tests/docker-compose.yml')

local_resource("🔁 tests",
  cmd = "go test -v ./...",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
)

local_resource("🔁 tests short",
  cmd = "go test -v ./...",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
)

local_resource("nuke integration env",
  cmd = "cd integration_tests && docker-compose stop teedy-server && docker-compose rm --force teedy-server && docker-compose up -d",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
)

local_resource("go mod tidy",
  cmd = "go mod tidy",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
)