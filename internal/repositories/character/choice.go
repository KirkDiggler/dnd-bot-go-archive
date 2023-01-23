package character

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeMultiple         OptionType = "multiple"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type Choice struct {
	Name     string    `json:"name"`
	Key      string    `json:"key"`
	Selected bool      `json:"selected"`
	Count    int       `json:"count"`
	Options  []*Option `json:"options"`
}
type Option struct {
	CountedReference *CountedReferenceOption `json:"counted_reference"`
	Reference        *ReferenceOption        `json:"reference"`
	Choice           *Choice                 `json:"choice"`
	Multiple         *MultipleOption         `json:"multiple"`
}

type ReferenceItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CountedReferenceOption struct {
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

type ReferenceOption struct {
	Reference *ReferenceItem `json:"reference"`
}

type MultipleOption struct {
	Selected bool      `json:"selected"`
	Items    []*Option `json:"items"`
}
