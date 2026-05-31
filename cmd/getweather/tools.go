package main

func ConvertCelciusToFahrenheit(temp float64) (float64, error) {
	converted := (temp * 9 / 5) + 32
	return converted, nil
}
