package action

import (
	"gameProject/world"
)

func checkAvailableDomains(d world.Domain) string {

	domainList := "можно пройти - "

	for index, value := range d.GetAvailableDomains() {

		n, al := value.GetNaming()
		if al != "" {
			domainList += al
		} else {
			domainList += n
		}

		if index != len(d.GetAvailableDomains())-1 {
			domainList += ", "
		}
	}

	return domainList
}

func checkContent(d world.Domain) string {

	var content string

	// Значение для корретности знаков препинания.
	var nonEmptyInventoryCount int

	for _, itemList := range d.GetInventory() {
		if len(itemList) != 0 {
			nonEmptyInventoryCount++
		}
	}

	if nonEmptyInventoryCount == 0 {
		return "здесь ничего нет"
	}

	keyList := keyPrioritySorter(d.GetInventory()) // Список инвентарей по ключам.

	for _, key := range keyList {

		currInv := d.GetInventory()[key]

		if len(currInv) != 0 {
			nonEmptyInventoryCount--

			content += key.KeyValue + ": "

			for index, item := range currInv {
				content += item.GetName()

				// Если текущий инвентарь не закончился или есть еще непустые инвентари.
				if index != len(currInv)-1 || nonEmptyInventoryCount != 0 {
					content += ", "
				}
			}

		}

	}

	return content
}

// Вспомогательная функция для обеспечения порядка вывода контента в соотвествии с приоритетами, заданными в ините.
func keyPrioritySorter(inv world.DomainInventory) []world.DomainInventoryKey {
	sortedKeys := make([]world.DomainInventoryKey, len(inv)) // Слайс длиной = количеству видов(ключей) инвентарей.

	for k := range inv {
		sortedKeys[k.Priority-1] = k
	}

	return sortedKeys
}

// Вспомогательная функция поиска имени локации в слайсе, учитывая элиасы.
func findDomainName(s string, arr []world.AvalibleDomainNaming) (bool, world.AvalibleDomainNaming) {
	for _, val := range arr {
		n, al := val.GetNaming()
		if n == s || al == s {
			return true, val
		}
	}
	return false, world.AvalibleDomainNaming{}
}

// Функция апдейта всех квестов.
func AllQuestsUpd(ch *world.Character) {
	for _, q := range ch.GetQuestbook() {
		q.Update()
	}
}
