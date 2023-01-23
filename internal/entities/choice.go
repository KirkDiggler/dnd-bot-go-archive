package entities

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeMultiple         OptionType = "multiple"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type Choice struct {
	Name     string   `json:"name"`
	Key      string   `json:"key"`
	Selected bool     `json:"selected"`
	Count    int      `json:"count"`
	Options  []Option `json:"options"`
}
type Option interface {
	GetOptionType() OptionType
	GetName() string
	GetKey() string
}

type CountedReferenceOption struct {
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

func (o *CountedReferenceOption) GetOptionType() OptionType {
	return OptionTypeCountedReference
}

func (o *CountedReferenceOption) GetName() string {
	return o.Reference.Name
}

func (o *CountedReferenceOption) GetKey() string {
	return o.Reference.Key
}

type ReferenceOption struct {
	Reference *ReferenceItem `json:"reference"`
}

func (o *ReferenceOption) GetOptionType() OptionType {
	return OptionTypeReference
}

func (o *ReferenceOption) GetName() string {
	return o.Reference.Name
}

func (o *ReferenceOption) GetKey() string {
	return o.Reference.Key
}

func (o *Choice) GetOptionType() OptionType {
	return OptionTypeChoice
}

func (o *Choice) GetName() string {
	return o.Name
}

func (o *Choice) GetKey() string {
	return o.Key
}

type MultipleOption struct {
	Selected bool     `json:"selected"`
	Items    []Option `json:"items"`
}

func (o *MultipleOption) GetOptionType() OptionType {
	return OptionTypeMultiple
}

func (o *MultipleOption) GetName() string {
	return ""
}

func (o *MultipleOption) GetKey() string {
	return ""
}
