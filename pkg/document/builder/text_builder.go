package builder

// TextBuilder constructs text commands with styling
type TextBuilder struct {
	parent  *DocumentBuilder
	content textContent
	label   *textLabel
	newLine *bool
}

type textContent struct {
	Text  string     `json:"text"`
	Style *textStyle `json:"content_style,omitempty"`
	Align *string    `json:"align,omitempty"`
}

type textLabel struct {
	Text      string     `json:"text,omitempty"`
	Style     *textStyle `json:"label_style,omitempty"`
	Separator *string    `json:"separator,omitempty"`
	Align     *string    `json:"align,omitempty"`
}

type textStyle struct {
	Bold      *bool   `json:"bold,omitempty"`
	Size      *string `json:"size,omitempty"`
	Underline *string `json:"underline,omitempty"`
	Inverse   *bool   `json:"inverse,omitempty"`
	Font      *string `json:"font,omitempty"`
}

type textCommand struct {
	Content textContent `json:"content"`
	Label   *textLabel  `json:"label,omitempty"`
	NewLine *bool       `json:"new_line,omitempty"`
}

func newTextBuilder(parent *DocumentBuilder, content string) *TextBuilder {
	return &TextBuilder{
		parent: parent,
		content: textContent{
			Text: content,
		},
	}
}

// Bold enables bold text
func (tb *TextBuilder) Bold() *TextBuilder {
	if tb.content.Style == nil {
		tb.content.Style = &textStyle{}
	}
	t := true
	tb.content.Style.Bold = &t
	return tb
}

// Size sets text size (e.g., "2x2", "3x3")
func (tb *TextBuilder) Size(size string) *TextBuilder {
	if tb.content.Style == nil {
		tb.content.Style = &textStyle{}
	}
	tb.content.Style.Size = &size
	return tb
}

// Underline sets underline mode ("1pt" or "2pt")
func (tb *TextBuilder) Underline(mode string) *TextBuilder {
	if tb.content.Style == nil {
		tb.content.Style = &textStyle{}
	}
	tb.content.Style.Underline = &mode
	return tb
}

// Inverse enables inverse mode (white on black)
func (tb *TextBuilder) Inverse() *TextBuilder {
	if tb.content.Style == nil {
		tb.content.Style = &textStyle{}
	}
	t := true
	tb.content.Style.Inverse = &t
	return tb
}

// Font sets the font (A or B)
func (tb *TextBuilder) Font(font string) *TextBuilder {
	if tb.content.Style == nil {
		tb.content.Style = &textStyle{}
	}
	tb.content.Style.Font = &font
	return tb
}

// Left aligns text to the left
func (tb *TextBuilder) Left() *TextBuilder {
	align := "left"
	tb.content.Align = &align
	return tb
}

// Center centers the text
func (tb *TextBuilder) Center() *TextBuilder {
	align := "center"
	tb.content.Align = &align
	return tb
}

// Right aligns text to the right
func (tb *TextBuilder) Right() *TextBuilder {
	align := "right"
	tb.content.Align = &align
	return tb
}

// WithLabel adds a label prefix
func (tb *TextBuilder) WithLabel(labelText string) *TextBuilder {
	tb.label = &textLabel{Text: labelText}
	return tb
}

// LabelSeparator sets the separator between label and content
func (tb *TextBuilder) LabelSeparator(sep string) *TextBuilder {
	if tb.label == nil {
		tb.label = &textLabel{}
	}
	tb.label.Separator = &sep
	return tb
}

// NoNewLine prevents automatic line feed
func (tb *TextBuilder) NoNewLine() *TextBuilder {
	f := false
	tb.newLine = &f
	return tb
}

// End finishes the text command and returns to document builder
func (tb *TextBuilder) End() *DocumentBuilder {
	cmd := textCommand{
		Content: tb.content,
		Label:   tb.label,
		NewLine: tb.newLine,
	}
	return tb.parent.addCommand("text", cmd)
}
