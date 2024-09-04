package main

import (
    "flag"
    "log"
    "net"
    "github.com/coredns/coredns/coremain"
    "github.com/coredns/coredns/plugin"
    "github.com/coredns/coredns/plugin/pkg/parse"
    "github.com/coredns/coredns/plugin/pkg/upstream"
)

func main() {
    flag.Parse()

    // Initialize CoreDNS with your plugin
    coremain.Run(plugin.NewServer(), plugin.New(), upstream.New())

    // Run CoreDNS server
    log.Fatal(coremain.RunAndServe())
}
