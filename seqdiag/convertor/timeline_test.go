package convertor

import (
	"testing"

	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag"
)

const testDataMessage = `
seqdiag {
  foo  -> bar;
  foo <-- bar
  foo --> bar --> baz;
  loop {
    foo  -> bar;
    foo <-- bar;
  }
  foo -> bar {
    baz -> qux;
  }
  foo -> foo;
}
`

const testDataMessageTrip = `
seqdiag {
  foo => bar;
  bar => baz --> qux {
	qux => quux;
  }
  foo => bar {
	bar => baz {
	  baz => qux;
	}
  }
  foo => foo;
}
`

const testDataMessageLabel = `
seqdiag {
  browser  -> web [label = "GET /options"];
  browser <-- web [label = "option list"];
  browser => web => db [label = "pass-through"];
  browser -> web [label = "POST /objects"] {
    web -> db [label = "INSERT INTO objects"];
  }
}
`

const testDataMessageNote = `
seqdiag {
  foo  -> bar [note = "Note"];
  foo <-- bar [leftnote = "LeftNote"];
  foo => bar [note = "Note on Trip Message"];
  foo -> bar [leftnote  = "Each side notes: Left",
              rightnote = "Each side notes: Right"];
  foo -> bar -> baz [note = "Note on Chained Messages"];
}
`

const testDataFragment = `
seqdiag {
  loop {
    foo -> bar;
	alt {
      bar -->> baz;
	}
    foo <-- bar;
  }
}
`

const testDataEmptyFragment = `
seqdiag {
  foo -> bar;
  loop { }
  bar -> baz;
}
`

const testDataSeparator = `
seqdiag {
  ... Sep1 ...
  foo1 -> foo1;
  === Sep2 ===
  foo2 -> foo2;
  === Sep3 ===
  foo3 -> foo3 {
    === Sep4 ===
    foo4 -> foo4;
    === Sep5 ===
  }
  === Sep6 ===
  foo5 -> foo5;
  === Sep7 ===
}
`

func parseDiagram(t *testing.T, testData string) *model.SequenceDiagram {
	d := seqdiag.ParseSeqdiag([]byte(testData))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	seq := &model.SequenceDiagram{Lifelines: lls}
	err = ScanTimeline(d, seq)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	return seq
}

func checkMessage(t *testing.T, msg *model.Message, idx int, from, to string, msgType model.MessageType) {
	if msg.Index != idx {
		t.Fatalf("Mismatches index of message [expect: %d, actual: %d]", idx, msg.Index)
	}
	if msg.From.Name != from {
		t.Fatalf("Mismatches lifeline name [expect: %s, actual: %s]", from, msg.From.Name)
	}
	if msg.To.Name != to {
		t.Fatalf("Mismatches lifeline name [expect: %s, actual: %s]", to, msg.To.Name)
	}
	if msg.Type != msgType {
		t.Fatalf("Mismatches message type [expect: %v, actual: %v]", msgType, msg.Type)
	}
}

func TestExtractMessages(t *testing.T) {
	seq := parseDiagram(t, testDataMessage)

	if len(seq.Messages) != 9 {
		t.Fatalf("Too many or few messages %d", len(seq.Messages))
	}

	checkMessage(t, seq.Messages[0], 0, "foo", "bar", model.Synchronous)
	checkMessage(t, seq.Messages[1], 1, "bar", "foo", model.Reply)
	checkMessage(t, seq.Messages[2], 2, "foo", "bar", model.Asynchronous)
	checkMessage(t, seq.Messages[3], 3, "bar", "baz", model.Asynchronous)
	checkMessage(t, seq.Messages[4], 4, "foo", "bar", model.Synchronous)
	checkMessage(t, seq.Messages[5], 5, "bar", "foo", model.Reply)
	checkMessage(t, seq.Messages[6], 6, "foo", "bar", model.Synchronous)
	checkMessage(t, seq.Messages[7], 7, "baz", "qux", model.Synchronous)
	checkMessage(t, seq.Messages[8], 8, "foo", "foo", model.SelfReference)
}

func TestExtractMessagesTrip(t *testing.T) {
	seq := parseDiagram(t, testDataMessageTrip)

	if len(seq.Messages) != 14 {
		t.Fatalf("Too many or few messages %d", len(seq.Messages))
	}

	checkMessage(t, seq.Messages[0], 0, "foo", "bar", model.Synchronous)
	checkMessage(t, seq.Messages[1], 1, "bar", "foo", model.Reply)
	checkMessage(t, seq.Messages[2], 2, "bar", "baz", model.Synchronous)
	checkMessage(t, seq.Messages[3], 3, "baz", "qux", model.Asynchronous)
	checkMessage(t, seq.Messages[4], 4, "qux", "quux", model.Synchronous)
	checkMessage(t, seq.Messages[5], 5, "quux", "qux", model.Reply)
	checkMessage(t, seq.Messages[6], 6, "baz", "bar", model.Reply)
	checkMessage(t, seq.Messages[7], 7, "foo", "bar", model.Synchronous)
	checkMessage(t, seq.Messages[8], 8, "bar", "baz", model.Synchronous)
	checkMessage(t, seq.Messages[9], 9, "baz", "qux", model.Synchronous)
	checkMessage(t, seq.Messages[10], 10, "qux", "baz", model.Reply)
	checkMessage(t, seq.Messages[11], 11, "baz", "bar", model.Reply)
	checkMessage(t, seq.Messages[12], 12, "bar", "foo", model.Reply)
	checkMessage(t, seq.Messages[13], 13, "foo", "foo", model.SelfReference)
}

func checkMessageLabel(t *testing.T, msg *model.Message, label string) {
	if msg.Text != label {
		t.Fatalf("Mismatches label of message [expect: %s, actual: %s]", label, msg.Text)
	}
}

func TestExtractMessagesLabel(t *testing.T) {
	seq := parseDiagram(t, testDataMessageLabel)

	checkMessageLabel(t, seq.Messages[0], "GET /options")
	checkMessageLabel(t, seq.Messages[1], "option list")
	checkMessageLabel(t, seq.Messages[2], "pass-through")
	checkMessageLabel(t, seq.Messages[3], "pass-through")
	checkMessageLabel(t, seq.Messages[4], "")
	checkMessageLabel(t, seq.Messages[5], "")
	checkMessageLabel(t, seq.Messages[6], "POST /objects")
	checkMessageLabel(t, seq.Messages[7], "INSERT INTO objects")
}

func checkNote(t *testing.T, note *model.Note, idx int, onLeft bool, text string) {
	if note.Assoc.Index != idx {
		t.Fatalf("Mismatches index of message associated note [expect: %d, actual: %d]", idx, note.Assoc.Index)
	}
	if note.OnLeft != onLeft {
		t.Fatalf("Invalid side [expect: onleft=%v, actual: onleft=%v]", onLeft, note.OnLeft)
	}
	if note.Text != text {
		t.Fatalf("Mismatches note text [expect: %s, actual: %s]", text, note.Text)
	}
}

func TestExtractMessagesNote(t *testing.T) {
	seq := parseDiagram(t, testDataMessageNote)

	checkNote(t, seq.Notes[0], 0, false, "Note")
	checkNote(t, seq.Notes[1], 1, true, "LeftNote")
	checkNote(t, seq.Notes[2], 2, false, "Note on Trip Message")
	checkNote(t, seq.Notes[3], 4, true, "Each side notes: Left")
	checkNote(t, seq.Notes[4], 4, false, "Each side notes: Right")
	checkNote(t, seq.Notes[5], 5, false, "Note on Chained Messages")
	checkNote(t, seq.Notes[6], 6, false, "Note on Chained Messages")
}

func checkFragment(t *testing.T, frag *model.Fragment, idx, begin, end int, fragType model.FragmentType) {
	if frag.Index != idx {
		t.Fatalf("Mismatches index of fragment [expect: %d, actual: %d]", idx, frag.Index)
	}
	if frag.Begin.Index != begin {
		t.Fatalf("Mismatches index of begin message of the fragment [expect: %d, actual: %d]", begin, frag.Begin.Index)
	}
	if frag.End.Index != end {
		t.Fatalf("Mismatches index of end message of the fragment [expect: %d, actual: %d]", end, frag.End.Index)
	}
	if frag.Type != fragType {
		t.Fatalf("Mismatches type of the fragment [expect: %v, actual: %v]", fragType, frag.Type)
	}
}

func TestExtractFragments(t *testing.T) {
	seq := parseDiagram(t, testDataFragment)

	checkFragment(t, seq.Fragments[0], 0, 0, 2, model.Loop)
	checkFragment(t, seq.Fragments[1], 1, 1, 1, model.Alt)
}

func TestExtractFragmentsEmpty(t *testing.T) {
	d := seqdiag.ParseSeqdiag([]byte(testDataEmptyFragment))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	seq := &model.SequenceDiagram{Lifelines: lls}
	err = ScanTimeline(d, seq)
	if err == nil {
		t.Fatalf("Expected error does not occure")
	}
}

func checkSeparator(t *testing.T, sep *model.Separator, text string, beforeFrom string) {
	if sep.Text != text {
		t.Fatalf("Mismatches text of separator [expect: %s, actual: %s]", text, sep.Text)
	}
	if (beforeFrom == "nil" && sep.Before != nil) || (beforeFrom != "nil" && sep.Before.From.Name != beforeFrom) {
		t.Fatalf("Mismatches message of before separator [expect: %s, actual: %s]", beforeFrom, sep.Before.From.Name)
	}
}

func TestExtractSeparator(t *testing.T) {
	seq := parseDiagram(t, testDataSeparator)

	checkSeparator(t, seq.Separators[0], "Sep1", "nil")
	checkSeparator(t, seq.Separators[1], "Sep2", "foo1")
	checkSeparator(t, seq.Separators[2], "Sep3", "foo2")
	checkSeparator(t, seq.Separators[3], "Sep4", "foo3")
	checkSeparator(t, seq.Separators[4], "Sep5", "foo4")
	checkSeparator(t, seq.Separators[5], "Sep6", "foo4")
	checkSeparator(t, seq.Separators[6], "Sep7", "foo5")
}
