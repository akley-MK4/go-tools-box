package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/bpf"
	"os"
)

func main() {
	LinkType := flag.Uint("link_type", uint(layers.LinkTypeEthernet), "link_type=1")
	captureLength := flag.Int("capture_length", 262144, "capture_length=262144")
	expr := flag.String("expr", "", "expr='tcp and port 8080'")
	flag.Parse()

	pcapPBFInstructions, compileErr := pcap.CompileBPFFilter(layers.LinkType(*LinkType), *captureLength, *expr)
	if compileErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to CompileBPFFilter, %v\n", compileErr)
		os.Exit(1)
		return
	}

	var retInstructions []bpf.RawInstruction
	for _, instruction := range pcapPBFInstructions {
		retInstructions = append(retInstructions, bpf.RawInstruction{
			Op: instruction.Code,
			Jt: instruction.Jt,
			Jf: instruction.Jf,
			K:  instruction.K,
		})
	}

	d, marshalErr := json.Marshal(retInstructions)
	if marshalErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to marshal, %v\n", marshalErr)
		os.Exit(1)
	}

	fmt.Println(string(d))
	os.Exit(0)
}
