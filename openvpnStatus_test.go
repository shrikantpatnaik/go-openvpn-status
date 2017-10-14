package openvpnStatus

import (
	"testing"
)

func checkEmpty(status *Status) bool {
	if len(status.ClientList) > 0 || len(status.RoutingTable) > 0 {
		return false
	}
	return true
}

func TestNonExistentFile(t *testing.T) {
	_, e := ParseFile("examples/nonExistentFile")
	if e == nil {
		t.Errorf("Parsing Non Existent File failed")
	}
}

func TestEmptyFile(t *testing.T) {
	s, e := ParseFile("examples/emptyFile.status")
	if s.IsUp == true || e.Error() != "Status File is empty" {
		t.Errorf("Incorrect error while parsing empty file")
	}
}

func TestBadFile(t *testing.T) {
	s, e := ParseFile("examples/badFile.status")
	if s.IsUp == true || e.Error() != "Unable to Parse Status file" {
		t.Errorf("Incorrect error while parsing bad file")
	}
}

func TestEmptyServer(t *testing.T) {
	s, e := ParseFile("examples/emptyServer.status")
	if e != nil {
		if !checkEmpty(s) {
			t.Errorf("Parsing Status file with no clients failed")
		}
	}
}

func TestFileWith2Clients(t *testing.T) {
	s, _ := ParseFile("examples/server.status")
	if len(s.ClientList) != 2 && len(s.RoutingTable) != 2 {
		t.Errorf("Parsing Status file with 2 clients failed")
	}
}
