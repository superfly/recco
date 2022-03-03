Recco extracts information from application source trees to inform and simplify deployment on supported platforms. This work is initially designed to inform Nix-based deployments on [Fly.io](https://fly.io).

## Concepts

Recco has one simple goal: given an application source tree, it will generate a *YAML deployment specification* to be used by build systems and deployment platforms. Recco should have few dependencies and be runnable at deploy time to catch changes made in the source tree.

Knowledge about runtimes, frameworks and common application services is encoded in YAML files called *scanners*. Here's a sample [scanner for Rails apps](scanners/rails.yml).  

Scanners can be of different types with different rules: runtime (ruby), framework (rails) or service (sidekiq). We know things about each of these, such as whether YJIT is available (Ruby 3.1.1 and above), whether to set `RAILS_MASTER_KEY` (if credentials are present) or whether to run a second process for the app (a Sidekiq worker).

## Implementation dieas 

Roughly, here's how it could work:

1. Run a source three through all scanners in `scanners/*.yml`
2. For matched runtimes, extract the version to inform package installers (like the Ruby version from `.ruby-version` or `Gemfile`)
3. For matched frameworks, extract versions to inform setting secrets (like `RAILS_MASTER_KEY` from `config/master.key`)
4. For matched services, set env vars and secrets like puma WEB_CONCURRENCY and 'puma -c config/puma.rb'
5. Export a YAML spec intended to be versioned in source trees and picked up by build systems

One could imagine this being extended to run commands for deploy preparation, for example to create a [Docker-based release for Phoenix](https://hexdocs.pm/phoenix/Mix.Tasks.Phx.Gen.Release.html).

## Why do we need this?

Today, information about configuring deployments is locked up in Dockerfiles, buildpacks, platform-specific tooling and developers heads. Developers are no longer comfortable saying "Just deploy on Heroku". We need a way to break out of that mindset with confidence and without a degree in Dockerfiles. You've probably run into a lot of problems unrelated to your app code that you've quietly erased from your memory. The goal is to reduce friction at deploy time while also educating developers about what's possible.

## What about Buildpacks? 

Buildpacks, at first glance, seem to offer a simple middleware-like API as a way out. But the reality is that buildpacks are fragmented across platforms and make no consistent guarantees about how they'll behave. It's a whole different experience running a standard deployment on Heroku and deploying with Heroku's buildpack-compatible Docker builder. The same runtimes and frameworks have multiple implementations of the same logic in bash, Golang, and so on.

## What, are some no-code zealot? Isn't this just reinventing the wheel?

Nah - let's look at the details of the problems developers have at deploy time. I think they come down to:

1. Lack of knowledge on what knobs to turn for deployment, and how
2. Lack of clarity around the actual build process: do we have caching, persistent storage, etc?
3. Lack of visibilty into the logic used to make decisions about the deployment environment (disparate dockerfiles, buildpacks)
4. Lack of confidence making changes to deployment configuration (packages going missing)
5. Lack of flexibility when trying to compose software (mismatched versions of Ruby, Node, etc)

A lot of this lack is related to details of the build system. But some is related to the possibilities in each domain being obscured.

## Isn't this just declarative versus imperative config?

I anticipate pushback here from those who shun declarative configuration over code. Here we're not going this far. We're simply extracting what is actually *data* to a digestible and extensible format. The more data we can extract, the simpler the code that consumes this information will need to be.

## How does this information get into a build system?

I think this information may only be useful to systems that offer fine-grained control over deployment configurations. Things like running multiple processes in-VMs, or installing specific versions of runtimes, can be hard to compose in inflexible systems like Dockerfiles or fragmented systems like buildpacks. At Fly, we're looking into Nix to help us here. More on this topic soon.
