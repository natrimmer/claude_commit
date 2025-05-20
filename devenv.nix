{
  pkgs,
  ...
}:

{
  #----------------------------------------------------------------------------
  # Basic Environment Setup
  #----------------------------------------------------------------------------
  env.GREET = "example project";

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
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  enterShell = ''
    hello
  '';

  enterTest = ''
    echo "Running tests"
    go --version
  '';

  #----------------------------------------------------------------------------
  # Tasks
  #----------------------------------------------------------------------------

  tasks."go:test" = {
    exec = ''
      go test ./...
    '';
  };

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
