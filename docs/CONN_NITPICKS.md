# Connection Package Nitpicks

## Architecture Issues

1.  **Platform-Specific Struct Definitions**: The `WindowsPrintConnector` struct is defined differently in `windows.go` and `windows_stub.go`. This prevents cross-platform unit testing of the connector logic, as the struct fields are not available on non-Windows platforms in the stub version.
2.  **Direct Syscall Usage**: `windows.go` directly calls `syscall` functions and global variables (DLL procs). This makes it impossible to unit test the `WindowsPrintConnector` without a real Windows environment and a real printer.
3.  **Global State**: The use of global `syscall.NewProc` makes the code harder to mock and test in isolation.

## Implementation Gaps

1.  **Read Not Implemented**: The `Read` method in `WindowsPrintConnector` returns an error and is not implemented. While Windows Spooler API is primarily for writing, status reading might be possible via other means or bi-directional mode, but it is currently just an error stub.
2.  **Hardcoded Strings**: String literals like "winspool.drv" and function names are hardcoded.
3.  **Error Handling**: Some error messages could be more descriptive.

## Refactoring Needs

-   Introduce a `PrinterService` interface to abstract the OS API.
-   Unified `WindowsPrintConnector` struct across platforms to allow logic testing on Linux/CI.
