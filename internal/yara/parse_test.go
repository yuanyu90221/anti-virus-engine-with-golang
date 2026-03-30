package yara

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYaraLine_WithMetadata(t *testing.T) {
	name, severity := parseYaraLine(`TestRule [description="...",severity="high",date="2026-01-01"] /some/file`)
	assert.Equal(t, "TestRule", name)
	assert.Equal(t, "high", severity)
}

func TestParseYaraLine_WithoutMetadata(t *testing.T) {
	name, severity := parseYaraLine("TestRule /some/file")
	assert.Equal(t, "TestRule", name)
	assert.Equal(t, "unknown", severity)
}

func TestParseYaraLine_NoSeverityInMeta(t *testing.T) {
	name, severity := parseYaraLine(`TestRule [description="something",date="2026-01-01"] /some/file`)
	assert.Equal(t, "TestRule", name)
	assert.Equal(t, "unknown", severity)
}

func TestParseYaraLine_EmptyLine(t *testing.T) {
	name, severity := parseYaraLine("")
	assert.Equal(t, "", name)
	assert.Equal(t, "unknown", severity)
}

func TestParseSeverityFromMeta_Found(t *testing.T) {
	assert.Equal(t, "high", parseSeverityFromMeta(`description="foo",severity="high",date="2026-01-01"`))
}

func TestParseSeverityFromMeta_NotFound(t *testing.T) {
	assert.Equal(t, "unknown", parseSeverityFromMeta(`description="foo",date="2026-01-01"`))
}

func TestParseSeverityFromMeta_Empty(t *testing.T) {
	assert.Equal(t, "unknown", parseSeverityFromMeta(""))
}

func TestParseSeverityFromMeta_OnlySeverity(t *testing.T) {
	assert.Equal(t, "medium", parseSeverityFromMeta(`severity="medium"`))
}
