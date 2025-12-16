#!/bin/bash

# Sense API - cURL Examples
# Server: http://174.138.14.66:8080
# Replace YOUR_TOKEN with actual JWT token after login

BASE_URL="http://174.138.14.66:8080"
TOKEN="YOUR_TOKEN_HERE"

echo "=== Health Check ==="
curl -X GET "${BASE_URL}/health"

echo -e "\n\n=== Authentication ==="

# Register new user
echo "Register new user:"
curl -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "securepass123",
    "phone": "+375291234567",
    "description": "Test user description"
  }'

# Login
echo -e "\n\nLogin:"
curl -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "securepass123"
  }'

# Check auth (requires token)
echo -e "\n\nCheck auth:"
curl -X GET "${BASE_URL}/auth/check" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== Publications ==="

# Create publication
echo "Create publication:"
curl -X POST "${BASE_URL}/publication" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "post",
    "content": "This is my first post!",
    "visibility": "public",
    "tags": ["photography", "art"]
  }'

# Get publication by ID (replace {id} with actual publication ID)
echo -e "\n\nGet publication:"
curl -X GET "${BASE_URL}/publication/{id}" \
  -H "Authorization: Bearer ${TOKEN}"

# Update publication
echo -e "\n\nUpdate publication:"
curl -X PUT "${BASE_URL}/publication/{id}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated content",
    "visibility": "public"
  }'

# Like publication
echo -e "\n\nLike publication:"
curl -X POST "${BASE_URL}/publication/{id}/like" \
  -H "Authorization: Bearer ${TOKEN}"

# Get likes
echo -e "\n\nGet publication likes:"
curl -X GET "${BASE_URL}/publication/{id}/likes?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

# Save publication
echo -e "\n\nSave publication:"
curl -X POST "${BASE_URL}/publication/{id}/save" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "note": "Interesting article"
  }'

# Delete publication
echo -e "\n\nDelete publication:"
curl -X DELETE "${BASE_URL}/publication/{id}" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== Feed ==="

# Get public feed
echo "Get public feed:"
curl -X GET "${BASE_URL}/feed?limit=20&offset=0"

# Get user feed (replace {id} with user ID)
echo -e "\n\nGet user feed:"
curl -X GET "${BASE_URL}/feed/user/{id}?limit=20&offset=0"

# Get my feed (requires auth)
echo -e "\n\nGet my feed:"
curl -X GET "${BASE_URL}/feed/me?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

# Get saved publications
echo -e "\n\nGet saved publications:"
curl -X GET "${BASE_URL}/feed/me/saved?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== Profile ==="

# Get my profile
echo "Get my profile:"
curl -X GET "${BASE_URL}/profile/me" \
  -H "Authorization: Bearer ${TOKEN}"

# Get user profile (replace {id} with user ID)
echo -e "\n\nGet user profile:"
curl -X GET "${BASE_URL}/profile/{id}" \
  -H "Authorization: Bearer ${TOKEN}"

# Update profile
echo -e "\n\nUpdate profile:"
curl -X PUT "${BASE_URL}/profile/me" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description",
    "phone": "+375291234568"
  }'

# Follow user
echo -e "\n\nFollow user:"
curl -X POST "${BASE_URL}/profile/{id}/follow" \
  -H "Authorization: Bearer ${TOKEN}"

# Unfollow user
echo -e "\n\nUnfollow user:"
curl -X DELETE "${BASE_URL}/profile/{id}/follow" \
  -H "Authorization: Bearer ${TOKEN}"

# Get followers
echo -e "\n\nGet followers:"
curl -X GET "${BASE_URL}/profile/{id}/followers?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

# Get following
echo -e "\n\nGet following:"
curl -X GET "${BASE_URL}/profile/{id}/following?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== Comments ==="

# Create comment
echo "Create comment:"
curl -X POST "${BASE_URL}/comment" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "publication_id": "{publication_id}",
    "content": "Great post!",
    "parent_id": null
  }'

# Get comments for publication
echo -e "\n\nGet comments:"
curl -X GET "${BASE_URL}/comment/publication/{publication_id}?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

# Update comment
echo -e "\n\nUpdate comment:"
curl -X PUT "${BASE_URL}/comment/{comment_id}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated comment"
  }'

# Delete comment
echo -e "\n\nDelete comment:"
curl -X DELETE "${BASE_URL}/comment/{comment_id}" \
  -H "Authorization: Bearer ${TOKEN}"

# Like comment
echo -e "\n\nLike comment:"
curl -X POST "${BASE_URL}/comment/{comment_id}/like" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== Media ==="

# Upload media
echo "Upload media:"
curl -X POST "${BASE_URL}/media" \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "file=@/path/to/image.jpg" \
  -F "type=image"

# Get media by ID
echo -e "\n\nGet media:"
curl -X GET "${BASE_URL}/media/{media_id}" \
  -H "Authorization: Bearer ${TOKEN}"

# Delete media
echo -e "\n\nDelete media:"
curl -X DELETE "${BASE_URL}/media/{media_id}" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n\n=== AI Features ==="

# Get recommendations
echo "Get recommendations:"
curl -X POST "${BASE_URL}/recommendations" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "limit": 10,
    "offset": 0
  }'

# Get recommendations feed
echo -e "\n\nGet recommendations feed:"
curl -X GET "${BASE_URL}/recommendations/feed?limit=20&offset=0" \
  -H "Authorization: Bearer ${TOKEN}"

# Hide recommendation
echo -e "\n\nHide recommendation:"
curl -X POST "${BASE_URL}/recommendations/{id}/hide" \
  -H "Authorization: Bearer ${TOKEN}"

# Purify text
echo -e "\n\nPurify text:"
curl -X POST "${BASE_URL}/purify" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Some text to moderate"
  }'

echo -e "\n\n=== Logout ==="
curl -X POST "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer ${TOKEN}"


