# Akeyless Gateway Migrator

Akeyless Gateway Migrator is a CLI tool that helps you migrate configuration from one Akeyless Gateway cluster to another. 

## Installation

You can install `agmigrator` using Homebrew with the following command:

```bash
brew install akeyless-community/agmigrator/agmigrator
```

## Usage

You can use the following command line flags:

- `--akeyless-source-token` (required): The token for the source Akeyless Gateway.
- `--akeyless-destination-token` (required): The token for the destination Akeyless Gateway.
- `--source-gateway-config-url` (required): The configuration URL for the source Akeyless Gateway.
- `--destination-gateway-config-url` (required): The configuration URL for the destination Akeyless Gateway.
- `--filter-config-file-path` (optional): The path to a file containing the name(s) of the k8s auth configurations to be migrated. [Exmaple found here](text-file-filter-example.txt) If not provided, all configurations will be migrated.
- `--debug` (optional): Enable debug mode to output all the argument values.

Alternatively, you can set the following environment variables:

- `AKEYLESS_SOURCE_TOKEN` (required): Equivalent to `--akeyless-source-token`.
- `AKEYLESS_DESTINATION_TOKEN` (required): Equivalent to `--akeyless-destination-token`.
- `SOURCE_GATEWAY_URL` (required): Equivalent to `--source-gateway-config-url`.
- `DESTINATION_GATEWAY_URL` (required): Equivalent to `--destination-gateway-config-url`.
- `FILTER_CONFIG_FILE_PATH` (optional): Equivalent to `--filter-config-file-path`.
- `DEBUG` (optional): Equivalent to `--debug`.

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
