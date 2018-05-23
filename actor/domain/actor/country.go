package actor

type (
	Country struct {
		Name          string   `json:"name"`
		Capital       string   `json:"capital"`
		ISO3166Alpha2 string   `json:"ISO_3166_2"`
		ISO3166Alpha3 string   `json:"ISO_3166_3"`
		CallingCode   string   `json:"calling_code"`
		Currency      Currency `json:"currency"`
		Image         Image    `json:"image"`
	}

	Currency struct {
		Name    string `json:"name"`
		ISO4217 string `json:"ISO_4217"`
	}

	Image struct {
		Flat  Size `json:"flat"`
		Shiny Size `json:"shiny"`
	}

	Size struct {
		Sixteen    string `json:"16"`
		TwentyFour string `json:"24"`
		ThirtyTwo  string `json:"32"`
		FortyEight string `json:"48"`
		SixtyFour  string `json:"64"`
	}
)
