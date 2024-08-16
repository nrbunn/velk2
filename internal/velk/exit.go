package velk

import "fmt"

type Vnum struct {
	ZoneId    int
	VirtualId int
}

func (r Vnum) ToString() string {
	return fmt.Sprintf("%d-%d", r.ZoneId, r.VirtualId)
}

func ParseVnum(str string) (Vnum, error) {
	r := &Vnum{}
	_, err := fmt.Sscanf(str, "%d-%d", &r.ZoneId, &r.VirtualId)
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
