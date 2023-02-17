package entities

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeMultiple         OptionType = "multiple"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type ChoiceStatus string

const (
	ChoiceStatusUnset    ChoiceStatus = ""
	ChoiceStatusActive   ChoiceStatus = "active"
	ChoiceStatusInactive ChoiceStatus = "inactive"
	ChoiceStatusSelected ChoiceStatus = "selected"
)

type ChoiceType string

const (
	ChoiceTypeUnset       ChoiceType = ""
	ChoiceTypeProficiency ChoiceType = "proficiency"
	ChoiceTypeLanguage    ChoiceType = "language"
	ChoiceTypeEquipment   ChoiceType = "equipment"
)

type Choice struct {
	Name     string       `json:"name"`
	Type     ChoiceType   `json:"type"`
	Key      string       `json:"key"`
	Status   ChoiceStatus `json:"status"`
	Selected bool         `json:"selected"`
	Count    int          `json:"count"`
	Options  []Option     `json:"options"`
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

func (o *Choice) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *Choice) SetStatus(status ChoiceStatus) {
	o.Status = status
}

// Select selects an option by key
// sets the current choice to Active if there are remaining opiotns to be chosen, otherwise marks the top level choice as Selected
func (o *Choice) Select(key string) Option {
	var selected Option
	selectedCount := 0

	for _, option := range o.Options {
		if option.GetStatus() == ChoiceStatusSelected {
			selectedCount++
			continue // if this option is selected we don't need to do anything
		}

		if option.GetKey() == key {
			option.SetStatus(ChoiceStatusActive)
			selected = option

			continue
		}

		if option.GetOptionType() == OptionTypeChoice {
			choiceOption := option.(*Choice)
			choice := choiceOption.Select(key)
			if choice != nil {
				selected = choice
				if selected.GetStatus() == ChoiceStatusSelected {
					selectedCount++
				}
			}
			// TODO: if the selected status is selected what should we do?
		}

		if option.GetOptionType() == OptionTypeMultiple {
			multipleOption := option.(*MultipleOption)
			totalCount := 0
			// Go through all the items and select the one that matches the key
			// keep track of how many have been selected
			for _, item := range multipleOption.Items {
				if item.GetKey() == key {
					if item.GetStatus() != ChoiceStatusSelected {
						selected = item
						item.SetStatus(ChoiceStatusSelected)
					}
					totalCount++
				}
			}
			// If they have all been selected we will mark the top level option as selected
			if totalCount == o.Count {
				option.SetStatus(ChoiceStatusSelected)
			}
		}

		if option.GetKey() == key {
			option.SetStatus(ChoiceStatusSelected)
			selected = option
			selectedCount++
		}
	}

	if selectedCount == o.Count {
		o.Status = ChoiceStatusSelected
	} else if selectedCount > 0 {
		o.Status = ChoiceStatusActive
	}

	return selected
}

type Option interface {
	GetOptionType() OptionType
	GetName() string
	GetKey() string
	GetStatus() ChoiceStatus
	SetStatus(ChoiceStatus)
}

type CountedReferenceOption struct {
	Status    ChoiceStatus   `json:"status"`
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

func (o *CountedReferenceOption) GetOptionType() OptionType {
	return OptionTypeCountedReference
}

func (o *CountedReferenceOption) GetName() string {
	return o.Reference.Name
}

func (o *CountedReferenceOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *CountedReferenceOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}

func (o *CountedReferenceOption) GetKey() string {
	return o.Reference.Key
}

type ReferenceOption struct {
	Status    ChoiceStatus   `json:"status"`
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

func (o *ReferenceOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *ReferenceOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}

type MultipleOption struct {
	Status ChoiceStatus `json:"status"`
	Key    string       `json:"key"`
	Name   string       `json:"name"`
	Items  []Option     `json:"items"`
}

func (o *MultipleOption) Select(key string) Option {
	totalCount := 0
	var selected Option

	// Go through all the items and select the one that matches the key
	// keep track of how many have been selected
	for _, item := range o.Items {
		switch item.GetOptionType() {
		case OptionTypeChoice:
			choiceOption := item.(*Choice)
			choice := choiceOption.Select(key)
			if choice != nil {
				selected = choice
				if selected.GetStatus() == ChoiceStatusSelected {
					totalCount++
				}

				break
			}
		default:
			if item.GetKey() == key {
				totalCount++

				if item.GetStatus() != ChoiceStatusSelected {
					selected = item
					item.SetStatus(ChoiceStatusSelected)

					break
				}
			}
		}
	}

	// If they have all been selected we will mark the top level option as selected
	if totalCount == len(o.Items) {
		o.Status = ChoiceStatusSelected
	}

	return selected
}
func (o *MultipleOption) GetOptionType() OptionType {
	return OptionTypeMultiple
}

func (o *MultipleOption) GetName() string {
	return o.Name
}

func (o *MultipleOption) GetKey() string {
	return o.Key
}

func (o *MultipleOption) GetStatus() ChoiceStatus {
	return o.Status
}

func (o *MultipleOption) SetStatus(status ChoiceStatus) {
	o.Status = status
}
