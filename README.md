# nuon-ext-starship

Starship prompt module that shows your current Nuon context: environment, org name, and install name.

## Install

```bash
nuon ext install nuonco/nuon-ext-starship
```

## Setup

Generate the Starship config snippet:

```bash
nuon starship init
```

Add the output to your `~/.config/starship.toml`.

## Usage

The prompt module activates automatically when you're in a directory containing `runner.toml` (a Nuon app directory).

Example prompt output:

```
☁️ stage 🏢 acme-corp 📦 us-west-prod
```

### How it works

1. Reads `~/.nuon` for the current `org_id`, `install_id`, and `api_url`
2. Parses the API URL to determine the environment (`prod`, `stage`, `dev`, `local`)
3. Resolves org/install IDs to human-readable names via the Nuon API
4. Caches resolved names at `~/.cache/nuon/starship.json` — API calls only happen when IDs change

### Commands

| Command | Description |
|---------|-------------|
| `nuon starship prompt` | Output the formatted prompt string (called by Starship) |
| `nuon starship init` | Print the `starship.toml` snippet to add |
| `nuon starship refresh` | Force-refresh the name cache |

### Graceful degradation

- **No org selected**: prompt is empty (module hidden by Starship)
- **No API token / expired token**: shows truncated IDs instead of names
- **API unreachable**: uses cached names if available, otherwise truncated IDs
- **No install selected**: only shows env + org
