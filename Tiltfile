docker_compose('integration_tests/docker-compose.yml')

local_resource("🔁 tests",
  cmd = "go test -v ./... -count=1",
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

local_resource("teedy-cli backup",
  cmd = "source .env_backup && rm -rf backup_test/ && go run cmd/teedy-cli/main.go backup --url 'http://localhost:8080' --destinationpath './backup_test' && tree backup_test/",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
)

local_resource("teedy-cli restore",
  cmd = "go run cmd/teedy-cli/main.go restore --url 'http://localhost:8080' --sourcepath './backup_test'",
  auto_init = False,
  trigger_mode = TRIGGER_MODE_MANUAL,
  env = {
    "TEEDY_USERNAME": "admin",
    "TEEDY_PASSWORD": "superSecure"
  },
)