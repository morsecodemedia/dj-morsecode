package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func EnrichmentCacheKey(
	artist string,
	title string,
) string {

	artist = normalizeEnrichmentKeyPart(
		artist,
	)

	title = normalizeEnrichmentKeyPart(
		title,
	)

	value := artist +
		"\x00" +
		title

	sum := sha256.Sum256(
		[]byte(value),
	)

	return hex.EncodeToString(
		sum[:],
	)

}

func normalizeEnrichmentKeyPart(
	value string,
) string {

	return strings.ToLower(
		strings.TrimSpace(value),
	)

}
