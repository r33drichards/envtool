# EnvTool

EnvTool is a command-line utility for managing environment variables in shell environments. It allows you to automatically load variables from `.env` files into your shell environment when you change directories or start a new shell session.

## Features

- Automatically load environment variables from `.env` files
- Track which variables are managed by the tool
- Intelligently unset variables that are no longer defined
- Support for both bash and zsh shells
- Support for user-specific or system-wide configuration

## Installation

### From Source

```bash
git clone https://github.com/username/envtool.git
cd envtool
go install
```

### Using Homebrew (coming soon)

```bash
brew install username/tap/envtool
```

## Quick Start

1. Initialize shell configuration:

```bash
# For system-wide configuration (requires sudo)
sudo envtool init

# For user-specific configuration
envtool init --user
```

2. Create a `.env` file in your project directory:

```bash
echo "DATABASE_URL=postgres://localhost:5432/mydb" > .env
echo "API_KEY=my_secret_key" >> .env
```

3. Navigate to your project directory and your environment variables will be automatically loaded!

## Usage

### Initialize Shell Configuration

The `init` command sets up your shell configuration files to automatically run EnvTool when your shell prompt is displayed:

```bash
# System-wide configuration (requires sudo)
sudo envtool init

# User-specific configuration
envtool init --user

# Specify custom configuration file paths
envtool init --bashrc ~/.bashrc --zshrc ~/.zshrc

# Update only bash or only zsh (defaults to both if neither specified)
# Only bash:
envtool init --bash
# Only zsh:
envtool init --zsh

# Positional args when selecting a single shell:
#   envtool init --bash <rc_path> [env_file]
#   envtool init --zsh <rc_path> [env_file]
# Examples:
#   Write bash hook to /etc/bashrc using /path/to/.env:
sudo envtool init --bash /etc/bashrc /etc/incus-env.env

#   Write zsh hook to ~/.zshrc using default .env from config:
envtool init --zsh ~/.zshrc
```

### Manually Load Environment Variables

You can also manually load environment variables from a specific `.env` file:

```bash
# Load variables from the default .env file in the current directory
eval "$(envtool env)"

# Load variables from a specific .env file
eval "$(envtool env --env-file /path/to/.env)"

# Specify shell type (bash is default)
eval "$(envtool env zsh --env-file /path/to/.env)"
```

## How It Works

EnvTool works by adding a hook to your shell prompt that executes the `envtool env` command every time your prompt is displayed. The command reads the `.env` file in your current directory, exports the variables, and keeps track of which variables it has set.

When you move to a different directory or the `.env` file changes, EnvTool will automatically update your environment, exporting new variables and unsetting variables that are no longer defined.

## Configuration

EnvTool can be configured through command-line flags or a configuration file. The default configuration file location is `~/.envtool.yaml`.

Example configuration file:

```yaml
env-file: .env
init:
  user: true
  bashrc: ~/.bashrc
  zshrc: ~/.zshrc
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.