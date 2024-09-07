package handwriting

import (
    "bytes"
    "encoding/binary"
    "github.com/google/uuid"
    "github.com/trymoose/handwriting2svg/pkg/handwriting/internal/handwriting"
    "google.golang.org/protobuf/proto"
    "time"
)

//go:generate protoc --go_out=./internal/handwriting --go_opt=paths=source_relative ./handwriting.proto

type Handwriting struct {
    CreatedAt time.Time       `json:"creation_date"`
    ID        uuid.UUID       `json:"identifier"`
    Size      Size            `json:"size"`
    Frame     Rectangle       `json:"frame"`
    Version   int             `json:"version"` // ???
    Strokes   [][]StrokePoint `json:"strokes"`
}

type StrokePoint struct {
    Point    Point
    Velocity Int16 `json:"velocity"` // ???
    Pressure Int16 `json:"pressure"` // ???
}

type Size struct {
    Width  Int16 `json:"width"`
    Height Int16 `json:"height"`
}

type Point struct {
    X Int16 `json:"x"`
    Y Int16 `json:"y"`
}

type Rectangle struct {
    Origin Point `json:"origin"`
    Size   Size  `json:"size"`
}

func Unmarshal(b []byte) (*Handwriting, error) {
    var pr handwriting.BaseMessage
    if err := proto.Unmarshal(b, &pr); err != nil {
        return nil, err
    }

    identifier, err := uuid.Parse(pr.ID)
    if err != nil {
        return nil, err
    }

    var size Size
    if err := binary.Read(bytes.NewReader(pr.Handwriting.Size), binary.LittleEndian, &size); err != nil {
        return nil, err
    }

    var frame Rectangle
    if err := binary.Read(bytes.NewReader(pr.Handwriting.Frame), binary.LittleEndian, &frame); err != nil {
        return nil, err
    }

    strokes, err := readStrokes(&pr)
    if err != nil {
        return nil, err
    }

    return &Handwriting{
        CreatedAt: time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC).Add(time.Millisecond * time.Duration(pr.CreatedAt)),
        ID:        identifier,
        Version:   int(pr.Handwriting.Version),
        Strokes:   strokes,
        Size:      size,
        Frame:     frame,
    }, nil
}
