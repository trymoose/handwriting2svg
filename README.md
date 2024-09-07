# handwriting2svg

## Install

```bash
go install github.com/trymoose/handwriting2svg@latest
```

## Usage

```
handwriting2svg [input] [output]

Usage of handwriting2svg:
  Converts a handwritten message retrieved by running this sql on the imessage database:
    SELECT payload_data FROM message WHERE balloon_bundle_id = 'com.apple.Handwriting.HandwritingProvider';
  Metadata is written to stderr.
  On non-zero exit data is invalid, and any files created will not be cleaned up.

  input (optional)
    Filename to read as input. If '-' or not specified stdin is used.
  output (optional)
    File to write to. If not specified stdout is used.
```

## Example Output

![example](https://github.com/user-attachments/assets/7b51f9dc-bfb4-45af-a797-e81249f27410)

## iMessage Output

### Screenshot

![Screenshot](https://github.com/user-attachments/assets/5f77c97f-1d9f-4ebc-921a-35908db0c722)

### SMS

![SMS](https://github.com/user-attachments/assets/80df4673-a610-4b4c-b61b-69c6f2e1d8e5)
