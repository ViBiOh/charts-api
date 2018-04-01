package conservatories

type conservatory struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	Department uint8   `json:"department"`
	Zip        string  `json:"zip"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
}
