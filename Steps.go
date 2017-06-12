package goblackholes

import (
	"math"
	//"sync/atomic"
	"time"
	//"fmt"
)

func countFitness(output chan *Agent, input *Agent) {
	if input.TypeOfFunction.Rastrigin == true {
		input.Fitness = 20 + math.Pow(input.X, 2) + math.Pow(input.Y, 2) - 10*(math.Cos(2*math.Pi*input.X)+math.Cos(2*math.Pi*input.Y))
	} else if input.TypeOfFunction.Rosenbrock == true {
		input.Fitness = math.Pow(float64(1)-input.X, 2) + 100*math.Pow(input.Y-math.Pow(input.X, 2), 2)
	} else if input.TypeOfFunction.Easom == true {
		input.Fitness = -math.Cos(input.X) * math.Cos(input.Y) * math.Exp(-(math.Pow(input.X-math.Pi, 2) + math.Pow(input.Y-math.Pi, 2)))
	} else if input.TypeOfFunction.McCormick == true {
		input.Fitness = math.Sin(input.X+input.Y) + math.Pow(input.X-input.Y, 2) - 1.5*input.X + 2.5*input.Y + 1
	} else if input.TypeOfFunction.StringEvaluation != "" {
		channelString := make(chan string, 1)
		channelEvaluate := make(chan float64, 1)
		go parseFunction(input.TypeOfFunction.StringEvaluation, channelString, input.X, input.Y)
		go evaluateFunction(<-channelString, channelEvaluate)
		input.Fitness = <-channelEvaluate
	}
	output <- input
}

func getBest(output chan *Agent, outputChan chan *Agent, input chan *Agent) {
	for {
		for i := 0; i < agentAmount; i++ {
			agentList[i] = <-input
		}
		for _, agent := range agentList {
			if bestAgent.fitness > agent.Fitness {
				bestAgent.x = agent.X
				bestAgent.y = agent.Y
				bestAgent.fitness = agent.Fitness
				bestAgent.step = agent.Times
				outputChan <- bestAgent.Convert()
				agent.newPosition()
			}
			agent.Times += 1
			output <- agent
			outputChan <- agent
		}
	}
}

func move(output chan *Agent, input *Agent) {
	time.Sleep(time.Duration(slowmotion) * time.Millisecond)
	input.X += <-randomBuffer * (bestAgent.x - input.X)
	input.Y += <-randomBuffer * (bestAgent.y - input.Y)
	output <- input
}

func eventHorizon(output chan *Agent, input *Agent) {
	var (
		xLen, yLen, len chan float64 = make(chan float64, 1), make(chan float64, 1), make(chan float64, 1)
	)
	go func() { xLen <- math.Pow(bestAgent.x-input.Y, 2.0) }()
	go func() { yLen <- math.Pow(bestAgent.y-input.Y, 2.0) }()
	go func() { len <- math.Pow(<-xLen+<-yLen, 0.5) }()
	if bestAgent.eventHorizon >= <-len {
		input.newPosition()
	}
	output <- input
}
