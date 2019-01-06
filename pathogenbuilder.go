package main

import(
  "fmt"
  "os"
)

//Writes a .PATHOGEN file specifying the name of the disease, it's base reproductive ratio, and it's mortality rate. This is
//a private function called in main.go
func BuildPathogen() string {
  fmt.Print("Please create a pathogen to infect your population. First, name your pathogen (no whitespaces please): ")

  //Read the pathogen name from user input
  pathName := ""
  fmt.Scanln(&pathName)

  fmt.Println("")
  fmt.Println("")

  fmt.Println("Specify a Basic Reproductive Ratio (Ro) for", pathName, ". The higher this decimal number, the more easily spreadable your pathogen is.")
  fmt.Println("For reference, here are some common Ro values for known pathogens [https://en.wikipedia.org/wiki/Basic_reproduction_number]. Note that a value below 1 cannot reach epidemic levels.")
  fmt.Println("Measles: 12-15 \t Smallpox, Polio, Rubella, Mumps: 5-7 \t HIV: 2-5 \t Influenze: 2-3 \t SARS, Ebola: 2-2.5")

  //Read the base reproductive ratio as a string
  roString := ""
  fmt.Scanln(&roString)



  fmt.Println("")
  fmt.Println("")

  fmt.Println("Please enter the mortality rate for", pathName, ", as a decimal number between 0 and 1, inclusive")
  fmt.Println("Sample mortality rates for illnesses:")
  fmt.Println("Ebola (Zaire Strain): 0.90 \t HIV (Untreated): 0.85 \t Influenza A H5N1: .60 \t Smallpox: 0.30 \t SARS: 0.11 \t Measles: .03 \t Spanish Flu Variant: .03")

  //Read the mortality rate as a string
  mortString := ""
  fmt.Scanln(&mortString)


  //now write to a .PATHOGEN file in the /pathogens/ directory
  outputPath := "pathogens/" + pathName + ".PATHOGEN"
  outFile, errF := os.Create(outputPath)
  if errF != nil {
    fmt.Println("Error writing to file")
    os.Exit(1)
  }

  defer outFile.Close()

  fmt.Fprintln(outFile, pathName)
  fmt.Fprintln(outFile, roString)
  fmt.Fprintln(outFile, mortString)

  return outputPath
}
