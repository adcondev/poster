package testutils

import (
	"errors"
	"testing"
)

func TestAssertions_SuccessCases(t *testing.T) {
	// AssertBytes
	AssertBytes(t, []byte{0x01, 0x02}, []byte{0x01, 0x02})

	// AssertError
	AssertError(t, nil, nil)
	customErr := errors.New("custom error")
	AssertError(t, customErr, customErr)

	// AssertErrorOccurred
	if !AssertErrorOccurred(t, customErr, true, "testMethod") {
		t.Error("AssertErrorOccurred should return true when error matches expectation")
	}
	if !AssertErrorOccurred(t, nil, false, "testMethod") {
		t.Error("AssertErrorOccurred should return true when no-error matches expectation")
	}

	// AssertContains
	AssertContains(t, []byte("abcdef"), []byte("cde"))

	// AssertNotEmpty
	AssertNotEmpty(t, []byte{0x01}, "test data")

	// AssertLength
	AssertLength(t, []byte{1, 2, 3}, 3)

	// AssertHasPrefix
	AssertHasPrefix(t, []byte("prefix_data"), []byte("prefix"))

	// AssertNumeric
	AssertNumeric(t, []byte("1234567890"))

	// AssertUppercase
	AssertUppercase(t, []byte("ABCXYZ123"))

	// AssertPrintableASCII
	AssertPrintableASCII(t, []byte("Hello, World! 123"))

	// AssertHasNullTerminator
	AssertHasNullTerminator(t, []byte{0x01, 0x02, 0x00})

	// AssertEvenLength
	AssertEvenLength(t, []byte{0x01, 0x02, 0x03, 0x04})

	// AssertContainsOnly
	AssertContainsOnly(t, []byte("abca"), []byte("abc"))

	// AssertValidLength
	AssertValidLength(t, []byte("test"), 2, 6)
}

func TestAssertions_FailureCases(t *testing.T) {
	mockT := new(testing.T)

	// AssertBytes failures
	AssertBytes(mockT, []byte{1}, []byte{2})
	AssertBytes(mockT, []byte{1}, []byte{2}, "custom msg")
	AssertBytes(mockT, []byte{1}, []byte{2}, "formatted %s", "msg")

	// AssertError failures
	AssertError(mockT, errors.New("unexpected"), nil)
	AssertError(mockT, errors.New("err1"), errors.New("err2"))

	// AssertErrorOccurred failure
	_ = AssertErrorOccurred(mockT, errors.New("err"), false, "method")
	_ = AssertErrorOccurred(mockT, nil, true, "method")

	// AssertContains failure
	AssertContains(mockT, []byte("abc"), []byte("xyz"))
	AssertContains(mockT, []byte("abc"), []byte("xyz"), "not found")
	AssertContains(mockT, []byte("abc"), []byte("xyz"), "not found %s", "needle")

	// AssertNotEmpty failure
	AssertNotEmpty(mockT, nil, "empty")

	// AssertLength failure
	AssertLength(mockT, []byte{1}, 3)
	AssertLength(mockT, []byte{1}, 3, "wrong len")
	AssertLength(mockT, []byte{1}, 3, "wrong len %d", 3)

	// AssertHasPrefix failure
	AssertHasPrefix(mockT, []byte("abc"), []byte("xyz"))
	AssertHasPrefix(mockT, []byte("abc"), []byte("xyz"), "no prefix")
	AssertHasPrefix(mockT, []byte("abc"), []byte("xyz"), "no prefix %s", "xyz")

	// AssertNumeric failure
	AssertNumeric(mockT, []byte("123a45"))

	// AssertUppercase failure
	AssertUppercase(mockT, []byte("ABCdEF"))

	// AssertPrintableASCII failure
	AssertPrintableASCII(mockT, []byte{0x01, 0x02})

	// AssertHasNullTerminator failure
	AssertHasNullTerminator(mockT, []byte{1, 2})
	AssertHasNullTerminator(mockT, nil)

	// AssertEvenLength failure
	AssertEvenLength(mockT, []byte{1, 2, 3})

	// AssertContainsOnly failure
	AssertContainsOnly(mockT, []byte("abcd"), []byte("abc"))

	// AssertValidLength failure
	AssertValidLength(mockT, []byte("toolong"), 1, 3)

	if !mockT.Failed() {
		t.Error("expected mockT to be marked failed")
	}
}

func TestFormatMessage(t *testing.T) {
	// Default message with no args
	msg1 := formatMessage("default")
	if msg1 != "default" {
		t.Errorf("formatMessage() = %q, want \"default\"", msg1)
	}

	// Single argument message
	msg2 := formatMessage("default", "custom message")
	if msg2 != "custom message" {
		t.Errorf("formatMessage() = %q, want \"custom message\"", msg2)
	}

	// Formatted message with multiple arguments
	msg3 := formatMessage("default", "code %d: %s", 404, "not found")
	if msg3 != "code 404: not found" {
		t.Errorf("formatMessage() = %q, want \"code 404: not found\"", msg3)
	}
}

func TestByteHelpers(t *testing.T) {
	// RepeatByte
	repeated := RepeatByte(5, 0xAA)
	if len(repeated) != 5 {
		t.Fatalf("RepeatByte length = %d, want 5", len(repeated))
	}
	for i, b := range repeated {
		if b != 0xAA {
			t.Errorf("repeated[%d] = %#x, want 0xAA", i, b)
		}
	}

	// GenerateString
	str := GenerateString(4, 'Z')
	if str != "ZZZZ" {
		t.Errorf("GenerateString = %q, want \"ZZZZ\"", str)
	}

	// BuildCommand
	cmd := BuildCommand(0x1B, '@', 0x01, 0x02)
	wantCmd := []byte{0x1B, '@', 0x01, 0x02}
	AssertBytes(t, cmd, wantCmd)

	// BufferBuilder
	builder := &BufferBuilder{}
	builder.Append([]byte{0x10, 0x20}).Append([]byte{0x30, 0x40})
	AssertBytes(t, builder.GetBuffer(), []byte{0x10, 0x20, 0x30, 0x40})
}

func TestFakeCapabilityBase(t *testing.T) {
	fc := &FakeCapabilityBase{
		Tracker:  &MockTracker{CallCount: make(map[string]int)},
		StateMap: make(map[string]interface{}),
	}

	fc.AppendCommand([]byte{0x1B, '@'}, "Initialize")
	AssertBytes(t, fc.GetBuffer(), []byte{0x1B, '@'})
	if count := fc.Tracker.GetCallCount("Initialize"); count != 1 {
		t.Errorf("GetCallCount = %d, want 1", count)
	}

	fc.SetState("align", "center")
	state := fc.GetState()
	if state["align"] != "center" {
		t.Errorf("GetState[\"align\"] = %v, want \"center\"", state["align"])
	}

	fc.Reset()
	if len(fc.GetBuffer()) != 0 {
		t.Errorf("GetBuffer after Reset = %v, want empty", fc.GetBuffer())
	}
	if len(fc.GetState()) != 0 {
		t.Errorf("GetState after Reset = %v, want empty", fc.GetState())
	}
	if fc.Tracker.GetCallCount("Initialize") != 0 {
		t.Errorf("Tracker call count after Reset = %d, want 0", fc.Tracker.GetCallCount("Initialize"))
	}
}

func TestMockTracker(t *testing.T) {
	mt := &MockTracker{
		CallCount: make(map[string]int),
	}

	mt.Track("Print")
	mt.Track("Print")
	mt.Track("Feed")

	if count := mt.GetCallCount("Print"); count != 2 {
		t.Errorf("GetCallCount(Print) = %d, want 2", count)
	}
	if count := mt.GetCallCount("Feed"); count != 1 {
		t.Errorf("GetCallCount(Feed) = %d, want 1", count)
	}
	if mt.LastCommand != "Feed" {
		t.Errorf("LastCommand = %q, want \"Feed\"", mt.LastCommand)
	}
	if len(mt.CallHistory) != 3 {
		t.Errorf("CallHistory length = %d, want 3", len(mt.CallHistory))
	}

	mt.Reset()
	if mt.GetCallCount("Print") != 0 {
		t.Errorf("GetCallCount after reset = %d, want 0", mt.GetCallCount("Print"))
	}
	if len(mt.CallHistory) != 0 {
		t.Errorf("CallHistory length after reset = %d, want 0", len(mt.CallHistory))
	}
	if mt.LastCommand != "" {
		t.Errorf("LastCommand after reset = %q, want empty", mt.LastCommand)
	}
}

func TestMockConnector(t *testing.T) {
	// Default methods
	mc := &MockConnector{}

	n, err := mc.Write([]byte("hello"))
	if err != nil || n != 5 {
		t.Errorf("Write() = (%d, %v), want (5, nil)", n, err)
	}

	if err := mc.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}

	rn, rerr := mc.Read([]byte{0})
	if rerr != nil || rn != 0 {
		t.Errorf("Read() = (%d, %v), want (0, nil)", rn, rerr)
	}

	// Custom WriteFunc and CloseFunc error injection
	injectedErr := errors.New("injected connection error")
	mc.WriteFunc = func(_ []byte) (int, error) {
		return 0, injectedErr
	}
	mc.CloseFunc = func() error {
		return injectedErr
	}

	_, err = mc.Write([]byte("data"))
	if !errors.Is(err, injectedErr) {
		t.Errorf("Write() error = %v, want injectedErr", err)
	}

	err = mc.Close()
	if !errors.Is(err, injectedErr) {
		t.Errorf("Close() error = %v, want injectedErr", err)
	}
}
