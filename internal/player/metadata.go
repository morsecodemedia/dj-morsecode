package player

func IsStationID(title string) bool {

	switch title {

	case "":
		return true

	case "97.5FM RadioX":
		return true

	}

	return false

}
