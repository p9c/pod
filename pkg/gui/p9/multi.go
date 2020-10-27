package p9

import (
	l "gioui.org/layout"
	"github.com/urfave/cli"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Multi struct {
	*Theme
	lines            *cli.StringSlice
	clickables       []*Clickable
	buttons          []*ButtonLayout
	input            *Input
	inputLocation    int
	addClickable     *Clickable
	removeClickables []*Clickable
	removeButtons    []*IconButton
}

func (th *Theme) Multiline(txt *cli.StringSlice, borderColorFocused, borderColorUnfocused string,
	size int, handle func(txt []string)) (m *Multi) {
	if handle == nil {
		handle = func(txt []string) {
			Debug(txt)
		}
	}
	addClickable := th.Clickable()
	m = &Multi{
		Theme:         th,
		lines:         txt,
		inputLocation: -1,
		addClickable:  addClickable,
	}
	handleChange := func(txt string) {
		Debug("handleChange", m.inputLocation)
		(*m.lines)[m.inputLocation] = txt
		// after submit clear the editor
		m.inputLocation = -1
		handle(*m.lines)
	}
	m.input = th.Input("", borderColorFocused, borderColorUnfocused, size, handleChange)
	m.clickables = append(m.clickables, (*Clickable)(nil))
	m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	m.removeClickables = append(m.removeClickables, (*Clickable)(nil))
	m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	for i := range *m.lines {
		Debug("making clickables")
		x := i
		clickable := m.Theme.Clickable().SetClick(
			func() {
				m.inputLocation = x
				Debug("button clicked", x, m.inputLocation)
			})
		if len(*m.lines) > len(m.clickables) {
			m.clickables = append(m.clickables, clickable)
		} else {
			m.clickables[i] = clickable
		}
		Debug("making button")
		btn := m.Theme.ButtonLayout(clickable).CornerRadius(0).Background("transparent").
			Embed(
				m.Theme.Flex().
					Rigid(
						m.Theme.Fill("Primary",
							m.Theme.Inset(0.5,
								m.Theme.Body2((*m.lines)[i]).Color("DocText").Fn,
							).Fn,
						).Fn,
					).Fn,
			)
		if len(*m.lines) > len(m.buttons) {
			m.buttons = append(m.buttons, btn)
		} else {
			m.buttons[i] = btn
		}
		Debug("making clickables")
		removeClickable := m.Theme.Clickable()
		if len(*m.lines) > len(m.removeClickables) {
			m.removeClickables = append(m.removeClickables, removeClickable)
		} else {
			m.removeClickables[i] = removeClickable
		}
		Debug("making remove button")
		removeBtn := m.Theme.IconButton(removeClickable).
			Icon(
				m.Theme.Icon().Scale(1.5).Color("DocText").Src(icons.ActionDelete),
			).
			Background("DocBg").
			SetClick(func() {
				Debug("remove button", i, "clicked")
				m.inputLocation = -1
				if i == len(*m.lines)-1 {
					*m.lines = (*m.lines)[:len(*m.lines)-1]
					m.clickables = m.clickables[:len(m.clickables)-1]
					m.buttons = m.buttons[:len(m.buttons)-1]
					m.removeClickables = m.removeClickables[:len(m.removeClickables)-1]
					m.removeButtons = m.removeButtons[:len(m.removeButtons)-1]
				} else {
					*m.lines = append((*m.lines)[:i], (*m.lines)[i+1:]...)
					m.clickables = append(m.clickables[:i], m.clickables[i+1:]...)
					m.buttons = append(m.buttons[:i], m.buttons[i+1:]...)
					m.removeClickables = append(m.removeClickables[:i], m.removeClickables[i+1:]...)
					m.removeButtons = append(m.removeButtons[:i], m.removeButtons[i+1:]...)
				}
			})
		if len(*m.lines) > len(m.removeButtons) {
			m.removeButtons = append(m.removeButtons, removeBtn)
		} else {
			m.removeButtons[i] = removeBtn
		}
	}
	return m
}

func (m *Multi) UpdateWidgets() *Multi {
	if len(m.clickables) < len(*m.lines) {
		Debug("allocating new clickables")
		m.clickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.buttons) < len(*m.lines) {
		Debug("allocating new buttons")
		m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	}
	if len(m.removeClickables) < len(*m.lines) {
		Debug("allocating new removeClickables")
		m.removeClickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.removeButtons) < len(*m.lines) {
		Debug("allocating new removeButtons")
		m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	}
	return m
}

func (m *Multi) PopulateWidgets() *Multi {
	added := false
	for i := range *m.lines {
		if m.clickables[i] == nil {
			added = true
			Debug("making clickables", i)
			x := i
			m.clickables[i] = m.Theme.Clickable().SetClick(
				func() {
					Debug("clicked", x, m.inputLocation)
					m.inputLocation = x
					m.input.Editor().SetText((*m.lines)[x])
					m.input.Editor().Focus()
				})
		}
		// m.clickables[i]
		if m.buttons[i] == nil {
			added = true
			btn := m.Theme.ButtonLayout(m.clickables[i]).CornerRadius(0).Background("transparent")
			m.buttons[i] = btn
		}
		m.buttons[i].Embed(
			m.Theme.Flex().
				Rigid(
					m.Theme.Inset(0.5,
						m.Theme.Body2((*m.lines)[i]).Color("DocText").Fn,
					).Fn,
				).Fn,
		)
		if m.removeClickables[i] == nil {
			added = true
			removeClickable := m.Theme.Clickable()
			m.removeClickables[i] = removeClickable
		}
		if m.removeButtons[i] == nil {
			added = true
			Debug("making remove button", i)
			x := i
			m.removeButtons[i] = m.Theme.IconButton(m.removeClickables[i].
				SetClick(func() {
					Debug("remove button", x, "clicked", len(*m.lines))
					m.inputLocation = -1
					if len(*m.lines)-1 == i {
						*m.lines = (*m.lines)[:len(*m.lines)-1]
					} else {
						*m.lines = append((*m.lines)[:x], (*m.lines)[x+1:]...)
					}
				})).
				Icon(
					m.Theme.Icon().Scale(1.5).Color("DocText").Src(icons.ActionDelete),
				).
				Background("DocBg")
		}
	}
	if added {
		Debug("clearing editor")
		m.input.editor.SetText("")
	}
	return m
}

func (m *Multi) Fn(gtx l.Context) l.Dimensions {
	m.UpdateWidgets()
	m.PopulateWidgets()
	addButton := m.Theme.IconButton(m.addClickable).Icon(
		m.Theme.Icon().Scale(1.5).Color("Primary").Src(icons.ContentAdd),
	)
	var widgets []l.Widget
	if m.inputLocation > 0 && m.inputLocation < len(*m.lines) {
		m.input.Editor().SetText((*m.lines)[m.inputLocation])
	}
	for i := range *m.lines {
		if m.buttons[i] == nil {
			x := i
			btn := m.Theme.ButtonLayout(m.clickables[i].SetClick(
				func() {
					Debug("button pressed", (*m.lines)[x], x, m.inputLocation)
					m.inputLocation = x
					m.input.editor.SetText((*m.lines)[x])
					m.input.editor.Focus()
				})).CornerRadius(0).Background("transparent").
				Embed(
					m.Theme.Flex().
						Rigid(
							m.Theme.Inset(0.5,
								m.Theme.Body2((*m.lines)[x]).Color("DocText").Fn,
							).Fn,
						).Fn,
				)
			m.buttons[i] = btn
		}
		if i == m.inputLocation {
			m.input.Editor().SetText((*m.lines)[i])
			input := m.Flex().
				Rigid(
					m.removeButtons[i].Fn,
				).
				Rigid(
					m.input.Fn,
				).
				Fn
			widgets = append(widgets, input)
		} else {
			x := i
			m.clickables[i].SetClick(
				func() {
					Debug("setting", x, m.inputLocation)
					m.inputLocation = x
					m.input.editor.SetText((*m.lines)[x])
					m.input.editor.Focus()
				})
			button := m.Flex().
				Rigid(
					m.removeButtons[i].Fn,
				).
				Rigid(
					m.buttons[i].Fn,
				).
				Fn
			widgets = append(widgets, button)
		}
	}
	widgets = append(widgets, addButton.SetClick(func() {
		Debug("clicked add")
		m.inputLocation = len(*m.lines)
		*m.lines = append(*m.lines, "")
		Debugs([]string(*m.lines))
		m.UpdateWidgets()
		m.PopulateWidgets()
		m.input.editor.SetText("")
		m.input.editor.Focus()
	}).Background("DocBg").Fn)
	m.UpdateWidgets()
	m.PopulateWidgets()
	// Debug(m.inputLocation)
	// if m.inputLocation > 0 {
	// 	m.input.Editor().Focus()
	// }
	out := m.Theme.VFlex()
	for i := range widgets {
		out.Rigid(widgets[i])
	}
	return out.Fn(gtx)
}
