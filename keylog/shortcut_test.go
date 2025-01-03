package keylog

import "testing"

func TestGetNewShortcutDetector(t *testing.T) {
	sgs := []ShortcutGroup{
		{
			Id: "1", Name: "test1", Shortcuts: []ShortcutCodes{
				{Id: "1", Name: "dummy1", Codes: []uint16{1, 2, 3}, Type: SequentialShortcutType},
				{Id: "2", Name: "dummy2", Codes: []uint16{4, 5}, Type: SequentialShortcutType},
			},
		},
		{
			Id: "2", Name: "test2", Shortcuts: []ShortcutCodes{
				{Id: "3", Name: "dummy3", Codes: []uint16{10, 11}, Type: HoldShortcutType},
			},
		},
	}
	MustGetNewShortcutsDetector(sgs)
}

func TestDuplicateShortcutGroupIds(t *testing.T) {
	sgs := []ShortcutGroup{
		{
			Id: "1", Name: "test1", Shortcuts: []ShortcutCodes{
				{Id: "1", Name: "dummy1", Codes: []uint16{1, 2, 3}, Type: SequentialShortcutType},
				{Id: "2", Name: "dummy2", Codes: []uint16{4, 5}, Type: SequentialShortcutType},
			},
		},
		{
			Id: "1", Name: "test2", Shortcuts: []ShortcutCodes{
				{Id: "3", Name: "dummy3", Codes: []uint16{10, 11}, Type: HoldShortcutType},
			},
		},
	}
	_, err := getShortcutsFromGroups(sgs)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestDuplicateShortcutIds(t *testing.T) {
	sgs := []ShortcutGroup{
		{
			Id: "1", Name: "test1", Shortcuts: []ShortcutCodes{
				{Id: "1", Name: "dummy1", Codes: []uint16{1, 2, 3}, Type: SequentialShortcutType},
				{Id: "2", Name: "dummy2", Codes: []uint16{4, 5}, Type: SequentialShortcutType},
			},
		},
		{
			Id: "2", Name: "test2", Shortcuts: []ShortcutCodes{
				{Id: "1", Name: "dummy3", Codes: []uint16{10, 11}, Type: HoldShortcutType},
			},
		},
	}
	_, err := getShortcutsFromGroups(sgs)
	if err == nil {
		t.Error("Expected error")
	}
}
