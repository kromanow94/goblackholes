package goblackholes

import (
	//"fmt"
	"math"
	//"runtime"
	"sync"
	//"time"
	//"os"
	//"time"
	//"os"
	//"fmt"
	//"runtime"
	//"time"
	"fmt"
	"time"
)

var (
	agentAmount         int
	singleServiceAmount int      = 50
	slowmotion          int      = 0
	border              Border_s// = Border_s{-1.5, -3.0, 4.0, 4.0, 0, 0, 0, 0}
	agentList           []*Agent     //= make([]*Agent, agentAmount, agentAmount)
	bestAgent           bestAgent_s  // = bestAgent_s{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, 0, 0, sync.Mutex{}}
	newAgentChannel     chan *Agent  // = make(chan *Agent, agentAmount)
	fitnessChannel      chan *Agent  // = make(chan *Agent, agentAmount)
	getBestChannel      chan *Agent  // = make(chan *Agent, agentAmount)
	moveChannel         chan *Agent  // = make(chan *Agent, agentAmount)
	eventHorizonChannel chan *Agent  // = make(chan *Agent, agentAmount)
	protoChannel        chan *Agent  //  = make(chan *Agent, agentAmount)
	randomBuffer        chan float64 // = make(chan float64, agentAmount*5)
	sendResults	bool = true
	typeOfFunction TypeOfFunction_s //  TypeOfFunction_s{McCormick: true}
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
	sendResults = false
	defer func(){
		if r := recover() ; r != nil{
			fmt.Println("Recovered: ",r)
		}
	}()
	//recover()
	//close(outputChan)
	return

}

func startComputing(outputChan chan *Agent) {
	for i := 0; i < singleServiceAmount; i++ {
		go func() {
			defer func(){
				if r := recover(); r != nil{
					//fmt.Println("service recovery: ", r)
				}
			}()
			for {
				getBest(getBestChannel, outputChan, <-fitnessChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
					//fmt.Println("service recovery: ", r)
				}
			}()
			for {
				move(moveChannel, <-getBestChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
					//fmt.Println("service recovery: ", r)
				}
			}()
			for {
				eventHorizon(eventHorizonChannel, <-moveChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
					//fmt.Println("service recovery: ", r)
				}
			}()
			for {
				countFitness(protoChannel, <-eventHorizonChannel)
			}
		}()
		go func() {
			defer func(){
				if r := recover(); r != nil{
					//fmt.Println("service recovery: ", r)
				}
			}()
			for {
				getProto := <-protoChannel
				fitnessChannel <- getProto
				if sendResults {
					outputChan <- getProto
				}
			}
		}()
	}
}

//func utils() {
/// reports
//go func() {
//	for {
//		select {
//		case <-endComputing:
//			return
//		default:
//			fmt.Println(bestAgent)
//			fmt.Println(runtime.NumGoroutine())
//			averageStepAmount := averageStepAmount()
//			fmt.Println("averageStepAmount: ", averageStepAmount)
//			time.Sleep(500 * time.Millisecond)
//		}
//	}
//}()

//// check if got the Best answer
//go func() {
//	<-maxAccuracy
//	fmt.Println("It can't be bether:")
//	for i := 0; i < 4*singleServiceAmount+10; i++ {
//		endComputing <- true
//	}
//	fmt.Println(bestAgent)
//	averageStepAmount := averageStepAmount()
//	fmt.Println("averageStepAmount: ", averageStepAmount)
//	exitProgram <- true
//	return
//}()
//}

//func averageStepAmount() uint64 {
//	var averageStepAmount uint64
//	for i := 0; i < agentAmount; i++ {
//		averageStepAmount += agentList[i].Times
//	}
//	averageStepAmount /= uint64(agentAmount)
//	return averageStepAmount
//}

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
