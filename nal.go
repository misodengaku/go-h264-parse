package h264parse

import (
	"bytes"
	"fmt"
	"math"
)

func Unmarshal(data []byte) (NALUs, error) {
	var nalus NALUs
	// var index int

	rawNALUs := bytes.Split(data, []byte{0x00, 0x00, 0x01})

	for _, rawNALU := range rawNALUs[1:] {
		// fmt.Printf("%#v\n", rawNALU)
		/*
			for index = 0; index < len(rawNALU); index++ {
				if index < 2 {
					if rawNALU[index] != 0 {
						return NALUs{}, fmt.Errorf("startcode is invalid. index=%d", naluIndex)
					}
				} else {
					if rawNALU[index] == 1 {
						index++
						break
					} else if rawNALU[index] != 0 {
						return NALUs{}, fmt.Errorf("startcode is invalid. index=%d", naluIndex)
					}
				}
			}*/

		nal, err := ParseNAL(rawNALU)
		if err != nil {
			return NALUs{}, err

		}
		nalus.Units = append(nalus.Units, nal)
		// n = NAL{}

	}

	return nalus, nil
}

func Marshal(ns NALUs) ([]byte, error) {
	b := make([]byte, 1, 1024)
	prefix := []byte{0x00, 0x00, 0x01}
	for _, n := range ns.Units {
		b = append(b, prefix...)
		b = append(b, n.HeaderBytes...)
		b = append(b, n.RBSPByte...)
	}
	return b, nil
}

func ParseNAL(data []byte) (NAL, error) {
	index := 0
	n := NAL{}
	// for index < len(data) {
	// fmt.Println("index", index)
	if data[index]>>7&0x01 != 0 {
		return NAL{}, fmt.Errorf("forbidden_zero_bit is not 0")
	}
	n.RefIDC = (data[index] >> 5) & 0x03
	n.UnitType = data[index] & 0x1f
	numBytesInRBSP := 0
	nalUnitHeaderBytes := 1
	if n.UnitType == 14 || n.UnitType == 20 || n.UnitType == 21 {
		fmt.Printf("[BUG] not implemented\n")
		// 	if( nal_unit_type ! = 21 )
		// 		svc_extension_flag All u(1)
		// 	else
		// 		avc_3d_extension_flag All u(1)
		// 		if( svc_extension_flag ) {
		// 			nal_unit_header_svc_extension( ) /* specified in Annex G */ All
		// 			nalUnitHeaderBytes += 3
		// 		} else if( avc_3d_extension_flag ) {
		// 			nal_unit_header_3davc_extension( ) /* specified in Annex J */
		// 			nalUnitHeaderBytes += 2
		// 		} else {
		// 			nal_unit_header_mvc_extension( ) /* specified in Annex H */ All
		// 			nalUnitHeaderBytes += 3
		// 		}
	}
	n.HeaderBytes = data[:nalUnitHeaderBytes]

	index += nalUnitHeaderBytes

	// RBSP
	n.RBSPByte = make([]byte, 0, 16)
	i := 0
	for i = index; i < len(data); i++ {
		if (i+2) < len(data) && (data[i] == 0x00 && data[i+1] == 0x00 && data[i+2] == 0x03) {
			n.RBSPByte = append(n.RBSPByte, data[i], data[i+1])
			i += 2
			numBytesInRBSP += 2
			// 0x03
		} else {
			n.RBSPByte = append(n.RBSPByte, data[i])
			numBytesInRBSP++
		}
	}
	index += numBytesInRBSP

	n.ParseRBSP()

	// end of unit
	// break
	// }
	return n, nil
}

func (n *NAL) ParseRBSP() error {

	switch n.UnitType {
	case SequenceParameterSet:
		err := n.parseSPS()
		if err != nil {
			return err
		}
	case PictureParameterSet:
		err := n.parsePPS()
		if err != nil {
			return err
		}
	case SupplementalEnhancementInformation:
		err := n.parseSEI()
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeExpGolombCode(d []byte, bitOffset int) (uint64, int) {
	leadingZeroBits := 0
	b := byte(0)
	startByte := bitOffset / 8
	shift := 7 - (bitOffset % 8)

	for _, v := range d[startByte:] {
		for ; shift >= 0; shift-- {
			b = (v >> shift) & 0x01
			if b != 0 {
				break
			}
			leadingZeroBits++
		}
		if b != 0 {
			break
		}
		shift = 7
	}
	if leadingZeroBits == 0 {
		return uint64(0), 1
	}
	readBits := uint64(0)
	startByte = (bitOffset + leadingZeroBits) / 8
	shift = 7 - ((bitOffset + leadingZeroBits) % 8) - 1
	i := 0
	for _, v := range d[startByte:] {
		for ; shift >= 0; shift-- {
			readBits = (readBits << 1) | uint64((v>>shift)&0x01)
			i++
			if i >= leadingZeroBits {
				break
			}
		}
		if i <= leadingZeroBits {
			break
		}
		shift = 7
	}
	return uint64(math.Pow(2.0, float64(leadingZeroBits))) - 1 + readBits, leadingZeroBits*2 + 1
}
