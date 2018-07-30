package stormtf

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"

	"github.com/golang/protobuf/proto"
)

const maskDelta uint32 = 0xa282ead8

func mask(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + maskDelta
}

func unmask(masked uint32) uint32 {
	rot := masked - maskDelta
	return ((rot >> 17) | (rot << 15))
}

func uint64ToBytes(x uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, x)
	return b
}

var crc32Table = crc32.MakeTable(crc32.Castagnoli)

func crc32Hash(data []byte) uint32 {
	return crc32.Checksum(data, crc32Table)
}

func checksum(data []byte) uint32 {
	crc := crc32.Checksum(data, crc32Table)
	return ((crc >> 15) | (crc << 17)) + maskDelta
}

func verifyChecksum(data []byte, crcMasked uint32) bool {
	rot := crcMasked - maskDelta
	unmaskedCrc := ((rot >> 17) | (rot << 15))

	crc := crc32.Checksum(data, crc32Table)

	return crc == unmaskedCrc
}

func writeTFRecordExample(w io.Writer, example *Example) (int, error) {
	//Format of a single record:
	//  uint64    length
	//  uint32    masked crc of length
	//  byte      data[length]
	//  uint32    masked crc of data

	payload, err := proto.Marshal(example)
	if err != nil {
		return 0, err
	}

	length := len(payload)
	header := make([]byte, 12)
	footer := make([]byte, 4)

	binary.LittleEndian.PutUint64(header[0:8], uint64(length))
	binary.LittleEndian.PutUint32(header[8:12], checksum(header[0:8]))
	binary.LittleEndian.PutUint32(footer[0:4], checksum(payload))

	in1, err := w.Write(header)
	if err != nil {
		return in1, err
	}
	in2, err := w.Write(payload)
	if err != nil {
		return in1 + in2, err
	}
	in3, err := w.Write(footer)
	if err != nil {
		return in1 + in2 + in3, err
	}

	return in1 + in2 + in3, nil
}

func readTFRecordExample(r io.Reader) (*Example, error) {
	header := make([]byte, 12)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}

	crc := binary.LittleEndian.Uint32(header[8:12])
	if !verifyChecksum(header[0:8], crc) {
		return nil, errors.New("Invalid crc for length")
	}

	length := binary.LittleEndian.Uint64(header[0:8])

	payload := make([]byte, length)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}

	footer := make([]byte, 4)
	_, err = io.ReadFull(r, footer)
	if err != nil {
		return nil, err
	}

	crc = binary.LittleEndian.Uint32(footer[0:4])
	if !verifyChecksum(payload, crc) {
		return nil, errors.New("Invalid crc for payload")
	}

	ex := &Example{}
	err = proto.Unmarshal(payload, ex)
	if err != nil {
		return nil, err
	}

	return ex, nil
}
