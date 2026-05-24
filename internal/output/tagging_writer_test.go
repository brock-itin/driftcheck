package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftcheck/internal/drift"
)

func makeTagMap() map[string]drift.TagSet {
	return map[string]drift.TagSet{
		"web::image::image":  {"critical", "image-related"},
		"api::env::API_KEY":  {"secret"},
	}
}

func TestWriteTagged_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteTagged(&buf, nil)
	if !strings.Contains(buf.String(), "No tagged findings") {
		t.Errorf("expected no-findings message, got: %s", buf.String())
	}
}

func TestWriteTagged_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	WriteTagged(&buf, makeTagMap())
	out := buf.String()
	if !strings.Contains(out, "FINDING") {
		t.Error("expected FINDING header")
	}
	if !strings.Contains(out, "TAGS") {
		t.Error("expected TAGS header")
	}
	if !strings.Contains(out, "critical") {
		t.Error("expected tag 'critical' in output")
	}
	if !strings.Contains(out, "secret") {
		t.Error("expected tag 'secret' in output")
	}
}

func TestWriteTagged_CountLine(t *testing.T) {
	var buf bytes.Buffer
	WriteTagged(&buf, makeTagMap())
	if !strings.Contains(buf.String(), "2 finding(s) tagged") {
		t.Errorf("expected count line, got: %s", buf.String())
	}
}

func TestTaggingStatusLine_NoTags(t *testing.T) {
	line := TaggingStatusLine(nil)
	if !strings.Contains(line, "no rules matched") {
		t.Errorf("unexpected: %s", line)
	}
}

func TestTaggingStatusLine_WithTags(t *testing.T) {
	line := TaggingStatusLine(makeTagMap())
	if !strings.Contains(line, "2 finding(s)") {
		t.Errorf("expected finding count in line: %s", line)
	}
	if !strings.Contains(line, "unique tag(s)") {
		t.Errorf("expected unique tag count in line: %s", line)
	}
}

func TestTruncateTag_Short(t *testing.T) {
	if truncateTag("hello", 10) != "hello" {
		t.Error("short string should not be truncated")
	}
}

func TestTruncateTag_Long(t *testing.T) {
	s := strings.Repeat("x", 60)
	result := truncateTag(s, 50)
	if len([]rune(result)) > 50 {
		t.Errorf("truncated string too long: %d", len(result))
	}
}
