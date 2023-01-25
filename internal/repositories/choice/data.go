package choice

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeMultiple         OptionType = "multiple"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type Type string

const (
	TypeUnset       Type = ""
	TypeProficiency Type = "proficiency"
	TypeLanguage    Type = "language"
	TypeEquipment   Type = "equipment"
)

type Data struct {
	CharacterID string
	Type        Type
	Choices     []*Choice
}

type Status string

const (
	StatusUnset    Status = ""
	StatusActive   Status = "active"
	StatusSelected Status = "selected"
)

type Choice struct {
	Name    string    `json:"name"`
	Key     string    `json:"key"`
	Type    Type      `json:"type"`
	Status  Status    `json:"status"`
	Count   int       `json:"count"`
	Options []*Option `json:"options"`
}
type Option struct {
	Type             OptionType              `json:"type"`
	CountedReference *CountedReferenceOption `json:"counted_reference,omitempty"`
	Reference        *ReferenceOption        `json:"reference,omitempty"`
	Choice           *Choice                 `json:"choice,omitempty"`
	Multiple         *MultipleOption         `json:"multiple,omitempty"`
}

type ReferenceItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CountedReferenceOption struct {
	Status    Status         `json:"status"`
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

type ReferenceOption struct {
	Status    Status         `json:"status"`
	Reference *ReferenceItem `json:"reference"`
}

type MultipleOption struct {
	Status  Status    `json:"status"`
	Options []*Option `json:"options"`
}
