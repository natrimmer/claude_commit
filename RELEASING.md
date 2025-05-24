# Release Process

This document describes how to create a new release of Claude Commit that will automatically build and publish binaries for different platforms.

## Quick Release (Recommended)

The development environment includes automated version increment scripts that handle the entire release process:

```bash
# For bug fixes and small changes
patch    # Increments v1.2.3 → v1.2.4

# For new features (backwards compatible)
minor    # Increments v1.2.3 → v1.3.0

# For breaking changes
major    # Increments v1.2.3 → v2.0.0
```

### What These Scripts Do

1. **Safety checks**: Ensure your working directory is clean (no uncommitted changes)
2. **Version calculation**: Automatically determine the next version number
3. **Confirmation prompt**: Show you what will happen and ask for confirmation
4. **Tag creation**: Create the new git tag locally
5. **Push tag**: Push the tag to trigger the automated build

Each script will show you:

- Current version
- Proposed new version
- Warning that this triggers a release build
- Confirmation prompt (defaults to "No" for safety)

## Manual Release Process

If you prefer to handle versioning manually:

1. Update the code and make any necessary changes
2. Commit and push your changes to the main branch
3. Create and push a new tag with a version number (following semver):

```bash
git tag v0.1.0  # Change version as appropriate
git push origin v0.1.0
```

## Automated Build Process

Once a tag is pushed (either via the scripts or manually), GitHub Actions will automatically:

- Build binaries for all supported platforms (Linux, macOS, Windows)
- Create a GitHub Release
- Upload the binaries to the release

## Supported Platforms

The GitHub Actions workflow builds binaries for:

- Linux (amd64, arm64)
- macOS (amd64 - Intel, arm64 - Apple Silicon)
- Windows (amd64)

## Version Strategy

We follow [Semantic Versioning (semver)](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Rollback

If you need to remove a tag:

```bash
# Delete local tag
git tag -d v1.2.3

# Delete remote tag
git push origin --delete v1.2.3
```

**Note**: If GitHub Actions has already created a release, you'll need to delete it manually from the GitHub web interface.
