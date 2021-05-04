package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	TFTPserver = "192.168.101.127" //Invader7 for now
)

type getFil struct {
	srv        string
	shouldFail bool
}

var GFtests = []getFil{
	{srv: "http://" + DfltSrv + "/" + DfltVer},
	{srv: "http://" + DfltSrv + "/" + "NONEXISTENT", shouldFail: true},
	//	{TFTPserver, DfltVer},
}

func TestGetFile(t *testing.T) {
	TmpDir = "/tmp"
	for _, pair := range GFtests {
		fmt.Printf("Opening %s\n", pair.srv)
		_, err := getFile(pair.srv, ArchiveName)
		if err != nil {
			if pair.shouldFail {
				fmt.Printf("Got expected error %s\n",
					err)
				continue
			}
			t.Errorf("HTTP: Error downloading: %v", err)
			return
		}
		if pair.shouldFail {
			t.Errorf("Did not get error opening %s", pair.srv)
		}
		fn := filepath.Join(TmpDir, ArchiveName)
		if _, err = os.Stat(fn); os.IsNotExist(err) {
			t.Errorf("HTTP: File did not get created, error: %v", err)
			return
		}
		if err = rmFile(ArchiveName); err != nil {
			t.Errorf("HTTP: File did not get removed, error: %v", err)
			return
		}
	}
}

type isVerNewer struct {
	curver string
	newver string
	result bool
}

var z = []byte{0xff, 0xff, 0xff, 0xff}

var IVNtests = []isVerNewer{
	{"v0.2", "v0.3", true},
	{"v0.3", "v0.2", false},
	{"v0.3", "20170901", true},
	{"20170901", "v0.2", false},
	{"20170901", "20170902", true},
	{"20170901", "20170830", false},
	{string(z), "20170830", true},
	{"20170830", string(z), false},
}

func TestIsVersionNewer(t *testing.T) {
	for _, pair := range IVNtests {
		v, err := isVersionNewer(pair.curver, pair.newver)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if v != pair.result {
			t.Error(
				"For", pair.curver, pair.newver,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}

func TestRmFile(t *testing.T) {
	fn := filepath.Join(TmpDir, "tempfile")
	f, err := os.Create(fn)
	if err != nil {
		t.Errorf("Could not create file, error: %v", err)
		return
	}
	f.Close()
	if _, err = os.Stat(fn); os.IsNotExist(err) {
		t.Errorf("File did not get created, error: %v", err)
		return
	}
	if err = rmFile("tempfile"); err != nil {
		t.Errorf("Error during remove file: %v", err)
		return
	}
	if _, err = os.Stat(fn); !os.IsNotExist(err) {
		t.Errorf("File did not get removed, error: %v", err)
		return
	}
}
