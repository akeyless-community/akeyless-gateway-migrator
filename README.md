# Akeyless Gateway Migrator

## Installation

You can install `agmigrator` using Homebrew with the following command:

```bash
brew install akeyless-community/agmigrator/agmigrator
```

## Usage

You can use the following command line flags:

- `--akeyless-source-token`: Akeyless source token
- `--akeyless-destination-token`: Akeyless destination token
- `--source-gateway-config-url`: Source gateway Config URL
- `--destination-gateway-config-url`: Destination gateway config URL
- `--filter-config-file-path`: Filter config file path
- `--debug`: Enable debug mode

Alternatively, you can set the following environment variables:

- `AKEYLESS_SOURCE_TOKEN`: Akeyless source token
- `AKEYLESS_DESTINATION_TOKEN`: Akeyless destination token
- `SOURCE_GATEWAY_URL`: Source gateway Config URL
- `DESTINATION_GATEWAY_URL`: Destination gateway config URL
- `FILTER_CONFIG_FILE_PATH`: Filter config file path
- `DEBUG`: Enable debug mode

## Examples

### Running the CLI with flags

```bash
agmigrator kubernetes \
  --akeyless-destination-token "t-fe9e3c72d2d50b2e19e7020c322239e3" \
  --akeyless-source-token "t-fe9e3c72d2d50b2e19e7020c322239e3" \
  --source-gateway-config-url "https://gw-config.old.akeyless.fans" \
  --destination-gateway-config-url "https://gw-config.new.akeyless.fans"
```

### Running the CLI with environment variables for mac/linux

```bash
export AKEYLESS_SOURCE_TOKEN="t-fe9e3c72d2d50b2e19e7020c322239e3"
export AKEYLESS_DESTINATION_TOKEN="t-fe9e3c72d2d50b2e19e7020c322239e3"
export SOURCE_GATEWAY_CONFIG_URL="https://gw-config.old.akeyless.fans"
export DESTINATION_GATEWAY_CONFIG_URL="https://gw-config.new.akeyless.fans"
agmigrator kubernetes
```
