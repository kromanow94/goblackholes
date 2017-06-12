package goblackholes

import (
	"math"
	"fmt"
	"time"
	"sync"
)

var (
	agentAmount         int
	singleServiceAmount int      = 50
	border              Border_s
	agentList           []*Agent
	bestAgent           bestAgent_s
	newAgentChannel     chan *Agent
	fitnessChannel      chan *Agent
	getBestChannel      chan *Agent
	moveChannel         chan *Agent
	eventHorizonChannel chan *Agent
		protoChannel        chan *Agent
	randomBuffer        chan float64
	sendResults	bool = true
	typeOfFunction TypeOfFunction_s
	counter uint64
)

func Start(outputChan chan *Agent, quitChan chan bool, initVariables InitVariables) {
	agentAmount = initVariables.AgentAmount
	border = initVariables.Border
	typeOfFunction = initVariables.TypeOfFucntion

	initialize()

	startComputing(outputChan)

	stopStatInfo := false
	go func() {
		for !stopStatInfo {
			fmt.Println(bestAgent)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	<-quitChan
	stopStatInfo = true
	defer func(){
		if r := recover() ; r != nil{
			fmt.Println("Recovered: ",r)
		}
	}()
	return

}

func startComputing(outputChan chan *Agent) {
	//go func() {
	//	defer func(){
	//		if r := recover(); r != nil{
	//		}
	//	}()
	//	getBest(getBestChannel, outputChan, fitnessChannel)
	//}()
	for i := 0; i < singleServiceAmount; i++ {
		go func() {

			defer func() {

				if r := recover(); r != nil {
				}
			}()
			for {
				getBest(getBestChannel, outputChan, <-fitnessChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
				}
			}()
			for {
				move(moveChannel, <-getBestChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
				}
			}()
			for {
				eventHorizon(eventHorizonChannel, <-moveChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
				}
			}()
			for {
				countFitness(protoChannel, <-eventHorizonChannel)
			}
		}()
	}
	go func() {
				defer func(){
					if r := recover(); r != nil{
					}
				}()
				for {
					getProto := <-protoChannel
					if sendResults {
						fitnessChannel <- getProto
						outputChan <- getProto
					}
				}
			}()
}

func initialize() {
	newAgentChannel = make(chan *Agent, agentAmount)
	fitnessChannel = make(chan *Agent, agentAmount)
	getBestChannel = make(chan *Agent, agentAmount)
	moveChannel = make(chan *Agent, agentAmount)
	eventHorizonChannel = make(chan *Agent, agentAmount)
	protoChannel = make(chan *Agent, agentAmount)
	randomBuffer = make(chan float64, agentAmount*5)

	bestAgent = bestAgent_s{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, 0, 0, sync.Mutex{}}

	go func() {
		defer func(){
			if r := recover(); r != nil{
				//fmt.Println("service recovery: ", r)
			}
		}()
		for {
			randomBuffer <- NextDouble()
		}
	}()

	border.setUp()
	agentList = make([]*Agent, agentAmount, agentAmount)

	/// create agents
	for i := 0; i < agentAmount; i++ {
		agent := Agent{
			X:              <-randomBuffer*border.HorizontalLength - border.HorizontalLength/2 + border.HorizontalCenter,
			Y:              <-randomBuffer*border.VerticalLength - border.VerticalLength/2 + border.VerticalCenter,
			Fitness:        math.MaxFloat64,
			Border:         border,
			TypeOfFunction: typeOfFunction}
		agentList[i] = &agent
	}

	/// send agents to channel
	for i := 0; i < agentAmount; i++ {
		newAgentChannel <- agentList[i]
	}

	/// set pre-values for agents
	for i := 0; i < agentAmount; i++ {
		countFitness(fitnessChannel, <-newAgentChannel)
	}

	bestAgent.eventHorizon = math.SmallestNonzeroFloat64

}
