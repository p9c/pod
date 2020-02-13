package components

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"time"
)

var (
	testLabel         = "testtopLabel"
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis: layout.Vertical,
		ScrollToEnd:true,
	}
)

func DuoUIconsoleWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		layout.Flex{}.Layout(duo.DuoUIcontext,
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround,
					}.Layout(duo.DuoUIcontext,
						layout.Flexed(1, func() {
							consoleOutputList.Layout(duo.DuoUIcontext, len(rc.CommandsHistory.Commands), func(i int) {
								t := rc.CommandsHistory.Commands[i]
								layout.Flex{
									Alignment: layout.End,
								}.Layout(duo.DuoUIcontext,
									layout.Rigid(func() {
										sat := duo.DuoUItheme.Body1("ds://" + t.ComID)
										sat.Font.Size = unit.Dp(16)
										sat.Layout(duo.DuoUIcontext)
									}),
								)
							})
						}),
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(8)).Layout(duo.DuoUIcontext, func() {
								e := duo.DuoUItheme.DuoUIeditor("Run command", "Run txt")
								e.Font.Style = text.Regular
								e.Font.Size = unit.Dp(16)
								e.Layout(duo.DuoUIcontext, consoleInputField)
								for _, e := range consoleInputField.Events(duo.DuoUIcontext) {
									if e, ok := e.(widget.SubmitEvent); ok {
										rc.CommandsHistory.Commands = append(rc.CommandsHistory.Commands, models.DuoUIcommand{
											ComID: e.Text,
											Time:  time.Time{},
										})
										consoleInputField.SetText("")
									}
								}
							})
						}))
				})
			}),
		)
	}
}