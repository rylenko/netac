# NetAC

A simple program that uses a multicast group to determine the number of running copies of an application on the network. IPv4 and IPv6 are supported.

# Install

Build:

```
$ make
```

Install. Most likely you will need `sudo`, `doas` or somethink like that before the command:

```
$ make install
```

# Usage

To see possible CLI flags:

```
$ netac -help
```

Example of launching on IPv4 multicast group 224.0.0.1 interface wlo1:

```
$ netac -iface wlo1 -ip 224.0.0.1
```

Example of launching on IPv6 multicast group 224.0.0.1 interface wlo1:

```
$ netac -iface wlo1 -ip ff01::1
```
