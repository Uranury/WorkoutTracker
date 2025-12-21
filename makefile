MIGRATIONS_DIR := migrations

# Usage: make new-migration NAME=create_users_table
new-migration:
ifndef NAME
	$(error NAME is not set. Usage: make new-migration NAME=create_users_table)
endif
	@echo "Creating new migration: $(NAME)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
	@echo "Migration files created:"
	@ls -1 $(MIGRATIONS_DIR)/*$(NAME).sql