# NAME

hei - A simple consistent command wrapper for nix

# SYNOPSIS

hei

```
[--debug|-D]
[--dry-run|-d]
[--flake|-f]=[value]
[--help|-h]
[--version|-v]
```

# DESCRIPTION

A simple consistent command wrapper for nix

**Usage**:

```
hei [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--debug, -D**: Enable debug logging

**--dry-run, -d**: Perform a dry run without making any changes

**--flake, -f**="": Path to flake. Will default to auto-detect

**--help, -h**: show help

**--version, -v**: print the version


# COMMANDS

## build

Run nix flake check on your flake

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## check

Run checks on the given flake paths or the default ones if none are provided

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## gc

Garbage collect & optimize nix store

**--all, -a**: Collect all garbage

**--help, -h**: show help

**--system, -s**: Collect system garbage

### help, h

Shows a list of commands or help for one command

## gen

Manage nix generations

**--help, -h**: show help

### list

List nix generations

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### delete

Build the given flake paths or the default ones if none are provided

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### diff

Build the given flake paths or the default ones if none are provided

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### switch

Switch generation

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## p

Shortcut for nix profile commands

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## rebuild

Rebuild your nix configuration

**--fast, -f**: Build in fast mode

**--help, -h**: show help

**--offline, -o**: Build in offline mode

**--rollback, -r**: Build in fast mode

**--update, -u**: Pull on nix flake before rebuilding

### help, h

Shows a list of commands or help for one command

## repl

open a repl in your nix config

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## rollback

Roll back to previous generation of nixos. See gen list for the current generations.

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## search

Search nixpkgs for packages

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## show

Run nix flake show on the given flake paths or the default ones if none are provided

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## ssh

Run a hei command on a remote NixOS system

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## swap

Recursively swap nix-store symlinks with copies (or back)

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## test

Run a nix check

**--help, -h**: show help

**--interactive, -i**: Run the test with the interctive test driver

### help, h

Shows a list of commands or help for one command

## upgrade

deprecated, use rebuild -u instead

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## update

Update the given flake paths or the default ones if none are provided

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command

## completion

Generate shell completion scripts

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command
