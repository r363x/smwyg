package status

type statusDetails struct {
    Left   string
    Center string
    Right  string
}

type Model struct {
    Content statusDetails
}

type Msg struct {
    Section int
    Message string
}

const (
    SecLeft  = iota
    SecCenter
    SecRight
)


