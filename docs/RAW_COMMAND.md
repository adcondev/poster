# Raw Command Documentation

## ⚠️ Important Limitations

1. **Write-Only Operation**: Raw commands are unidirectional. Commands that expect responses
   (like status queries) will NOT return data through this interface.

2. **No State Tracking**: Raw commands bypass all state management. The library cannot track
   changes made by raw commands. It is up to the printer and you to ensure correct state.

3. **Buffer Responsibility**: You are responsible for any data left in read buffers by
   query commands.

4. **Encoding Management**: Raw commands do not handle text encoding. Ensure any text data is
   correctly encoded before sending.

## Security Modes

### Standard Mode (safe_mode: false)

- Commands execute without validation
- Full responsibility on the developer
- Suitable for production with tested commands

### Safe Mode (safe_mode: true)

- Known dangerous commands are BLOCKED
- Execution stops with error
- Recommended for development/testing

## Common Escape Sequences

### Cash Drawer

```json
{
  "type": "raw",
  "data": {
    "hex": "1B 70 00 32 64",
    "comment": "Open drawer pin 2, 100ms pulse"
  }
}
```

### Beeper/Buzzer

```json
{
  "type": "raw",
  "data": {
    "hex": "07",
    "comment": "Single beep"
  }
}
```

### Chinese Printer Compatibility

```json
{
  "type": "raw",
  "data": {
    "hex": "1B 21 00",
    "comment": "Reset text style - workaround for GP-5890X"
  }
}
```

## Blocked Commands in Safe Mode

| Command | Hex        | Risk Level | Description                       |
|---------|------------|------------|-----------------------------------|
| ESC @   | `1B 40`    | HIGH       | Full reset - clears ALL settings  |
| ESC = n | `1B 3D 00` | HIGH       | Disable printer                   |
| DLE ENQ | `10 05`    | MEDIUM     | Status query - leaves unread data |
| ESC p   | `1B 70`    | LOW        | Cash drawer - physical activation |

## Builder Usage Examples

```go
// Safe development usage
builder.AddRawSafe("1B 40", "Reset attempt") // Will be BLOCKED

// Production usage  
builder.AddRaw("1B 70 00 32 64", "Open drawer") // Executes

// Convenience methods (bypass safety)
builder.AddPulse() // Opens drawer
builder.AddBeep(3) // Three beeps
```
