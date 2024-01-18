package action

import (
	"gameProject/world"
	"strings"
)

const wrongArgs string = "некорретный набор аргументов команды"
const StopGame string = "игра завершена"

func NewCommandBase() world.CommandBase {
	var cBase = world.CommandBase{

		"осмотреться": lookAround,

		"идти": goNext,

		"надеть": putOnItem,

		"взять": takeItem,

		"применить": useItem,

		"задачи": quests,

		"помощь": func(args []string) string {
			cmds := "ДОСТУПНЫЕ КОМАНДЫ:\n\nосмотреться\nидти $domain\nнадеть $item\nвзять $item\nприменить $item $interactive odject\nзадачи"
			return cmds
		},

		"стоп": func(args []string) string { return StopGame },
	}

	return cBase
}

func lookAround(args []string) string {

	if len(args) > 0 {
		return wrongArgs
	}

	var answer string

	if !world.CurrCharacter.InfoAnnounced {
		answer += world.CurrCharacter.CurrLocation.GetInfo() + ", "
		world.CurrCharacter.InfoAnnounced = true
	}

	answer += checkContent(*world.CurrCharacter.CurrLocation) + ". "

	answer += checkAvailableDomains(*world.CurrCharacter.CurrLocation)

	return answer
}

func goNext(args []string) string {

	if len(args) != 1 {
		return wrongArgs
	}

	var answer string

	// Проверяем, есть ли локация среди доступных.
	contains, dNaming := findDomainName(args[0], world.CurrCharacter.CurrLocation.GetAvailableDomains())
	var loc *world.Domain

	if contains {
		n, _ := dNaming.GetNaming()
		loc = world.CurrWorld.GetDomain(n) // Проверяем, существует ли локация в мире.
	}

	if args[0] == world.CurrCharacter.CurrLocation.GetName() {
		return "вы уже в данной локации"
	} else if loc != nil {
		world.CurrCharacter.CurrLocation = loc
	} else {
		return "нет пути в " + args[0]
	}

	answer = world.CurrCharacter.CurrLocation.GetInfo() + ". "
	world.CurrCharacter.InfoAnnounced = true

	answer += checkAvailableDomains(*world.CurrCharacter.CurrLocation)

	AllQuestsUpd(world.CurrCharacter)

	return answer
}

func putOnItem(args []string) string {

	if len(args) != 1 {
		return wrongArgs
	}

	var answer string

	currInventory := world.CurrCharacter.CurrLocation.GetInventory()

	for _, itemSet := range currInventory {
		for _, item := range itemSet {
			if item.GetName() == args[0] {

				if item.IsWearable {

					answer += "вы надели: " + args[0]

					world.CurrCharacter.AddItem(item)
					world.CurrCharacter.CurrLocation.DeleteItem(item)

					world.CurrCharacter.QuestControler()

					if item.IsInventory { // Если надет предмет-инвентарь - разбаниваем собирание предметов.
						world.CurrCharacter.InventoryUnlocked = true
					}

				} else {
					answer += "данный предмет невозможно надеть"
				}

				return answer

			}
		}
	}

	answer += "нет такого"

	AllQuestsUpd(world.CurrCharacter)

	return answer
}

func takeItem(args []string) string {
	if len(args) != 1 {
		return wrongArgs
	}

	var answer string

	currInventory := world.CurrCharacter.CurrLocation.GetInventory()

	for _, itemSet := range currInventory {
		for _, item := range itemSet {
			if item.GetName() == args[0] {

				if item.IsGetable {

					if world.CurrCharacter.InventoryUnlocked {

						answer += "предмет добавлен в инвентарь: " + args[0]

						world.CurrCharacter.AddItem(item)
						world.CurrCharacter.CurrLocation.DeleteItem(item)

						world.CurrCharacter.QuestControler()

					} else {
						answer += "некуда класть"
					}
				} else {
					answer += "данный предмет невозможно взять"
				}

				return answer
			}
		}
	}

	answer += "нет такого"

	AllQuestsUpd(world.CurrCharacter)

	return answer
}

func useItem(args []string) string {

	if len(args) != 2 {
		return wrongArgs
	}

	var answer string
	var counter int

	currInv := world.CurrCharacter.GetInventory() // Слайс айтемов.
	ioList := world.CurrWorld.GetInteractList()   // Слайс указателей на объекты.

	// Блок проверки инвентаря персонажа.
	for _, item := range currInv {

		if item.GetName() == args[0] {
			break
		}
		counter++
	}

	if counter == len(currInv) || len(currInv) == 0 {
		answer += "нет предмета в инвентаре - " + args[0]
		return answer
	}

	// Блок проверки интерактивных объектов в текущей локации.
	for _, obj := range ioList {

		if obj.GetName() == args[1] && obj.GetLocationName() == world.CurrCharacter.CurrLocation.GetName() && obj.GetKey() == args[0] {
			// Блок активации интерактивного объекта.
			obj.Activate()

			answer += obj.GetStateInfo()

			return answer
		}
	}

	answer += "применение невозможно"

	AllQuestsUpd(world.CurrCharacter)

	return answer
}

func quests(args []string) string {
	var answer string

	answer += "надо "

	questsList := world.CurrCharacter.GetQuestbook() // Список квестов.

	qDone := 0
	qContainer := make([]string, 0)

	for _, quest := range questsList {
		if !quest.CheckDone() { // Если квест еще не выполнен...
			qContainer = append(qContainer, quest.GetName(), ", ")
		} else {
			qDone++
		}
	}

	if len(questsList) == qDone {
		return "нет активных задач"
	}

	if len(qContainer) > 2 { // Размер контейнера больше 1 квеста с учетом запятой.
		qContainer[len(qContainer)-3] = " и "
	}

	answer += strings.Join(qContainer[:len(qContainer)-1], "")

	return answer
}
