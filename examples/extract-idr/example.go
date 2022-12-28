package main

import (
	"fmt"
	"os"

	h264parse "github.com/misodengaku/go-h264-parse"
)

func main() {
	b, err := os.ReadFile("valid.h264")
	if err != nil {
		panic(err)
	}

	nalus, err := h264parse.Unmarshal(b)
	if err != nil {
		panic(err)
	}
	filteredNALUs := make([]h264parse.NAL, 0, len(nalus.Units))
	var lastUnitType h264parse.NALUnitType
	for i, nal := range nalus.Units {
		msg := ""

		// drop non-IDR frame
		if nal.UnitType != h264parse.CodedSliceNonIDRPicture {
			isPass := true
			if nal.UnitType == h264parse.SupplementalEnhancementInformation {
				fmt.Printf("%#v\n", nal.SEI)
			} else if nal.UnitType == h264parse.AccessUnitDelimiter {
				if lastUnitType == h264parse.AccessUnitDelimiter {
					isPass = false
					msg = " => reject"
				}
			}

			if isPass {
				filteredNALUs = append(filteredNALUs, nal)
				msg = " => pass"
			}
			lastUnitType = nal.UnitType
		}
		fmt.Printf("unit %d %s%s\n", i, nal.UnitType.String(), msg)
	}
	nalus.Units = filteredNALUs
	nb, err := h264parse.Marshal(nalus)
	if err != nil {
		panic(err)
	}
	os.WriteFile("dump_valid.h264", nb, 0755)
}
