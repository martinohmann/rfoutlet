rfoutlet
========

Outlet control for Raspberry PI 2/3.

TODO

Installation
------------

`node` executable is required to build the frontend. rfoutlet was tested with
`node v10.11.0+` and `go v1.11+` but may also work with older versions.

```sh
go get github.com/martinohmann/rfoutlet
cd $GOPATH/github.com/martinohmann/rfoutlet
make build
make install
```

Raspberry PI Setup
------------------

TODO

Usage example
-------------

`sudo` is required in order to access `/dev/gpiomem`. If not provided,
rfoutlet looks in `/etc/rfoutlet/config.yml` for its config. Check
[dist/config.yml](dist/config.yml) for an example config file with all
available config values.

Start the server:

```sh
sudo rfoutlet --listen-address 0.0.0.0:3333 --config /etc/rfoutlet/config.yml
```

Browse to `locahost:3333`.

Development / Testing
---------------------

rfoutlet is meant to run on a Raspberry PI 2/3 to work properly. However, for
development purposes you can also run it on your local machine. In this can the
transmission of the rf codes is simulated.

Run `make` without arguments to see other available commands.

License
-------

rfoutlet is released under the MIT License. See the bundled LICENSE file for details.
