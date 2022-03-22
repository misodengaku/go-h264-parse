package h264parse

func (n *NAL) parseSEI() error {

	numBits := 0
	byteOffset := 0
	// numResultBits := 0

	n.SEI.PayloadType = 0
	n.SEI.PayloadSize = 0
	nextBits := n.RBSPByte[byteOffset]

	for {
		if nextBits == 0xff {
			n.PayloadType += 255
			numBits += 8
			byteOffset += numBits / 8
			numBits = numBits % 8
			nextBits = n.RBSPByte[byteOffset]
			continue
		}
		break
	}

	n.PayloadType += int(nextBits)
	numBits += 8
	byteOffset += numBits / 8
	numBits = numBits % 8
	nextBits = n.RBSPByte[byteOffset]

	// read size
	for {
		if nextBits == 0xff {
			n.PayloadSize += 255
			numBits += 8
			byteOffset += numBits / 8
			numBits = numBits % 8
			nextBits = n.RBSPByte[byteOffset]
			continue
		}
		break
	}

	n.PayloadSize += int(nextBits)
	numBits += 8
	byteOffset += numBits / 8
	numBits = numBits % 8
	nextBits = n.RBSPByte[byteOffset]

	return nil
}
