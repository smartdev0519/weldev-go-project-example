# Go Project Example

This is an example for Go project.

The motivation behind this project is to learn and widen my limited knowledge about programming, project design, and concepts implementation. In this project, I will try to implement business logic/flow into Go program for various use-cases.

Some of them might not follow existing specs/standards, feel free to open issues, and please let me know.

## Designing Project For Industrial Programming

What is industrial programming? [Peter Bourgon](https://peter.bourgon.org/go-for-industrial-programming/) explain the terms, as:

- In a startup or corporate environment.
- Within a team where engineers come and go.
- On code that outlives any single engineer.
- Serving highly mutable business requirements.

## The Project

The project theme is `Property`. I will try to build a Property application, where people able to search and book the property.

### Use-cases

1. Users were able to register and log in.
2. Users were able to register their Properties.
    - Register the property detail
    - Upload the property image
3. Users were able to book a Property.
4. Users were able to receive notifications and have a notification inbox.

### Project Stack

This project is using:

1. PostgreSQL for the main database
2. Redis for k/v, user session management, caching. 

## Getting Started

This is guide to get started with this project and installing dependencies required for running this project locally

### Requirements

1. Make
2. Docker
3. Soda CLI from [gobuffalo](https://gobuffalo.io/en/docs/db/toolbox). You can install this by using `make install-deps`

### Create Database And Migrate

To create and migrate the database, we will use `soda CLI` created by `gobuffalo`. The command is wrapper by this [script](/database/setup.sh).

Use this command to fully create and migrate the databse schema:

`make dbup`

### Flags

The following flags is avaiable to help the project configuration and debug parameters.

- `--config_file` define where the configuration file is located
- `--env_file` define where environment variable file is located, this is a helper file to set environment variable.
- `--log` define the log configuration for the project. The flag contains comma separated value:
    - `--log=file=./to_your_file.log` the location of log file.
    - `--log=level=debug|info|warning|error|fatal` the level of log.
    - `--log=color=1` to set if color to console/terminal is enabled.
- `--debug` define the debug configuration for the project. The flag contains comma separated value:
    - `--debug=server=1` to turn on the debug server.
    - `--debug=testconfig=1` to test the configuration of the project.

### Configuration

For applicatoin configuration, config-file is used and located on the root project directory. The [configuration-file](./project.config.toml) is written in `toml` grammar.

The configuration value in the `configuration-file` is embed as environment-variable value. The value will be replaced by the environment-variable value in runtime or when program started. To help the process, the project use the help of [environment-variable](./project.env.toml) file. The `environment-variable` files is choosen because it is simpler as we can have multiple files and we can hide/ignore the file for specific use-case, for example, secret value of `client_id` and `client_secret` of some cloud vendor.

The mixed of `configuration-variable` and `environment-variable` is used to help people in project to see what configuration structure is exists within the project, and able to dynamically changed depends on the environment variables value.

### Environment State

The project have no environment state. Different flags and configuration value is used in different environment.

Environment state like `dev`, `staging`, and `production` is usually used to check in what environment the program/application is running. From experience, this considered harmful for the program itself, as developer tempted to abuse the state for many things. Developer tempted to abuse the state because the function is available, and sometimes it is the easiest way to accomplish some goals. By using the state, people in the project are cutting edges and create conditional expression for various use-cases. This leads to broken mental model, bugs, and edge-cases to the product which make life harder for the maintainers.

For example, in code:

```go
if env.IsDevelopment() {
    // do something only in dev
}

if env.IsStaging() {
    // do something only in staging
}

if env.IsProduction() {
    // do something only in production
}
```

This environment state, sometimes also used for configuration directive. When the configuration directive is gated by the environment state, another problem occur. Because, configuration for each environment might have different variables and value, and different configuration file can mean different things.

For example, in configuration with spec `project.{environment_name}.config.toml`:

- project.dev.config.toml
- project.staging.config.toml
- project.production.config.toml

Or, imagine if you have many different configurations(with various reasons/decision) using this kind of directive:

- project.dev.config1.toml
- project.staging.config1.toml
- project.production.config1.toml
- project.dev.config2.toml
- project.staging.config2.toml
- project.production.config2.toml

Things got very messy indeed.

Multiple configuration with environment state directive, usually used to address different configurations in each environments. For example, when a database is pointing to one instance in `dev` but not in `staging`, which completely different. Or, when doing doing some migration we want to get rid of some configuration variables in some environment. This all are valid use-cases, and the given solution by using the environment state for configuration directive works. Usually, until the configuration is become too long and different for each environments, then turning into problems for the maintainers.

As sometimes we need to run with some special configuration in non-production or in production environment, this might be able to achieved by using the combination of flags and configuration-file. Variables from flags and configuration is more clear and straightforward than `IsEnvrionment`, and can be used to checked the design choices, do we have too many hacks? Why? For whatever reason, the flags/configuration variables between environments should stay the same, to maintain consistency.

But, in the end, it depends on each project policies and governance.

## Project Structure

### Cmd

All Go main programs is located in `go_project_example/cmd/*` folder.

### Internal

To be added

## Code Structure

As stated at the top of this document, the design contains several layers and components and very much similar to onion ring or clean architecture attempt.

But, let's talk about the components first.

### Components

1. Server
2. Usecase
3. Repository
4. Entity

#### Server

Is where all the `http` handler exists. This layer is responsible to hold all the `http` handlers and request validation.

**Main Server**

Main server for serving main application. All business use-case handler exists within this server.

**Admin Server**

Admin server is a server for administrational purpose. On by default and should not open to public.

Use-case for admin server:

- `/metrics` endpoint
- check current configuration value

**Debug Server**

Debug server is a server for experimental purpose, and should be enabled with a spesific flag. This server should not be triggered in production environment.

Use-case for debug server:

- Login bypass
- Serve fileserver for local object storage

#### Usecase

To be added

#### Repository

To be added

#### Entity

To be added

#### Layers

To be added

## Acknowledgement & References

To be added