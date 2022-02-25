package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cmp interface {
}

type Strategy struct {
	guess    string
	depth    int
	children map[cmp]Strategy
}

func loadTextFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

func main() {
	words := loadTextFile("wordle_legal_words.txt")
	legal := loadTextFile("wordle_legal_guesses.txt")
	// Take the first 10 words
	fmt.Println("Solving...")
	sol := solve(words, legal)
	fmt.Println("Solved with at most", sol.depth, "guesses!")
	for {
		tree := sol
		fmt.Println()
		for {
			fmt.Println("Guess", tree.guess)
			var pattern string
			fmt.Print("Enter the pattern: ")
			fmt.Scanln(&pattern)
			fmt.Println()
			if pattern == "eeeee" {
				fmt.Println("You win!")
				break
			}
			tree = tree.children[pattern]
		}
	}
}

func compare(guess, word string) cmp {
	//fmt.Println(guess, word)
	ans := []string{}

	letters := make(map[rune]int)

	for _, c := range word {
		if _, ok := letters[c]; ok {
			letters[c]++
		} else {
			letters[c] = 1
		}
	}
	for i, c := range word {
		if c == rune(guess[i]) {
			ans = append(ans, "e")
			letters[c]--
		} else {
			ans = append(ans, ".")
		}
	}

	for i, c := range guess {
		if x, ok := letters[c]; ok && x > 0 {
			ans[i] = "c"
			letters[c]--
		}
	}

	return strings.Join(ans, "")
}

func solveHelper(words, legal []string, limit int) (int, string) {
	if len(words) <= 2 {
		return len(words), words[0]
	}
	if limit <= 1 {
		return 100, ""
	}
	// fmt.Print(len(words), " ")

	bestScore := limit
	bestGuess := ""
	for _, guess := range legal {
		m := make(map[cmp][]string)
		p := false
		for _, w := range words {
			if w == guess {
				p = true
				continue
			}

			c := compare(guess, w)

			if _, ok := m[c]; !ok {
				m[c] = []string{w}
			} else {
				m[c] = append(m[c], w)
			}
		}
		if len(m) == 1 && !p {
			continue
		}
		score := 0
		for _, val := range m {
			var d int
			if len(val) >= 3 {
				d = len(val)
			} else {
				z, _ := solveHelper(val, legal, bestScore-1)
				d = z + 1
			}
			if d > score {
				score = d
			}

			if score >= limit {
				continue
			}
		}

		if score < bestScore {
			if score <= 2 {
				return score, guess
			}
			bestScore = score
			bestGuess = guess

		}
	}
	// if bestGuess == "" {
	//   panic("no guess")
	// }
	return bestScore, bestGuess
}

func solve(words, legal []string) (strat Strategy) {
	depth, guess := solveHelper(words, legal, 1000)
	strat.depth = 1
	strat.guess = guess
	strat.children = make(map[cmp]Strategy)
	if depth == 1 {
		return
	}

	m := make(map[cmp][]string)

	for _, w := range words {
		if w == guess {
			continue
		}
		c := compare(guess, w)
		if _, ok := m[c]; !ok {
			m[c] = []string{w}
		} else {
			m[c] = append(m[c], w)
		}
	}

	for key, val := range m {
		strat.children[key] = solve(val, legal)
		if strat.children[key].depth >= strat.depth {
			strat.depth = strat.children[key].depth + 1
		}
	}

	return
}
