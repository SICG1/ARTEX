package traffic

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestGetFallaSinCuerpoLegacy(t *testing.T) {
	dir := t.TempDir()
	db, err := sql.Open("sqlite", "file:"+filepath.Join(dir, "idx.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	tr := &Traffic{dir: dir, db: db, closed: make(chan struct{})}
	if err := tr.initIndex(); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES('1',1,'h','GET','/','http://h/',200,'text/plain',0,0,'legacy/x')`); err != nil {
		t.Fatal(err)
	}
	if _, _, err := tr.Get("1"); err == nil {
		t.Fatal("sin archivos legacy el cuerpo no esta disponible: debe devolver error, no texto vacio")
	}
}
