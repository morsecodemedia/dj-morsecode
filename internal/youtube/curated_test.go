func TestCuratedCatalog(
	t *testing.T,
) {

	catalog := NewCatalog(
		Curated,
	)

	items, err := catalog.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) == 0 {
		t.Fatal(
			"expected curated YouTube items",
		)
	}

	for i, item := range items {

		if !item.Valid() {

			t.Fatalf(
				"expected item %d to be valid",
				i,
			)

		}

	}

}