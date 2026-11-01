package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var help bool
	flag.BoolVar(&help, "help", false, "показать использование программы")
	flag.BoolVar(&help, "h", false, "показать использование программы")

	var count int
	flag.IntVar(&count, "n", 1, "количество итераций")
	flag.IntVar(&count, "count", 1, "количество итераций")

	seed := flag.Int64("seed", -1, "зерно")

	var sum bool
	flag.BoolVar(&sum, "sum", false, "выводить только сумму значений костей")
	flag.BoolVar(&sum, "s", false, "выводить только сумму значений костей")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "выводить подробную информацию по каждому броску внутри итерации")
	flag.BoolVar(&verbose, "verbose", false, "выводить подробную информацию по каждому броску внутри итерации")
	flag.Parse()

	if help {
		PrintHelp()
		return
	}

	if sum && verbose {
		fmt.Println("Нельзя использовать одновременно флаг sum и verbose")
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Нужно писать в формате NdM")
		os.Exit(1)
	}

	if count < 1 {
		fmt.Println("count не может быть меньше единицы")
		os.Exit(1)
	}

	dice := flag.Args()[0]
	_, _, err := Parse(dice)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var r *rand.Rand
	if *seed != -1 {
		r = rand.New(rand.NewSource(*seed))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if count == 1 {
		if sum {
			summ, err := RollWithOnlySum(r, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Sum: " + strconv.Itoa(summ))
		} else if verbose {
			res, err := VerboseRoll(r, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		} else {
			res, err := Roll(r, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		}
	} else {
		if sum {
			res, err := RollSumWithIterations(r, count, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		} else if verbose {
			res, err := VerboseRollWithIterations(r, count, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		} else {
			res, err := RollWithIterations(r, count, dice)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		}
	}
}

func Roll(r *rand.Rand, dice string) (string, error) {
	N, M, err := Parse(dice)
	if err != nil {
		return "", err
	}
	summ := 0
	res := "Rolls: "
	for i := range N {
		number := r.Intn(M) + 1
		summ += number
		if i > 0 {
			res += "+"
		}
		res += strconv.Itoa(number)
	}
	res += "=" + strconv.Itoa(summ)

	return res, nil
}

func VerboseRoll(r *rand.Rand, dice string) (string, error) {
	res := "Rolling " + dice + ":\n"
	N, M, err := Parse(dice)
	if err != nil {
		return "", err
	}
	summ := 0
	for i := range N {
		number := r.Intn(M) + 1
		summ += number
		res += "Dice " + strconv.Itoa(i+1) + ": " + strconv.Itoa(number) + "\n"
	}
	summStr := strconv.Itoa(summ)
	avgStr := fmt.Sprint(float64(summ) / float64(N))
	res += "Sum: " + summStr + "\n"
	res += "Average: " + avgStr + "\n"

	return res, nil
}

func VerboseRollWithIterations(r *rand.Rand, iterations int, dice string) (string, error) {
	N, M, err := Parse(dice)
	if err != nil {
		return "", err
	}
	res := "Rolling " + dice + " (" + strconv.Itoa(iterations) + " iterations):\n"
	for j := range iterations {
		res += "Iteration " + strconv.Itoa(j+1) + ":\n"
		summ := 0
		for i := range N {
			number := r.Intn(M) + 1
			summ += number
			res += "  Dice " + strconv.Itoa(i+1) + ": " + strconv.Itoa(number) + "\n"
		}
		summStr := strconv.Itoa(summ)
		avgStr := fmt.Sprint(float64(summ) / float64(N))
		res += "  Sum: " + summStr + "\n"
		res += "  Average: " + avgStr + "\n"
	}

	return res, nil
}

func RollWithIterations(r *rand.Rand, iterations int, dice string) (string, error) {
	res := ""
	for i := range iterations {
		res += "Iteration " + strconv.Itoa(i+1) + ":\n"
		rollStr, err := Roll(r, dice)
		if err != nil {
			return "", err
		}
		res += "  " + rollStr + "\n"
	}

	return res, nil
}

func RollWithOnlySum(r *rand.Rand, dice string) (int, error) {
	N, M, err := Parse(dice)
	if err != nil {
		return 0, err
	}
	summ := 0
	for range N {
		number := r.Intn(M) + 1
		summ += number
	}

	return summ, nil
}

func RollSumWithIterations(r *rand.Rand, iterations int, dice string) (string, error) {
	N, M, err := Parse(dice)
	if err != nil {
		return "", err
	}
	res := ""
	for i := range iterations {
		summ := 0
		for range N {
			number := r.Intn(M) + 1
			summ += number
		}
		res += "Iteration " + strconv.Itoa(i+1) + ":\n"
		res += "  Sum: " + strconv.Itoa(summ) + "\n"
	}

	return res, nil
}

func Seed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func Parse(dice string) (int, int, error) {
	parts := strings.Split(dice, "d")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Неверный формат ввода")
	}
	N, errorN := strconv.Atoi(parts[0])
	M, errorM := strconv.Atoi(parts[1])
	if errorN != nil || N <= 0 || N > 100 {
		return 0, 0, fmt.Errorf("Неверное число костей")
	}
	if errorM != nil || M <= 0 || M > 100 {
		return 0, 0, fmt.Errorf("Неверное число граней")
	}

	return N, M, nil
}

func PrintHelp() {
	fmt.Println("CLI-утилита для броска игральных костей")
	fmt.Println()
	fmt.Println("Флаги:")
	fmt.Println("-n, --count <число> - выполняет указанное количество итераций бросков,")
	fmt.Println("выводит информацию отдельно по каждой итерации")
	fmt.Println()
	fmt.Println("--seed <число> - задаёт зерно (seed) для генератора случайных чисел;")
	fmt.Println("при одинаковом seed результаты бросков полностью воспроизводимы")
	fmt.Println()
	fmt.Println("-s, --sum - выводит для каждой итерации только сумму значений костей, несовместим с флагом -v")
	fmt.Println()
	fmt.Println("-v, --verbose - выводит подробную информацию по каждому броску внутри итерации,")
	fmt.Println("а также среднее значение среди всех итераций, несовместим с флагом -s")
	fmt.Println()
	fmt.Println("-h, --help - показать справку")
	fmt.Println()
}
