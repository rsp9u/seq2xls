package seqdiag

import (
	"testing"

	"github.com/rsp9u/seq2xls/model"
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
	d := ParseSeqdiag([]byte(testDataMessage))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}
	msgs, _, err := ExtractMessages(d, lls)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(msgs) != 9 {
		t.Fatalf("Too many or few messages %d", len(msgs))
	}

	checkMessage(t, msgs[0], 0, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[1], 1, "bar", "foo", model.Reply)
	checkMessage(t, msgs[2], 2, "foo", "bar", model.Asynchronous)
	checkMessage(t, msgs[3], 3, "bar", "baz", model.Asynchronous)
	checkMessage(t, msgs[4], 4, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[5], 5, "bar", "foo", model.Reply)
	checkMessage(t, msgs[6], 6, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[7], 7, "baz", "qux", model.Synchronous)
	checkMessage(t, msgs[8], 8, "foo", "foo", model.SelfReference)
}

func TestExtractMessagesTrip(t *testing.T) {
	d := ParseSeqdiag([]byte(testDataMessageTrip))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}
	msgs, _, err := ExtractMessages(d, lls)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(msgs) != 14 {
		t.Fatalf("Too many or few messages %d", len(msgs))
	}

	checkMessage(t, msgs[0], 0, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[1], 1, "bar", "foo", model.Reply)
	checkMessage(t, msgs[2], 2, "bar", "baz", model.Synchronous)
	checkMessage(t, msgs[3], 3, "baz", "qux", model.Asynchronous)
	checkMessage(t, msgs[4], 4, "qux", "quux", model.Synchronous)
	checkMessage(t, msgs[5], 5, "quux", "qux", model.Reply)
	checkMessage(t, msgs[6], 6, "baz", "bar", model.Reply)
	checkMessage(t, msgs[7], 7, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[8], 8, "bar", "baz", model.Synchronous)
	checkMessage(t, msgs[9], 9, "baz", "qux", model.Synchronous)
	checkMessage(t, msgs[10], 10, "qux", "baz", model.Reply)
	checkMessage(t, msgs[11], 11, "baz", "bar", model.Reply)
	checkMessage(t, msgs[12], 12, "bar", "foo", model.Reply)
	checkMessage(t, msgs[13], 13, "foo", "foo", model.SelfReference)
}

func checkMessageLabel(t *testing.T, msg *model.Message, label string) {
	if msg.Text != label {
		t.Fatalf("Mismatches label of message [expect: %s, actual: %s]", label, msg.Text)
	}
}

func TestExtractMessagesLabel(t *testing.T) {
	d := ParseSeqdiag([]byte(testDataMessageLabel))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}
	msgs, _, err := ExtractMessages(d, lls)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	checkMessageLabel(t, msgs[0], "GET /options")
	checkMessageLabel(t, msgs[1], "option list")
	checkMessageLabel(t, msgs[2], "pass-through")
	checkMessageLabel(t, msgs[3], "pass-through")
	checkMessageLabel(t, msgs[4], "")
	checkMessageLabel(t, msgs[5], "")
	checkMessageLabel(t, msgs[6], "POST /objects")
	checkMessageLabel(t, msgs[7], "INSERT INTO objects")
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
	d := ParseSeqdiag([]byte(testDataMessageNote))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}
	_, notes, err := ExtractMessages(d, lls)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	checkNote(t, notes[0], 0, false, "Note")
	checkNote(t, notes[1], 1, true, "LeftNote")
	checkNote(t, notes[2], 2, false, "Note on Trip Message")
	checkNote(t, notes[3], 4, true, "Each side notes: Left")
	checkNote(t, notes[4], 4, false, "Each side notes: Right")
	checkNote(t, notes[5], 5, false, "Note on Chained Messages")
	checkNote(t, notes[6], 6, false, "Note on Chained Messages")
}
