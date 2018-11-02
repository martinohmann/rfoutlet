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

TODO

Raspberry PI Setup
------------------

TODO

Outlets
-------

TODO

Usage example
-------------

`sudo` is required in order to access `/dev/gpiomem`. If not provided,
rfoutlet looks in `/etc/rfoutlet/config.yml` for its config. Check
[dist/config.yml](dist/config.yml) for an example config file with all
available config values.

Start the server:

```sh
sudo rfoutlet -listen-address 0.0.0.0:3333 -config /etc/rfoutlet/config.yml
```

Browse to `locahost:3333`.

By default rfoutlet uses gpio pin 17 (physical 11 / wiringPi 0) for
transmission of the rf codes. A different pin can be use by providing the
`-gpio-pin` flag. Check out the [Raspberry Pi pinouts](https://pinout.xyz/) for
reference.

Code transmission
-----------------

This repo provides a tool called `rftransmit` to send rf codes. You can use
this for testing or wrap it with your own outlet control tool.

```sh
rftransmit -help
```

Development / Testing
---------------------

rfoutlet is meant to run on a Raspberry PI 2/3 to work properly. However, for
development purposes you can also run it on your local machine. In this can the
transmission of the rf codes is simulated.

Run `make` without arguments to see other available commands.

Todo
----

- [x] port [codesend](https://github.com/ninjablocks/433Utils/blob/master/RPi_utils/codesend.cpp) to golang
- [ ] port [RFSniffer](https://github.com/ninjablocks/433Utils/blob/master/RPi_utils/RFSniffer.cpp) to golang
- [ ] persist outlet state across server restarts
- [ ] time switch: switch outlets on/off using user defined rules
- [ ] make transmission code available as library

License
-------

rfoutlet is released under the MIT License. See the bundled LICENSE file for details.

Resources
---------

- [Raspberry Pi pinouts](https://pinout.xyz/)
- [Wireless Power Outlets](https://timleland.com/wireless-power-outlets/)
- [ninjablocks 433Utils](https://github.com/ninjablocks/433Utils)
- [WiringPi](https://projects.drogon.net/raspberry-pi/wiringpi/download-and-install/)
