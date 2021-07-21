package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var numbersDictionary = map[int]string{
	0: "",
	1: "один",
	2: "два",
	3: "три",
	4: "четыре",
	5: "пять",
	6: "шесть",
	7: "семь",
	8: "восемь",
	9: "девять",

	10: "десять",
	11: "одиннадцать",
	12: "двенадцать",
	13: "тринадцать",
	14: "четырнадцать",
	15: "пятнадцать",
	16: "шестнадцать",
	17: "семнадцать",
	18: "восемнадцать",
	19: "девятнадцать",

	20: "двадцать",
	30: "тридцать",
	40: "сорок",
	50: "пятьдесят",
	60: "шестьдесят",
	70: "семьдесят",
	80: "восемьдесят",
	90: "девяносто",

	100: "сто",
	200: "двести",
	300: "триста",
	400: "четыреста",
	500: "пятьсот",
	600: "шестьсот",
	700: "семьсот",
	800: "восемьсот",
	900: "девятьсот",
}

var periodDictionary = [][]string{
	{"тысяча", "тысячи", "тысяч"},
	{"миллион", "миллиона", "миллионов"},
	{"миллиард", "миллиарда", "миллиардов"},
	{"триллион", "триллиона", "триллионов"},
	{"квадриллион", "квадриллиона", "квадриллионов"},
	{"квинтиллион", "квинтиллиона", "квинтиллионов"},
}

var (
	// ErrNumbersDictionary indicates that a number missing in dictionary.
	ErrNumbersDictionary = errors.New("missing in numbers dictionary")
	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("parameter length should be 1 <number>")
)

func getTripletName(n int, dictionary map[int]string) (string, error) {
	words := make([]string, 0, 3)
	if n >= 100 {
		key := n % 1000 / 100 * 100
		word, ok := dictionary[key]
		if !ok {
			return "", fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
		n %= 100
	}
	if n > 0 && n < 20 {
		key := n
		word, ok := dictionary[key]
		if !ok {
			return "", fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
		n = 0
	}
	if n >= 20 {
		key := n / 10 * 10
		word, ok := dictionary[key]
		if !ok {
			return "", fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
		n %= 10
	}
	if n >= 0 {
		key := n
		word, ok := dictionary[key]
		if !ok {
			return "", fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
	}

	return strings.TrimSpace(strings.Join(words, " ")), nil
}

func parseTriplets(number int64) []int {
	t := []int{}
	for number > 0 {
		t = append(t, int(number%1000))
		number /= 1000
	}

	return t
}

func getNumberName(number string, nd map[int]string, pd [][]string) (string, error) {
	words := []string{}
	n, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parsing int:%w", err)
	}
	if n == 0 {
		return "нуль", nil
	}
	if n < 0 {
		words = append(words, "минус")
		n *= -1
	}
	triplets := parseTriplets(n)
	for i := len(triplets) - 1; i > 0; i-- {
		tripletName, err := getTripletName(triplets[i], nd)
		if err != nil {
			return "", err
		}
		if i == 1 {
			if strings.HasSuffix(tripletName, "ин") {
				tripletName = strings.TrimSuffix(tripletName, "ин") + "на"
			}
			if strings.HasSuffix(tripletName, "ва") {
				tripletName = strings.Replace(tripletName, "ва", "ве", 1)
			}
		}
		words = append(words, tripletName, pd[i-1][getPluralIndex(triplets[i])])
	}
	tripletName, err := getTripletName(triplets[0], nd)
	if err != nil {
		return "", err
	}
	words = append(words, tripletName)

	return strings.TrimSpace(strings.Join(words, " ")), nil
}

// Task spell params number in words.
func Task(w io.Writer, args []string) error {
	// ErrSize indicates that a value does not have the right syntax for the size type.
	if len(args) != 1 {
		return ErrParameters
	}
	numberName, err := getNumberName(args[0], numbersDictionary, periodDictionary)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s - %s\n", args[0], numberName)

	return nil
}

func getPluralIndex(n int) int {
	n %= 100
	if n > 4 && n < 20 {
		return 2
	}
	n %= 10
	if n > 1 && n <= 4 {
		return 1
	}
	if n == 1 {
		return 0
	}

	return 2
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print number name\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <number>", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(0)
	}
}
