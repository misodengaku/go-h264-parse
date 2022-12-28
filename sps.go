package h264parse

import "fmt"

func (n *NAL) parseSPS() error {
	var err error
	bbr := BitByteReader{}
	bbr.New(n.RBSPByte)

	n.SPS.ProfileIDC, err = bbr.ReadByte()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet0Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet1Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet2Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet3Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet4Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.SPS.ConstraintSet5Flag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {

		bit, err := bbr.ReadBit()
		if err != nil {
			return err
		}
		if bit != 0 {
			return fmt.Errorf("reserved_zero_2bits is not 0")
		}
	}
	n.SPS.LevelIDC, err = bbr.ReadByte()
	if err != nil {
		return err
	}
	n.SPS_ID, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}
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

		n.ChromaFormatIDC, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}
		if n.ChromaFormatIDC == 3 {
			n.SeparateColourPlaneFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}
		}
		n.BitDepthLumaMinus8, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}
		n.BitDepthChromaMinus8, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		n.QPPrimeYZeroTransformBypassFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		n.SeqScalingMatrixPresentFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		if n.SeqScalingMatrixPresentFlag {
			fmt.Printf("[WARN] SeqScalingMatrixPresentFlag is 1. not implemented\n")
			// not implemented
		}

		n.Log2MaxFrameNumMinus4, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		n.PicOrderCntType, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		if n.PicOrderCntType == 0 {
			n.Log2MaxPicOrderCntLSBMinus4, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}
		} else if n.PicOrderCntType == 1 {
			n.DeltaPicOrderAlwaysZeroFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			n.OffsetForNonRefPic, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			n.OffsetForTopToBottomField, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			n.NumRefFramesInPicOrderCntCycle, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			for i := 0; i < int(n.NumRefFramesInPicOrderCntCycle); i++ {
				offset, err := bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.OffsetForRefFrames = append(n.OffsetForRefFrames, offset)
			}
		}

		n.MaxNumRefFrames, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		n.GapsInFrameNumValueAllowedFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		n.PicWidthInMBSMinus1, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		n.PicHeightInMapUnitsMinus1, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		n.FrameMBSOnlyFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		if !n.FrameMBSOnlyFlag {
			n.MBAdaptiveFrameFieldFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}
		}

		n.Direct8x8InferenceFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		n.FrameCroppingFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		if n.FrameCroppingFlag {
			n.FrameCropLeftOffset, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			n.FrameCropRightOffset, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			n.FrameCropTopOffset, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}

			n.FrameCropBottomOffset, err = bbr.ReadExpGolombCode()
			if err != nil {
				return err
			}
		}

		n.VUIParametersPresentFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		if n.VUIParametersPresentFlag {
			n.VUI.AspectRatioInfoPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.AspectRatioInfoPresentFlag {
				n.VUI.AspectRatioIDC, err = bbr.ReadByte()
				if err != nil {
					return err
				}

				if n.VUI.AspectRatioIDC == AspectRatioIDC_ExtendedSAR {
					n.VUI.SARWidth, err = bbr.ReadUInt16()
					if err != nil {
						return err
					}
					n.VUI.SARHeight, err = bbr.ReadUInt16()
					if err != nil {
						return err
					}
				}
			}

			n.VUI.OverscanInfoPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.OverscanInfoPresentFlag {
				n.VUI.OverscanAppropriateFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}
			}

			n.VUI.VideoSignalTypePresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.VideoSignalTypePresentFlag {
				r, err := bbr.ReadBits(3)
				if err != nil {
					return err
				}
				n.VUI.VideoFormat = byte(r & 0x07)

				n.VUI.VideoFullRangeFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}

				n.VUI.ColourDescriptionPresentFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}

				if n.VUI.ColourDescriptionPresentFlag {
					n.VUI.ColourPrimaries, err = bbr.ReadByte()
					if err != nil {
						return err
					}

					n.VUI.TransferCharacteristics, err = bbr.ReadByte()
					if err != nil {
						return err
					}

					n.VUI.MatrixCoefficients, err = bbr.ReadByte()
					if err != nil {
						return err
					}
				}
			}

			n.VUI.ChromaLocInfoPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.ChromaLocInfoPresentFlag {
				n.VUI.ChromaSampleLocTypeTopField, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}

				n.ChromaSampleLocTypeBottomField, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
			}

			n.VUI.TimingInfoPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.TimingInfoPresentFlag {
				n.VUI.NumUnitsInTick, err = bbr.ReadUInt32()
				if err != nil {
					return err
				}

				n.VUI.TimeScale, err = bbr.ReadUInt32()
				if err != nil {
					return err
				}

				n.VUI.FixedFrameRateFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}
			}

			n.VUI.NALHRDParametersPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

		if n.VUI.NALHRDParametersPresentFlag {
			n.VUI.NALHRDParameters, err = n.readHRDParams(&bbr)
			if err != nil {
				return err
			}

		}

			n.VUI.VCLHRDParametersPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

		if n.VUI.VCLHRDParametersPresentFlag {
			n.VUI.VCLHRDParameters, err = n.readHRDParams(&bbr)
			if err != nil {
				return err
			}
		}

			if n.VUI.NALHRDParametersPresentFlag || n.VUI.VCLHRDParametersPresentFlag {
				n.VUI.LowDelayHRDFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}
			}

			n.VUI.PicStructPresentFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			n.VUI.BitstreamRestrictionFlag, err = bbr.ReadBool()
			if err != nil {
				return err
			}

			if n.VUI.BitstreamRestrictionFlag {
				n.VUI.MotionVectorsOverPicBoundariesFlag, err = bbr.ReadBool()
				if err != nil {
					return err
				}
				n.MaxBytesPerPicDenom, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}

				n.MaxBitsPerMBDenom, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.Log2MaxMVLengthHorizontal, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.Log2MaxMVLengthVertical, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.MaxNumReorderFrames, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.MaxDecFrameBuffering, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
			}
		}

		// var trail byte
		// for numBits%8 != 0 {
		// 	byteOffset += numBits / 8
		// 	numBits = numBits % 8
		// 	trail = trail<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
		// 	numBits++
		// }
		// byteOffset += numBits / 8
		// numBits = numBits % 8
		// fmt.Printf("SPS trail: %02x numBits: %d\n", trail, numBits)
		// fmt.Printf("SPS remain(%d): %#v\n", len(n.RBSPByte)-byteOffset, n.RBSPByte[byteOffset:])
	}

	return nil
}

func (n *NAL) readHRDParams(bbr *BitByteReader) (HRDParameters, error) {
	var err error
	hrdparam := HRDParameters{}

	hrdparam.CPBCntMinus1, err = bbr.ReadExpGolombCode()
	if err != nil {
		return HRDParameters{}, err
	}
	hrdparam.BitRateScale, err = bbr.ReadBitsAsByte(4)
	if err != nil {
		return HRDParameters{}, err
	}
	hrdparam.CPBSizeScale, err = bbr.ReadBitsAsByte(4)
	if err != nil {
		return HRDParameters{}, err
	}
	for i := 0; i <= int(hrdparam.CPBCntMinus1); i++ {
		a := AlternativeCPBSpecification{}
		a.BitRateValueMinus1, err = bbr.ReadExpGolombCode()
		if err != nil {
			return HRDParameters{}, err
		}
		a.CPBSizeValueMinus1, err = bbr.ReadExpGolombCode()
		if err != nil {
			return HRDParameters{}, err
		}
		a.CBRFlag, err = bbr.ReadBool()
		if err != nil {
			return HRDParameters{}, err
		}
		hrdparam.AlternativeCPBSpecifications = append(hrdparam.AlternativeCPBSpecifications, a)
	}
	hrdparam.InitialCPBRemovalDelayLengthMinus1, err = bbr.ReadBitsAsByte(5)
	if err != nil {
		return HRDParameters{}, err
	}
	hrdparam.CPBRemovalDelayLengthMinus1, err = bbr.ReadBitsAsByte(5)
	if err != nil {
		return HRDParameters{}, err
	}
	hrdparam.DPBOutputDelayLengthMinus1, err = bbr.ReadBitsAsByte(5)
	if err != nil {
		return HRDParameters{}, err
	}
	hrdparam.TimeOffsetLength, err = bbr.ReadBitsAsByte(5)
	if err != nil {
		return HRDParameters{}, err
	}
	return hrdparam, nil
}
