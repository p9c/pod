package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	buttonCornerOne = new(widget.Button)
	buttonCornerTwo = new(widget.Button)
	buttonSettings  = new(widget.Button)
	buttonNetwork   = new(widget.Button)
	buttonBlocks    = new(widget.Button)
	buttonConsole   = new(widget.Button)
	buttonHelp      = new(widget.Button)
	cornerNav       = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func DuoUIfooter(duo *models.DuoUI, rc *rcd.RcVar) func() {
	return func() {
		cs := duo.DuoUIcontext.Constraints
		helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 64, "ff303030", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		var (
			width             float32 = 48
			height            float32 = 48
			iconSize          int     = 32
			paddingVertical   float32 = 8
			paddingHorizontal float32 = 8
		)
		settingsIcon, _ := parallel.NewDuoUIicon(icons.ActionSettings)
		blocksIcon, _ := parallel.NewDuoUIicon(icons.ActionExplore)
		networkIcon, _ := parallel.NewDuoUIicon(icons.ActionFingerprint)
		consoleIcon, _ := parallel.NewDuoUIicon(icons.ActionInput)
		helpIcon, _ := parallel.NewDuoUIicon(icons.NavigationArrowDropDown)
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
					cornerButtons := []func(){
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
								var networkMeniItem parallel.DuoUIbutton
								networkMeniItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
								for buttonCornerOne.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Network"
								}
								networkMeniItem.Layout(duo.DuoUIcontext, buttonCornerOne)
							})
						},
						func() {
							var settingsMenuItem parallel.DuoUIbutton
							settingsMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

							for buttonCornerTwo.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Settings"
							}
							settingsMenuItem.Layout(duo.DuoUIcontext, buttonCornerTwo)
						},
					}
					cornerNav.Layout(duo.DuoUIcontext, len(cornerButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, cornerButtons[i])
					})
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
					navButtons := []func(){
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
								var networkMeniItem parallel.DuoUIbutton
								networkMeniItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, networkIcon)
								for buttonNetwork.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Network"
								}
								networkMeniItem.Layout(duo.DuoUIcontext, buttonNetwork)
							})
						},
						func() {
							var blocksMenuItem parallel.DuoUIbutton
							blocksMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, blocksIcon)
							for buttonBlocks.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Explorer"
							}
							blocksMenuItem.Layout(duo.DuoUIcontext, buttonBlocks)
						},
						func() {
							var helpMenuItem parallel.DuoUIbutton
							helpMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, helpIcon)
							for buttonHelp.Clicked(duo.DuoUIcontext) {
								rc.IsNotificationRun = true
							}
							helpMenuItem.Layout(duo.DuoUIcontext, buttonHelp)
						},
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
								var consoleMenuItem parallel.DuoUIbutton
								consoleMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, consoleIcon)
								for buttonConsole.Clicked(duo.DuoUIcontext) {
									duo.CurrentPage = "Console"
								}
								consoleMenuItem.Layout(duo.DuoUIcontext, buttonConsole)
							})
						},
						func() {
							var settingsMenuItem parallel.DuoUIbutton
							settingsMenuItem = duo.DuoUItheme.DuoUIbutton("", "", "ff303030", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal, settingsIcon)

							for buttonSettings.Clicked(duo.DuoUIcontext) {
								duo.CurrentPage = "Settings"
							}
							settingsMenuItem.Layout(duo.DuoUIcontext, buttonSettings)
						},
					}
					footerNav.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
					})
				})
			}),
		)
	}
}
