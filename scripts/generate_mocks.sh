#!/bin/bash

# Generate mocks for all repository interfaces and services

set -e

MOCKS_DIR="internal/usecase/mocks"

# Create mocks directory
mkdir -p "$MOCKS_DIR"

echo "Generating mocks..."

# Generate mocks for domain repositories
go run go.uber.org/mock/mockgen@latest -source=internal/domain/user_repository.go -destination="$MOCKS_DIR/mock_user_repository.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/domain/publication_repository.go -destination="$MOCKS_DIR/mock_publication_repository.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/domain/comment_repository.go -destination="$MOCKS_DIR/mock_comment_repository.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/domain/media_repository.go -destination="$MOCKS_DIR/mock_media_repository.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/domain/recommendation_repository.go -destination="$MOCKS_DIR/mock_recommendation_repository.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/domain/tag_repository.go -destination="$MOCKS_DIR/mock_tag_repository.go" -package=mocks

# Generate mocks for infrastructure services
go run go.uber.org/mock/mockgen@latest -source=internal/infrastructure/jwt/token_interface.go -destination="$MOCKS_DIR/mock_token_service.go" -package=mocks
go run go.uber.org/mock/mockgen@latest -source=internal/infrastructure/ai/client_interface.go -destination="$MOCKS_DIR/mock_ai_client.go" -package=mocks

echo "Mocks generated successfully!"

