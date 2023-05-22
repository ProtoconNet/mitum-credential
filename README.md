### mitum-credential

*mitum-credential* is a [mitum](https://github.com/ProtoconNet/mitum2)-based contract model and is a service that provides credentials to be issued and updated using DID.

#### Installation

```sh
$ git clone https://github.com/ProtoconNet/mitum-credential

$ cd mitum-credential

$ go build -o ./mc ./main.go
```

#### Run

```sh
$ ./mc init --design=<config file> <genesis config file>

$ ./mc run --design=<config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
[genesis-design.yml](genesis-design.yml) is a sample of `genesis config file`.