# Viper style example

In this example, the viper.Viper config registry permeates all
the components of the app (in the `pkg` folder), which retrieve their settings
using the `Getxxxx` methods.

Although not mandatory, we maintain a central reference for config keys
in `pkg/config-keys`. This reference mainly serves documentation purposes.

## Trade-offs

**Benefits**:
* flexibility of a decentralized configuration: each module in the app
  knows where to get its settings.
* app modules may define their own defaults
* The central reference is advisory only and is not strictly required:
  modules may evolve and possibly be refactored later to expose which
  keys are in use.

**Shortcomings**:
* settings are not strongly typed or globally validated
* when many modules evolve separately, it may become hard to understand
  the full configuration or the impact of changing a given setting in
  the config.

## Running the example
