package h264parse

import "fmt"

func (n *NAL) parseSPS() error {

	numBits := 0
	byteOffset := 0
	numResultBits := 0

	n.SPS.ProfileIDC = n.RBSPByte[0]
	n.SPS.ConstraintSet0Flag = (n.RBSPByte[1] >> 7) & 0x01
	n.SPS.ConstraintSet1Flag = (n.RBSPByte[1] >> 6) & 0x01
	n.SPS.ConstraintSet2Flag = (n.RBSPByte[1] >> 5) & 0x01
	n.SPS.ConstraintSet3Flag = (n.RBSPByte[1] >> 4) & 0x01
	n.SPS.ConstraintSet4Flag = (n.RBSPByte[1] >> 3) & 0x01
	if (n.RBSPByte[1] & 0x03) != 0 {
		return fmt.Errorf("reserved_zero_2bits is not 0")
	}
	n.SPS.LevelIDC = n.RBSPByte[2]
	n.SPS.ID, numBits = decodeExpGolombCode(n.RBSPByte[3:], 0)
	byteOffset = 3
	if n.SPS.ProfileIDC == 100 ||
		n.SPS.ProfileIDC == 110 ||
		n.SPS.ProfileIDC == 122 ||
		n.SPS.ProfileIDC == 244 ||
		n.SPS.ProfileIDC == 44 ||
		n.SPS.ProfileIDC == 83 ||
		n.SPS.ProfileIDC == 86 ||
		n.SPS.ProfileIDC == 118 ||
		n.SPS.ProfileIDC == 128 ||
		n.SPS.ProfileIDC == 138 ||
		n.SPS.ProfileIDC == 139 ||
		n.SPS.ProfileIDC == 134 ||
		n.SPS.ProfileIDC == 135 {
		n.ChromaFormatIDC, numResultBits = decodeExpGolombCode(n.RBSPByte[3:], numBits)
		numBits += numResultBits
		if n.ChromaFormatIDC == 3 {
			// separate_colour_plane_flag 0 u(1)
			fmt.Printf("[BUG] separate_colour_plane_flag is 1\n")
		}
		n.BitDepthLumaMinus8, numResultBits = decodeExpGolombCode(n.RBSPByte[3:], numBits)
		numBits += numResultBits
		n.BitDepthChromaMinus8, numResultBits = decodeExpGolombCode(n.RBSPByte[3:], numBits)
		numBits += numResultBits
		byteOffset += numBits / 8
		numBits = numBits % 8
		n.QPPrimeYZeroTransformBypassFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		n.SeqScalingMatrixPresentFlag = ((n.RBSPByte[byteOffset] >> (6 - numBits)) & 0x01) == 1
		numBits += 2
		if n.SeqScalingMatrixPresentFlag {
			fmt.Printf("[BUG] SeqScalingMatrixPresentFlag is 1\n")
			// not implemented
		}
		byteOffset += numBits / 8
		numBits = numBits % 8
		n.Log2MaxFrameNumMinus4, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.PicOrderCntType, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits

		if n.PicOrderCntType == 0 {
			// not implemented
			fmt.Printf("[BUG] PicOrderCntType is 0\n")
		} else if n.PicOrderCntType == 1 {
			// not implemented
			fmt.Printf("[BUG] PicOrderCntType is 1\n")
		}

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.MaxNumRefFrames, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.GapsInFrameNumValueAllowedFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.PicWidthInMBSMinus1, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.PicHeightInMapUnitsMinus1, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
		numBits += numResultBits

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.FrameMBSOnlyFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++

		if n.FrameMBSOnlyFlag {
			byteOffset += numBits / 8
			numBits = numBits % 8
			n.MBAdaptiveFrameFieldFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++
		}

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.Direct8x8InferenceFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.FrameCroppingFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++

		if n.FrameCroppingFlag {

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.FrameCropLeftOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.FrameCropRightOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.FrameCropTopOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.FrameCropBottomOffset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits
		}

		byteOffset += numBits / 8
		numBits = numBits % 8
		n.VUIParametersPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
		numBits++

		if n.VUIParametersPresentFlag {
			// not implemented
			fmt.Printf("[WARN] VUIParametersPresentFlag is not implemented\n")
		}

		// byteOffset += numBits / 8
		// numBits = numBits % 8

		var trail byte
		for numBits%8 != 0 {

			byteOffset += numBits / 8
			numBits = numBits % 8
			trail = trail<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
			numBits++
		}
		byteOffset += numBits / 8
		numBits = numBits % 8
		fmt.Printf("SPS trail: %02x\n", trail)
		fmt.Printf("SPS remain: %#v\n", n.RBSPByte[byteOffset:])
	}

	return nil
}
