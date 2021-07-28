package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// NumberDictionary represent one word number dictionary.
type NumberDictionary map[int]string

// PeriodDictionary represent period dictionary.
type PeriodDictionary [][]string

var nd = NumberDictionary{
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

var pd = PeriodDictionary{
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
	// ErrPeriodDictionary indicates that a period missing in dictionary.
	ErrPeriodDictionary = errors.New("missing in period dictionary")

	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("parameter length should be 1 <number>")
)

func appendFromDictionary(key int, words *[]string, d NumberDictionary) error {
	word, ok := d[key]
	if !ok {
		return fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
	}
	*words = append(*words, word)
	return nil
}

func getTripletName(n int, d NumberDictionary) (string, error) {
	words := make([]string, 0, 3)
	if n >= 100 {
		if err := appendFromDictionary(n/100*100, &words, d); err != nil {
			return "", fmt.Errorf("get triplet name:%w", err)
		}
		n %= 100
	}
	if n > 0 && n < 20 {
		if err := appendFromDictionary(n, &words, d); err != nil {
			return "", fmt.Errorf("get triplet name:%w", err)
		}
		n = 0
	}
	if n >= 20 {
		if err := appendFromDictionary(n/10*10, &words, d); err != nil {
			return "", fmt.Errorf("get triplet name:%w", err)
		}
		n %= 10
	}
	if n > 0 {
		if err := appendFromDictionary(n, &words, d); err != nil {
			return "", fmt.Errorf("get triplet name:%w", err)
		}
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

func getPeriodName(idx, n int, pd PeriodDictionary) (string, error) {
	if idx > len(pd)-1 {
		return "", ErrPeriodDictionary
	}
	periodNames := pd[idx]
	pluralIndex := getPluralIndex(n)

	if pluralIndex > len(periodNames)-1 {
		return "", ErrPeriodDictionary
	}
	return periodNames[pluralIndex], nil
}

func appendTripletsNames(tr []int, words *[]string, nd NumberDictionary, pd PeriodDictionary) error {
	for i := len(tr) - 1; i > 0; i-- {
		tripletName, err := getTripletName(tr[i], nd)
		if err != nil {
			return fmt.Errorf("append triplet name:%e", err)
		}
		if i == 1 {
			if strings.HasSuffix(tripletName, "ин") {
				tripletName = strings.TrimSuffix(tripletName, "ин") + "на"
			}
			if strings.HasSuffix(tripletName, "ва") {
				tripletName = strings.Replace(tripletName, "ва", "ве", 1)
			}
		}
		pn, err := getPeriodName(i-1, tr[i], pd)
		if err != nil {
			return fmt.Errorf("append triplet name:%e", err)
		}
		*words = append(*words, tripletName, pn)
	}
	tripletName, err := getTripletName(tr[0], nd)
	if err != nil {
		return fmt.Errorf("append triplet name:%e", err)
	}
	*words = append(*words, tripletName)
	return nil
}
func getNumberName(n int64, nd NumberDictionary, pd PeriodDictionary) (string, error) {
	words := []string{}
	if n == 0 {
		return "нуль", nil
	}
	if n < 0 {
		words = append(words, "минус")
		n *= -1
	}
	triplets := parseTriplets(n)
	if err := appendTripletsNames(triplets, &words, nd, pd); err != nil {
		return "", fmt.Errorf("get number name:%w", err)
	}
	return strings.TrimSpace(strings.Join(words, " ")), nil
}

// Task spell params number in words.
func Task(w io.Writer, args []string) error {
	// ErrSize indicates that a value does not have the right syntax for the size type.
	if len(args) != 1 {
		return ErrParameters
	}
	n, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("parsing int:%w", err)
	}

	numberName, err := getNumberName(n, nd, pd)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s - %s\n", args[0], numberName)

	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print number name\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <number>", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
		}
		fmt.Println(err)
	}
}
