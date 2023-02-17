package entities

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type suiteChoice struct {
	suite.Suite

	choice *Choice
}

func (s *suiteChoice) SetupTest() {
	s.choice = &Choice{
		Name:  "(a) a martial weapon and a shield or (b) two martial weapons",
		Count: 1,
		Type:  ChoiceTypeEquipment,
		Options: []Option{
			&MultipleOption{
				Key:  "option-1",
				Name: "a martial weapon and a shield",
				Items: []Option{
					&Choice{
						Name:  "a martial weapon",
						Count: 1,
						Type:  ChoiceTypeEquipment,
						Options: []Option{
							&ReferenceOption{
								Reference: &ReferenceItem{
									Key:  "battleaxe",
									Name: "Battleaxe",
								},
							},
							&ReferenceOption{
								Reference: &ReferenceItem{
									Key:  "flail",
									Name: "Flail",
								},
							},
						},
					},
					&CountedReferenceOption{
						Count: 1,
						Reference: &ReferenceItem{
							Key:  "shield",
							Name: "Shield",
						},
					},
				},
			},
			&Choice{
				Name:  "two martial weapons",
				Count: 2,
				Type:  ChoiceTypeEquipment,
				Options: []Option{
					&ReferenceOption{
						Reference: &ReferenceItem{
							Key:  "handaxe",
							Name: "Handaxe",
						},
					},
					&ReferenceOption{
						Reference: &ReferenceItem{
							Key:  "longsword",
							Name: "Longsword",
						},
					},
					&ReferenceOption{
						Reference: &ReferenceItem{
							Key:  "rapier",
							Name: "Rapier",
						},
					},
				},
			},
		},
	}
}

func (s *suiteChoice) TestSelectsMultiple() {
	actual := s.choice.Select("option-1")
	s.Equal("a martial weapon and a shield", actual.GetName())

}

func TestChoice(t *testing.T) {
	suite.Run(t, new(suiteChoice))
}
