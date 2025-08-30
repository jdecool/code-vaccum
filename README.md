Code Vaccum
===========

Scrape code from coding platforms (GitHub/GitLab). Download repositories from organizations or individual users.

## Usage

```bash
  $ ./github-vacuum --provider [github|gitlab] --output [filesystem|nil|repo] [--org org-name] [--username username]
```

### Examples

```bash
# Download all repositories from a GitHub organization
./github-vacuum --provider github --output filesystem --org myorg

# Download all repositories from a specific user
./github-vacuum --provider github --output filesystem --username someuser

# Download from multiple users and organizations
./github-vacuum --provider github --output filesystem --org myorg --username user1 --username user2

# Use with GitLab and custom endpoint
./github-vacuum --provider gitlab --provider-endpoint https://gitlab.example.com --provider-access-token TOKEN --username someuser --output filesystem

# Use specific SSH key for authentication
./github-vacuum --provider github --output filesystem --username someuser --ssh-key ~/.ssh/id_rsa
```

## Options

### Providers options

* `--provider`: provider to use (could be `github` or `gitlab`)
* `--provider-endpoint`: if use a self-hosted instance, you can specify the endpoint to use
* `--provider-access-token`: access token for authenticated queries

### Filtering options

* `--org`: filter by organization name (can be used multiple times)
* `--username`: filter by username to download all repositories of a user (can be used multiple times)

### Output options

* `--output`: output format to use (could be `filesystem`, `nil`, or `repo`)
  - `filesystem`: Clone repositories to local filesystem
  - `nil`: No-op output for dry-run/testing
  - `repo`: Repository-based output format
* `--output-folder`: available for `filesystem`. Output folder where projects will be cloned (default: current path)
* `--ssh-key`: path to SSH private key file for Git authentication (e.g., `~/.ssh/id_rsa`)

### General options

* `--debug`: Enable debug logging to see detailed processing information
* `--quiet`: Enable quiet mode (only show warnings and errors)
