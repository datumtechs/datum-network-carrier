## Metis Carrier

Official Golang implementation of the Metis-Carrier protocol.

## Building the source

Building `carrier` requires both a Go (version 1.16 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make carrier
```

or, to build the full suite of utilities:

```shell
make all
```

## Executables

The Metis-Carrier project comes with several wrappers/executables found in the `cmd` directory.

|    Command    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  **`carrier`**   | Our main Metis-Carrier CLI client. It is the entry point into the Metis-Carrier network, capable of running as a full node (default). `carrier --help` |
|   `keytool`    | Stand-alone signing tool, which can be used to generate a key pair.  |
|  `enr-calculator`   | ENR generator tool, which can be used to generate the address record information of the corresponding node。 |                                                                                                                                                                                                                                                                 |

## Running `carrier`

Going through all the possible command line flags is out of scope here (please use `carrier -help` to find out),
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `carrier` instance.

### Full node

The start of the node needs to know the address information of the data center service in advance, and the program will take the initiative to conduct 
a link communication test with the data service after it is started. If the connection fails, the program fails to start, and vice versa.

```shell
$ carrier --datadir ./data --grpc-datacenter-host 192.168.10.146 --grpc-datacenter-port 9098
```

* `--grpc-datacenter-host` Data Service listening interface (default: `192.168.112.32`)
* `--grpc-datacenter-port` Data Service listening port (default: `9099`)

*Note: For more flags settings, please use the `carrier --help` command to view*


### Docker quick start

One of the quickest ways to get Carrier up and running on your machine is by using
Docker:

```shell
docker run -d --name carrier-node -v /Users/alice/carrier:/root \
           -p 3500:3500 -p 13000:13000 -p 12000:12000 \
           carrier/metis-carrier
```

This will start `carrier` with a DB memory allowance of 1GB just as the
above command does.  It will also create a persistent volume in your home directory for
saving your data as well as map the default ports. 

Do not forget `--rpc-host 0.0.0.0`, if you want to access GRPC from other containers
and/or hosts. 


### Programmatically interfacing `carrier` nodes

As a developer, sooner rather than later you'll want to start interacting with `carrier`. To aid
this, `carrier` has built-in support for a GRPC based APIs.

HTTP based GRPC API options:

* `--rpc-host` GRPC server listening interface (default: `localhost`)
* `--rpc-port` HTTP-RPC server listening port (default: `8545`)


### Starting with specified node

With the bootnode operational and externally reachable (you can try
`telnet <ip> <port>` to ensure it's indeed reachable), start every subsequent `carrier`
node pointed to the bootnode for peer discovery via the `--bootstrap-node` flag. It will
probably also be desirable to keep the data directory of your private network separated, so
do also specify a custom `--datadir` flag.

```shell
$ carrier --datadir=path/to/custom/data/folder --bootstrap-node=<bootnode-enode-url-from-above>
```

## License

The Metis-Carrier library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The Metis-Carrier binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
