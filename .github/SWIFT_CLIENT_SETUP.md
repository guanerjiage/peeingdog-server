# Swift Client Auto-Generation Setup

This workflow automatically generates a Swift client SDK from the OpenAPI spec and pushes it to a separate repository.

## Prerequisites

1. **Create a separate Swift client repository** (e.g., `peeingdog-swift-client`)
   - Initialize with a basic structure or empty repo
   - Recommended structure after first generation will be auto-created

2. **Generate a Personal Access Token (PAT)**
   - Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Select scopes: `repo` (full control of private repositories)
   - Copy the token

## GitHub Secrets Configuration

Add these secrets to the **peeingdog-server** repository (Settings → Secrets and variables → Actions):

| Secret Name | Value |
|-------------|-------|
| `SWIFT_CLIENT_REPO` | `username/peeingdog-swift-client` |
| `SWIFT_CLIENT_REPO_TOKEN` | Your Personal Access Token from above |

## How It Works

The workflow triggers on:
- ✅ Push to `main` branch with changes to `openapi.yaml` or this workflow file
- ✅ Manual trigger via "Run workflow" button in Actions tab

### Steps:
1. Checks out peeingdog-server repo
2. Installs OpenAPI Generator via Homebrew
3. Generates Swift5 client from `openapi.yaml`
4. Checks out the separate Swift client repo
5. Replaces old generated files with new ones
6. Auto-commits and pushes changes (with `[skip ci]` to prevent recursion)

## Generated Files Structure

```
peeingdog-swift-client/
├── OpenAPIClient/
│   ├── Classes/
│   │   └── OpenAPIs/
│   │       ├── APIs/        # API endpoint classes (UsersAPI, MessagesAPI, etc.)
│   │       ├── Models/      # Data models (User, Message, etc.)
│   │       └── *.swift      # Helper files (Configuration, Extensions, etc.)
│   └── *.podspec           # CocoaPods configuration
├── docs/                    # API documentation
├── Package.swift           # Swift Package Manager config
└── openapi.yaml           # Source API spec
```

## Usage in iOS Projects

### CocoaPods
```ruby
# Podfile
pod 'OpenAPIClient', :git => 'https://github.com/YOUR-USERNAME/peeingdog-swift-client.git'
```

### Swift Package Manager
Add to your Xcode project:
```
https://github.com/YOUR-USERNAME/peeingdog-swift-client.git
```

## Troubleshooting

- **Workflow not triggering**: Ensure the branch is `main` and `openapi.yaml` was modified
- **Permission denied on push**: Verify `SWIFT_CLIENT_REPO_TOKEN` has `repo` scope
- **Merge conflicts**: If manual edits were made to the Swift repo, pull before the next auto-generation

## Related Tasks

See also:
- `python:generate` - Generates Python client (same pattern can be adapted)
- `typescript:generate` - Generates TypeScript client (same pattern can be adapted)
