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
	ErrPeriodDictionary = errors.New("missing in number period dictionary")
	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("should have one int64 parameter <number>")
)

func convertTensToWords(t int, d NumberDictionary) ([]string, error) {
	words := make([]string, 0, 2)
	word, ok := "", false
	if t > 0 && t < 20 {
		if word, ok = d[t]; !ok {
			return nil, fmt.Errorf("%d:%w", t, ErrNumbersDictionary)
		}
		return append(words, word), nil
	}
	if t >= 20 {
		key := t / 10 * 10
		if word, ok = d[key]; !ok {
			return nil, fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
		t -= key
	}
	if t > 0 {
		if word, ok = d[t]; !ok {
			return nil, fmt.Errorf("%d:%w", t, ErrNumbersDictionary)
		}
		words = append(words, word)
	}
	return words, nil
}

func convertTripletToWords(t int, d NumberDictionary) ([]string, error) {
	words := make([]string, 0, 3)
	if t >= 100 {
		key := t / 100 * 100
		word, ok := d[key]
		if !ok {
			return nil, fmt.Errorf("%d:%w", key, ErrNumbersDictionary)
		}
		words = append(words, word)
		t -= key
	}
	if t > 0 {
		hw, err := convertTensToWords(t, d)
		if err != nil {
			return nil, fmt.Errorf("get tens words:%w", err)
		}
		words = append(words, hw...)
	}
	return words, nil
}

func convertTripletToWord(n int, d NumberDictionary) (string, error) {
	words, err := convertTripletToWords(n, d)
	if err != nil {
		return "", fmt.Errorf("get triplet name:%w", err)
	}
	return strings.TrimSpace(strings.Join(words, " ")), nil
}

func getPeriodName(idx, n int, pd PeriodDictionary) (string, error) {
	if idx > len(pd)-1 || len(pd[idx]) != 3 {
		return "", ErrPeriodDictionary
	}
	n %= 100
	if n > 4 && n < 20 {
		return pd[idx][2], nil
	}
	n %= 10
	if n > 1 && n <= 4 {
		return pd[idx][1], nil
	}
	if n == 1 {
		return pd[idx][0], nil
	}
	return pd[idx][2], nil
}

func fixThousandsSuffix(tn *string) {
	if strings.HasSuffix(*tn, "ин") {
		*tn = strings.TrimSuffix(*tn, "ин") + "на"
	}
	if strings.HasSuffix(*tn, "ва") {
		*tn = strings.TrimSuffix(*tn, "ва") + "ве"
	}
}

func convertTripletsToWords(tr []int, nd NumberDictionary, pd PeriodDictionary) ([]string, error) {
	words := make([]string, 0)

	for i := len(tr) - 1; i >= 0; i-- {
		if tr[i] == 0 {
			continue
		}
		tn, err := convertTripletToWord(tr[i], nd)
		if err != nil {
			return nil, fmt.Errorf("get triplet name:%e", err)
		}

		if i == 1 {
			fixThousandsSuffix(&tn)
		}

		words = append(words, tn)

		if i > 0 {
			pn, err := getPeriodName(i-1, tr[i], pd)
			if err != nil {
				return nil, fmt.Errorf("get triplet period:%e", err)
			}

			words = append(words, pn)
		}
	}
	return words, nil
}

func parseTriplets(n int64) []int {
	t := make([]int, 0)
	for n > 0 {
		t = append(t, int(n%1000))
		n /= 1000
	}
	return t
}

func convertNumberToWord(n int64, nd NumberDictionary, pd PeriodDictionary) (string, error) {
	words := []string{}
	if n == 0 {
		return "нуль", nil
	}

	if n < 0 {
		words = append(words, "минус")
		n *= -1
	}

	t := parseTriplets(n)

	tw, err := convertTripletsToWords(t, nd, pd)
	if err != nil {
		return "", fmt.Errorf("get number name:%w", err)
	}

	words = append(words, tw...)
	return strings.TrimSpace(strings.Join(words, " ")), nil
}

// Task converts params number in words.
func Task(w io.Writer, parameters []string) error {
	if len(parameters) != 1 {
		return ErrParameters
	}
	n, err := strconv.ParseInt(parameters[0], 10, 64)
	if err != nil {
		return fmt.Errorf("number should be int64:%w", err)
	}

	word, err := convertNumberToWord(n, nd, pd)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s - %s\n", parameters[0], word)

	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print number converted to words\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <number>\n", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
		}
		fmt.Println(err)
	}
}
