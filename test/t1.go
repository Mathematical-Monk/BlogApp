package main

import "fmt"

// 1. Define a "base" struct with some common properties.
type Engine struct {
	Horsepower int
	IsElectric bool
}

// 2. Define another struct that "embeds" Engine.
// Notice there is no field name, only the type `Engine`.
type Car struct {
	Make  string
	Model string
	Engine // The Engine struct is embedded here.
}

func main() {
	// Create an instance of the Car struct.
	myCar := Car{
		Make:  "Tesla",
		Model: "Model 3",
		Engine: Engine{ // Initialize the embedded struct by its type name.
			Horsepower: 350,
			IsElectric: true,
		},
	}

	// 3. Access the embedded fields DIRECTLY from the Car instance.
	// This is called "field promotion".
	fmt.Printf("Make: %s, Model: %s\n", myCar.Make, myCar.Model)
	fmt.Printf("Horsepower: %d\n", myCar.Horsepower) // Accessing Engine's field
	fmt.Printf("Is Electric: %t\n", myCar.IsElectric) // Accessing Engine's field

	fmt.Println("---")

	// You can also access the embedded struct itself as a field.
	fmt.Printf("Accessing via the nested struct: %d HP\n", myCar.Engine.Horsepower)

	fmt.Printf("accessing the nested struct: %+v\n", myCar.Engine)
}