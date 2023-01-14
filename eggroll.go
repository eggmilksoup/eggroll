package main

import "crypto/rand"
import "fmt"
import "math/big"
import "os"
import "strconv"
import "strings"
import "time"
import "unicode"

import "github.com/bwmarrin/discordgo"

func solve(expr []string) string {
	if len(expr) == 0 {
		return "Error: empty expression"
	}

	if expr[len(expr) - 1] == ")" {
		i := len(expr) - 1
		for j := i - 1; j > -1; j -- {
			if expr[j] == "(" {
				var pre []string
				if j != 0 {
					pre = expr[:j - 1]
				}
				var post []string
				if i + 1 < len(expr) {
					post = expr[i + 2:]
				}
				expr = append(pre, solve(expr[j+1:i]))
				for _, str := range post {
					expr = append(expr, str)
				}
				return solve(expr)
			}
			if expr[j] == ")" {
				i = j
			}
		}
		return "Error: mismatched parentheses"
	}
	if unicode.IsDigit(rune(expr[len(expr) - 1][0])) {
		if len(expr) == 1 {
			return expr[0]
		}
		switch expr[len(expr) - 2] {
			case "+":
				if len(expr) == 2 {
					return "Error: + takes 2 arguments"
				}

				second, err := strconv.Atoi(expr[len(expr) - 1])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}
				str := solve(expr[:len(expr) - 2])
				first, err := strconv.Atoi(str)
				if err != nil {
					return str
				}
				return fmt.Sprintf("%d", first + second)

			case "-":
				if len(expr) == 2 {
					num, err := strconv.Atoi(expr[1])
					if err != nil {
						return "Error: strconv.Atoi: " + err.Error()
					}

					return fmt.Sprintf("%d", -1 * num)
				}

				second, err := strconv.Atoi(expr[len(expr) - 1])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}

				str := solve(expr[:len(expr) - 2])
				first, err := strconv.Atoi(str)
				if err != nil {
					return str
				}
				return fmt.Sprintf("%d", first - second)

			case "*":
				if len(expr) == 2 {
					return "Error: * takes 2 arguments"
				}

				if len(expr) > 3 {
					if expr[len(expr) - 4] == "+" ||
							expr[len(expr) - 4] == "-" {
						
						first, err := strconv.Atoi(expr[len(expr) - 3])
						if err != nil {
							return "Error: strconv.Atoi: " + err.Error()
						}

						second, err := strconv.Atoi(expr[len(expr) - 1])
						if err != nil {
							return "Error: strconv.Atoi: " + err.Error()
						}

						return solve( append(
							expr[:len(expr) - 3],
							fmt.Sprintf("%d", first * second)))
					}
				}

				second, err := strconv.Atoi(expr[len(expr) - 1])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}

				str := solve(expr[:len(expr) - 2])
				first, err := strconv.Atoi(str)
				if err != nil {
					return str
				}
				return fmt.Sprintf("%d", first * second)

			case "/":
				if len(expr) == 2 {
					return "Error: / takes 2 arguments"
				}

				if len(expr) > 3 {
					if expr[len(expr) - 4] == "+" ||
							expr[len(expr) - 4] == "-" {
						
						first, err := strconv.Atoi(expr[len(expr) - 3])
						if err != nil {
							return "Error: strconv.Atoi: " + err.Error()
						}

						second, err := strconv.Atoi(expr[len(expr) - 1])
						if err != nil {
							return "Error: strconv.Atoi: " + err.Error()
						}

						return solve(append(
							expr[:len(expr) - 3],
							fmt.Sprintf("%d", first / second)))
					}
				}

				second, err := strconv.Atoi(expr[len(expr) - 1])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}

				str := solve(expr[:len(expr) - 2])
				first, err := strconv.Atoi(str)
				if err != nil {
					return str
				}
				return fmt.Sprintf("%d", first / second)

			case "d":
				if len(expr) == 2 {
					num, err := strconv.Atoi(expr[1])
					if err != nil {
						return "Error: strconv.Atoi: " + err.Error()
					}

					res, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
					if err != nil {
						return "Error: rand.Int: " + err.Error()
					}

					res.Add(res, big.NewInt(1))

					return res.Text(10)
				}

				if len(expr) > 3 && expr[len(expr) - 3] == ")" {
					j := len(expr) - 3
					for i := j - 1; i > -1; i -- {
						if expr[i] == "(" {
							var pre []string
							if j != 0 {
								pre = expr[:j - 1]
							}
							var post []string
							if i + 1 != len(expr) {
								post = expr[i + 2:]
							}
							expr = append(pre, solve(expr[i+1:j]))
							for _, str := range post {
								expr = append(expr, str)
							}
							return solve(expr)
						}
						if expr[i] == ")" {
							j = i
						}
					}
					return "Error: unmatched parentheses"
				}

				count, err := strconv.Atoi(expr[len(expr) - 3])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}

				num, err := strconv.Atoi(expr[len(expr) - 1])
				if err != nil {
					return "Error: strconv.Atoi: " + err.Error()
				}

				sum := big.NewInt(0)
				for i := 1; i <= count; i ++ {
					res, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
					if err != nil {
						return "Error: rand.Int: " + err.Error()
					}
					res.Add(res, big.NewInt(1))
					sum.Add(sum, res)
				}

				return sum.Text(10)
		}
	}

	return "Error: could not parse expression"
}

func main() {

	discord, err := discordgo.New("Bot " + os.Getenv("key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "discordgo.New: %s\n", err.Error())
		os.Exit(1)
	}

	err = discord.Open()

	if err != nil {
		fmt.Fprintf(os.Stderr, "discordgo.Session.Open: %s\n", err.Error())
		os.Exit(1)
	}

	discord.AddHandler( func(
			discord *discordgo.Session,
			event *discordgo.MessageCreate) {

		words := strings.Fields(event.Message.Content)
		if len(words) > 0 && words[0] == "/roll" {

			var expr []string
			for i := 6; i < len(event.Message.Content); i++ {

				switch rune(event.Message.Content[i]) {
					case ' ':
					case '0': fallthrough
					case '1': fallthrough
					case '2': fallthrough
					case '3': fallthrough
					case '4': fallthrough
					case '5': fallthrough
					case '6': fallthrough
					case '7': fallthrough
					case '8': fallthrough
					case '9':
						j := i + 1
						for ; j < len(event.Message.Content) && 
							unicode.IsDigit(rune(event.Message.Content[j]));
							j ++ {}

						expr = append(expr, event.Message.Content[i:j])
						i = j - 1
					case '(': fallthrough
					case ')': fallthrough
					case '+': fallthrough
					case '-': fallthrough
					case '*': fallthrough
					case '/': fallthrough
					case 'd':
						expr = append(expr, event.Message.Content[i:i + 1])
					default:
						discord.ChannelMessageSend(
							event.Message.ChannelID,
							"Unrecognized symbol: " +
							event.Message.Content[i:i + 1])
						return
				}

			}

			discord.ChannelMessageSend(event.Message.ChannelID, solve(expr))
		}

	})

	for true {
		time.Sleep(1<<63 - 1)
	}
}
