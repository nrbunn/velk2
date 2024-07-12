package velk

import "fmt"

type Rnum struct {
	ZoneId int
	RoomId int
}

func (r Rnum) ToString() string {
	return fmt.Sprintf("%d-%d", r.ZoneId, r.RoomId)
}

func ParseRnum(str string) (Rnum, error) {
	r := &Rnum{}
	_, err := fmt.Sscanf(str, "%d-%d", &r.ZoneId, &r.RoomId)
	return *r, err
}

func getOppositeDirection(dir string) string {
	switch dir {
	case "north":
		return "south"
	case "east":
		return "west"
	case "south":
		return "north"
	case "west":
		return "east"
	case "up":
		return "down"
	case "down":
		return "up"
	default:
		return dir
	}
}
