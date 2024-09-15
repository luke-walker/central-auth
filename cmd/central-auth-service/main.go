package main

import (
    "flag"

    "central-auth/internal/server"
)

func main() {
    addrFlag := flag.String("addr", "127.0.0.1", "specify the address")
    portFlag := flag.Int("port", 6060, "specify the port")
    flag.Parse()

    authServer := server.NewAuthServer(*addrFlag, *portFlag)
    authServer.Start()
}
