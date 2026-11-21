package notify

import "testing"

// Line breaks in sender or recipients would let a value add SMTP commands or headers such as Bcc.
func TestEmailRejectsLineBreaksInAddresses(t *testing.T) {
	bad := []string{"a@b.cl\r\nBcc: x@evil.tld", "a@b.cl\nRCPT TO:<x@evil.tld>", "not an address"}
	for _, b := range bad {
		if err := checkMailAddresses(b, []string{"ok@b.cl"}); err == nil {
			t.Fatalf("sender %q was accepted", b)
		}
		if err := checkMailAddresses("ok@b.cl", []string{b}); err == nil {
			t.Fatalf("recipient %q was accepted", b)
		}
	}
	if err := checkMailAddresses("Sebastian <s@b.cl>", []string{"c@d.cl"}); err != nil {
		t.Fatalf("a normal sender with a display name was rejected: %v", err)
	}
}
