# cronlocker

cronlocker is a commandline tool to allow running cronjobs on multiple hosts while ensuring that it only runs once at a time.

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
  -minlocktime int
        Configures the minimum time in milliseconds a lock is held (default 5000)
```

## Packaging

Execute `make package` to package the application as a debian package.

Note: [fpm](https://github.com/jordansissel/fpm) is required

## Bugs and Contribution

For bugs and feature requests open an issue on Github. For code contributions fork the repo, make your changes and create a pull request.

## License

[LICENSE] (MIT)
