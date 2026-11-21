package intercept

import (
	"reflect"
	"testing"
)

// An unset intercept_enabled_tools setting uses the defaults. A malformed one enables no tool,
// as in the base; the warning goes to the log. The defaults are never applied to a broken setting.
func TestEnabledToolNamesDefaultsOnlyWhenUnset(t *testing.T) {
	if got := enabledToolNames("", false); !reflect.DeepEqual(got, defaultEnabledTools) {
		t.Fatalf("unset setting gave %v, want the defaults", got)
	}
	if got := enabledToolNames("{not json", true); len(got) != 0 {
		t.Fatalf("malformed setting gave %v, want no tools as in the base", got)
	}
	if got := enabledToolNames(`["Bash"]`, true); !reflect.DeepEqual(got, []string{"Bash"}) {
		t.Fatalf("valid setting gave %v, want [Bash]", got)
	}
}
