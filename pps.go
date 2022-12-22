package h264parse

func (n *NAL) parsePPS() error {
	var err error
	bbr := BitByteReader{}
	bbr.New(n.RBSPByte)

	n.PPS.ID, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.SPS_ID, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.EntropyCodingModeFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	n.PPS.BottomFieldPicOrderInFramePresentFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}
	n.PPS.NumSliceGroupsMinus1, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	if n.PPS.NumSliceGroupsMinus1 > 0 {

		n.PPS.SliceGroupMapType, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}

		if n.SliceGroupMapType == 0 {
			for i := 0; i < int(n.PPS.NumSliceGroupsMinus1); i++ {
				rl, err := bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}

				n.RunLengthMinus1 = append(n.RunLengthMinus1, rl)
			}
		} else if n.SliceGroupMapType == 2 {
			for i := 0; i < int(n.PPS.NumSliceGroupsMinus1); i++ {
				rl, err := bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
				n.TopLeft = append(n.RunLengthMinus1, rl)

				rl, err = bbr.ReadExpGolombCode()
				if err != nil {
					return err
				}
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

	n.PPS.NumRefIndexL0DefaultActiveMinus1, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.NumRefIndexL1DefaultActiveMinus1, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.WeightedPredFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	// 1/2
	highBit, err := bbr.ReadBit()
	if err != nil {
		return err
	}

	// 2/2
	lowBit, err := bbr.ReadBit()
	if err != nil {
		return err
	}
	n.PPS.WeightedBipredIdc = highBit<<1 | lowBit

	n.PPS.PicInitQPMinus26, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.PicInitQSMinus26, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.ChromaQPIndexOffset, err = bbr.ReadExpGolombCode()
	if err != nil {
		return err
	}

	n.PPS.DeblockingFilterControlPresentFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	n.PPS.ConstrainedIntraPredFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	n.PPS.RedundantPicCntPresentFlag, err = bbr.ReadBool()
	if err != nil {
		return err
	}

	if bbr.RemainBytes() > 0 || bbr.CurrentBit() > 0 {
		n.PPS.Transform8x8ModeFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

		n.PPS.PicScalingMatrixPresentFlag, err = bbr.ReadBool()
		if err != nil {
			return err
		}

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
				f, err := bbr.ReadBool()
				if err != nil {
					return err
				}

				n.PPS.SeqScalingListPresentFlags[i] = f
			}
		}

		n.PPS.SecondChromaQPIndexOffset, err = bbr.ReadExpGolombCode()
		if err != nil {
			return err
		}
	}

	// bbr.GetTrailBytes()
	// byteOffset += numBits / 8
	// numBits = numBits % 8
	// fmt.Printf("PPS trail: %02x numBits: %d\n", trail, numBits)
	// fmt.Printf("PPS remain(%d): %#v\n", len(n.RBSPByte)-byteOffset, n.RBSPByte[byteOffset:])

	return nil
}
