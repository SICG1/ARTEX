package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// Rendering for the findings page "Export": renders a batch of findings table rows
// into a summary Markdown, a single-finding Markdown, or CSV. JSON is serialized
// directly from DTOs at the server layer, not here.

// sortFindingsForExport sorts by severity descending, then by time descending, consistent with the summary report's grouping.
func sortFindingsForExport(fs []*db.DBFinding) {
	sort.SliceStable(fs, func(i, j int) bool {
		ri, rj := severityRank(fs[i].Severity), severityRank(fs[j].Severity)
		if ri != rj {
			return ri < rj // smaller sevRank means more severe
		}
		return fs[i].CreatedAt.After(fs[j].CreatedAt)
	})
}

// findingTitle returns the vulnerability's readable title: name → class → "Uncategorized".
func findingTitle(f *db.DBFinding) string {
	return nz(f.Name, nz(f.VulnClass, "Uncategorized"))
}

// summaryKey maps a severity to a summary row: the four known names, lower case; anything else is info.
func summaryKey(sev string) string {
	switch k := strings.ToLower(strings.TrimSpace(sev)); k {
	case "critical", "high", "medium", "low":
		return k
	}
	return "info"
}

// FindingsMarkdown combines a batch of findings into a summary report (summary + grouped by severity,
// each entry with class/status/owning task/evidence/detailed report).
func FindingsMarkdown(fs []*db.DBFinding, generatedAt time.Time) string {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var b strings.Builder
	b.WriteString("# Vulnerability Findings Summary Report\n\n")
	fmt.Fprintf(&b, "- **Generated at**: %s\n", generatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **Total findings**: %d\n\n", len(items))

	// Summary: counts per severity.
	counts := map[string]int{}
	for _, f := range items {
		counts[summaryKey(f.Severity)]++
	}
	b.WriteString("## Summary\n\n")
	b.WriteString("| Severity | Count |\n| --- | --- |\n")
	for _, s := range []struct{ key, label string }{
		{"critical", "Critical"}, {"high", "High"}, {"medium", "Medium"}, {"low", "Low"}, {"info", "Info / other"},
	} {
		fmt.Fprintf(&b, "| %s | %d |\n", s.label, counts[s.key])
	}
	b.WriteString("\n")

	if len(items) == 0 {
		b.WriteString("_No matching vulnerabilities._\n")
		return b.String()
	}

	b.WriteString("## Vulnerability Details\n\n")
	for i, f := range items {
		fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
		if f.VulnClass != "" {
			fmt.Fprintf(&b, "- **Class**: %s\n", f.VulnClass)
		}
		fmt.Fprintf(&b, "- **Status**: %s\n", nz(f.Status, "pending"))
		fmt.Fprintf(&b, "- **Report stage**: %s\n", nz(f.ReportState, "none"))
		if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
			fmt.Fprintf(&b, "- **Task**: %s\n", desc)
		}
		fmt.Fprintf(&b, "- **Found at**: %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		q := qualityFromDBFinding(f)
		b.WriteString(qualityBullets(q, false))
		if s := strings.TrimSpace(f.Summary); s != "" {
			fmt.Fprintf(&b, "%s\n\n", s)
		}
		if e := strings.TrimSpace(f.Evidence); e != "" {
			fence := codeFence(e)
			fmt.Fprintf(&b, "**Evidence:**\n\n%s\n%s\n%s\n\n", fence, e, fence)
		}
		b.WriteString(qualityNarrative(q))
		b.WriteString(detailedReportFull(f.Report, "**Detailed report:**"))
		b.WriteString(findingTrafficMarkdown(f, false))
		b.WriteString("---\n\n")
	}
	return b.String()
}

// codeFence returns a backtick fence longer than any run of backticks in text, so the text
// cannot close the code block early. It is never shorter than three backticks.
func codeFence(text string) string {
	longest, run := 0, 0
	for i := 0; i < len(text); i++ {
		if text[i] == '`' {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	return strings.Repeat("`", max(3, longest+1))
}

// SingleFindingMarkdown renders a single vulnerability as a standalone Markdown (used for "one file per vulnerability" packaging).
func SingleFindingMarkdown(f *db.DBFinding, generatedAt time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# [%s] %s\n\n", strings.ToUpper(nz(f.Severity, "info")), findingTitle(f))
	if f.VulnClass != "" {
		fmt.Fprintf(&b, "- **Class**: %s\n", f.VulnClass)
	}
	fmt.Fprintf(&b, "- **Severity**: %s\n", nz(f.Severity, "info"))
	fmt.Fprintf(&b, "- **Status**: %s\n", nz(f.Status, "pending"))
	fmt.Fprintf(&b, "- **Report stage**: %s\n", nz(f.ReportState, "none"))
	if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
		fmt.Fprintf(&b, "- **Task**: %s\n", desc)
	}
	fmt.Fprintf(&b, "- **Found at**: %s\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "- **Generated at**: %s\n", generatedAt.Format("2006-01-02 15:04:05"))
	q := qualityFromDBFinding(f)
	b.WriteString(qualityBullets(q, false))
	if s := strings.TrimSpace(f.Summary); s != "" {
		fmt.Fprintf(&b, "## Summary\n\n%s\n\n", s)
	}
	if e := strings.TrimSpace(f.Evidence); e != "" {
		fence := codeFence(e)
		fmt.Fprintf(&b, "## Evidence\n\n%s\n%s\n%s\n\n", fence, e, fence)
	}
	b.WriteString(qualityNarrative(q))
	b.WriteString(detailedReportFull(f.Report, "## Detailed Report"))
	b.WriteString(findingTrafficMarkdown(f, true))
	return b.String()
}

var unsafeFilenameChars = regexp.MustCompile(`[^\p{Han}\p{L}\p{N}._-]+`)

// FindingFilename generates a safe .md filename for "one file per vulnerability", e.g.
// `critical_SQL_injection_#123.md`. Strips path separators and control characters to avoid illegal paths inside the zip.
func FindingFilename(f *db.DBFinding) string {
	sev := nz(f.Severity, "info")
	title := findingTitle(f)
	name := fmt.Sprintf("%s_%s_#%d", sev, title, f.ID)
	name = unsafeFilenameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = fmt.Sprintf("finding_%d", f.ID)
	}
	// Defensive: strip one more path layer to rule out zip slip.
	name = path.Base(name)
	if r := []rune(name); len(r) > 120 { // cut by runes: a byte cut can split a UTF-8 character
		name = string(r[:120])
	}
	return name + ".md"
}

// FindingsCSV renders a batch of findings as CSV (with a UTF-8 BOM so Excel detects the encoding correctly).
// It omits the long report/evidence full text and keeps only summary-level fields; use Markdown/JSON export for full text.
func FindingsCSV(fs []*db.DBFinding) []byte {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "Name", "Class", "Severity", "Status", "Report Stage", "Task", "Found At", "Summary", "Traffic Evidence Count", "Traffic Evidence IDs"})
	for _, f := range items {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			csvCell(findingTitle(f)),
			csvCell(f.VulnClass),
			csvCell(nz(f.Severity, "info")),
			csvCell(nz(f.Status, "pending")),
			csvCell(nz(f.ReportState, "none")),
			csvCell(f.TaskDescription),
			csvCell(f.CreatedAt.Format("2006-01-02 15:04:05")),
			csvCell(strings.TrimSpace(f.Summary)),
			fmt.Sprint(len(f.TrafficBindings)),
			csvCell(findingTrafficIDs(f)),
		})
	}
	w.Flush()
	return buf.Bytes()
}

// csvCell prefixes cells that spreadsheet apps treat as formulas. encoding/csv
// only quotes delimiters; Excel/LibreOffice still evaluate a leading
// '=', '+', '-', '@', tab, or CR when the file is opened.
// csvSafe is the earlier name of csvCell, kept for the local tests.
func csvSafe(s string) string { return csvCell(s) }

func csvCell(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}

func findingTrafficIDs(f *db.DBFinding) string {
	ids := make([]string, 0, len(f.TrafficBindings))
	for _, b := range f.TrafficBindings {
		ids = append(ids, fmt.Sprint(b.ID))
	}
	return strings.Join(ids, ",")
}

func findingTrafficMarkdown(f *db.DBFinding, attachments bool) string {
	stale := reportStale(f.Report, f.EvidenceVersion, f.ReportEvidenceVersion)
	if len(f.TrafficBindings) == 0 && !stale {
		return ""
	}
	var out strings.Builder
	out.WriteString("\n## Related Traffic Evidence\n\n")
	fmt.Fprintf(&out, "Evidence version: %d; bindings: %d.\n\n", f.EvidenceVersion, len(f.TrafficBindings))
	out.WriteString(staleReportNotice(stale))
	for i, b := range f.TrafficBindings {
		fmt.Fprintf(&out, "%d. **Evidence #%d · %s** — `%s %s`, status %d\n", i+1, b.ID, b.Role, b.Snapshot.Method, strings.ReplaceAll(b.Snapshot.URL, "`", "%60"), b.Snapshot.Status)
		if b.Note != "" {
			fmt.Fprintf(&out, "   %s\n", strings.ReplaceAll(b.Note, "\n", "\n   "))
		}
		if attachments {
			fmt.Fprintf(&out, "   [Request message](evidence/%d/%d/request.http) · [Response message](evidence/%d/%d/response.http)\n", f.ID, b.ID, f.ID, b.ID)
		}
	}
	out.WriteString("\n")
	return out.String()
}
