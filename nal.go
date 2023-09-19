package h264parse

import (
	"bytes"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

func Unmarshal(data []byte) (NALUs, error) {
	var nalus NALUs
	// var index int

	rawNALUs := bytes.Split(data, []byte{0x00, 0x00, 0x01})

	for index, rawNALU := range rawNALUs[1:] {
		nal, err := ParseNAL(rawNALU)
		if err != nil {
			return NALUs{}, errors.Wrap(err, fmt.Sprintf("ParseNAL failed at packet %d", index))
		}
		err = nal.ParseRBSP()
		if err != nil {
			return NALUs{}, errors.Wrap(err, fmt.Sprintf("ParseRBSP failed at packet %d", index))
		}
		nalus.Units = append(nalus.Units, nal)
		// n = NAL{}

	}

	return nalus, nil
}

func Marshal(ns NALUs) ([]byte, error) {
	b := make([]byte, 0, 1024)
	prefix := []byte{0x00, 0x00, 0x01}
	for _, n := range ns.Units {
		b = append(b, prefix...)
		b = append(b, n.HeaderBytes...)
		b = append(b, n.GetEBSP()...)
	}
	return b, nil
}

func ParseNAL(data []byte) (NAL, error) {
	if len(data) == 0 {
		return NAL{}, fmt.Errorf("data is empty")
	}
	n := NAL{}
	if data[0]>>7&0x01 != 0 {
		return NAL{}, fmt.Errorf("forbidden_zero_bit is not 0")
	}
	n.RefIDC = (data[0] >> 5) & 0x03
	n.UnitType = NALUnitType(data[0] & 0x1f)
	numBytesInRBSP := 0
	nalUnitHeaderBytes := 1
	if n.UnitType == PrefixNALUnit || n.UnitType == CodedSliceExtension || n.UnitType == CodedSliceExtensionDepthViewComponentOr3DAVCTextureView {
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

	index := nalUnitHeaderBytes

	// RBSP
	n.RBSPByte = make([]byte, 0, 16)
	i := 0
	for i = index; i < len(data); i++ {
		// unescape EBSP
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
	return n, nil
}

func (n *NAL) GetEBSP() []byte {
	rbsp := n.RBSPByte
	state := 0
	ebsp := []byte{}
	for _, b := range rbsp {
		switch state {
		case 0:
			fallthrough
		case 1:
			if b == 0 {
				state++
			} else {
				state = 0
			}
			ebsp = append(ebsp, b)
		case 2:
			if b == 0 || b == 1 || b == 2 || b == 3 {
				ebsp = append(ebsp, 0x03, b)
			} else {
				ebsp = append(ebsp, b)
			}
			if b == 0 {
				state = 1
			} else {
				state = 0
			}
		}
	}
	return ebsp
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

func decodeSignedExpGolombCode(d []byte, bitOffset int) (int64, int) {
	codeNum, l := decodeExpGolombCode(d, bitOffset)
	sev := int64(codeNum & 0x01)
	sev += int64(codeNum >> 1)

	if codeNum >= 2 && codeNum%2 == 0 {
		sev *= -1
	}
	return sev, l
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
	if shift == 0 {
		startByte++
		shift = 7
	} else {
		shift = 7 - ((bitOffset + leadingZeroBits) % 8) - 1
	}
	i := 0

ShiftLoop:
	for _, v := range d[startByte:] {
		for ; shift >= 0; shift-- {
			readBits = (readBits << 1) | uint64((v>>shift)&0x01)
			i++
			if i >= leadingZeroBits {
				break ShiftLoop
			}
		}
		shift = 7
	}
	return uint64(math.Pow(2.0, float64(leadingZeroBits))) - 1 + readBits, leadingZeroBits*2 + 1
}
