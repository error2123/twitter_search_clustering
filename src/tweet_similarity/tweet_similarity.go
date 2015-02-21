package tweet_similarity


import (
	"math"
	"strings"
)

/**************************************************
 * Cosine takes two strings str1 and str2
 * (sentences delimited by spaces) and computes
 * the cosine similarity between the words.
 * It returns a float value between 0 to 1 with
 * 0 being no similarity and 1 being the highest
 * similarity.
 ***************************************************/
func Cosine(str1, str2 string) float64 {
	// delimit the strings by space
	temp1 := strings.Fields(strings.ToLower(str1))
	temp2 := strings.Fields(strings.ToLower(str2))

	// to convert string to vector
	// to associate each character with its frequency values in string
	// vect1 contains the frequencies of each character
	vect1 := make(map[string]int)
	for _, elem := range temp1 {
		vect1[elem]++
	}

	// vect2 contains the frequencies of each character
	vect2 := make(map[string]int)
	for _, elem := range temp2 {
		vect2[elem]++
	}

	// intersection contains common characters
	intersection := []string{}
	// traverse keys and add what is common to both.
	// (keys, in Hash/Map, are like indices in array)
	for key := range vect1 {
		if _, exist := vect2[key]; exist {
			intersection = append(intersection, key)
		}
	}

	// If all the vector elements are equal, cos will be 1.
	// Equal texts return the value 1.
	// We need to traverse the intersection(common) characters of texts.
	// In doing so, we can expect two same texts to return 1, cos 0°
	// to calculate A·B
	sum := 0.0
	for _, elem := range intersection {
		sum += float64(vect1[elem]) * float64(vect2[elem])
	}
	numerator := sum

	// to calculate |A|*|B|
	sum1 := 0.0
	for key := range vect1 {
		sum1 += math.Pow(float64(vect1[key]), 2)
	}
	sum2 := 0.0
	for key := range vect2 {
		sum2 += math.Pow(float64(vect2[key]), 2)
	}
	denominator := math.Sqrt(sum1) * math.Sqrt(sum2)

	// smoothing because we can't divide by 0
	if numerator == 0.0 || denominator == 0.0 {
		return 0.0001
	}
	return float64(numerator) / denominator
}
