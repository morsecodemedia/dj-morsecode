package metadata

import "strings"

type Source string

const (
	SourceGeneric       Source = ""
	SourceIHeart        Source = "iheart"
	SourceLautFM        Source = "laut.fm"
	SourceWXPN          Source = "wxpn"
	SourceAfterHoursFM  Source = "afterhoursfm"
	SourceRadioParadise Source = "radioparadise"
	SourceSomaFM        Source = "somafm"
)

func DetectSource(
	fields map[string]string,
) Source {

	url := strings.ToLower(
		fields["icy-url"],
	)

	name := strings.ToLower(
		fields["icy-name"],
	)

	switch {

	case strings.Contains(
		url,
		"iheartradio.com",
	):
		return SourceIHeart

	case strings.Contains(
		url,
		"laut.fm",
	):
		return SourceLautFM

	case strings.Contains(
		url,
		"xpn.org",
	):
		return SourceWXPN

	case strings.Contains(
		url,
		"ah.fm",
	):
		return SourceAfterHoursFM

	case strings.Contains(
		url,
		"radioparadise.com",
	):
		return SourceRadioParadise

	case strings.Contains(
		url,
		"somafm.com",
	):
		return SourceSomaFM

	case strings.Contains(
		name,
		"somafm",
	):
		return SourceSomaFM

	default:
		return SourceGeneric

	}

}
