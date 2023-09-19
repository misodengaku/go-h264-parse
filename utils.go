package h264parse

import (
	"errors"
	"fmt"
)

type BitByteReader struct {
	data       []byte
	numBits    int
	byteOffset uint64
}

var EOS = errors.New("End of stream reached")

func (b *BitByteReader) New(data []byte) {
	b.data = data
	b.numBits = 0
	b.byteOffset = 0
}

// Rec. ITU-T H.264 (08/2021) p.41
func (b *BitByteReader) MoreDataInByteStream() bool {
	return b.RemainBytes() > 0 || b.CurrentBit() > 0
}

// Rec. ITU-T H.264 (08/2021) p.41
func (b *BitByteReader) MoreRBSPData() bool {
	if !b.MoreDataInByteStream() {
		return false
	}
	return !b.CheckRBSPTrailingBitsOnly()
}

// Rec. ITU-T H.264 (08/2021) p.41
func (b *BitByteReader) CheckRBSPTrailingBitsOnly() bool {
	byteOffset := b.byteOffset
	numBits := b.numBits
	isStopBitFound := false
	for {
		result, err := b.PeekBit(byteOffset, numBits)
		if err != nil {
			if errors.Is(err, EOS) {
				break
			}
			panic(err)
		}
		if !isStopBitFound {
			if result == 1 {
				isStopBitFound = true
			} else {
				return false
			}
		} else {
			if result != 0 {
				return false
			}
		}
		numBits++
		byteOffset += uint64(numBits / 8)
		numBits = numBits % 8
	}

	return isStopBitFound
}

// Rec. ITU-T H.264 (08/2021) p.211
func (b *BitByteReader) ReadSignedExpGolombCode() (int64, error) {
	if len(b.data) <= int(b.byteOffset) {
		return 0, fmt.Errorf("ReadExpGolombCode: not enough data")
	}
	result, resultBits := decodeSignedExpGolombCode(b.data[b.byteOffset:], b.numBits)
	b.numBits += resultBits
	b.byteOffset += uint64(b.numBits / 8)
	b.numBits = b.numBits % 8

	return result, nil
}

// Rec. ITU-T H.264 (08/2021) pp.209-210
func (b *BitByteReader) ReadExpGolombCode() (uint64, error) {
	if len(b.data) <= int(b.byteOffset) {
		return 0, fmt.Errorf("not enough data")
	}
	result, resultBits := decodeExpGolombCode(b.data[b.byteOffset:], b.numBits)
	b.numBits += resultBits
	b.byteOffset += uint64(b.numBits / 8)
	b.numBits = b.numBits % 8

	return result, nil
}

func (b *BitByteReader) PeekBit(offset uint64, numBits int) (byte, error) {
	if len(b.data) <= int(offset) {
		return 0, EOS
	}
	return (b.data[offset] >> (7 - numBits)) & 0x01, nil
}

func (b *BitByteReader) ReadBit() (byte, error) {
	result, err := b.PeekBit(b.byteOffset, b.numBits)
	if err != nil {
		return 0, err
	}
	b.numBits++
	b.byteOffset += uint64(b.numBits / 8)
	b.numBits = b.numBits % 8

	return result, nil
}

func (b *BitByteReader) ReadBitsAsByte(bits int) (byte, error) {
	r, err := b.ReadBits(bits)
	if err != nil {
		return 0, err
	}
	return byte(r), err
}

func (b *BitByteReader) ReadBits(bits int) (uint64, error) {
	if len(b.data) <= int(b.byteOffset) {
		return 0, fmt.Errorf("not enough data")
	}
	if bits > 64 {
		return 0, fmt.Errorf("ReadBits can only read up to 64 bits")
	}
	var r uint64
	for i := 0; i < bits; i++ {
		bit, err := b.ReadBit()
		if err != nil {
			return 0, err
		}
		r = r<<1 | uint64(bit)
	}
	return r, nil
}

func (b *BitByteReader) ReadByte() (byte, error) {
	result, err := b.ReadBits(8)
	if err != nil {
		return 0, err
	}

	return byte(result), nil
}

func (b *BitByteReader) ReadBytes(size int) ([]byte, error) {
	if b.RemainBytes() < size {
		return nil, fmt.Errorf("not enough data")
	}

	r := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		rb, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		r = append(r, rb)
	}

	return r, nil
}

func (b *BitByteReader) ReadUInt16() (uint16, error) {
	result, err := b.ReadBits(16)
	if err != nil {
		return 0, err
	}

	return uint16(result), nil
}

func (b *BitByteReader) ReadUInt32() (uint32, error) {
	result, err := b.ReadBits(32)
	if err != nil {
		return 0, err
	}

	return uint32(result), nil
}

func (b *BitByteReader) ReadBool() (bool, error) {
	r, err := b.ReadBit()
	return r == 1, err
}

func (b *BitByteReader) RemainBytes() int {
	return len(b.data[b.byteOffset:])
}

func (b *BitByteReader) CurrentBit() int {
	return b.numBits
}

// func (b *BitByteReader) GetTrailBytes() ([]byte, int) {

// 	var trail byte
// 	readBit := 0
// 	for b.numBits%8 != 0 {
// 		b.byteOffset += uint64(b.numBits / 8)
// 		if b.byteOffset > uint64(len(b.data)) {
// 			return []byte{}, readBit
// 		}
// 		b.numBits = b.numBits % 8
// 		trail = trail<<1 | ((n.RBSPByte[byteOffset] >> (7 - numBits)) & 0x01)
// 		numBits++
// 	}
// }
