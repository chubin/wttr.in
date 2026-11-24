package v1

func unitTemp() map[bool]string {
	return map[bool]string{
		false: "C",
		true:  "F",
	}
}

func localizedRain() map[string]map[bool]string {
	return map[string]map[bool]string{
		"en": {
			false: "mm",
			true:  "in",
		},
		"be": {
			false: "мм",
			true:  "in",
		},
		"ru": {
			false: "мм",
			true:  "in",
		},
		"uk": {
			false: "мм",
			true:  "in",
		},
	}
}

func localizedVis() map[string]map[bool]string {
	return map[string]map[bool]string{
		"en": {
			false: "km",
			true:  "mi",
		},
		"be": {
			false: "км",
			true:  "mi",
		},
		"ru": {
			false: "км",
			true:  "mi",
		},
		"uk": {
			false: "км",
			true:  "mi",
		},
	}
}

func localizedWind() map[string]map[int]string {
	return map[string]map[int]string{
		"en": {
			0: "km/h",
			1: "mph",
			2: "m/s",
		},
		"az": {
			0: "km/saat",
			1: "mph",
			2: "m/s",
		},
		"be": {
			0: "км/г",
			1: "mph",
			2: "м/c",
		},
		"ru": {
			0: "км/ч",
			1: "mph",
			2: "м/c",
		},
		"tr": {
			0: "km/sa",
			1: "mph",
			2: "m/s",
		},
		"uk": {
			0: "км/год",
			1: "mph",
			2: "м/c",
		},
	}
}

func unitWind(unit int, lang string) string {
	translation, ok := localizedWind()[lang]
	if !ok {
		translation = localizedWind()["en"]
	}

	return translation[unit]
}

func unitVis(unit bool, lang string) string {
	translation, ok := localizedVis()[lang]
	if !ok {
		translation = localizedVis()["en"]
	}

	return translation[unit]
}

func unitRain(unit bool, lang string) string {
	translation, ok := localizedRain()[lang]
	if !ok {
		translation = localizedRain()["en"]
	}

	return translation[unit]
}
