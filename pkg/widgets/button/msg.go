package button


type MsgType int

const (
    ButtonPressed MsgType = iota
    ButtonReleased
    ButtonSelected
    ButtonUnselected
)

type Msg struct {
    Type MsgType
}

