package converter

func ConvertLength(value float64, from string, to string) float64 {
	toMeter := map[string]float64{
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1.0,
		"kilometer":  1000.0,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.34,
	}

	meters := value * toMeter[from]
	return meters / toMeter[to]
}

func ConvertTemp(val float64, from, to string) float64 {
	if from == to {
		return val
	}

	var celsius float64
	switch from {
	case "farenheit":
		celsius = (val - 32) * 5 / 9
	case "kelvin":
		celsius = val - 273.15
	default:
		celsius = val
	}

	switch to {
	case "farenheit":
		return (celsius * 9 / 5) + 32
	case "kelvin":
		return celsius + 273.15
	default:
		return celsius
	}
}
