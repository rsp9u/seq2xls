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
