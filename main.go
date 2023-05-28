package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"

	"github.com/julioguillermo/OXO/oxo"
	"github.com/julioguillermo/staticneurogenetic"
)

const (
	PopulationSize = 100
	Survivors      = 10
	MutSize        = 0.1
	MutRate        = 0.1

	MutType   = staticneurogenetic.MultiMutation
	CrossType = staticneurogenetic.DivPointCross

	Activation = staticneurogenetic.Relu

	PopSizeChRandom = false

	ChanBuff   = 100
	RecordSize = 300

	PopMul     = 5
	Iterations = 5
)

type Record struct {
	Fitness     []float64
	Wins        []float64
	BestWins    []float64
	Generations []uint64
	Changes     []uint64
}

var (
	layers = []int{9, 50, 30, 9}
	record = Record{}
	wins   []int
	//connections []net.Conn
)

func onePlay(ia *staticneurogenetic.SNG, p1 int, p2 int) {
	game := oxo.NewOXO()
	moves1 := 0
	moves2 := 0
	for !game.End() {
		var player int
		if game.GetPlayer() == 1 {
			player = p1
			moves1++
		} else {
			player = p2
			moves2++
		}

		_, pos := ia.MaxOutput(player, game.State())
		valid := game.Play(pos)
		if !valid {
			ia.AddFitness(player, -1)
			for !valid {
				pos = (pos + 1) % 9
				valid = game.Play(pos)
			}
		}
	}
	if game.Winner() == 1 {
		ia.AddFitness(p1, 3/float64(moves1))
		ia.AddFitness(p2, -1/float64(moves1))
		wins[p1]++
	} else if game.Winner() == -1 {
		ia.AddFitness(p2, 3/float64(moves2))
		ia.AddFitness(p1, -1/float64(moves2))
		wins[p2]++
	}
}

func oneRandomPlay(ia *staticneurogenetic.SNG, p int, player int8) {
	game := oxo.NewOXO()
	moves := 0
	n_moves := 0
	for !game.End() {
		if game.GetPlayer() == player {
			moves++
			_, pos := ia.MaxOutput(p, game.State())
			valid := game.Play(pos)
			if !valid {
				ia.AddFitness(p, -1)
				for !valid {
					pos = (pos + 1) % 9
					valid = game.Play(pos)
				}
			}
		} else {
			n_moves++
			valid := false
			for !valid {
				valid = game.Play(rand.Intn(9))
			}
		}
	}
	if game.Winner() == player {
		ia.AddFitness(p, 3/float64(moves))
		wins[p]++
	} else {
		ia.AddFitness(p, -1/float64(n_moves))
	}
}

func playOne(ia *staticneurogenetic.SNG, p1, p2 int) {
	if p2 < 0 || p2 >= PopulationSize {
		oneRandomPlay(ia, p1, 1)
	} else if p1 < 0 || p1 >= PopulationSize {
		oneRandomPlay(ia, p2, -1)
	} else {
		onePlay(ia, p1, p2)
	}
}

func playWorker(ia *staticneurogenetic.SNG, wg *sync.WaitGroup, ctl chan [2]int) {
	for ps := range ctl {
		p1 := ps[0]
		p2 := ps[1]
		playOne(ia, p1, p2)
		wg.Done()
	}
}

func playConcurrently(ia *staticneurogenetic.SNG) {
	ps := make(chan [2]int, ChanBuff)
	var wg sync.WaitGroup

	for i := 0; i < ChanBuff; i++ {
		go playWorker(ia, &wg, ps)
	}

	for i := 0; i < Iterations; i++ {
		for p1 := 0; p1 < PopulationSize*PopMul; p1++ {
			for p2 := 0; p2 < PopulationSize*PopMul; p2++ {
				if p1 == p2 {
					continue
				}
				if p1 >= PopulationSize && p2 >= PopulationSize {
					break
				}
				wg.Add(1)
				ps <- [2]int{p1, p2}
			}
		}
	}

	wg.Wait()
	close(ps)
}

/*func play(ia *staticneurogenetic.SNG) {
	for p1 := 0; p1 < PopulationSize*2; p1++ {
		for p2 := 0; p2 < PopulationSize*2; p2++ {
			if p1 == p2 {
				continue
			}
			if p1 >= PopulationSize && p2 >= PopulationSize {
				break
			}
			playOne(ia, p1, p2)
		}
	}
}*/

func playGeneration(ia *staticneurogenetic.SNG) {
	start := time.Now()
	wins = make([]int, PopulationSize)
	ia.ResetFitness()

	playConcurrently(ia)
	//play(ia)

	ia.Sort()
	fitness := ia.GetFitness(0) * 50 / PopulationSize / PopMul / Iterations
	generation := ia.GetGeneration()

	max_win := 0
	for _, v := range wins {
		if max_win < v {
			max_win = v
		}
	}
	win_pro := float64(max_win) * 50 / PopulationSize / PopMul / Iterations
	best_win_pro := float64(wins[ia.GetLastBestIndex()]) * 50 / PopulationSize / PopMul / Iterations

	record.Fitness = append(record.Fitness, fitness)
	record.Generations = append(record.Generations, generation)
	record.Wins = append(record.Wins, win_pro)
	record.BestWins = append(record.BestWins, best_win_pro)
	for len(record.Fitness) > RecordSize {
		record.Fitness = record.Fitness[1:]
		record.Generations = record.Generations[1:]
		record.Wins = record.Wins[1:]
		record.BestWins = record.BestWins[1:]
	}

	best_ind := ia.GetLastBestIndex()
	var best_ind_str string
	if best_ind != 0 {
		best_ind_str = fmt.Sprintf("\033[32m => %d", best_ind)
		record.Changes = append(record.Changes, generation)
		if len(record.Changes) > 0 {
			for record.Changes[0] < record.Generations[0] {
				record.Changes = record.Changes[1:]
			}
		}
	}

	buf, err := json.Marshal(record)
	if err == nil {
		ioutil.WriteFile("record.json", buf, 0777)
		//go sendData(buf)
	}

	ia.NextGeneration()
	ia.SaveAsBin("oxo.bin")

	dTime := time.Since(start)
	fmt.Printf("\033[33mGeneration:\033[0m %d \033[31m[%.4f s]%s\033[0m\n    \033[36mBest:\033[0m     %f%%\033[0m\n    \033[35mBestWins:\033[0m %f%%\n    \033[32mMaxWins:\033[0m  %f%%\n    >> \033[34mNEXT:\033[0m  %s\n", generation, dTime.Seconds(), best_ind_str, fitness, best_win_pro, win_pro, time.Now().Add(dTime).String())
}

/*func sendData(buf []byte) {
	for i, c := range connections {
		_, err := c.Write(buf)
		if err != nil {
			c.Close()
			connections = append(connections[:i], connections[i+1:]...)
		}
	}
}

func initSocketServer() {
	server, err := net.Listen("tcp", "0.0.0.0:9988")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		server.Close()
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			connections = append(connections, conn)
		}
	}
}*/

func main() {
	//go initSocketServer()
	buf, err := ioutil.ReadFile("record.json")
	if err == nil {
		json.Unmarshal(buf, &record)
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	ia, err := staticneurogenetic.LoadFromBin("oxo.bin")
	if err != nil {
		ia = staticneurogenetic.NewSNG(layers, Activation, PopulationSize, Survivors, MutRate, MutSize, MutType, CrossType)
	} else {
		ia.Survivors = Survivors
		ia.MutSize = MutSize
		ia.MutRate = MutRate
		ia.CrossType = CrossType
		ia.MutationType = MutType
		ia.SetPopulationSize(PopulationSize, PopSizeChRandom)
	}

	for {
		playGeneration(ia)
	}
}
