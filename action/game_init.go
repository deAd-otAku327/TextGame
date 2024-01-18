package action

import (
	"gameProject/world"
)

func InitGame() {
	// Создаем мир и персонажа.
	world.CurrCharacter = new(world.Character)
	world.CurrWorld = new(world.World)

	// Поднимаем базу команд.
	world.CurrWorld.CommandList = NewCommandBase()

	var arrInv []world.Item                        // Шелл для айтемов.
	var arrAvlDomains []world.AvalibleDomainNaming // Шелл для доступных локаций.
	var arrBanned []string                         // Шелл для забаненных действий.

	// Инициализация кухни.
	var kitchen world.Domain

	kitchen.SetName("кухня")

	kitchen.SetInventory(world.DomainInventory{world.DomainInventoryKey{KeyValue: "на столе", Priority: 1}: append(arrInv, world.NewItem("чай", false, false, false, false))})

	kitchen.SetAvailableDomains(append(arrAvlDomains, world.NewAvlDomain("коридор", "")))
	kitchen.SetInfo("ты находишься на кухне")

	world.CurrWorld.LocationBase = append(world.CurrWorld.LocationBase, kitchen) // Закинули в список локаций мира.

	// Инициализация комнаты.
	var room world.Domain

	room.SetName("комната")

	room.SetInventory(world.DomainInventory{world.DomainInventoryKey{KeyValue: "на столе", Priority: 1}: append(arrInv, world.NewItem("ключи", false, true, false, true), world.NewItem("конспекты", false, true, false, false)),
		world.DomainInventoryKey{KeyValue: "на стуле", Priority: 2}: append(arrInv,
			world.NewItem("рюкзак", true, false, true, false))},
	)

	room.SetAvailableDomains(append(arrAvlDomains, world.NewAvlDomain("коридор", "")))
	room.SetInfo("ты в своей комнате")

	world.CurrWorld.LocationBase = append(world.CurrWorld.LocationBase, room)

	// Инициализация холла
	var hall world.Domain

	hall.SetName("коридор")
	hall.SetAvailableDomains(append(arrAvlDomains, world.NewAvlDomain("кухня", ""), world.NewAvlDomain("комната", ""), world.NewAvlDomain("улица", "")))
	hall.SetInfo("ничего интересного")

	world.CurrWorld.LocationBase = append(world.CurrWorld.LocationBase, hall)

	// Инициализация улицы
	var outdoors world.Domain

	outdoors.SetName("улица")
	outdoors.SetInfo("на улице весна")
	outdoors.SetAvailableDomains(append(arrAvlDomains, world.NewAvlDomain("коридор", "дом")))

	world.CurrWorld.LocationBase = append(world.CurrWorld.LocationBase, outdoors)

	// Иницализация персонажа
	world.CurrCharacter.CurrLocation = world.CurrWorld.GetDomain("кухня")
	world.CurrCharacter.InfoAnnounced = false
	world.CurrCharacter.InventoryUnlocked = false
	world.CurrCharacter.SetQuestbook([]*world.Quest{world.NewQuest("собрать рюкзак",
		func(ch *world.Character) bool {
			currInv := ch.GetInventory()

			itemFlag1 := false
			itemFlag2 := false

			for _, item := range currInv {
				if item.GetName() == "рюкзак" {
					itemFlag1 = true
				}

				if item.GetName() == "конспекты" {
					itemFlag2 = true
				}
			}

			if itemFlag1 && itemFlag2 {
				return true
			}
			return false
		}),
		world.NewQuest("идти в универ",
			func(ch *world.Character) bool {
				if ch.CurrLocation.GetName() == "улица" && ch.GetQuest("собрать рюкзак").CheckDone() {
					return true
				}
				return false
			}),
	})

	// Инициализация интерактивных объектов.
	world.CurrWorld.IOList = append(world.CurrWorld.IOList, world.NewInteractiveObject("дверь", "ключи", "коридор",
		append(arrBanned, "идти улица"), "дверь открыта", "дверь закрыта", false))
}
