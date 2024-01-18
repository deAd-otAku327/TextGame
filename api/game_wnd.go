package api

import (
	"gameProject/action"
	"gameProject/world"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func StartGame() {
	img, _ := fyne.LoadResourceFromPath("icon.ico")

	a := app.New()
	game := a.NewWindow("Game")

	game.CenterOnScreen()
	game.Resize(fyne.Size{1000, 400})
	game.SetIcon(img)

	сomm := widget.NewEntry()

	gameField := widget.NewLabelWithStyle("ДОБРО ПОЖАЛОВАТЬ! ВВЕДИТЕ КОМАНДУ...\n"+`ДЛЯ ПРОСМОТРА ДОСТУПНЫХ КОМАНД ВВЕДИТЕ "помощь"`,
		fyne.TextAlign(fyne.TextAlignCenter), fyne.TextStyle{Bold: true})

	execButton := widget.NewButton("Выполнить", func() {
		answer := HandleCommand(сomm.Text)

		if сomm.Text != "помощь" {
			answer = strings.ToUpper(answer)
		}
		gameField.SetText(answer)

		if strings.ToLower(answer) == action.StopGame {
			сomm.Hide()
		} else {
			сomm.SetText("")
		}
	})

	game.SetContent(container.NewVBox(
		layout.NewSpacer(),
		gameField,
		layout.NewSpacer(),
		сomm,
		execButton,
	))
	game.ShowAndRun()
}

func HandleCommand(command string) string {

	parsedCom := strings.Fields(command)

	answer := executor(parsedCom)

	return answer
}

func executor(parsedCom []string) string {

	if len(parsedCom) == 0 {
		return "поле ввода команды пустое"
	}
	action, ok := world.CurrWorld.CommandList[parsedCom[0]] // Первый элемент всегда команда, чекаем наличие в базе команд.

	if !ok {
		return "неизвестная команда"
	}

	for _, obj := range world.CurrWorld.IOList { // Проверка команды на возможность выполнения.
		if obj.CheckBanned(parsedCom) {
			return obj.GetStateInfo()
		}
	}

	return action(parsedCom[1:])
}
