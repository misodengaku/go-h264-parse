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
			byteOffset += numBits / 8
			numBits = numBits % 8
			n.SeparateColourPlaneFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++
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
			fmt.Printf("[WARN] SeqScalingMatrixPresentFlag is 1. not implemented\n")
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
			byteOffset += numBits / 8
			numBits = numBits % 8
			n.Log2MaxPicOrderCntLSBMinus4, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits
		} else if n.PicOrderCntType == 1 {
			// not implemented
			fmt.Printf("[WARN] PicOrderCntType is 1. not implemented\n")

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.DeltaPicOrderAlwaysZeroFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.OffsetForNonRefPic, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits
			byteOffset += numBits / 8

			numBits = numBits % 8
			n.OffsetForTopToBottomField, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits
			byteOffset += numBits / 8

			numBits = numBits % 8
			n.NumRefFramesInPicOrderCntCycle, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
			numBits += numResultBits

			for i := 0; i < int(n.NumRefFramesInPicOrderCntCycle); i++ {
				var offset uint64
				numBits = numBits % 8
				offset, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
				n.OffsetForRefFrames = append(n.OffsetForRefFrames, offset)
			}
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

		if !n.FrameMBSOnlyFlag {
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
			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.AspectRatioInfoPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.AspectRatioInfoPresentFlag {
				n.VUI.AspectRatioIDC = 0
				for i := 0; i < 8; i++ {
					byteOffset += numBits / 8
					numBits = numBits % 8
					n.VUI.AspectRatioIDC = (n.VUI.AspectRatioIDC << 1) | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
					numBits++
				}
				if n.VUI.AspectRatioIDC == AspectRatioIDC_ExtendedSAR {
					n.VUI.SARWidth = 0
					n.VUI.SARHeight = 0
					for i := 0; i < 16; i++ {
						byteOffset += numBits / 8
						numBits = numBits % 8
						n.VUI.SARWidth = (n.VUI.SARWidth << 1) | uint16((n.RBSPByte[byteOffset]>>(7-numBits))&0x01)
						numBits++
					}
					for i := 0; i < 16; i++ {
						byteOffset += numBits / 8
						numBits = numBits % 8
						n.VUI.SARHeight = (n.VUI.SARHeight << 1) | uint16((n.RBSPByte[byteOffset]>>(7-numBits))&0x01)
						numBits++
					}
				}
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.OverscanInfoPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.OverscanInfoPresentFlag {
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.OverscanAppropriateFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.VideoSignalTypePresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.VideoSignalTypePresentFlag {
				n.VUI.VideoFormat = 0
				for i := 0; i < 3; i++ {
					byteOffset += numBits / 8
					numBits = numBits % 8
					n.VUI.VideoFormat = (n.VUI.VideoFormat << 1) | (n.RBSPByte[byteOffset]>>(7-numBits))&0x01
					numBits++
				}

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.VideoFullRangeFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.ColourDescriptionPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++

				if n.VUI.ColourDescriptionPresentFlag {

					n.VUI.ColourPrimaries = 0
					n.VUI.TransferCharacteristics = 0
					n.VUI.MatrixCoefficients = 0
					for i := 0; i < 8; i++ {
						byteOffset += numBits / 8
						numBits = numBits % 8
						n.VUI.ColourPrimaries = (n.VUI.ColourPrimaries << 1) | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
						numBits++
					}
					for i := 0; i < 8; i++ {
						byteOffset += numBits / 8
						numBits = numBits % 8
						n.VUI.TransferCharacteristics = (n.VUI.TransferCharacteristics << 1) | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
						numBits++
					}
					for i := 0; i < 8; i++ {
						byteOffset += numBits / 8
						numBits = numBits % 8
						n.VUI.MatrixCoefficients = (n.VUI.MatrixCoefficients << 1) | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
						numBits++
					}
				}
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.ChromaLocInfoPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.ChromaLocInfoPresentFlag {
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.ChromaSampleLocTypeTopField, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.ChromaSampleLocTypeBottomField, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.TimingInfoPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.TimingInfoPresentFlag {
				n.VUI.NumUnitsInTick = 0
				n.VUI.TimeScale = 0
				for i := 0; i < 32; i++ {
					byteOffset += numBits / 8
					numBits = numBits % 8
					n.VUI.NumUnitsInTick = (n.VUI.NumUnitsInTick << 1) | uint32((n.RBSPByte[byteOffset]>>(7-numBits))&0x01)
					numBits++
				}
				for i := 0; i < 32; i++ {
					byteOffset += numBits / 8
					numBits = numBits % 8
					n.VUI.TimeScale = (n.VUI.TimeScale << 1) | uint32((n.RBSPByte[byteOffset]>>(7-numBits))&0x01)
					numBits++
				}
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.FixedFrameRateFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.NALHRDParametersPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.NALHRDParametersPresentFlag {
				// not implemented
				fmt.Printf("[WARN] NALHRDParametersPresentFlag is not implemented\n")
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.VCLHRDParametersPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.VCLHRDParametersPresentFlag {
				// not implemented
				fmt.Printf("[WARN] VCLHRDParametersPresentFlag is not implemented\n")
			}

			if n.VUI.NALHRDParametersPresentFlag || n.VUI.VCLHRDParametersPresentFlag {
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.LowDelayHRDFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++
			}

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.PicStructPresentFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			byteOffset += numBits / 8
			numBits = numBits % 8
			n.VUI.BitstreamRestrictionFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
			numBits++

			if n.VUI.BitstreamRestrictionFlag {
				byteOffset += numBits / 8
				numBits = numBits % 8
				n.VUI.MotionVectorsOverPicBoundariesFlag = ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01) == 1
				numBits++

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.MaxBytesPerPicDenom, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.MaxBitsPerMBDenom, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.Log2MaxMVLengthHorizontal, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.Log2MaxMVLengthVertical, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.MaxNumReorderFrames, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits

				byteOffset += numBits / 8
				numBits = numBits % 8
				n.MaxDecFrameBuffering, numResultBits = decodeExpGolombCode(n.RBSPByte[byteOffset:], numBits)
				numBits += numResultBits
			}
		}

		var trail byte
		for numBits%8 != 0 {
			byteOffset += numBits / 8
			numBits = numBits % 8
			trail = trail<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
			numBits++
		}
		// byteOffset += numBits / 8
		// numBits = numBits % 8
		// fmt.Printf("SPS trail: %02x numBits: %d\n", trail, numBits)
		// fmt.Printf("SPS remain(%d): %#v\n", len(n.RBSPByte)-byteOffset, n.RBSPByte[byteOffset:])
	}

	return nil
}
