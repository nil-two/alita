package alita

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	DIGIT_ONLY             = regexp.MustCompile(`^\d+$`)
	COLON_SEPARATED_DIGITS = regexp.MustCompile(`^(\d+):(\d+)$`)
)

type Margin struct {
	left  int
	right int
}

func NewMargin() *Margin {
	return &Margin{
		left:  1,
		right: 1,
	}
}

func (m *Margin) String() string {
	return fmt.Sprint(*m)
}

func (m *Margin) Set(format string) error {
	if DIGIT_ONLY.MatchString(format) {
		n, err := strconv.Atoi(format)
		if err != nil {
			return err
		}
		m.left, m.right = n, n
		return nil
	}
	if COLON_SEPARATED_DIGITS.MatchString(format) {
		a := COLON_SEPARATED_DIGITS.FindAllStringSubmatch(format, -1)
		left, err := strconv.Atoi(a[0][1])
		if err != nil {
			return err
		}
		right, err := strconv.Atoi(a[0][2])
		if err != nil {
			return err
		}
		m.left, m.right = left, right
		return nil
	}
	return fmt.Errorf("margin: invalid format: %s", format)
}

func (m *Margin) Join(a []string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := m.left * (len(a) / 2)
	if len(a)%2 != 0 {
		n += m.left * (len(a) / 2)
	}
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	lm, rm := strings.Repeat(" ", m.left), strings.Repeat(" ", m.right)

	b := make([]byte, n)
	bp := copy(b, a[0])
	for i := 2; i <= len(a); i += 2 {
		bp += copy(b[bp:], lm)
		bp += copy(b[bp:], a[i-1])
		if i != len(a) {
			bp += copy(b[bp:], rm)
			bp += copy(b[bp:], a[i])
		}
	}
	return string(b)
}
