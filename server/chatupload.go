package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
)

const (
	// maxChatUpload caps the combined file data in one attachment upload.
	maxChatUpload = 512 << 20 // 512 MiB
	// Allow multipart headers and boundaries in addition to the file data.
	chatUploadMultipartOverhead = 1 << 20 // 1 MiB
	chatUploadLimitMessage      = "files exceed the 512 MiB upload limit"
)

// safeChatID guards the {id} path segment against traversal — task ids are numeric,
// session ids are alnum/_/- ; anything with "/" or ".." is rejected.
var safeChatID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// chatAttachment is one uploaded file as the frontend + agent see it. Path is relative
// to the chat's working dir (e.g. "uploads/report.txt"), which is the agent's CWD, so
// it can Read/Bash the file directly; Name/Size drive the UI card.
type chatAttachment struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	// Abs is the absolute on-disk path (m.dir is already absolute). For pre-task staging
	// (scope=staging) the frontend uses it to write the prompt into the description; task/session use composeAgentMessage, which builds the path in the backend.
	Abs string `json:"abs,omitempty"`
}

// chatUpload implements method-1 file support: it saves one or more files into a chat's
// working dir under uploads/, so the agent opens them with its existing Read/Bash tools
// and the sent message carries their paths. No LLM-layer change, no multimodal.
//
// POST /api/chat/upload?scope=task|session|staging&id=<id>, multipart field "file"
// (repeatable). Returns {attachments:[{name,path,size,abs}]}. Target dir mirrors the
// agent CWD layout (the run directory is created like the agent creates it):
//
//	scope=task    → <workDir>/tasks/<id>/main/uploads/   (the main agent's run directory)
//	scope=session → <workDir>/sessions/<id>/uploads/
//	scope=staging → <workDir>/drafts/<id>/uploads/   (pre-task staging: the task has no ID yet;
//	                files land here first and the frontend writes their abs paths into the task description)
func (s *Server) chatUpload(w http.ResponseWriter, r *http.Request) {
	var sub string
	taskScoped := false
	switch r.URL.Query().Get("scope") {
	case "task":
		sub = "tasks"
		taskScoped = true
	case "session":
		sub = "sessions"
	case "staging":
		sub = "drafts"
	default:
		writeErr(w, 400, "scope must be task / session / staging")
		return
	}
	id := r.URL.Query().Get("id")
	if !safeChatID.MatchString(id) {
		writeErr(w, 400, "invalid id")
		return
	}
	base := filepath.Join(s.m.dir, sub, id) // the directory uploads/ goes in
	if taskScoped {
		t := s.m.ResolveTask(id)
		if t == nil {
			writeErr(w, 404, "task not found")
			return
		}
		if !s.engine.beginTaskOperation(id) {
			writeErr(w, http.StatusConflict, "Task is being deleted; cannot upload attachments")
			return
		}
		defer s.engine.decInflight(id)
		taskID, _ := strconv.ParseInt(t.ID, 10, 64)
		base = agent.EnsureRunDir(s.m.dir, taskID, 0)
	} else if sub == "sessions" {
		base = agent.EnsureSessionDir(s.m.dir, id)
	}
	dir := filepath.Join(base, "uploads")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		writeErr(w, 500, "failed to create directory: "+publicErr(err))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxChatUpload+chatUploadMultipartOverhead)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			writeErr(w, http.StatusRequestEntityTooLarge, chatUploadLimitMessage)
			return
		}
		writeErr(w, 400, "failed to parse the upload: "+err.Error())
		return
	}
	defer r.MultipartForm.RemoveAll()
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		writeErr(w, 400, "missing upload file (form field file)")
		return
	}
	var totalSize int64
	for _, hdr := range files {
		totalSize += hdr.Size
	}
	if totalSize > maxChatUpload {
		writeErr(w, http.StatusRequestEntityTooLarge, chatUploadLimitMessage)
		return
	}
	for _, hdr := range files {
		if n := filepath.Base(hdr.Filename); len(n) > maxUploadNameBytes || strings.IndexByte(n, 0) >= 0 {
			writeErr(w, 400, "nombre de archivo demasiado largo o no valido")
			return
		}
	}
	out := make([]chatAttachment, 0, len(files))
	for _, hdr := range files {
		name := filepath.Base(hdr.Filename) // strip any path component
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
			continue
		}
		f, dest, err := createUpload(dir, name)
		if err != nil {
			writeErr(w, 500, "Save failed: "+publicErr(err))
			return
		}
		if err := copyUpload(f, hdr); err != nil {
			os.Remove(dest)
			writeErr(w, 500, "Save failed: "+publicErr(err))
			return
		}
		base := filepath.Base(dest)
		att := chatAttachment{Name: base, Path: "uploads/" + base, Size: hdr.Size}
		if sub == "drafts" { // only pre-task staging reads the absolute path (see chatAttachment.Abs)
			att.Abs = dest
		}
		out = append(out, att)
	}
	writeJSON(w, 200, map[string]any{"attachments": out})
}

// maxUploadNameBytes keeps a name plus its "-N" suffix under the 255-byte limit of common file systems.
const maxUploadNameBytes = 200

// maxUniqueUploadTries bounds the "-N" suffix search in createUpload.
const maxUniqueUploadTries = 1000

// createUpload creates a new attachment file in dir: name, or name-1, name-2... when taken.
// Creation is exclusive, so two simultaneous uploads of one name never share a file.
func createUpload(dir, name string) (*os.File, string, error) {
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	for i := 0; i < maxUniqueUploadTries; i++ {
		cand := filepath.Join(dir, name)
		if i > 0 {
			cand = filepath.Join(dir, fmt.Sprintf("%s-%d%s", stem, i, ext))
		}
		f, err := os.OpenFile(cand, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if os.IsExist(err) {
			continue
		}
		if err != nil { // e.g. name too long or NUL byte: retrying cannot help
			return nil, "", err
		}
		return f, cand, nil
	}
	return nil, "", fmt.Errorf("demasiados archivos con el nombre %q", name)
}

// composeAgentMessage appends an attachment manifest to the user's message so the agent
// knows which files were uploaded and where to Read them. baseDir is the agent's working
// dir (its CWD); we emit ABSOLUTE paths (baseDir + relative) so the agent can Read/Bash
// them unambiguously regardless of how it interprets relative paths.
func composeAgentMessage(msg string, atts []chatAttachment, baseDir string) string {
	if len(atts) == 0 {
		return msg
	}
	var b strings.Builder
	b.WriteString(msg)
	b.WriteString("\n\n[User-uploaded attachments] (absolute paths; use Read/Bash to view them):")
	for _, a := range atts {
		abs, ok := attachmentAbs(baseDir, a.Path)
		if !ok {
			log.Printf("chat: attachment path outside uploads dropped: %q", a.Path)
			continue
		}
		fmt.Fprintf(&b, "\n- %s (%s)", abs, humanBytes(a.Size))
	}
	return b.String()
}

// userActivityWithAttachments builds the persisted 'user' activity. With attachments,
// Detail holds JSON {text, attachments} so the transcript renders text + attachment
// cards; Summary stays the plain text (the list payload omits Detail, lazy-loaded).
func userActivityWithAttachments(worker, text string, atts []chatAttachment) db.Activity {
	a := db.Activity{Worker: worker, Kind: "user", Summary: text}
	if len(atts) > 0 {
		blob, _ := json.Marshal(map[string]any{"text": text, "attachments": atts})
		a.Detail = string(blob)
	}
	return a
}

// attachmentAbs maps an attachment path sent by the client to an absolute path. Only
// "uploads/<name>" with a plain file name qualifies: ".." or an absolute path would let
// the agent read files outside the session.
func attachmentAbs(baseDir, rel string) (string, bool) {
	name := filepath.Base(rel)
	if name == "." || name == ".." || rel != "uploads/"+name {
		return "", false
	}
	return filepath.Join(baseDir, "uploads", name), true
}

func humanBytes(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
