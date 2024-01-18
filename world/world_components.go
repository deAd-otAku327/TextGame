package world

type CommandBase map[string]func(args []string) string

type DomainInventory map[DomainInventoryKey][]Item

type DomainInventoryKey struct {
	KeyValue string
	Priority int
}

type InteractiveObject struct {
	objName      string
	key          string
	locationName string

	bannedUserRequests []string

	stateInfoOn  string
	stateInfoOff string

	isActivated bool
}

type AvalibleDomainNaming struct {
	name  string // Название домейна.
	alias string // Псевдоним домейна, для элементарных локаций оставлять пустым, используется для комплексов.
}

type Domain struct {
	name string // Название локации.

	inventory DomainInventory // Инвентарь локации.

	domainInfo string // Первичная информация, получаемая игроком при переходе или осмотре локации.

	availableDomains []AvalibleDomainNaming // Список локаций, доступных для перехода из данной.
}

type Character struct {
	InfoAnnounced     bool // Флаг уведомленности юзера о его местоположении.
	InventoryUnlocked bool // Доступен ли инвентарь для сбора предметов (необходимо наличие инвентарного предмета).

	CurrLocation *Domain

	questBook []*Quest // Задачи формата: квест - набор айтемов для завершения.

	inventory []Item
}

type Quest struct {
	name string
	done bool

	checkUpdate func(ch *Character) bool
}

type Item struct {
	name string

	IsWearable  bool
	IsGetable   bool
	IsInventory bool
	IsUsable    bool
}

type World struct {
	LocationBase []Domain // Список локаций в мире.

	IOList []InteractiveObject // Список интерактивных объектов в мире.

	CommandList CommandBase // База доступных юзеру команд.
}

var CurrWorld *World         // Текущее состояние игрового мира.
var CurrCharacter *Character // Текущее состояние игрока.
