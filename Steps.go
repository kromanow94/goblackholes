package goblackholes

import (
	//"fmt"
	"math"
	"sync"
	"time"
)

type Agent struct {
	x, y, fitness  float64
	border         Border
	typeOfFunction TypeOfFunction_S
	times          uint64
}
type BestAgent struct {
	x, y, fitness, eventHorizon float64
	step                        uint64
	mutex                       sync.Mutex
}
type Border struct {
	X1, Y1, X2, Y2                                                     float64
	HorizontalLength, VerticalLength, HorizontalCenter, VerticalCenter float64
}

func (b *Border) SetUp() {
	b.HorizontalLength = b.X2 - b.X1
	b.HorizontalCenter = (b.X2 + b.X1) / 2
	b.VerticalLength = b.Y2 - b.Y1
	b.VerticalCenter = (b.Y2 + b.Y1) / 2
}

type TypeOfFunction_S struct {
	StringEvaluation string
	Rastrigin        bool
	Rosenbrock       bool
	Easom            bool
	McCormick        bool
}

func countFitness(output chan *Agent, input *Agent) {
	if input.typeOfFunction.Rastrigin == true {
		input.fitness = 20 + math.Pow(input.x, 2) + math.Pow(input.y, 2) - 10*(math.Cos(2*math.Pi*input.x)+math.Cos(2*math.Pi*input.y))
	} else if input.typeOfFunction.Rosenbrock == true {
		input.fitness = math.Pow(float64(1)-input.x, 2) + 100*math.Pow(input.y-math.Pow(input.x, 2), 2)
	} else if input.typeOfFunction.Easom == true {
		//input.x = 3.14
		//input.y = 3.14
		input.fitness = -math.Cos(input.x) * math.Cos(input.y) * math.Exp(-(math.Pow(input.x-math.Pi, 2) + math.Pow(input.y-math.Pi, 2)))
		//fmt.Println(input.fitness)
		//time.Sleep(10*time.Second)
	} else if input.typeOfFunction.McCormick == true {
		input.fitness = math.Sin(input.x+input.y) + math.Pow(input.x-input.y, 2) - 1.5*input.x + 2.5*input.y + 1
	} else if input.typeOfFunction.StringEvaluation != "" {
		channelString := make(chan string, 1)
		channelEvaluate := make(chan float64, 1)
		go ParseFunction(input.typeOfFunction.StringEvaluation, channelString, input.x, input.y)
		go EvaluateFunction(<-channelString, channelEvaluate)
		input.fitness = <-channelEvaluate
	}
	output <- input
}

func getBest(output chan *Agent, input *Agent) {
	if bestAgent.fitness > input.fitness {
		bestAgent.mutex.Lock()
		bestAgent.x = input.x
		bestAgent.y = input.y
		bestAgent.fitness = input.fitness
		bestAgent.step = input.times
		bestAgent.mutex.Unlock()
		// ToDo channel
		input.newPosition()
	}
	input.times += 1
	output <- input
	i++
}
func (agent *Agent) newPosition() {
	go func() {
		agent.x = <-randomBuffer*agent.border.HorizontalLength - agent.border.HorizontalLength/2 + agent.border.HorizontalCenter
	}()
	go func() {
		agent.y = <-randomBuffer*agent.border.VerticalLength - agent.border.VerticalLength/2 + agent.border.VerticalCenter
	}()
}

func move(output chan *Agent, input *Agent) {
	time.Sleep(time.Duration(slowmotion) * time.Millisecond)
	input.x += <-randomBuffer * (bestAgent.x - input.x)
	input.y += <-randomBuffer * (bestAgent.y - input.y)
	output <- input
}

func countEventHorizon() {
	var fSum float64
	for i := 0; i < agentAmount; i++ {
		fSum = agentList[i].fitness
	}
	bestAgent.mutex.Lock()
	bestAgent.eventHorizon = bestAgent.fitness / fSum
	bestAgent.mutex.Unlock()
	//if math.IsNaN(bestAgent.eventHorizon) || bestAgent.eventHorizon == 0 {
	//	maxAccuracy <- true
	//}
}

func eventHorizon(output chan *Agent, input *Agent) {
	var (
		xLen, yLen, len chan float64 = make(chan float64, 1), make(chan float64, 1), make(chan float64, 1)
	)
	go func() { xLen <- math.Pow(bestAgent.x-input.y, 2.0) }()
	go func() { yLen <- math.Pow(bestAgent.y-input.y, 2.0) }()
	go func() { len <- math.Pow(<-xLen + <-yLen, 0.5) }()
	if bestAgent.eventHorizon >= <-len {
		input.newPosition()
	}
	output <- input
}
