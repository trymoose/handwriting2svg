package svg

import (
    "encoding/json"
    svg "github.com/ajstarks/svgo"
    "github.com/trymoose/handwriting2svg/pkg/handwriting"
    "io"
)

type SVG struct {
    canvas              *svg.SVG
    boxWidth, boxHeight int
    maxX, maxY          int
}

func New(w io.Writer) *SVG {
    g := SVG{}
    g.canvas = svg.New(w)
    return &g
}

func (g *SVG) Close() error {
    g.canvas.End()
    return nil
}

func (g *SVG) Start(hw *handwriting.Handwriting) error {
    g.boxWidth = hw.Frame.Size.Width.Int()
    g.boxHeight = hw.Frame.Size.Height.Int()
    for _, stroke := range hw.Strokes {
        for _, point := range stroke {
            g.maxX = max(g.maxX, point.Point.X.Int())
            g.maxY = max(g.maxY, point.Point.Y.Int())
        }
    }

    g.canvas.Start(g.boxWidth+5, g.boxHeight+5)
    g.canvas.Title(hw.ID.String())

    b, err := json.Marshal(hw)
    if err != nil {
        return err
    }
    g.canvas.Desc(string(b))
    return nil
}

func (g *SVG) WriteStrokes(strokes [][]handwriting.StrokePoint) {
    for _, stroke := range strokes {
        g.WriteStroke(stroke)
    }
}

func (g *SVG) WriteStroke(stroke []handwriting.StrokePoint) {
    x, y := make([]int, len(stroke)), make([]int, len(stroke))
    var prevX, prevY int
    for i, point := range stroke {
        // Scale to bounding box
        xs, ys := point.Point.X.Int(), point.Point.Y.Int()
        xs, ys = (xs*g.boxWidth)/g.maxX, (ys*g.boxHeight)/g.maxY

        // Dots don't move, simulate movement for line to be drawn
        if i == 0 {
            prevX, prevY = xs, ys
        } else if prevX == xs && prevY == ys {
            xs, ys = xs+1, ys+1
        } else {
            prevX, prevY = xs, ys
        }
        x[i], y[i] = xs, ys
    }
    g.canvas.Polyline(x, y, `fill="none"`, `stroke="black"`, `stroke-width="5"`, `stroke-linecap="round"`, `stroke-linejoin="round"`)
}
