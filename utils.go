package goblackholes

import (
	"fmt"
	"sync"
)

type Agent struct {
	X, Y, Fitness  float64
	Border         Border_s
	TypeOfFunction TypeOfFunction_s
	Times          uint64
	Best           bool
}
type bestAgent_s struct {
	x, y, fitness, eventHorizon float64
	step                        uint64
	mutex                       sync.Mutex
}
type Border_s struct {
	X1, Y1, X2, Y2                                                     float64
	HorizontalLength, VerticalLength, HorizontalCenter, VerticalCenter float64
}

func (b *Border_s) setUp() {
	b.HorizontalLength = b.X2 - b.X1
	b.HorizontalCenter = (b.X2 + b.X1) / 2
	b.VerticalLength = b.Y2 - b.Y1
	b.VerticalCenter = (b.Y2 + b.Y1) / 2
}

func(b *Border_s) ToStr() string{
	return fmt.Sprintf("%f < x < %f ; %f < y < %f", b.X1, b.X2, b.Y1, b.Y2)

}

type TypeOfFunction_s struct {
	StringEvaluation string
	Rastrigin        bool
	Rosenbrock       bool
	Easom            bool
	McCormick        bool
}

type InitVariables struct {
	AgentAmount         int
	SingleServiceAmount int
	TypeOfFucntion      TypeOfFunction_s
	Border              Border_s
}

func (this *bestAgent_s) Convert() *Agent {
	agent := &Agent{
		this.x,
		this.y,
		this.fitness,
		Border_s{},
		TypeOfFunction_s{},
		this.step,
		true,
	}
	return agent
}

func (agent *Agent) newPosition() {
	go func() {
		agent.X = <-randomBuffer*agent.Border.HorizontalLength - agent.Border.HorizontalLength/2 + agent.Border.HorizontalCenter
	}()
	go func() {
		agent.Y = <-randomBuffer*agent.Border.VerticalLength - agent.Border.VerticalLength/2 + agent.Border.VerticalCenter
	}()
	agent.Best = false

}
