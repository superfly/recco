Recco extracts information from application source trees to inform and simplify deployment on supported platforms. This work is initially designed to inform Nix-based deployments on [Fly.io](https://fly.io).

## Concepts

In Recco, knowledge about runtimes, frameworks and common services are encoded in YAML files called *scanners*. Here's a [scanner for Rails apps](scanners/rails.yml). A lightweight Golang library and CLI will take a source tree as an input, and spit out a YAML *deployment spec*. This spec should be picked up by deployment systems. The idea is for this tool to be runnable continuously - every time it's deployed.

Scanners can be of type: runtime, framework, service, e.g. ruby, rails, sidekiq.

Roughly, here's how it should work:

1. Run a source three through all scanners in scanners/*.yml
2. For matched runtimes, extract the version to inform package installers (like the Ruby version from `.ruby-version` or `Gemfile`)
3. For matched frameworks, extract versions to inform setting secrets (like `RAILS_MASTER_KEY` from `config/master.key`)
4. For matched services, set env vars and secrets like puma WEB_CONCURRENCY and 'puma -c config/puma.rb'

This is the main scope. One could imagine this being extended to run commands for deploy preparation, for example to create a [Docker-based release for Phoenix](https://hexdocs.pm/phoenix/Mix.Tasks.Phx.Gen.Release.html).
## Reasoning

Today, information about configuring deployments is locked up in Dockerfiles, buildpacks, platform-specific tooling and developers heads. We're no longer comfortable saying "Just deploy on Heroku", and need a way to break out of that mindset with confidence.

Buildpacks, at first glance, seem to offer a simple middleware-like API as a way out. But the reality is that buildpacks are fragmented across platforms and make no consistent guarantees about how they'll behave. It's a whole different experience running a standard deployment on Heroku and deploying with Heroku's buildpack-compatible Docker builder. The same runtimes and frameworks have multiple implementations of the same logic in bash, Golang, and so on.

While we can't go all no-code, we can look at the details of the problem. Most of the problems come down to unpredictability around the build process itself, and lack of visibilty into the logic used to make decisions about the deployment environment.

So if we compress some of this information into a low-calorie standard, we can make better decisions about improving build systems.

I anticipate pushback here from those who shun declarative configurtion over code. The idea here is not to reinvent buildpacks in YAML. It's to simplify - to extract what is actually *data* to a digestible and extensible format, and keep implementation details in code. The more data we can extract, the simpler the code will need to be.

Finally, this sort of information may only be useful to systems that offer fine-grained control over deployment configurations. Things like running multiple processes in-VMs, or installing specific versions of runtimes, can be hard to compose in inflexible systems like Dockerfiles or fragmented systems like buildpacks. We're get more into this topic as things progress.
