# ffprogress

ffprogress provides elapsed time, ETA and progress percentage based on your ffmpeg call.  e.g.

```shell
Elapsed 00:16:59, ETA 02:26:24, Progress 19.79%
```

## Installation

Using `go get`:

```shell
go get -u github.com/ashmckenzie/ffprogress
```

Or, download a [release](https://github.com/ashmckenzie/ffprogress/releases).

## Usage

```shell
ffmpeg -y -i in.mkv -o out.mkv 2>&1 | ${GOPATH}/bin/ffprogress
```

NOTE:

* You must redirect STDERR to STDOUT (`2>&1`), otherwise `ffprogress` can not parse the necessary output
* As STDERR to STDOUT, you will need to add `-y` to overwrite (if exists) as you won't see the `Overwrite ? [y/N]` prompt
* If you're adjusting the output verbosity (`-v`), you will need at minimum `32` for `ffprogress` to work (`32` is the default)

## Contributing

1. Fork it ( https://github.com/ashmckenzie/ffprogress/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create a new Pull Request
