{
  pkgs,
  ...
}:

{
  #----------------------------------------------------------------------------
  # Basic Environment Setup
  #----------------------------------------------------------------------------
  env.GREET = "Claude Commit";

  #----------------------------------------------------------------------------
  # Languages and Packages
  #----------------------------------------------------------------------------
  # Go environment
  languages.go.enable = true;
  packages = [
    pkgs.golangci-lint
    pkgs.pkgsite
    pkgs.git
  ];

  #----------------------------------------------------------------------------
  # Scripts and Shell Hooks
  #----------------------------------------------------------------------------
  scripts = {
    hello.exec = ''
      echo hello from $GREET
    '';

    build.exec = ''
      VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0-dev")
      BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
      COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

      echo "Building Claude Commit $VERSION"
      go build -ldflags "-X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.commitSHA=$COMMIT_SHA" -o claude_commit .
    '';

    build-release.exec = ''
      VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0-dev")
      BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
      COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

      echo "Building Claude Commit $VERSION (release)"
      CGO_ENABLED=0 go build -ldflags "-w -s -X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.commitSHA=$COMMIT_SHA" -o claude_commit .
    '';

    version.exec = ''
      VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0-dev")
      BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
      COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

      echo "Version: $VERSION"
      echo "Build Date: $BUILD_DATE"
      echo "Commit: $COMMIT_SHA"
    '';

    test-code.exec = ''
      go test ./... -v
    '';

    test-coverage.exec = ''
      go test ./... -cover -coverprofile=coverage.out
      go tool cover -html=coverage.out -o coverage.html
      echo "Coverage report generated: coverage.html"
    '';

    test-race.exec = ''
      go test ./... -race
    '';

    bench.exec = ''
      go test ./... -bench=. -benchmem
    '';

    lint.exec = ''
      golangci-lint run
    '';

    fmt.exec = ''
      go fmt ./...
    '';

    vet.exec = ''
      go vet ./...
    '';

    clean.exec = ''
      rm -f claude_commit
      rm -f coverage.out coverage.html
      go clean -testcache
    '';

    ci.exec = ''
      echo "Running CI checks..."
      golangci-lint run
      go vet ./...
      go test ./... -race
      go test ./... -cover
      echo "All CI checks passed!"
    '';

    test-binary.exec = ''
      ./build
      echo "Testing built binary:"
      ./claude_commit --version
      ./claude_commit --help
    '';
  };

  enterShell = ''
    echo ""
    echo ""
    hello
    echo ""
    echo "Available commands:"
    echo "  build          - Build with version info"
    echo "  build-release  - Build optimized release binary"
    echo "  test-code      - Run tests"
    echo "  test-coverage  - Run tests with coverage"
    echo "  test-race      - Run tests with race detection"
    echo "  lint           - Run linter"
    echo "  fmt            - Format code"
    echo "  version        - Show version info"
    echo "  clean          - Clean build artifacts"
    echo "  ci             - Run all CI checks"
    echo ""
  '';

  enterTest = ''
    echo "Running tests"
    go --version
  '';

  #----------------------------------------------------------------------------
  # Tasks
  #----------------------------------------------------------------------------

  #----------------------------------------------------------------------------
  # Git Hooks
  #----------------------------------------------------------------------------
  git-hooks.hooks = {
    #----------------------------------------
    # Formatting Hooks - Run First
    #----------------------------------------
    beautysh.enable = true; # Format shell files
    gofmt.enable = true; # Format Go code
    nixfmt-rfc-style.enable = true; # Format Nix code

    #----------------------------------------
    # Linting Hooks - Run After Formatting
    #----------------------------------------
    shellcheck.enable = true; # Lint shell scripts
    golangci-lint.enable = true; # Lint Go code
    statix.enable = true; # Lint Nix code
    deadnix.enable = true; # Find unused Nix code

    #----------------------------------------
    # Security & Safety Hooks
    #----------------------------------------
    detect-private-keys.enable = true; # Prevent committing private keys
    check-added-large-files.enable = true; # Prevent committing large files
    check-case-conflicts.enable = true; # Check for case-insensitive conflicts
    check-merge-conflicts.enable = true; # Check for merge conflict markers
    check-executables-have-shebangs.enable = true; # Ensure executables have shebangs
    check-shebang-scripts-are-executable.enable = true; # Ensure scripts with shebangs are executable
  };
}
