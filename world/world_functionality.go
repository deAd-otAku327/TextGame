package world

import (
	"strings"
)

// Конструктор ключа локации.
func NewDomainInventoryKey(val string, priority int) DomainInventoryKey {
	var dmk DomainInventoryKey

	dmk.KeyValue = val
	dmk.Priority = priority

	return dmk
}

// Конструктор интерактивных объектов.
func NewInteractiveObject(name string, key string, loc string, banned []string, on string, off string, mode bool) InteractiveObject {
	var io InteractiveObject

	io.objName = name
	io.key = key

	io.bannedUserRequests = banned

	io.locationName = loc

	io.stateInfoOn = on
	io.stateInfoOff = off

	io.isActivated = mode

	return io
}

// Конструктор квестов.
func NewQuest(name string, upd func(ch *Character) bool) *Quest {
	var q Quest

	q.SetName(name)
	q.checkUpdate = upd
	return &q
}

// Конструктор предметов.
func NewItem(name string, isWearable, isGetable, isInventory, isUsable bool) Item {
	var it Item

	it.name = name

	it.IsWearable = isWearable
	it.IsGetable = isGetable
	it.IsInventory = isInventory
	it.IsUsable = isUsable

	return it
}

// Конструктор AvailableDomainNaming.
func NewAvlDomain(name string, alies string) AvalibleDomainNaming {
	avd := AvalibleDomainNaming{name, alies}
	return avd
}

// Методы AvailableDomainNaming.
func (avd AvalibleDomainNaming) GetNaming() (name string, alies string) {
	return avd.name, avd.alias
}

// Методы Domain.
func (d Domain) GetName() string {
	return d.name
}

func (d *Domain) SetName(n string) {
	d.name = n
}

func (d Domain) GetInventory() DomainInventory {
	return d.inventory
}

func (d *Domain) SetInventory(inv DomainInventory) {
	d.inventory = inv
}

func (d Domain) GetAvailableDomains() []AvalibleDomainNaming {
	return d.availableDomains
}

func (d *Domain) SetAvailableDomains(dList []AvalibleDomainNaming) {
	d.availableDomains = dList
}

func (d Domain) GetInfo() string {
	return d.domainInfo
}

func (d *Domain) SetInfo(info string) {
	d.domainInfo = info
}

// Удалить предмет из инвентаря, возвращает успешно ли прошло удаление (для общности, если программист сам не проконтролил).
func (d *Domain) DeleteItem(it Item) bool {
	key, index := findItem(it, *d)

	if index == -1 {
		return false
	}

	d.inventory[key] = append(d.inventory[key][:index], d.inventory[key][index+1:]...)
	return true
}

// Методы Character.
func (ch Character) GetInventory() []Item {
	return ch.inventory
}

func (ch Character) GetQuestbook() []*Quest {
	return ch.questBook
}

func (ch *Character) SetQuestbook(qb []*Quest) {
	ch.questBook = qb
}

func (ch *Character) AddItem(it Item) {
	ch.inventory = append(ch.inventory, it)
}

func (ch Character) GetQuest(n string) *Quest {
	for _, q := range ch.GetQuestbook() {
		if q.GetName() == n {
			return q
		}
	}

	return nil
}

func (ch *Character) QuestControler() {
	currQbook := ch.GetQuestbook()

	for _, quest := range currQbook { // Цикл по квестам.
		quest.Update()
	}
}

// Методы World.
func (w *World) GetDomain(tag string) *Domain {
	for _, loc := range w.LocationBase {
		if loc.name == tag {
			return &loc
		}
	}

	return nil
}

func (w World) GetInteractList() []*InteractiveObject {
	var ioPointers = make([]*InteractiveObject, len(w.IOList)) // Сила преаллокейтинга.
	for index := range w.IOList {
		ioPointers[index] = &w.IOList[index]
	}
	return ioPointers
}

// Методы Interactiveobject.
func (io InteractiveObject) GetStateInfo() string {
	if io.isActivated {
		return io.stateInfoOn
	} else {
		return io.stateInfoOff + ", требуется: " + io.GetKey()
	}
}

func (io *InteractiveObject) Activate() {
	io.isActivated = true
}

func (io InteractiveObject) GetName() string {
	return io.objName
}

func (io InteractiveObject) GetKey() string {
	return io.key
}

func (io InteractiveObject) GetLocationName() string {
	return io.locationName
}

func (io InteractiveObject) CheckBanned(parsedCom []string) bool { // Забанен ли запрос юзера неактивным объектом, с учетом состояния.
	for _, banned := range io.bannedUserRequests {
		// Собираем команду после парса обратно в строку, чтобы сопоставить форматом bannedUserRequests.
		if !io.isActivated && banned == strings.Join(parsedCom, " ") {
			return true
		}
	}

	return false
}

// Методы Quest.
func (q Quest) CheckDone() bool {
	return q.done
}

func (q *Quest) Update() {
	q.done = q.checkUpdate(CurrCharacter)
}

func (q Quest) GetName() string {
	return q.name
}

func (q *Quest) SetName(n string) {
	q.name = n
}

// Методы Item.
func (it Item) GetName() string {
	return it.name
}

// Вспомогательная функция поиска айтема в инвентаре помещения.
func findItem(it Item, d Domain) (key DomainInventoryKey, index int) {
	for key, itemList := range d.GetInventory() {
		for index, item := range itemList {

			if item == it {
				return key, index
			}
		}
	}

	return DomainInventoryKey{KeyValue: "", Priority: 0}, -1
}
