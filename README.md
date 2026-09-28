# Peon

Peon is a simple Git Worktree management tool designed for humans

## Installation

Install the package through one of the methods listed below and put this into your `.bashrc`:

```bash
eval "$(peon hook bash)"
```

Similarly `zsh` and `fish` hooks are also available.

### Nix

Add my nix repo to your config as shown here: [atomicptr/nix](https://github.com/atomicptr/nix#usage)

## Usage

The available commands are limited to the following actions

### Switch - Switch to a worktree, create a new one if needed

You can switch to an existing workspace using:

```bash
$ peon switch WORKTREE_NAME

# or short
$ peon sw WORKTREE_NAME
```

If you want to create a new worktree, this is also part of the `switch` command:

```bash
$ peon switch --create NEW_WORKTREE_NAME

# or short
$ peon sw -c NEW_WORKTREE_NAME
```

### List - List all available worktrees

This one is very self explanatory:

```bash
$ peon list

# or short
$ peon ls
```

This will list your worktrees

### Remove - Remove worktree

With

```bash
$ peon remove

# or short
$ peon rm
```

you can remove your current worktree.

**Note**: This will not work if you have unstaged changes available unless you add the `--force` flag

In addition you can also specify a worktree name to delete:

```bash
$ peon rm WORKTREE_NAME
```

## Motivation

While git worktrees are cool, I found them to be kind of a pain to use in general. While lots of other tools already
exist they all seem to be very focused on AI driven workflows, which I'm not interested in.

## License

GPLv3
