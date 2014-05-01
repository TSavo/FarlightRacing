// +build ignore

package main

import "github.com/TSavo/GreatSpaceRace"

func main() {
	server := greatspacerace.NewServer();
	server.Listen(8080);
}