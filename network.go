package main

import (
  "math/rand"
  "math"
)

type Node struct {
  id int
  vulnerability float64
  status string
  connections []*Node
}

type Network []*Node



//IsIn returns true if a given integer appears in a target []int
func IsIn(arr []int, k int) bool {
  for i := range arr {
    if arr[i] == k {
      return true
    }
  }
  return false
}


//MeanDegree is a network method that returns the mean degree of the network
func (n Network) MeanDegree() float64 {
  //calculate the mean degree of the network
  k := 0
  //sum the degree (length of connections slice) across all nodes
  for i := range n {
    k += len(n[i].connections)
  }
  meanDegree := float64(k) / float64(len(n))
  return meanDegree
}

//MeanResDegree is a network method that returns the mean original degree of the residual network
func (n Network) MeanResDegree() float64 {
  //Now we calculate the mean original degree of the residual nodes
  kRes := 0
  numRes := 0
  //sum the degree of all residual nodes
  for i := range n {
    if n[i].status == "S" {
      numRes++
      kRes += len(n[i].connections)
    }
  }
  meanResDegree := float64(kRes) / float64(numRes)
  return meanResDegree
}

//ResResDegree is a network method that returns the mean residual degree of the residual network
func (n Network) ResResDegree() float64 {
  kRes := 0
  numRes := 0
  //sum the residual degree of all residual nodes
  for i := range n {
    if n[i].status == "S" {
      numRes++
      for j := range n[i].connections {
        //Sum all the nodes connected to n[i] that are also residual nodes
        if n[i].connections[j].status == "S" {
          kRes++
        }
      }
    }
  }

  resResDeg := float64(kRes) / float64(numRes)
  return resResDeg
}

//MeanSquaredDegree is a network method that returns the mean squared-degree of the network
func (n Network) MeanSquaredDegree() float64 {
  k2 := 0
  //sum the square of the degree (length of connections slice) across all nodes
  for i := range n {
    deg := len(n[i].connections)
    k2 += deg * deg
  }
  meanSquaredDegree := float64(k2) / float64(len(n))
  return meanSquaredDegree
}

//NetworkFrailty takes an already infected network as an input and returns the frailty parameter
func NetworkFrailty(n Network) float64 {

  k := n.MeanDegree()
  kr := n.MeanResDegree()

  frailty := (k - kr) / k

  return frailty
}

//NetworkInterference takes an already infected network as an input and returns the interference parameter
func NetworkInterference(n Network) float64 {
  k := n.MeanDegree()
  kr := n.MeanResDegree()
  krr := n.ResResDegree()

  interf := (kr - krr) / k

  return interf
}


//Sample from the power-law distribution, using Newton's method to solve for the transcendental equation
func PowerLaw(alpha, kappa, C float64) int {
  //first, we sample from the uniform distribution over [0,1]
  seed := rand.Float64()

  //From Meyers et al. We are given that an appropriate power-law equation for epidemiology problems is
  //p_k = Ck^(-alpha) exp(-k/kappa)

  //Using algebra, this can be rewritten as the following:
  //(alpha)(kappa)log(k) + k - (kappa)(log(C/p_k)) = 0
  //With p_k = seed


  //We use Newton's method to solve for k (we have to approximate because log(k) + k is a transcendental function)
  //Newton's method takes a target 'guess' (We will start with k_0 = 1) and iterates through subtracting f(k_n)/f'(k_n) to generate
  //k_n+1. We set an arbitrary threshold of convergence delta, and when this threshold is reached, we return the closest integer. We don't need
  //an extreme degree of accuracy since we are looking for integer degree outputs only

  //So, initialize k to 1, and the convergence delta to 0.01
  k := 1.0
  delta := 0.0001

  nextK := 0.0

  //while not converged, implement Newton's method
  for true {
    //f(k) is simply the equation shown above.
    fk := (alpha * kappa * math.Log(k)) + k - (kappa * math.Log((C / seed)))

    //df(k), its derivative w.r.t. k, is equal to [(alpha)(kappa) + k]/k
    dfk := ((alpha * kappa) + k) / k

    nextK = k - (fk / dfk)
    if math.Abs(nextK - k) <= delta {
      k = nextK
      break
    }

    k = nextK

  }

  return int(math.Round(k))

}

//InitializeNetwork takes an empty network and initializes nodes with Gaussian vulnerability multiplier, status of susceptible,
//and an empty list of connections (slice of pointers to nodes)
func (n Network) InitializeNetwork() {
  for i := range n {
    c := make([]*Node, 0)
    vuln := GaussianVuln()
    n[i] = &Node{i, vuln, "S", c}
  }
}

//ConnectNetwork takes a network and connects each edge to random edges in the network such that the degree of each
//node is sampled from the power-law distribution outlined in Meyers et al.
func (n Network) ConnectNetwork(alpha, kappa, C float64) {
  for i := range n {
    edges := make([]*Node, 0)
    //The degree of node n[i] is taken from the Power-Law distribution used in Meyers et al.
    c := PowerLaw(alpha, kappa, C)
    if c > len(n) {
      c = len(n)
    }

    //alreadyConnected stores the id's of nodes already connected to our current node.
    //Used to eliminate repeats
    alreadyConnected := make([]int, 0)
    //Now connect node n[i] to c random nodes in the network n
    for c > 0 {
      target := rand.Intn(len(n))

      //A node cannot point to itself, so continue generating until a non-self number is reached
      //We also select a new target if the target is alreadyConnected
      for (target == i || IsIn(alreadyConnected, target)) {
        newTarget := rand.Intn(len(n))
        target = newTarget
      }

      var newEdge *Node
      newEdge = n[target]

      edges = append(edges, newEdge)
      alreadyConnected = append(alreadyConnected, target)
      c--
    }

    n[i].connections = edges
  }
}

//Vaccinate takes a vaccination rate as a float64 input and vaccinates every Node in network n with probability (rate)
func (n Network) Vaccinate(rate float64) {
  for i := range n {
    vaccineChance := rand.Float64()
    if vaccineChance <= rate {
      n[i].status = "V"
    }
  }
}

//IsInfected returns true if any node in a network is currently infected, false otherwise. This is how the
//algorithm knows when to stop iterating.
func (n Network) IsInfected() bool {
  for i := range n {
    if n[i].status == "I" {
      return true
    }
  }
  return false
}
