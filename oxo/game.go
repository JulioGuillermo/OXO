package oxo

type OXO struct {
	end    bool
	winner int8

	grid   [9]int8
	player int8
}

func NewOXO() *OXO {
	oxo := &OXO{}
	oxo.Reset()
	return oxo
}

func (p *OXO) Reset() {
	p.end = false
	p.winner = 0
	for i := range p.grid {
		p.grid[i] = 0
	}
	p.player = 1
}

func (p *OXO) GetPlayer() int8 {
	return p.player
}

func (p *OXO) End() bool {
	return p.end
}

func (p *OXO) Winner() int8 {
	return p.winner
}

func (p *OXO) getSums() int8 {
	// Rows
	s := p.grid[0] + p.grid[1] + p.grid[2]
	if s == 3 {
		return 1
	}
	if s == -3 {
		return -1
	}
	s = p.grid[3] + p.grid[4] + p.grid[5]
	if s == 3 {
		return 2
	}
	if s == -3 {
		return -2
	}
	s = p.grid[6] + p.grid[7] + p.grid[8]
	if s == 3 {
		return 3
	}
	if s == -3 {
		return -3
	}

	// Cols
	s = p.grid[0] + p.grid[3] + p.grid[6]
	if s == 3 {
		return 4
	}
	if s == -3 {
		return -4
	}
	s = p.grid[1] + p.grid[4] + p.grid[7]
	if s == 3 {
		return 5
	}
	if s == -3 {
		return -5
	}
	s = p.grid[2] + p.grid[5] + p.grid[8]
	if s == 3 {
		return 6
	}
	if s == -3 {
		return -6
	}

	// Diags
	s = p.grid[0] + p.grid[4] + p.grid[8]
	if s == 3 {
		return 7
	}
	if s == -3 {
		return -7
	}
	s = p.grid[2] + p.grid[4] + p.grid[6]
	if s == 3 {
		return 8
	}
	if s == -3 {
		return -8
	}

	return 0
}

func (p *OXO) checkEnd() bool {
	for _, e := range p.grid {
		if e == 0 {
			return false
		}
	}
	return true
}

func (p *OXO) CheckWin() {
	s := p.getSums()
	if s > 0 {
		p.end = true
		p.winner = 1
		return
	}
	if s < 0 {
		p.end = true
		p.winner = -1
		return
	}
	if p.checkEnd() {
		p.end = true
		p.winner = 0
	}
}

func (p *OXO) Play(pos int) bool {
	if p.end {
		return true
	}

	if p.grid[pos] != 0 {
		return false
	}

	p.grid[pos] = p.player
	p.player = -p.player

	p.CheckWin()
	return true
}

func (p *OXO) PlayXY(x, y int) bool {
	return p.Play(y*3 + x)
}

func (p *OXO) String() string {
	get_char := func(i int8) string {
		if i == 1 {
			return "x"
		}
		if i == -1 {
			return "o"
		}
		return " "
	}

	str := get_char(p.grid[0]) + "|" + get_char(p.grid[1]) + "|" + get_char(p.grid[2]) + "\n"
	str += "-----\n"
	str += get_char(p.grid[3]) + "|" + get_char(p.grid[4]) + "|" + get_char(p.grid[5]) + "\n"
	str += "-----\n"
	str += get_char(p.grid[6]) + "|" + get_char(p.grid[7]) + "|" + get_char(p.grid[8]) + "\n"

	return str
}

func (p *OXO) State() []float64 {
	state := make([]float64, 9)
	for i, v := range p.grid {
		state[i] = float64(v) * float64(p.player)
	}
	return state
}

func (p *OXO) Bytes() []byte {
	state := make([]byte, 9)
	for i, v := range p.grid {
		state[i] = byte(v)
	}
	return state
}
