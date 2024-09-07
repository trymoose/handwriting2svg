package handwriting

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "github.com/trymoose/errors"
    "github.com/trymoose/handwriting2svg/pkg/handwriting/internal/handwriting"
    "io"
)

func readStrokes(pr *handwriting.BaseMessage) ([][]StrokePoint, error) {
    data, err := decompress(pr)
    if err != nil {
        return nil, err
    }

    a := make([][]StrokePoint, 0, pr.Handwriting.StrokesCount)
    r := bytes.NewReader(data)
    for {
        var length uint16
        err := binary.Read(r, binary.LittleEndian, &length)
        if errors.Is(err, io.EOF) {
            if len(a) != int(pr.Handwriting.StrokesCount) {
                return nil, fmt.Errorf("expected strokes count of %d, got %d", pr.Handwriting.StrokesCount, len(a))
            }
            return a, nil
        } else if err != nil {
            return nil, err
        }

        aa := make([]StrokePoint, int(length))
        for i := range aa {
            if err := binary.Read(r, binary.LittleEndian, &aa[i]); err != nil {
                return nil, err
            }
        }
        a = append(a, aa)
    }
}
