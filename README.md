# Cobra and Viper Extended Example

An experimentation in Go's `cobra` and `viper` with some extras.

## What it does

- A `Config` struct with some hierarchy.
- An `init` command that takes the defaults from the `Config` struct and writes them into a `config.[toml, js, yaml, yml]` file.
- A `serve` command that exemplifies reading the configuration with the following order of importance:
  1. The command line flags.
  2. The environment variables.
  3. A `.env` file.
  4. The configuration file set on `viper`.
  5. Defaults.

This example is loosely based on a scenario where a web server uses a web `API` to get/set data.

## Deviations

### Global `initConfig`

On many examples, we see a call to `cobra.OnInitialize(initConfig)`. The `initConfig` function is called for all commands, since it's registered with the `cobra` environment, not per command.

This is OK when you don't have an `init` command. In this case, the `init` command does not need any general initialization of the config. It will do said initialization and write a new config file based on the defaults from `config.Config`.

In the case you don;t have a command like `init`, and all your commands depend on the `viper` config being setup, then please make the necessary changes to have that `initConfig`.

In my example, on the `serve` command, I do this initialization on it's `init` function.