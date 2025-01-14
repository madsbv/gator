# Gator
RSS feed aggregator with Postgres backend.

## Installation
To use gator you need a running Postgres server for it to connect to and use as backend (this can be a local instance of postgres). Currently you also need the Go toolchain for installation:

``` sh
go install https://github.com/madsbv/gator
```

Gator uses the configuration file `.gator.json` in your systems default configuration directory--on Linux this is your XDG config home, on MacOS this will be in `$HOME/Library/Application Support`.
The configuration file has the following structure:

``` json

{
    "db_url": "postgres://dev@localhost:5432/gator?sslmode=disable",
    "current_user_name": "username"
}
```

`db_url` must point to the postgres database you intend to use. For a local instance, you probably don't care about setting up SSL, so you can disable it with `sslmode=disable`.
You can leave `current_user_name` empty or omit it entirely; gator will manage that itself.

## Usage
Gator is a CLI application with a number of commands, described below. Each command is run as `gator command args`, where each command takes a varying number of arguments.

Gator is intended to be a long running service that fetches new posts in the background. The background service is invoked with the "agg" command (see below), and interacting with the stored posts happens on the CLI with other commands.

Gator supports tracking RSS feeds for multiple users from the same running instance. You must first register a user with the "register" command, which also "logs in" that user, i.e. sets `current_user_name` in the configuration file. You can then add feeds to track, follow feeds that other users have already added, view feeds and their posts, and so on.

- `gator login username`: Log in as user `username`.
- `gator register username`: Register `username` as a new user.
- `gator reset`: Deletes all saved user and feed data. **DANGEROUS**
- `gator agg duration`: Runs the gator service to aggregate posts from all feeds tracked by gator, for every user. Gator retrieves posts from one RSS feed at a time, scraping the most outdated feed every `duration`. Here `duration` can be anything that can be parsed by Go's `time.ParseDuration`. Examples include `5s, 1m, 3h`. For example, running `gator agg 1m` means that one feed is fetched every minute. If you are following 5 feeds, it takes 5 minutes to retrieve posts from all 5 feeds.
- `gator addfeed name url`: Add the RSS feed at location `url` to the collection of tracked feeds, under the name `name`.
- `gator feeds`: List all feeds in the database.
- `gator follow url`: As the current user (as set by `gator login`), follow the feed at `url` if it's already tracked by gator for another user.
- `gator following`:  List all feeds that the current user follows.
- `gator unfollow url`: Remove the feed at `url` from the list of feeds the current user follows. 
- `gator browse num`: The argument `num` is optional. Lists the `num` (or 2 if omitted) most recent posts from the feeds that the current user follows.
