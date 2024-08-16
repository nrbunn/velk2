package velk

type MODE int

const (
	Inactive  MODE = 0
	ReditMode MODE = 1
	OeditMode MODE = 2
	MeditMode MODE = 3
	ZeditMode MODE = 4
)

type Olc struct {
	Mode  MODE
	Redit *ReditData
	Medit *MeditData
}

func HandleOlc(player *Player) {
	switch player.Olc.Mode {
	case ReditMode:
		handleRedit(player)
	case MeditMode:
		handleMedit(player)
	}
}
