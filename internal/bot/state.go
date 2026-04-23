package bot

type Step int

const (
	StepFrom      Step = iota // ждём "откуда"
	StepTo                    // ждём "куда"
	StepDeparture             // ждём дату вылета
	StepReturn                // ждём дату возврата
	StepTransfers             // ждём кнопку
)

type UserState struct {
	Step      Step
	From      string
	To        string
	Departure string
	Return    string
}

var states = map[int64]*UserState{}
