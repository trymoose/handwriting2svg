package handwriting

import (
    "bytes"
    "fmt"
    "github.com/trymoose/errors"
    "github.com/trymoose/handwriting2svg/pkg/handwriting/internal/handwriting"
    "github.com/ulikunitz/xz"
    "io"
)

func decompress(pr *handwriting.BaseMessage) ([]byte, error) {
    switch pr.Handwriting.Compression {
    case handwriting.Compression_None:
        return pr.Handwriting.Strokes, nil
    case handwriting.Compression_XZ:
        if pr.Handwriting.DecompressedLength == nil {
            return nil, errors.New("decompressed length is not set for compressed payload")
        }

        data, err := io.ReadAll(errors.Get(xz.NewReader(bytes.NewReader(pr.Handwriting.Strokes))))
        if err != nil {
            return nil, err
        } else if len(data) != int(*pr.Handwriting.DecompressedLength) {
            return nil, fmt.Errorf("expected decompressed length of %d, got %d", *pr.Handwriting.DecompressedLength, len(data))
        }
        return data, nil
    default:
        return nil, errors.New("unknown compression algorithm")
    }
}
