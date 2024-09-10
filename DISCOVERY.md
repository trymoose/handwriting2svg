# Reading Handwritten Message `payload_data` from the iMessage database

## Analyzing the Raw Data

I started by generating a variety of handwritten messages in Messages. Some lines, dots, squiggles, some written text, and a few preloaded examples from iOS. After making a backup of the database, I dumped all these messages into individual files labeled for easy reference.

To analyze the data, I wrote a program that converted the messages into `.xlsx` files, where each byte of the message was displayed in its own cell. I also set it up to highlight matching bytes with colors across columns. From this, I made a few key observations:

- Every message started with the byte `0x11`.

    - The next eight bytes varied based on whether the message was a built-in one or a custom message.

    - Built-in messages ended with various amounts of `0xFF`, while user-created messages ended with `0x00 0x00 0x00.`

- UUIDs were consistently prefixed with `0x1A 0x24`.

    - Interestingly, `0x24` equals `$`, but it also corresponds to the length of the UUID string.

    - If a handwritten message was resent, it retained the same UUID.

- Certain byte offsets were identical across different messages.

- The XZ compression magic number showed up in some messages but not all, always at the same location.

## Protobuf Structure

At this point, I hypothesized the data was using a tag-length-value (TLV) structure. After some digging, I found that Apple often uses protobufs, so I tried decoding it with the `protoscope` utility. Running this gave me the complete structure of the data:

```
protoscope -explicit-wire-types handwriting.bin
```

```
2:I64 nnnnnnnnnnnni64
3:LEN {"########-####-####-####-############"}
4:LEN {
  2:LEN {`bbbbbbbb`}
  3:LEN {`bbbbbbbbbbbbbbbb`}
  4:VARINT n
  5:VARINT n
  6:VARINT n
  7:VARINT n
  8:LEN {
    `<raw/compressed payload>`
  }
}
```

However, protobufs don't store the field names or types directly in the message, and each structure had three byte arrays that needed decoding.

## Decoding Field Names

The first part of the message was relatively straightforward:

- __Field 2__: Interpreting this as an Apple epoch timestamp produced values that matched the times the messages were sent.
- __Field 3__: This was clearly a unique identifier for each handwritten message.

From there, I dove into the handwritten message struct, starting with the `VARINT` fields:

- __Field 4__: This value varied between messages. By looking at the images, I figured out it represented the number of lines in the message.
- __Field 5__: This value differed between compressed and uncompressed messages. Since protobuf stores enums as numbers, it seemed to fit, identifying the message type (compressed or not).
- __Field 6__: This field only appeared in compressed messages and matched the decompressed length of the payload.
- __Field 7__: This was always present and always had the value 4. I didn’t know what it did, so I called it version for now.

## Skeleton Protobuf Schema

Based on that, I wrote up a preliminary protobuf schema:

```
syntax = 'proto3';
package handwriting;

message BaseMessage {
  sfixed64 CreatedAt = 2;  // Milliseconds since 2001-01-01 UTC.
  string ID = 3;
  Handwriting Handwriting = 4;
}

enum Compression {
  Unknown = 0;
  None = 1;
  XZ = 4;
}

message Handwriting {
  bytes unknown_size_4 = 2;
  bytes unknown_size_8 = 3;
  int64 StrokesCount = 4;  // Number of arrays in Strokes.
  Compression Compression = 5;
  optional int64 DecompressedLength = 6;
  int64 always_4 = 7;
  bytes unknown_strokes_data = 8;
}
```

## Decoding the Raw Bytes

Next, I focused on decoding the three `bytes` fields. To decipher the `unknown_strokes_data` field, I returned to the spreadsheet and identified a repeating pattern: `0x?? 0x?? 0x?? 0x?? 0x?? 0x80 0xFF 0x7F`. This pattern occasionally diverged by two bytes. At the start of the data, there were always two bytes that varied between samples. By analyzing the single dot examples, I realized these initial bytes represented the count of 8-byte blocks. After this, another two-byte count would follow, representing the number of blocks corresponding to the strokes. This structure provided a 2D array of strokes, with each block presumably representing points on those strokes.

Decoding the 8-byte blocks was more challenging. I experimented with different formats (e.g., two 32-bit integers, floats, 16-bit numbers, and 8-bit numbers), but none made sense immediately. I set this aside and moved on to other fields.

## Brute-Forcing Byte Interpretation

In the `unknown_size_4` field, I consistently noticed `0x80` appearing at odd-numbered indexes. Since I was already working with `uint16` for the stroke points, I eventually XORed the values with `0x8000` while playing with the high bit. This process yielded values around ~250 x ~1200, which appeared to represent the size of the handwriting "ribbon" in iMessage. While I couldn’t test across identical devices, the size values differed between messages created on different devices, suggesting the dimensions were influenced by the originating device.

## Adjusting Numbers When Initial Results Don't Look Right

When applying the 16-bit XOR to `unknown_size_8`, the resulting values gave me widths and heights that seemed to match the image dimensions (e.g., wide images produced wide width values). However, the first two numbers were sometimes inexplicably large. After converting the uint16 values to int16, the numbers appeared more reasonable and aligned with where the image started on the canvas.

> After decoding the stroke points and mapping them to the width and height from the  `unknown_size_8` field, I found that these values were sufficient to properly size the output image. By mapping the points to match the width and height from `unknown_size_8` and then setting the output SVG to the width and height from `unknown_size_4`, the resulting image resembled what is sent when handwritten messages are transmitted via SMS.

## Reconstructing the Drawing

After decoding the stroke points using the same technique, I got x and y values that made sense for drawing the line. However, the values were much higher than the frame size, and there were two additional numbers. One seemed to relate to drawing speed, while the other was always `-1`. Since iMessage allows users to replay the message drawing, I assumed one of these values related to drawing velocity. The `-1` might represent pressure, though I couldn't confirm that since my device doesn't support pressure sensitivity.

> Interestingly, when I ignored the frame and size data and simply plotted the points using their max x and y values as the image height and width, my image viewer struggled with the rendering.

In the end, after scaling the points to fit the size of the frame rectangle, I was able to recreate the handwritten message as an image. For exporting the data, I chose the SVG format, as it supports smooth line drawing and worked well with the decoded point data.

## Final Protobuf

```
syntax = 'proto3';
package handwriting;

message BaseMessage {
  sfixed64 CreatedAt = 2;  // Milliseconds since 2001-01-01 UTC.
  string ID = 3;
  Handwriting Handwriting = 4;
}

enum Compression {
  Unknown = 0;
  None = 1;
  XZ = 4;
}

message Handwriting {
  bytes Size = 2;
  bytes Frame = 3;
  int64 StrokesCount = 4;  // Number of arrays in Strokes.
  Compression Compression = 5;
  optional int64 DecompressedLength = 6;
  int64 always_4 = 7;
  bytes Strokes = 8;
}
```
