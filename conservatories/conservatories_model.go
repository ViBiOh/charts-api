package conservatories

type conservatory struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	Department int     `json:"department"`
	Zip        string  `json:"zip"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lng"`
}

func findConservatories(page, pageSize int64, sortKey string, ascending bool, query string) (int64, []*conservatory, error) {
	count, err := countConservatories(query)
	if err != nil {
		return 0, nil, err
	}

	list, err := searchConservatories(page, pageSize, sortKey, ascending, query)

	return count, list, err

}
