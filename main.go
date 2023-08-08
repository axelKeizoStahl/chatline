package main

import (
    "flag"
    "chatline/client"
    "chatline/server"
)

func main() {
    conntype := flag.String("type", "client", "server or client type")
    flag.Parse()

    switch *conntype {
        case "server":
            server.Server()
        case "client":
            client.Client()
    }
}
