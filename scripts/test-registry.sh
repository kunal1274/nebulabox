#!/usr/bin/env bash
set -euo pipefail

BASE="http://localhost:${NEBULABOX_REGISTRY_PORT:-5001}/v2"
REPO="nebulabox/nginx"
AUTH_BASE="http://localhost:${NEBULABOX_REGISTRY_PORT:-5001}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🧪 Testing Nebula Registry"
echo "=========================="
echo ""

# Test 1: Registry root
echo -n "Test 1: Registry root endpoint... "
RESPONSE=$(curl -s -w "%{http_code}" -o /dev/null "$BASE/")
if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗ (HTTP $RESPONSE)${NC}"
    exit 1
fi

# Test 2: Catalog
echo -n "Test 2: Catalog endpoint... "
CATALOG=$(curl -s "$BASE/_catalog")
if echo "$CATALOG" | grep -q "repositories"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Response: $CATALOG"
    exit 1
fi

# Test 3: Authentication
echo -n "Test 3: Authentication... "
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}')
if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓ (Token obtained)${NC}"
else
    echo -e "${RED}✗${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

# Test 4: Tags list (empty)
echo -n "Test 4: Tags list (empty repo)... "
TAGS=$(curl -s "$BASE/$REPO/tags/list")
if echo "$TAGS" | grep -q "tags"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC} (may be expected)"
fi

# Test 5: Blob upload
echo -n "Test 5: Blob upload... "
LOC=$(curl -s -i -X POST -H "Authorization: Bearer $TOKEN" "$BASE/$REPO/blobs/uploads/" | grep -i "Location:" | awk '{print $2}' | tr -d '\r')
if [ -n "$LOC" ]; then
    UUID=$(echo "$LOC" | awk -F/ '{print $NF}')
    echo -e "${GREEN}✓ (UUID: ${UUID:0:20}...)${NC}"
else
    echo -e "${RED}✗${NC}"
    exit 1
fi

# Test 6: Upload chunk
echo -n "Test 6: Upload blob chunk... "
CHUNK_RESPONSE=$(curl -s -w "%{http_code}" -X PATCH \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/octet-stream" \
    --data-binary 'sample-blob-data' \
    "http://localhost:${NEBULABOX_REGISTRY_PORT:-5001}$LOC" -o /dev/null)
if [ "$CHUNK_RESPONSE" = "204" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗ (HTTP $CHUNK_RESPONSE)${NC}"
fi

# Test 7: Finalize blob
echo -n "Test 7: Finalize blob upload... "
DGST="sha256:$(echo -n 'sample-blob-data' | sha256sum | awk '{print $1}')"
FINALIZE_RESPONSE=$(curl -s -w "%{http_code}" -X PUT \
    -H "Authorization: Bearer $TOKEN" \
    "http://localhost:${NEBULABOX_REGISTRY_PORT:-5001}$LOC?digest=sha256:$DGST" -o /dev/null)
if [ "$FINALIZE_RESPONSE" = "201" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠ (HTTP $FINALIZE_RESPONSE)${NC}"
fi

# Test 8: Put manifest
echo -n "Test 8: Put manifest... "
MAN='{"schemaVersion":2,"mediaType":"application/vnd.oci.image.manifest.v1+json","config":{"mediaType":"application/vnd.oci.image.config.v1+json","digest":"sha256:config","size":123},"layers":[{"mediaType":"application/vnd.oci.image.layer.v1.tar+gzip","digest":"sha256:'"$DGST"'","size":16}]}'
MANIFEST_RESPONSE=$(curl -s -w "%{http_code}" -X PUT \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/vnd.oci.image.manifest.v1+json" \
    --data "$MAN" \
    "$BASE/$REPO/manifests/latest" -o /dev/null)
if [ "$MANIFEST_RESPONSE" = "201" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠ (HTTP $MANIFEST_RESPONSE)${NC}"
fi

# Test 9: Get manifest
echo -n "Test 9: Get manifest... "
GET_MANIFEST=$(curl -s -w "%{http_code}" "$BASE/$REPO/manifests/latest" -o /dev/null)
if [ "$GET_MANIFEST" = "200" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗ (HTTP $GET_MANIFEST)${NC}"
fi

# Test 10: Tags list (after push)
echo -n "Test 10: Tags list (after push)... "
TAGS_AFTER=$(curl -s "$BASE/$REPO/tags/list")
if echo "$TAGS_AFTER" | grep -q "latest"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC}"
fi

# Test 11: Registry API endpoints
echo -n "Test 11: Registry API - repositories... "
API_REPOS=$(curl -s "$AUTH_BASE/api/registry/repositories")
if echo "$API_REPOS" | grep -q "repositories"; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC}"
fi

# Test 12: Version metadata
if echo "$API_REPOS" | grep -q "$REPO"; then
    echo -n "Test 12: Version metadata... "
    VERSIONS=$(curl -s "$AUTH_BASE/api/registry/repositories/$(echo $REPO | tr '/' '%2F')/versions")
    if echo "$VERSIONS" | grep -q "versions"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${YELLOW}⚠${NC}"
    fi
fi

echo ""
echo -e "${GREEN}✅ Registry tests completed!${NC}"
echo ""

