MIGRATIONS_DIR = "./db/migrations"

# Extract the name from the arguments
# This takes everything after the 'create' command
MIGRATION_NAME := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

migration-create:
	goose -dir $(MIGRATIONS_DIR) create $(MIGRATION_NAME) sql 

%:
	@:
