package overlay

import (
    "strings"
    "regexp"

    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/x/ansi"
)

var ansiStyleRegexp = regexp.MustCompile(`\x1b[[\d;]*m`)

const (
    DefaultWidth = 60
    DefaultHeight = 20
)

type dimensions struct {
    width  int
    height int
}

// The Bubble Tea model for this overlay element
type ModelBase struct {

    // Store errors here
    Err         error

    // The contents of the overlay background
    bg          string

    // Whether the overlay should be shown or not
    Show        bool

    // The current styling to use
    style       gloss.Style

    // The width of the overlay
    width       int

    // The height of the overlay
    height      int

    // The column the cursor's on
    col         int

    // The row the cursor's on
    row         int

    // The dimensions of the whole window
    dimensions  dimensions

}

// Creates a new model with default settings
func NewBase() ModelBase {
    style := gloss.NewStyle().
        Border(gloss.RoundedBorder()).
        BorderStyle(gloss.DoubleBorder()).
        Bold(true).
        Padding(2)


    m := ModelBase{
        style: style,
        col: 0,
        row: 0,
        Show: false,
    }

    m.SetWidth(DefaultWidth)
    m.SetHeight(DefaultHeight)

    return m
}

// Sets the width of the overlay
func (m *ModelBase) SetWidth(width int) {
    m.width = width
    m.style = m.style.Width(m.width)
}

// Sets the height of the overlay
func (m *ModelBase) SetHeight(height int) {
    m.height = height
    m.style = m.style.Height(m.height)
}

// Sets the contents of the overlay background
func (m *ModelBase) SetBackground(content string) {
    m.bg = content
}

func (m ModelBase) View(content string) string {

    row := (m.dimensions.height - m.height) / 2
    col := (m.dimensions.width - m.width) / 2

	overlay := m.style.Render(content)

	bgLines := strings.Split(m.bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, overlayLine := range overlayLines {
		bgLine := bgLines[i+row] // TODO: index handling
		if len(bgLine) < col {
			bgLine += strings.Repeat(" ", col-len(bgLine)) // add padding
		}

		bgLeft := ansi.Truncate(bgLine, col, "")
		bgRight := truncateLeft(bgLine, col+ansi.StringWidth(overlayLine))

		bgLines[i+row] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(bgLines, "\n")
}

func (m *ModelBase) SetDimensions(width int, height int) {
    m.dimensions.width = width
    m.dimensions.height = height
}

func truncateLeft(line string, padding int) string {
	if strings.Contains(line, "\n") {
		panic("line must not contain newline")
	}

	// NOTE: line has no newline, so [strings.Join] after [strings.Split] is safe.
	wrapped := strings.Split(ansi.Hardwrap(line, padding, true), "\n")
	if len(wrapped) == 1 {
		return ""
	}

	var ansiStyle string
	ansiStyles := ansiStyleRegexp.FindAllString(wrapped[0], -1)
	if l := len(ansiStyles); l > 0 {
		ansiStyle = ansiStyles[l-1]
	}

	return ansiStyle + strings.Join(wrapped[1:], "")
}
