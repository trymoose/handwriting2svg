//   handwriting2svg [input] [output]
//
//   Usage of handwriting2svg:
//     Converts a handwritten message retrieved by running this sql on the imessage database:
//         SELECT payload_data FROM message WHERE balloon_bundle_id = 'com.apple.Handwriting.HandwritingProvider';
//     Metadata is written to stderr.
//     On non-zero exit data is invalid, and any files created will not be cleaned up.
//
//     input (optional)
//       Filename to read as input. If '-' or not specified stdin is used.
//     output (optional)
//       File to write to. If not specified stdout is used.
//
package main

import (
    "github.com/trymoose/errors"
    "github.com/trymoose/handwriting2svg/pkg/handwriting"
    "github.com/trymoose/handwriting2svg/pkg/svg"
    "io"
    "os"
)

func main() {
    pr, pw := io.Pipe()
    go GenerateSVG(pw)
    WriteOutput(pr)
}

func GenerateSVG(pw io.WriteCloser) {
    defer errors.Do(pw.Close)
    w := svg.New(pw)
    defer errors.Do(w.Close)

    input := GetInputProto()
    errors.Check(w.Start(input))
    w.WriteStrokes(input.Strokes)
}

func GetInputProto() *handwriting.Handwriting {
    return errors.Get(handwriting.Unmarshal(GetInputBytes()))
}

func GetInputBytes() []byte {
    if len(os.Args) < 2 {
        return errors.Get(io.ReadAll(os.Stdin))
    } else if os.Args[1] == "-" {
        return errors.Get(io.ReadAll(os.Stdin))
    }
    return errors.Get(os.ReadFile(os.Args[1]))
}

func WriteOutput(r io.Reader) {
    if len(os.Args) < 3 {
        errors.Get(io.Copy(os.Stdout, r))
    } else {
        f := errors.Get(os.Create(os.Args[2]))
        defer errors.Do(f.Close)
        errors.Get(io.Copy(f, r))
    }
}
