package main

import (
  "fmt"
  "os"
  "bufio"
  "math/rand"
  "math"
  "time"
  "strconv"
  "image"
  "log"
)

type Pathogen struct {
  name  string
  Ro  float64
  lethality float64
}



//Returns the transmissibility as a function of the infectivity of a given pathogen (Ro) and a given network
func Transmissibility(Ro float64, n Network) float64 {
  k := n.MeanDegree()

  k2 := n.MeanSquaredDegree()


  T := (Ro / k2) * (k - 1.0)
  return T

}

//Finally, some individuals are more susceptible to disease and mortality (elderly, immunocompromised people, children), and some are less. This is
//modeled by adding a Gaussian term "vulnerability" that will function as a multiplier of the probabilities of infection and mortality.
//A vulnerability of 1 would represent the average of the population, and a SD of 0.556 (this is based on U.S. national percentages of elderly and immunocompromised individuals)
func GaussianVuln() float64 {
  base := rand.NormFloat64()
  //Mean of 1, SD of 0.556
  v := (base * 0.556) + 1

  //It does not make sense to have vulnerability coefficient below 0, so we set a lower bound on v
  if v <= 0.1 {
    v = 0.1
  }

  return v
}


//Run one pass of infection through a given network. For each infected node, its neighbors are also infected with a probability T equal
//to the transmissibility of the pathogen in that network
func InfectOnce(n Network, p Pathogen) Network {
  //First, compute the transmissibility of p in this network
  transmitRate := Transmissibility(p.Ro, n)

  //Now, range over infected nodes in n
  for i := range n {
    if n[i].status == "I" {
      neighbors := n[i].connections

      //Infect each susceptible neighbor with probability transmitRate. If the neighbor is either dead or immune ("D" or "R"), they cannot be infected
      for k := range neighbors {
        infectChance := rand.Float64()
        infectChance *= n[i].vulnerability
        if infectChance <= transmitRate && neighbors[k].status == "S" {
          neighbors[k].status = "I"
        }
      }

      //Now we update the status of the infected node to either dead "D" or immune "R" with probability of death based on the lethality of the pathogen
      deathChance := rand.Float64()
      deathChance *= n[i].vulnerability
      if deathChance <= p.lethality {
        n[i].status = "D"
      } else {
        n[i].status = "R"
      }
    }
  }

  return n

}

//Initially used for debugging and readability, just converts the single character status into a single word description for
//a given Node.
 func ReadStatus(n *Node) string {
  if n.status == "S" {
    return "susceptible"
  } else if n.status == "I" {
    return "infected"
  } else if n.status == "V" {
    return "immune"
  } else if n.status == "R" {
    return "recovered"
  } else if n.status == "D" {
    return "dead"
  } else {
    return "ERROR READING STATUS"
  }
}


//Reads the name, ro, and death rate from a .PATHOGEN file specified
func ReadPathogenFromFile(filePath string) (string, float64, float64) {
  file, errF := os.Open(filePath)

  if errF != nil {
    fmt.Println("Error reading .PATHOGEN file")
    os.Exit(1)
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)

  //First, the name
  scanner.Scan()
  pathName := scanner.Text()

  //next, RO
  scanner.Scan()
  roString := scanner.Text()
  ro, err2 := strconv.ParseFloat(roString, 64)
  if err2 != nil {
    fmt.Println("Unable to Parse Ro.")
    os.Exit(2)
    } else if ro <= 0.0 {
      fmt.Println("Invalid input. Please enter a decimal number greater than 0.")
      os.Exit(2)
    }

  //finally, the death rate
  scanner.Scan()
  mortString := scanner.Text()
  deathRate, err3 := strconv.ParseFloat(mortString, 64)
  if err3 != nil {
      fmt.Println("Unable to Parse Mortality Rate.")
      os.Exit(3)
      } else if deathRate < 0.0 || deathRate > 1.0 {
        fmt.Println("Invalid input. Please enter a decimal number between 0 and 1, inclusive.")
        os.Exit(3)
      }
  fmt.Println("Successfully loaded", pathName)
  return pathName, ro, deathRate
}


func main() {
  //Seed the random generator
  rand.Seed(time.Now().UTC().UnixNano())


  //Prompt the user for the .PATHOGEN file, then load it into path.
  fmt.Println("Please enter the full .PATHOGEN filepath for your disease, or enter \"CUSTOM\" to create your own (case-sensitive):")
  disInput := ""
  fmt.Scanln(&disInput)

  //If the input is "CUSTOM", redirect to pathogenbuilder.go
  if disInput == "CUSTOM" {
    customPath := BuildPathogen()
    disInput = customPath
  } else {
    disInput = "pathogens/" + disInput
  }


  //Read the pathogen from the file indicated.
  pathName, ro, deathRate := ReadPathogenFromFile(disInput)


  //Prompt the user for the population info
  fmt.Print("Enter Population:")
  popString := ""
  fmt.Scanln(&popString)

  pop, err1 := strconv.Atoi(popString)
  if err1 != nil {
    fmt.Println("Unable to Parse population input.")
    os.Exit(1)
  } else if pop <= 0 {
    fmt.Println("Invalid input. Please enter an integer greater than 0.")
    os.Exit(1)
  } else {
    fmt.Println("Creating Network with population", pop)
  }

  fmt.Println("Network successfully generated!")

  fmt.Println("")
  fmt.Println("")


  //Prompt user for vaccine rate
  fmt.Println("What percentage of your population is vaccinated against", pathName, "?")
  vacString := ""
  fmt.Scanln(&vacString)

  vaccineRate, err4 := strconv.ParseFloat(vacString, 64)
  if err4 != nil {
    fmt.Println("Unable to Parse Vaccination Rate.")
    os.Exit(4)
  } else if vaccineRate < 0.0 || vaccineRate > 100.0 {
    fmt.Println("Invalid input. Please enter a number between 0 and 100, inclusive.")
    os.Exit(4)
  }


  //We are using vaccineRate as a probability, not a percentage, so divide by 100.0
  vaccineRate = vaccineRate / 100.0


  //Prompt user for patient(s) zero information
  fmt.Println("Specify the number of patients to start with the infection (a value of 1 corresponds to a single patient 0 and a value of 0 means no one is infected.)")
  pZeroString := ""
  fmt.Scanln(&pZeroString)

  pZero, err5 := strconv.Atoi(pZeroString)
  if err5 != nil {
    fmt.Println("Unable to Parse patient(s) zero.")
    os.Exit(5)
  } else if pZero < 0 {
    fmt.Println("Invalid input. Please enter an integer greater than or equal to 0.")
    os.Exit(5)
  }


  //Initialize an empty Network and an empty slice of images for visualization
  net := make(Network, pop)

  progression := make([]image.Image, 0)


  //Now initialize the network and Connect it using the parameters given by Meyers et al.
  net.InitializeNetwork()
  net.ConnectNetwork(2, 94.2, float64(pop)/10.0)

  //Initialize our Pathogen object based on the values from the .PATHOGEN file
  p1 := Pathogen{pathName, ro, deathRate}

  net.Vaccinate(vaccineRate)

  //Initialize the patient(s) zero
  for i := 0; i < pZero; i++ {
    //If everyone is vaccinated, no one is getting infected
    if vaccineRate == 1.0 {
      break
    }

    //Otherwise pick a random person from the network
    patientZeroID := rand.Intn(len(net))
    //Prevent repeats
    for net[patientZeroID].status == "I" || net[patientZeroID].status == "V" {
      patientZeroID = rand.Intn(len(net))
    }

    //Now set them to infected
    net[patientZeroID].status = "I"
  }

  //Now draw our initial infected network to '0.png'
  progression = append(progression, DrawNetwork(net, 10, 0))

  //numEpochs is used to keep track of what timestep we are in for the purposes of writing the progression image files.
  numEpochs := 1

  //Keep infecting until the network is no longer infected
  for true {
    net = InfectOnce(net, p1)
    progression = append(progression, DrawNetwork(net, 10, numEpochs))

    if net.IsInfected() == false {
      break
    }

    numEpochs++
  }

  //We write a death map which maps status strings to counts from our network, then it is
  //Passed through to WriteEpidemicToFile().
  deathMap := make(map[string]int)

  //Read the statuses so they are easier to interpret for debugging purposes.
  for i := range net {
    deathMap[ReadStatus(net[i])]++
  }

  fmt.Println("Processing images...")
  Process(progression, pathName)
  fmt.Println("done!")

  //Now write our epidemic to file
  fmt.Println("Writing Epidemic Statistics to", pathName + ".txt")
  WriteEpidemicToFile(deathMap, p1, net, vaccineRate * 100)
}

//WriteEpidemicToFile writes all the statistics of our epidemic to a file
//called [PATHOGEN_NAME].txt
func WriteEpidemicToFile(m map[string]int, p Pathogen, n Network, vacRate float64) {
  //Standard Go I/O code. Lots of Fprint statements so we print exactly what we want.
  file, err := os.Create(p.name + ".txt")
  if err != nil {
  log.Fatal("Cannot create file", err)
  }

  defer file.Close()
  fmt.Fprint(file, "Out of a total population of ")
  fmt.Fprint(file, len(n))
  fmt.Fprint(file, "\r\n")
  fmt.Fprint(file, p.name + " killed ")
  fmt.Fprint(file, m["dead"])
  if m["dead"] == 1 {
    fmt.Fprint(file, " person. ")
  } else {
    fmt.Fprint(file, " people. ")
  }

  fmt.Fprint(file, m["recovered"])
  if m["recovered"] == 1 {
    fmt.Fprint(file, " was infected, but survived. ")
  } else {
    fmt.Fprint(file, " were infected, but survived. ")
  }

  fmt.Fprint(file, m["immune"])
  if m["immune"] == 1 {
    fmt.Fprint(file, " was vaccinated and did not contract " + p.name + ", and ")
  } else {
    fmt.Fprint(file, " were vaccinated and did not contract " + p.name + ", and ")
  }


  fmt.Fprint(file, m["susceptible"])
  if m["recovered"] == 1 {
    fmt.Fprint(file, " were susceptible to the disease but was not exposed. \r\n")
  } else {
  fmt.Fprint(file, " were susceptible to the disease but were not exposed. \r\n")
  }



  fmt.Fprint(file, vacRate)
  fmt.Fprint(file, "% of the population was vaccinated, and the pathogen had a base reproductive ratio of ")
  fmt.Fprint(file, p.Ro)
  fmt.Fprint(file, " and a mortality rate of ")
  fmt.Fprint(file, p.lethality * 100)
  fmt.Fprintln(file, "%. \r\n")


  //Call the frailty and interference methods then print
  frailty := NetworkFrailty(n)
  interference := NetworkInterference(n)

  fmt.Fprintln(file, "Network Frailty Statistics: \r\n")
  fmt.Fprint(file, "Frailty: ", frailty, " \t ", "Interference: ", interference, "\r\n")

}



//DrawNetwork is an adaptation of the drawing code from Cellular Automata, rewritten slightly
//To write a network.
func DrawNetwork(n Network, cellWidth int, i int) image.Image {
  sqrt := int(math.Sqrt(float64(len(n))))

  height := (sqrt + 1) * cellWidth
	width := sqrt * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
  //darkGray := MakeColor(10, 10, 10)
	//black := MakeColor(0, 0, 0)
	blue := MakeColor(0, 107, 225)
	red := MakeColor(211, 10, 10)
	//green := MakeColor(0, 255, 0)
	yellow := MakeColor(170, 175, 8)
	//magenta := MakeColor(255, 0, 255)
	white := MakeColor(255, 255, 255)
	//cyan := MakeColor(0, 255, 255)

	// fill in colored squares. S and V are white, I is yellow, D is red, and R is blue
	for i := 0; i <= sqrt; i++ {
		for j := 0; j < sqrt; j++ {
      index := (i * sqrt) + j
			if index >= len(n) {
        c.SetFillColor(white)
      } else if n[index].status == "S" {
				c.SetFillColor(white)
			} else if n[index].status == "I" {
				c.SetFillColor(yellow)
			} else if n[index].status == "V" {
				c.SetFillColor(white)
			} else if n[index].status == "R" {
				c.SetFillColor(blue)
			} else if n[index].status == "D" {
				c.SetFillColor(red)
			} else {
        c.SetFillColor(white)
      }

			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}

  st := strconv.Itoa(i)

  //Save to the progression directory
  c.SaveToPNG("progression/" + st + ".png")


	return c.img
}
