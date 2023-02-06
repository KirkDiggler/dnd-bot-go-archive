package characters

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

func (m *manager) characterFromData(ctx context.Context, data *character.Data) (*entities.Character, error) {
	if data == nil {
		return nil, dnderr.NewMissingParameterError("data")
	}

	g, _ := errgroup.WithContext(ctx)

	var race *entities.Race
	var class *entities.Class

	g.Go(func() (err error) {

		race, err = m.client.GetRace(data.RaceKey)
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() (err error) {
		class, err = m.client.GetClass(data.ClassKey)
		if err != nil {
			return err
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	char := &entities.Character{
		ID:               data.ID,
		Name:             data.Name,
		OwnerID:          data.OwnerID,
		Speed:            data.Speed,
		AC:               data.AC,
		MaxHitPoints:     data.MaxHitPoints,
		CurrentHitPoints: data.CurrentHitPoints,
		HitDie:           data.HitDie,
		Experience:       data.Experience,
		Level:            data.Level,
		NextLevel:        data.NextLevel,
		Race:             race,
		Class:            class,
		Attribues:        attributDataToAttributes(data.Attributes),
		Rolls:            rollDatasToRollResults(data.Rolls),
	}

	for _, prof := range data.Proficiencies {
		prof := prof

		g.Go(func() (err error) {
			char.AddProficiency(dataToProficiency(prof))
			return nil
		})
	}

	for _, item := range data.Inventory {
		item := item
		g.Go(func() (err error) {
			equip, err := m.client.GetEquipment(item.Key)
			if err != nil {
				return err
			}

			char.AddInventory(equip)

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	for _, slot := range data.EquippedSlots {
		char.Equip(slot.Key)
	}

	return char, nil
}

func dataToProficiency(data *character.Proficiency) *entities.Proficiency {
	if data == nil {
		return nil
	}

	return &entities.Proficiency{
		Name: data.Name,
		Key:  data.Key,
		Type: entities.ProficiencyType(data.Type),
	}
}

func attributDataToAttributes(data *character.AttributeData) map[entities.Attribute]*entities.AbilityScore {
	if data == nil {
		data = &character.AttributeData{}
	}

	return map[entities.Attribute]*entities.AbilityScore{
		entities.AttributeStrength:     abilityScoreDataToAbilityScore(data.Str),
		entities.AttributeDexterity:    abilityScoreDataToAbilityScore(data.Dex),
		entities.AttributeConstitution: abilityScoreDataToAbilityScore(data.Con),
		entities.AttributeIntelligence: abilityScoreDataToAbilityScore(data.Int),
		entities.AttributeWisdom:       abilityScoreDataToAbilityScore(data.Wis),
		entities.AttributeCharisma:     abilityScoreDataToAbilityScore(data.Cha),
	}
}

func abilityScoreDataToAbilityScore(data *character.AbilityScoreData) *entities.AbilityScore {
	if data == nil {
		return nil
	}

	return &entities.AbilityScore{
		Score: data.Score,
		Bonus: data.Bonus,
	}
}
func rollDataToRollResult(data *character.RollData) *dice.RollResult {
	if data == nil {
		return nil
	}

	return &dice.RollResult{
		Used:    data.Used,
		Total:   data.Total,
		Highest: data.Highest,
		Lowest:  data.Lowest,
		Rolls:   data.Rolls,
	}
}

func rollDatasToRollResults(data []*character.RollData) []*dice.RollResult {
	results := make([]*dice.RollResult, len(data))
	for i, d := range data {
		results[i] = rollDataToRollResult(d)
	}

	return results
}

func characterToData(input *entities.Character) *character.Data {
	if input == nil {
		return nil
	}

	return &character.Data{
		ID:            input.ID,
		Name:          input.Name,
		OwnerID:       input.OwnerID,
		RaceKey:       input.Race.Key,
		ClassKey:      input.Class.Key,
		Attributes:    attributesToAttributeData(input.Attribues),
		Rolls:         rollResultsToRollDatas(input.Rolls),
		Proficiencies: proficienciesToDatas(input.Proficiencies),
		Inventory:     equipmentsToDatas(input.Inventory),
		EquippedSlots: equippedSlotsToDatas(input.EquippedSlots),
	}
}
func equippedSlotsToDatas(input map[entities.Slot]entities.Equipment) map[entities.Slot]*character.Equipment {
	datas := make(map[entities.Slot]*character.Equipment)

	for k, v := range input {
		datas[k] = equipmentToData(v)
	}

	return datas
}

func proficienciesToDatas(input map[entities.ProficiencyType][]*entities.Proficiency) []*character.Proficiency {
	datas := make([]*character.Proficiency, 0)

	for _, v := range input {
		for _, p := range v {
			datas = append(datas, proficiencyToData(p))
		}
	}

	return datas
}

func proficiencyToData(input *entities.Proficiency) *character.Proficiency {
	if input == nil {
		return nil
	}

	return &character.Proficiency{
		Name: input.Name,
		Key:  input.Key,
		Type: string(input.Type),
	}
}

func rollResultsToRollDatas(input []*dice.RollResult) []*character.RollData {
	datas := make([]*character.RollData, len(input))
	for i, r := range input {
		datas[i] = rollResultToRollData(r)
	}

	return datas
}
func rollResultToRollData(input *dice.RollResult) *character.RollData {
	if input == nil {
		return nil
	}

	return &character.RollData{
		Used:    input.Used,
		Total:   input.Total,
		Highest: input.Highest,
		Lowest:  input.Lowest,
		Rolls:   input.Rolls,
	}
}

func abilityScoreToData(input *entities.AbilityScore) *character.AbilityScoreData {
	if input == nil {
		return nil
	}

	return &character.AbilityScoreData{
		Score: input.Score,
		Bonus: input.Bonus,
	}
}
func attributesToAttributeData(input map[entities.Attribute]*entities.AbilityScore) *character.AttributeData {
	if input == nil {
		return nil
	}

	return &character.AttributeData{
		Str: abilityScoreToData(input[entities.AttributeStrength]),
		Dex: abilityScoreToData(input[entities.AttributeDexterity]),
		Con: abilityScoreToData(input[entities.AttributeConstitution]),
		Int: abilityScoreToData(input[entities.AttributeIntelligence]),
		Wis: abilityScoreToData(input[entities.AttributeWisdom]),
		Cha: abilityScoreToData(input[entities.AttributeCharisma]),
	}
}
func equipmentToData(input entities.Equipment) *character.Equipment {
	return &character.Equipment{
		Key:  input.GetKey(),
		Name: input.GetName(),
		Type: input.GetEquipmentType(),
	}
}

func equipmentsToDatas(input map[entities.EquipmentType][]entities.Equipment) []*character.Equipment {
	datas := make([]*character.Equipment, 0)

	for _, v := range input {
		for _, e := range v {
			datas = append(datas, equipmentToData(e))
		}
	}

	return datas
}
