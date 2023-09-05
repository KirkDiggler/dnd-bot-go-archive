package rooms

import (
	"context"
	"fmt"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/monster"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/room"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
)

const defaultMonster = "goblin"

type Implementation struct {
	client           dnd5e.Client
	characterManager characters.Manager
	roomRepo         room.Repository
	monsterRepo      monster.Interface
	uuider           types.UUIDGenerator
}

type Config struct {
	Client           dnd5e.Client
	CharacterManager characters.Manager
	RoomRepo         room.Repository
	MonsterRepo      monster.Interface
}

func New(cfg *Config) (*Implementation, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.CharacterManager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharacterManager")
	}

	if cfg.RoomRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.RoomRepo")
	}

	if cfg.MonsterRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.MonsterRepo")
	}

	return &Implementation{
		client:           cfg.Client,
		characterManager: cfg.CharacterManager,
		roomRepo:         cfg.RoomRepo,
		monsterRepo:      cfg.MonsterRepo,
		uuider:           &types.GoogleUUID{},
	}, nil
}

func (m *Implementation) LoadRoom(ctx context.Context, input *LoadRoomInput) (*LoadRoomOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	rooms, err := m.roomRepo.ListByPlayer(ctx, &room.ListByPlayerInput{
		PlayerID: input.PlayerID,
		Limit:    1,
		Offset:   0,
		Reverse:  true,
	})
	if err != nil {
		return nil, err
	}

	if len(rooms) == 0 || rooms[0].Status == room.StatusInactive {
		out, err := m.createRoom(ctx, input.PlayerID)
		if err != nil {
			return nil, err
		}
		return &LoadRoomOutput{
			Room: out,
		}, nil
	}

	out, err := m.hydrateRoom(ctx, rooms[0])
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate room: %w", err)
	}

	return &LoadRoomOutput{
		Room: out,
	}, nil
}

func (m *Implementation) HasActiveRoom(ctx context.Context, input *HasActiveRoomInput) (*HasActiveRoomOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	rooms, err := m.roomRepo.ListByPlayer(ctx, &room.ListByPlayerInput{
		PlayerID: input.PlayerID,
		Limit:    1,
		Offset:   0,
		Reverse:  true,
	})
	if err != nil {
		return nil, err
	}

	if len(rooms) == 0 || rooms[0].Status == room.StatusInactive {
		return &HasActiveRoomOutput{
			HasActiveRoom: false,
		}, nil
	}

	return &HasActiveRoomOutput{
		HasActiveRoom: rooms[0].Status == room.StatusActive,
	}, nil
}

func (m *Implementation) Attack(ctx context.Context, playerID string) (string, error) {
	rm, err := m.LoadRoom(ctx, &LoadRoomInput{PlayerID: playerID})
	if err != nil {
		return "", err
	}

	player, err := rm.Room.Character.Attack()
	if err != nil {
		return "", err
	}

	mon, err := rm.Room.Monster.Attack()
	if err != nil {
		return "", err
	}

	var msgBuilder strings.Builder

	if rm.Room.CharacterInitiative == 0 {
		rollResult, err := dice.Roll(1, 20, rm.Room.Character.Attribues[entities.AttributeDexterity].Bonus)
		if err != nil {
			return "", err
		}
		rm.Room.CharacterInitiative = rollResult.Total
	}

	if rm.Room.MonsterInitiative == 0 {
		rollResult, err := dice.Roll(1, 20, 0)
		if err != nil {
			return "", err
		}
		rm.Room.MonsterInitiative = rollResult.Total
	}

	msgBuilder.WriteString(fmt.Sprintf("You rolled a %d for initiative\n", rm.Room.CharacterInitiative))
	msgBuilder.WriteString(fmt.Sprintf("The monster rolled a %d for initiative\n", rm.Room.MonsterInitiative))

	msgBuilder.WriteString(fmt.Sprintf("You rolled a %d for attack, the monster has AC: %d\n", player[0].AttackRoll, rm.Room.Monster.Template.ArmorClass))
	msgBuilder.WriteString(fmt.Sprintf("The monster rolled a %d for attack, you have AC: %d\n", mon[0].AttackRoll, rm.Room.Character.AC))
	if rm.Room.CharacterInitiative > rm.Room.MonsterInitiative {
		if player[0].AttackRoll >= rm.Room.Monster.Template.ArmorClass {
			msgBuilder.WriteString(fmt.Sprintf("You hit the monster for %d damage\n", player[0].DamageRoll))
			rm.Room.Monster.CurrentHP -= player[0].DamageRoll
			if rm.Room.Monster.CurrentHP <= 0 {
				rm.Room.Status = entities.RoomStatusWon
				msgBuilder.WriteString("you defeated the monster")
				return msgBuilder.String(), nil
			}
		}

		if mon[0].AttackRoll >= rm.Room.Character.AC {
			msgBuilder.WriteString(fmt.Sprintf("The monster hit you for %d damage\n", mon[0].DamageRoll))
			rm.Room.Character.CurrentHitPoints -= mon[0].DamageRoll
			if rm.Room.Character.CurrentHitPoints <= 0 {
				rm.Room.Status = entities.RoomStatusLost
				msgBuilder.WriteString("you were defeated")
				return msgBuilder.String(), nil
			}
		}
		rm.Room.Character.CurrentHitPoints -= mon[0].DamageRoll
	} else {
		if mon[0].AttackRoll >= rm.Room.Character.AC {
			msgBuilder.WriteString(fmt.Sprintf("The monster hit you for %d damage\n", mon[0].DamageRoll))
			rm.Room.Character.CurrentHitPoints -= mon[0].DamageRoll
			if rm.Room.Character.CurrentHitPoints <= 0 {
				rm.Room.Status = entities.RoomStatusLost
				msgBuilder.WriteString("you were defeated")
				return msgBuilder.String(), nil
			}
		}

		if player[0].AttackRoll >= rm.Room.Monster.Template.ArmorClass {
			rm.Room.Monster.CurrentHP -= player[0].DamageRoll
			if rm.Room.Monster.CurrentHP <= 0 {
				rm.Room.Status = entities.RoomStatusWon
				msgBuilder.WriteString("you defeated the monster")
				return msgBuilder.String(), nil
			}
		}
	}

	_, err = m.roomRepo.Update(ctx, room.EntityToData(rm.Room))
	if err != nil {
		return "", err
	}

	err = m.monsterRepo.PutMonster(ctx, rm.Room.Monster)
	if err != nil {
		return "", err
	}

	return msgBuilder.String(), nil
}

func statusToEntity(status room.Status) entities.RoomStatus {
	switch status {
	case room.StatusActive:
		return entities.RoomStatusActive
	case room.StatusInactive:
		return entities.RoomStatusInactive
	default:
		return entities.RoomStatusUnset
	}
}

func (m *Implementation) createRoom(ctx context.Context, playerID string) (*entities.Room, error) {
	monsterTemplate, err := m.client.GetMonster(defaultMonster)
	if err != nil {
		return nil, err
	}

	hp, err := dice.RollString(monsterTemplate.HitDice)
	if err != nil {
		return nil, err
	}

	mon := &entities.Monster{
		ID:          m.uuider.New(),
		CharacterID: playerID,
		Key:         monsterTemplate.Key,
		CurrentHP:   hp.Total,
	}

	err = m.monsterRepo.PutMonster(ctx, mon)
	if err != nil {
		return nil, err
	}

	mon.Template = monsterTemplate

	character, err := m.characterManager.Get(ctx, playerID)
	if err != nil {
		return nil, err
	}

	data, err := m.roomRepo.Create(ctx, &room.Data{
		PlayerID:  playerID,
		MonsterID: mon.ID,
		Status:    room.StatusActive,
	})
	if err != nil {
		return nil, err
	}

	return &entities.Room{
		ID:        data.ID,
		Status:    statusToEntity(data.Status),
		Character: character,
		Monster:   mon,
	}, nil
}

func (m *Implementation) hydrateRoom(ctx context.Context, room *room.Data) (*entities.Room, error) {
	if room == nil {
		return nil, dnderr.NewMissingParameterError("room")
	}

	out := &entities.Room{
		ID:     room.ID,
		Status: statusToEntity(room.Status),
	}

	character, err := m.characterManager.Get(ctx, room.PlayerID)
	if err != nil {
		return nil, err
	}

	out.Character = character

	mon, err := m.monsterRepo.GetMonster(ctx, room.MonsterID)
	if err != nil {
		return nil, err
	}

	monsterTemplate, err := m.client.GetMonster(mon.Key)
	if err != nil {
		return nil, err
	}

	mon.Template = monsterTemplate
	out.Monster = mon

	return out, nil
}
