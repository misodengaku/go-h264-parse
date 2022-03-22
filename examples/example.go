package main

import (
	"fmt"
	"os"

	h264parse "github.com/misodengaku/go-h264-parse"
)

func main() {
	b, _ := os.ReadFile("test.h264")

	nalus, err := h264parse.Unmarshal(b)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\n", nalus)
	filteredNALUs := make([]h264parse.NAL, 0, len(nalus.Units))
	for _, nal := range nalus.Units {
		fmt.Printf("%#v\n", nal.UnitType)
		if nal.UnitType == h264parse.AccessUnitDelimiter {
			continue
		}
		if nal.UnitType == h264parse.SupplementalEnhancementInformation {
			continue
		}
		filteredNALUs = append(filteredNALUs, nal)
	}
	nalus.Units = filteredNALUs
	nb, _ := h264parse.Marshal(nalus)
	os.WriteFile("dump.h264", nb, 0755)
}
