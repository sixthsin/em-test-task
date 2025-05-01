package clients

type PersonAgeData struct {
	Count uint   `json:"count"`
	Name  string `json:"name"`
	Age   uint   `json:"age"`
}

type PersonGenderData struct {
	Count       uint    `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type PersonNationalityData struct {
	Count   uint          `json:"count"`
	Name    string        `json:"name"`
	Country []CountryData `json:"country"`
}

type CountryData struct {
	CointryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
