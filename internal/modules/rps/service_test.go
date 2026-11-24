package rps

import "testing"

func TestRPSValidation(t *testing.T) {
	s := NewService(nil)
	if e := s.Save(1, 0, 1, RPSInput{}); e == nil {
		t.Fatal("expected required RPS validation")
	}
	if e := s.SaveCPMK(1, 1, 0, CPMKInput{}); e == nil {
		t.Fatal("expected required CPMK validation")
	}
}
