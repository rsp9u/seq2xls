package seqdiag

import (
	"testing"

	"github.com/rsp9u/seq2xls/model"
)

const testDataLifeline = `
seqdiag {
  webserver; database; browser;
  browser  -> webserver [label = "GET /index.html"];
  browser <-- webserver;
  browser  -> webserver [label = "POST /blog/comment"];
              webserver  -> database [label = "INSERT comment"];
              webserver <-- database;
  browser <-- webserver;
}
`

func TestExtractLifelines(t *testing.T) {
	d := ParseSeqdiag([]byte(testDataLifeline))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(lls) > 3 {
		t.Fatalf("Too many lifelines %d", len(lls))
	}

	for _, ll := range lls {
		switch ll.Name {
		case "browser":
			if ll.Index != 2 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		case "webserver":
			if ll.Index != 0 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		case "database":
			if ll.Index != 1 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		default:
			t.Fatalf("Unknown lifeline %s", ll.Name)
		}
	}
}

const testDataMessage = `
seqdiag {
  foo -> bar;
  foo -> bar --> baz;
  foo => bar => baz;
  loop {
    foo  -> bar;
    foo <-- bar;
  }
  foo -> bar {
    baz => qux => quux;
  }
  foo -> foo;
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
	msgs, err := ExtractMessages(d, lls)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(msgs) != 15 {
		t.Fatalf("Too many or few messages %d", len(msgs))
	}

	checkMessage(t, msgs[0], 0, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[1], 1, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[2], 2, "bar", "baz", model.Asynchronous)
	checkMessage(t, msgs[3], 3, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[4], 4, "bar", "baz", model.Synchronous)
	checkMessage(t, msgs[5], 5, "baz", "bar", model.Synchronous)
	checkMessage(t, msgs[6], 6, "bar", "foo", model.Synchronous)
	checkMessage(t, msgs[7], 7, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[8], 8, "bar", "foo", model.Reply)
	checkMessage(t, msgs[9], 9, "foo", "bar", model.Synchronous)
	checkMessage(t, msgs[10], 10, "baz", "qux", model.Synchronous)
	checkMessage(t, msgs[11], 11, "qux", "quux", model.Synchronous)
	checkMessage(t, msgs[12], 12, "quux", "qux", model.Synchronous)
	checkMessage(t, msgs[13], 13, "qux", "baz", model.Synchronous)
	checkMessage(t, msgs[14], 14, "foo", "foo", model.SelfReference)
}
