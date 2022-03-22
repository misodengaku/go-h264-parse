package h264parse

import "fmt"

func (n *NAL) parsePPS() error {

	numBits := 0
	byteOffset := 0
	numResultBits := 0

	n.PPS.ID, numResultBits = decodeExpGolombCode(n.RBSPByte[0:], 0)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.SPS_ID, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.EntropyCodingModeFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.BottomFieldPicOrderInFramePresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.NumSliceGroupsMinus1, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	if n.PPS.NumSliceGroupsMinus1 > 0 {
		n.PPS.SliceGroupMapType, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits
		byteOffset += numBits / 8
		numBits = numBits % 8

		if n.SliceGroupMapType == 0 {
			for i := 0; i < int(n.PPS.NumSliceGroupsMinus1); i++ {
				var rl uint64
				rl, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.RunLengthMinus1 = append(n.RunLengthMinus1, rl)
			}
		} else if n.SliceGroupMapType == 2 {
			for i := 0; i < int(n.PPS.NumSliceGroupsMinus1); i++ {
				var rl uint64
				rl, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.TopLeft = append(n.RunLengthMinus1, rl)

				rl, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.BottomRight = append(n.RunLengthMinus1, rl)
			}
		} else if n.SliceGroupMapType == 3 ||
			n.SliceGroupMapType == 4 ||
			n.SliceGroupMapType == 5 {
			// not implemented
		} else if n.SliceGroupMapType == 6 {
			// not implemented
		}
	}

	n.PPS.NumRefIndexL0DefaultActiveMinus1, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.NumRefIndexL1DefaultActiveMinus1, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.WeightedPredFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	// 1/2
	n.PPS.WeightedBipredIdc = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	// 2/2
	n.PPS.WeightedBipredIdc = n.PPS.WeightedBipredIdc<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.PicInitQPMinus26, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.PicInitQSMinus26, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.ChromaQPIndexOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
	numBits += numResultBits
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.DeblockingFilterControlPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.ConstrainedIntraPredFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	n.PPS.RedundantPicCntPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
	numBits++
	byteOffset += numBits / 8
	numBits = numBits % 8

	// for numBits%8 != 0 {
	// 	n.PPS.WeightedBipredIdc = n.PPS.WeightedBipredIdc<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
	// 	numBits++
	// 	byteOffset += numBits / 8
	// 	numBits = numBits % 8
	// }
	if len(n.RBSPByte[byteOffset:]) > 0 || numBits > 0 {
		n.PPS.Transform8x8ModeFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++
		byteOffset += numBits / 8
		numBits = numBits % 8

		n.PPS.PicScalingMatrixPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++
		byteOffset += numBits / 8
		numBits = numBits % 8

		if n.PPS.PicScalingMatrixPresentFlag {
			t := 0
			if n.PPS.Transform8x8ModeFlag {
				t = 2 // FIXME
				// if n.SPS.ChromaFormatIDC != 3 {
				// 	t = 2
				// } else {
				// 	t = 6
				// }
			}
			n.PPS.SeqScalingListPresentFlags = make([]bool, 6+t)
			for i := 0; i < 6+t; i++ {
				f := ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++
				byteOffset += numBits / 8
				numBits = numBits % 8

				n.PPS.SeqScalingListPresentFlags[i] = f
			}
		}

		n.PPS.SecondChromaQPIndexOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits
		byteOffset += numBits / 8
		numBits = numBits % 8
	}

	var trail byte
	for numBits%8 != 0 {
		byteOffset += numBits / 8
		if byteOffset > len(n.RBSPByte) {
			return fmt.Errorf("RBSP out of range")
		}
		numBits = numBits % 8
		trail = trail<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
		numBits++
	}
	byteOffset += numBits / 8
	numBits = numBits % 8
	fmt.Printf("trail: %02x\n", trail)
	fmt.Printf("remain: %#v\n", n.RBSPByte[byteOffset:])

	return nil
}
