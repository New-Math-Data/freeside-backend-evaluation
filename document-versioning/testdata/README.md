# Test Data

This directory contains sample JSON documents for testing the Document Versioning API.

## Files

- `sample_document.json` - Initial document (version 1)
- `sample_document_v2.json` - Updated document (version 2) - adds fields, modifies capacity
- `sample_document_v3.json` - Final document (version 3) - marks as decommissioned

## Usage

You can use these files to test your API implementation:

```bash
# Create a document
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Asset",
    "content": '"$(cat testdata/sample_document.json)"'
  }'

# Update the document (replace {id} with the returned ID)
curl -X PUT http://localhost:8080/api/v1/documents/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "content": '"$(cat testdata/sample_document_v2.json)"'
  }'

# Get version history
curl http://localhost:8080/api/v1/documents/{id}/versions

# Get document at version 1
curl "http://localhost:8080/api/v1/documents/{id}?version=1"

# Revert to version 1
curl -X POST http://localhost:8080/api/v1/documents/{id}/revert \
  -H "Content-Type: application/json" \
  -d '{"version": 1}'
```
