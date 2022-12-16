package main

import (
	"fmt"
	"os"

	h264parse "github.com/misodengaku/go-h264-parse"
)

func main() {
	b, err := os.ReadFile("valid-1755.h264")
	if err != nil {
		panic(err)
	}

	nalus, err := h264parse.Unmarshal(b)
	if err != nil {
		panic(err)
	}
	filteredNALUs := make([]h264parse.NAL, 0, len(nalus.Units))
	for i, nal := range nalus.Units {
		msg := ""

		// drop non-IDR frame
		if nal.UnitType != h264parse.CodedSliceNonIDRPicture {
			filteredNALUs = append(filteredNALUs, nal)
			msg = " => pass"
		}
		fmt.Printf("unit %d %s%s\n", i, nal.UnitType.String(), msg)
	}
	nalus.Units = filteredNALUs
	nb, err := h264parse.Marshal(nalus)
	if err != nil {
		panic(err)
	}
	os.WriteFile("dump.h264", nb, 0755)
}
