# cronlocker

![GitHub](https://img.shields.io/github/license/viafintech/cronlocker) ![Build Status](https://github.com/viafintech/cronlocker/actions/workflows/test.yml/badge.svg)  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/viafintech/cronlocker/master) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/viafintech/cronlocker)

cronlocker is a commandline tool to allow running cronjobs on multiple hosts while ensuring that it only runs once at a time.
cronlocker utilizes the [consul](https://www.consul.io/) [lock](https://www.consul.io/docs/commands/lock.html) feature to ensure that.

## Usage

cronlocker can be easily executed by passing a key to lock as well as a command to be executed if the lock could be obtained successfully.

```
cronlocker -key=<key/to/lock> <command>
```

Use `cronlocker --help` to get the following output:

```
Usage of ./cronlocker:
  -endpoint string
      endpoint (default "http://localhost:8500")
  -key string
      key to monitor, e.g. cronjobs/any_service/cron_name (default "none")
  -lockwaittime int
      Configures the wait time for a lock in milliseconds (default 500)
  -maxexecutiontime int
      Configures the maximum time in milliseconds the execution of the given command can take
  -minlocktime int
      Configures the minimum time in milliseconds a lock is held (default 5000)
```

## Packaging

Execute `make package` to package the application as a debian package.

Note: [fpm](https://github.com/jordansissel/fpm) is required

## Bugs and Contribution

For bugs and feature requests open an issue on Github. For code contributions fork the repo, make your changes and create a pull request.

## License

[LICENSE](LICENSE)
