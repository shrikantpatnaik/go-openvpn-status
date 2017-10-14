# go-openvpn-status

A simple parser for parsing OpenVPN status files

## Usage

```go
status, _ := openvpnStatus.ParseFile("examples/server.status")

fmt.Printf(status.UpdateAt)

fmt.Printf(status.ClientList)

```

For more you can see the openvpnStatus.go file and look at the structs in the top of the file
