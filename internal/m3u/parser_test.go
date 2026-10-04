package m3u

import (
	"os"
	"strings"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestParseBasicPlaylist(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/minimal.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "testdata/minimal.m3u",
	}

	entries, err := Parse(
		file,
		origin,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 3 {

		t.Fatalf(
			"expected 3 entries, got %d",
			len(entries),
		)

	}

	if entries[0].Ref.URI !=
		"https://example.com/one" {

		t.Errorf(
			"unexpected first URI %q",
			entries[0].Ref.URI,
		)

	}

	if entries[2].Ref.URI !=
		"/Users/example/Music/song.mp3" {

		t.Errorf(
			"unexpected local path %q",
			entries[2].Ref.URI,
		)

	}

	if entries[0].Ref.Source.Kind !=
		source.KindM3U {

		t.Errorf(
			"expected M3U provenance, got %q",
			entries[0].Ref.Source.Kind,
		)

	}

}

func TestParseExtendedPlaylist(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/minimal_extended.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "testdata/minimal_extended.m3u",
	}

	entries, err := Parse(
		file,
		origin,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 3 {

		t.Fatalf(
			"expected 3 entries, got %d",
			len(entries),
		)

	}

	if entries[0].Ref.Name != "Z100" {

		t.Errorf(
			"expected Z100, got %q",
			entries[0].Ref.Name,
		)

	}

	if entries[0].Duration != -1 {

		t.Errorf(
			"expected indefinite duration -1, got %d",
			entries[0].Duration,
		)

	}

	if entries[2].Duration != 252 {

		t.Errorf(
			"expected duration 252, got %d",
			entries[2].Duration,
		)

	}

	if entries[2].Ref.Name !=
		"Stone Temple Pilots - Interstate Love Song" {

		t.Errorf(
			"unexpected track name %q",
			entries[2].Ref.Name,
		)

	}

}

func TestParseClearsPendingMetadata(
	t *testing.T,
) {

	content := strings.NewReader(
		"#EXTM3U\n" +
			"#EXTINF:-1,Z100\n" +
			"https://example.com/z100\n" +
			"https://example.com/plain\n",
	)

	entries, err := Parse(
		content,
		source.Source{
			Kind: source.KindM3U,
			URI:  "memory.m3u",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {

		t.Fatalf(
			"expected 2 entries, got %d",
			len(entries),
		)

	}

	if entries[1].Ref.Name != "" {

		t.Errorf(
			"expected no inherited name, got %q",
			entries[1].Ref.Name,
		)

	}

	if entries[1].Duration != 0 {

		t.Errorf(
			"expected no inherited duration, got %d",
			entries[1].Duration,
		)

	}

}

func TestParseToleratesMalformedEXTINF(
	t *testing.T,
) {

	content := strings.NewReader(
		"#EXTM3U\n" +
			"#EXTINF:potato,Example FM\n" +
			"https://example.com/live\n",
	)

	entries, err := Parse(
		content,
		source.Source{
			Kind: source.KindM3U,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf(
			"expected 1 entry, got %d",
			len(entries),
		)
	}

	if entries[0].Duration != 0 {

		t.Errorf(
			"expected unknown duration, got %d",
			entries[0].Duration,
		)

	}

	if entries[0].Ref.Name !=
		"Example FM" {

		t.Errorf(
			"expected preserved name, got %q",
			entries[0].Ref.Name,
		)

	}

}

func TestParseRealBasicPlaylist(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/basic.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "testdata/basic.m3u",
	}

	entries, err := Parse(
		file,
		origin,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) == 0 {
		t.Fatal(
			"expected playlist entries",
		)
	}

	for i, entry := range entries {

		if !entry.Ref.Valid() {

			t.Fatalf(
				"expected entry %d to have valid provenance and identity",
				i,
			)

		}

	}

}
func TestParseRealExtendedPlaylist(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/extended.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	origin := source.Source{
		Kind: source.KindM3U,
		URI:  "testdata/extended.m3u",
	}

	entries, err := Parse(
		file,
		origin,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) == 0 {
		t.Fatal(
			"expected playlist entries",
		)
	}

	named := 0
	indefinite := 0

	for i, entry := range entries {

		if !entry.Ref.Valid() {

			t.Fatalf(
				"expected entry %d to be valid",
				i,
			)

		}

		if entry.Ref.Name != "" {
			named++
		}

		if entry.Duration == -1 {
			indefinite++
		}

	}

	if named == 0 {
		t.Fatal(
			"expected EXTINF names to be preserved",
		)
	}

	if indefinite == 0 {
		t.Fatal(
			"expected indefinite stream durations",
		)
	}

}
