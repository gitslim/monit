package entities

import (
	"fmt"
)

// ExampleGetMetricType получение типа метрики по названию
func ExampleGetMetricType() {
	gauge, _ := GetMetricType("gauge")
	counter, _ := GetMetricType("counter")
	wrongMetricType, _ := GetMetricType("badmetrictype")

	fmt.Printf("%s %s %s\n", gauge, counter, wrongMetricType)
	// Output: gauge counter
}
