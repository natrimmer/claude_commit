# Release Process

This document describes how to create a new release of Claude Commit that will automatically build and publish binaries for different platforms.

## Creating a New Release

1. Update the code and make any necessary changes
2. Commit and push your changes to the main branch
3. Create and push a new tag with a version number (following semver):

```bash
git tag v0.1.0  # Change version as appropriate
git push origin v0.1.0
```

4. GitHub Actions will automatically:
   - Build binaries for all supported platforms (Linux, macOS, Windows)
   - Create a GitHub Release
   - Upload the binaries to the release

## Supported Platforms

The GitHub Actions workflow builds binaries for:

- Linux (amd64, arm64)
- macOS (amd64 - Intel, arm64 - Apple Silicon)
- Windows (amd64)
