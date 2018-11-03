rfoutlet
========

Outlet control for Raspberry PI 2/3.

TODO

Installation
------------

### Using `go get`

`node` executable is required to build the frontend. rfoutlet was tested with
`node v10.11.0+` and `go v1.11+` but may also work with older versions.

```sh
go get github.com/martinohmann/rfoutlet
cd $GOPATH/github.com/martinohmann/rfoutlet
make build
make install
```

### Using docker

Build the image for armv7:

```sh
make image-rfoutlet-armv7
```

This will create an image called `mohmann/rfoutlet:armv7`.

Raspberry PI Setup
------------------

TODO

Outlets
-------

TODO

Usage
-----

`sudo` is required in order to access `/dev/mem` and `/dev/gpiomem`. If not
provided, rfoutlet looks in `/etc/rfoutlet/config.yml` for its config (you can
change that by providing the `-config` flag). Check
[dist/config.yml](dist/config.yml) for an example config file with all
available config values.

By default rfoutlet uses gpio pin 17 (physical 11 / wiringPi 0) for
transmission of the rf codes. A different pin can be use by providing the
`-gpio-pin` flag. Check out the [Raspberry Pi pinouts](https://pinout.xyz/) for
reference.

rfoutlet listens on `0.0.0.0:3333` by default but you can change the listen
address by providing the `-listen-address` flag.

### Installed locally

Start the server and browse to `<raspberry-ip-address>:3333`:

```sh
sudo rfoutlet -listen-address 0.0.0.0:3333 -config /etc/rfoutlet/config.yml
```

### Via docker

Start the container and browse to `<raspberry-ip-address>:3333`:

```sh
docker run --rm \
    --privileged \
    -p 3333:3333 \
    -v $(pwd)/dist/config.yml:/etc/rfoutlet/config.yml \
    -v /dev/mem:/dev/mem \
    -v /dev/gpiomem:/dev/gpiomem \
    mohmann/rfoutlet:armv7
```

The container has to run in privileged mode in order to be able to access
`/dev/mem` and `/dev/gpiomem`.

Code sniffing
-------------

TODO

```sh
sudo rfsniff -help
```

Code transmission
-----------------

This repo provides a tool called `rftransmit` to send rf codes. You can use
this for testing or wrap it with your own outlet control tool. Check the help
for available options:

```sh
sudo rftransmit -help
```

Or using docker:

```sh
make image-rftransmit-armv7
docker run --rm \
    --privileged \
    -v /dev/mem:/dev/mem \
    -v /dev/gpiomem:/dev/gpiomem \
    mohmann/rftransmit:armv7 \
    -help
```

Development / Testing
---------------------

rfoutlet is meant to run on a Raspberry PI 2/3 to work properly. However, for
development purposes you can also run it on your local machine. In this can the
transmission of the rf codes is simulated.

Run `make` without arguments to see other available commands.

Todo
----

- [x] port
  [codesend](https://github.com/ninjablocks/433Utils/blob/master/RPi_utils/codesend.cpp)
  to golang (see [`cmd/rftransmit`](cmd/rftransmit))
- [x] port
  [RFSniffer](https://github.com/ninjablocks/433Utils/blob/master/RPi_utils/RFSniffer.cpp)
  to golang  (see [`cmd/rfsniff`](cmd/rfsniff))
- [x] make transmitter/receiver code available as library below `pkg/`
- [ ] persist outlet state across server restarts
- [ ] use receiver to detect outlet state changes (e.g. via remote control)?
- [ ] time switch: switch outlets on/off using user defined rules

License
-------

rfoutlet is released under the MIT License. See the bundled LICENSE file for details.

Resources
---------

- [Raspberry Pi pinouts](https://pinout.xyz/)
- [Wireless Power Outlets](https://timleland.com/wireless-power-outlets/)
- [ninjablocks 433Utils](https://github.com/ninjablocks/433Utils)
- [WiringPi](https://projects.drogon.net/raspberry-pi/wiringpi/download-and-install/)
