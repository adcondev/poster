# PROFILE Nitpicks

This document contains nitpicks and observations found while analyzing and testing the `pkg/profile` package.

## escpos_encoding.go

1.  **Side effect in `getEncoding`**: The `getEncoding` method logs a warning using `log.Printf` when an unsupported encoding is encountered. Libraries should generally avoid direct logging and instead return errors or allow the user to configure a logger.
2.  **Silent Fallback**: `getEncoding` falls back to `Windows-1252` when the code table is not found. While this prevents a crash, it might lead to incorrect output being printed without the caller knowing (other than the log).
3.  **Encoder Instantiation**: `getEncoding` calls `.NewEncoder()` on every call. Depending on the frequency of calls, this might be slightly inefficient, although `encoding.Encoder` creation is usually cheap.
4.  **Missing Mappings**: As noted in the TODOs, there are many ESC/POS code tables without direct Go encoding mappings.

## escpos_profile.go

1.  **Hardcoded Models**: The factory functions create specific profiles (e.g., `CreatePt210`). As the library grows, this might become hard to maintain. A configuration-based approach or a builder pattern might be more flexible in the future.
