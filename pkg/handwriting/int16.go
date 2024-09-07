package handwriting

import (
    "encoding/json"
    "strconv"
)

type Int16 int16

func (f Int16) Int() int {
    return int(uint16(f) ^ uint16(0x8000))
}

func (f Int16) MarshalJSON() ([]byte, error) {
    return json.Marshal(f.Int())
}

func (f *Int16) UnmarshalJSON(data []byte) error {
    var n int
    if err := json.Unmarshal(data, &n); err != nil {
        return err
    }
    *f = Int16(uint16(int16(n)) ^ uint16(0x8000))
    return nil
}

func (f Int16) String() string {
    return strconv.Itoa(f.Int())
}
