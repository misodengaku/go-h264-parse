package h264parse

func (n *NAL) parseSEI() error {
	bbr := BitByteReader{}
	bbr.New(n.RBSPByte)

	n.SEI.PayloadType = 0
	n.SEI.PayloadSize = 0

	nextBits, err := bbr.ReadByte()
	if err != nil {
		return err
	}
	for nextBits == 0xff {
		n.SEI.PayloadType += 255
		nextBits, err = bbr.ReadByte()
		if err != nil {
			return err
		}
	}
	n.SEI.PayloadType += int(nextBits) // last_payload_type_byte

	// read size
	nextBits, err = bbr.ReadByte()
	if err != nil {
		return err
	}
	for nextBits == 0xff {
		n.SEI.PayloadSize += 255
		nextBits, err = bbr.ReadByte()
		if err != nil {
			return err
		}
	}
	n.SEI.PayloadSize += int(nextBits) // last_payload_size_byte

	if n.SEI.PayloadSize >= 8 {
		n.SEI.PayloadBytes, err = bbr.ReadBytes(n.SEI.PayloadSize / 8)
		if err != nil {
			return err
		}
	}
	if n.SEI.PayloadSize%8 != 0 {
		r, err := bbr.ReadBits(n.SEI.PayloadSize % 8)
		if err != nil {
			return err
		}
		n.SEI.PayloadBytes = append(n.SEI.PayloadBytes, byte(r))
	}

	// FIXME: implement SEI.parsePayload
	// n.SEI.parsePayload()

	return nil
}
