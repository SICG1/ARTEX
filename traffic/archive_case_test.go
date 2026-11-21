package traffic

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestExportHostsIncluyeHostEnMayusculas(t *testing.T) {
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
VALUES('1',1,'Example.COM','GET','/','http://Example.COM/',200,'text/plain',0,0,'')`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO exchange_bodies(id,req_head,req_body,req_blob,resp_head,resp_body,resp_blob)
VALUES('1','GET / HTTP/1.1','',NULL,'HTTP 200','',NULL)`); err != nil {
		t.Fatal(err)
	}
	n, err := tr.ExportHosts([]string{"example.com"}, t.TempDir())
	if err != nil || n != 1 {
		t.Fatalf("ExportHosts = %d, %v; se esperaba 1 intercambio", n, err)
	}
}
