package libs

type Color struct {
	Name      string
	ColorCode string
}

var colorMap = map[int]Color{
	0: Color{
		Name:      "Normal",
		ColorCode: "\x1B[0;0m",
	},
	1: Color{
		Name:      "Red",
		ColorCode: "\x1B[0;31m",
	},
	2: Color{
		Name:      "Green",
		ColorCode: "\x1B[0;32m",
	},
	3: Color{
		Name:      "Yellow",
		ColorCode: "\x1B[0;33m",
	},
	4: Color{
		Name:      "Blue",
		ColorCode: "\x1B[0;34m",
	},
	5: Color{
		Name:      "Magenta",
		ColorCode: "\x1B[0;35m",
	},
	6: Color{
		Name:      "Cyan",
		ColorCode: "\x1B[0;36m",
	},
	7: Color{
		Name:      "White",
		ColorCode: "\x1B[0;37m",
	},

	//Bold
	8: Color{
		Name:      "Bold Red",
		ColorCode: "\x1B[1;31m",
	},
	9: Color{
		Name:      "Bold Green",
		ColorCode: "\x1B[1;32m",
	},
	10: Color{
		Name:      "Bold Yellow",
		ColorCode: "\x1B[1;33m",
	},
	11: Color{
		Name:      "Bold Blue",
		ColorCode: "\x1B[1;34m",
	},
	12: Color{
		Name:      "Bold Magenta",
		ColorCode: "\x1B[1;35m",
	},
	13: Color{
		Name:      "Bold Cyan",
		ColorCode: "\x1B[1;36m",
	},
	14: Color{
		Name:      "Bold White",
		ColorCode: "\x1B[1;37m",
	},

	//Background
	15: Color{
		Name:      "Bold Red",
		ColorCode: "\x1B[41m",
	},
	16: Color{
		Name:      "Bold Green",
		ColorCode: "\x1B[42m",
	},
	17: Color{
		Name:      "Bold Yellow",
		ColorCode: "\x1B[43m",
	},
	18: Color{
		Name:      "Bold Blue",
		ColorCode: "\x1B[44m",
	},
	19: Color{
		Name:      "Bold Magenta",
		ColorCode: "\x1B[45m",
	},
	20: Color{
		Name:      "Bold Cyan",
		ColorCode: "\x1B[46m",
	},
	21: Color{
		Name:      "Bold White",
		ColorCode: "\x1B[47m",
	},
}

func ProcessString(s string) string {

	outString := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			colorID := getColorCode(s[i+1])
			if colorID >= 0 {
				outString += colorMap[colorID].ColorCode
				i++
			}
		} else {
			outString += string(s[i])
		}
	}

	return outString
}

func getColorCode(code byte) int {
	switch code {
	/* Normal colours */
	case 'k':
		return 25
		break /* Black */
	case 'r':
		return 1
		break /* Red */
	case 'g':
		return 2
		break /* Green */
	case 'y':
		return 3
		break /* Yellow */
	case 'b':
		return 4
		break /* Blue */
	case 'm':
		return 5
		break /* Magenta */
	case 'c':
		return 6
		break /* Cyan */
	case 'w':
		return 7
		break /* White */

	/* Bold colours */
	case 'K':
		return 29
		break
	/* Bold black (Just for completeness) */
	case 'R':
		return 8
		break /* Bold red */
	case 'G':
		return 9
		break /* Bold green */
	case 'Y':
		return 10
		break /* Bold yellow */
	case 'B':
		return 11
		break /* Bold blue */
	case 'M':
		return 12
		break /* Bold magenta */
	case 'C':
		return 13
		break /* Bold cyan */
	case 'W':
		return 14
		break /* Bold white */

	/* Background colours */
	case '0':
		return 24
		break /* Black background */
	case '1':
		return 15
		break /* Red background */
	case '2':
		return 16
		break /* Green background */
	case '3':
		return 17
		break /* Yellow background */
	case '4':
		return 18
		break /* Blue background */
	case '5':
		return 19
		break /* Magenta background */
	case '6':
		return 20
		break /* Cyan background */
	case '7':
		return 21
		break /* White background */

	/* Misc characters */
	case '&':
		return 22
		break /* The & character */
	case '\\':
		return 23
		break /* The \ character */

	/* Special codes */
	case 'n':
		return 0
		break /* Normal */
	case 'f':
		return 26
		break /* Flash */
	case 'v':
		return 27
		break /* Reverse video */
	case 'u':
		return 28
		break /* Underline (Only for mono screens) */

	default:
		return -1
		break
	}
	return -1
}
