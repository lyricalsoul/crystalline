# Crystalline: a SeedLink client written in Go

## Running
You can run the CLI right now.
```bash
go run cli.go host:port
```
Where `host:port` is the address of the SeedLink server you want to connect to. For example:
```bash
go run cli.go seisrequest.iag.usp.br:18000
```
This will connect to the Universidade de São Paulo's SeedLink server, ran by the Instituto de Astronomia, Geofísica e Ciências Atmosféricas (IAG/USP).

Then, proceed as normal: send a `HELLO` command to the server to initiate the handshake. After that, consult the [SeedLink protocol documentation](https://docs.fdsn.org/projects/seedlink/en/latest/protocol.html) to see what to do next.
An example of a `HELLO` command is:
```bash
[renan@kingston crystalline]$ go run cli.go seisrequest.iag.usp.br:18000
Connecting to server seisrequest.iag.usp.br on port 18000...
Crystalline Interactive Shell. Type 'exit' to quit.
> hello
SeedLink v3.3 (2024.020)
Centro de Sismologia da USP

>
```

## TODO
- [ ] write common behaviors (such as a `Initiate` function that sends HELLO and USERAGENT on v4, HELLO on v3);
- [ ] miniSEED support for data mode;
- [ ] automatic station discovery (via CAT on v3, INFO STATIONS on v4. btw CAT isn't available on USGS servers);
- [ ] abstract behaviors so both v3 and v4 servers are supported;
- [ ] distribute alongside a GUI to quickly connect and visualize real-time data.
