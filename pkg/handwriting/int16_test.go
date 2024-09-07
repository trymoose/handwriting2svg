package handwriting

import "testing"

func TestInt16_UnmarshalJSON(t *testing.T) {
    v := Int16(-32292)
    b, err := v.MarshalJSON()
    if err != nil {
        t.Fatal(err)
    }
    var vv Int16
    if err := vv.UnmarshalJSON(b); err != nil {
        t.Fatal(err)
    }

    if vv != v {
        t.Fatalf("want %+v got %+v", v, vv)
    }
}
