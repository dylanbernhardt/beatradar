package config

var GenreMap = map[string]string{
	"140 / DEEP DUBSTEP / GRIME":           "95",
	"AFRO HOUSE":                           "89",
	"AMAPIANO":                             "98",
	"BASS / CLUB":                          "85",
	"BASS HOUSE":                           "91",
	"BREAKS / BREAKBEAT / UK BASS":         "9",
	"DANCE / ELECTRO POP":                  "39",
	"DEEP HOUSE":                           "12",
	"DJ TOOLS":                             "16",
	"DRUM & BASS":                          "1",
	"DUBSTEP":                              "18",
	"ELECTRO (CLASSIC / DETROIT / MODERN)": "94",
	"ELECTRONICA":                          "3",
	"FUNKY HOUSE":                          "81",
	"HARD DANCE / HARDCORE":                "8",
	"HARD TECHNO":                          "2",
	"HOUSE":                                "5",
	"INDIE DANCE":                          "37",
	"JACKIN HOUSE":                         "97",
	"MAINSTAGE":                            "96",
	"MELODIC HOUSE & TECHNO":               "90",
	"MINIMAL / DEEP TECH":                  "14",
	"NU DISCO / DISCO":                     "50",
	"ORGANIC HOUSE / DOWNTEMPO":            "93",
	"PROGRESSIVE HOUSE":                    "15",
	"PSY-TRANCE":                           "13",
	"TECH HOUSE":                           "11",
	"TECHNO (PEAK TIME / DRIVING)":         "6",
	"TECHNO (RAW / DEEP / HYPNOTIC)":       "92",
	"TRANCE (MAIN FLOOR)":                  "7",
	"TRANCE (RAW / DEEP / HYPNOTIC)":       "99",
	"TRAP / WAVE":                          "38",
	"UK GARAGE / BASSLINE":                 "86",
}

func GetGenreID(name string) (string, bool) {
	id, ok := GenreMap[name]
	return id, ok
}

func GetGenreName(id string) string {
	for name, genreID := range GenreMap {
		if genreID == id {
			return name
		}
	}
	return "Unknown Genre"
}
