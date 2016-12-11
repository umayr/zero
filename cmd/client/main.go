package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/umayr/zero"
)

const (
	Add   = "ADD"
	Push  = "PUSH"
	Pop   = "POP"
	Show  = "SHOW"
	Keys  = "KEYS"
	Count = "COUNT"
	Exit  = "EXIT"
	Del   = "DEL"
)

type call func(string, interface{}, interface{}) error
type fn func(call, []string) (string, error)

var (
	regNum = regexp.MustCompile(`^(-?\d+)$`)
	regArr = regexp.MustCompile(`^(\[.*\])$`)
)

func add(c call, args []string) (ln string, err error) {
	if len(args) <= 1 {
		return "", fmt.Errorf("must provide a key and value")
	}

	var kind string
	v := strings.Join(args[1:], " ")

	if regNum.MatchString(v) {
		kind = zero.Number
	} else if regArr.MatchString(v) {
		kind = zero.Array
	} else {
		kind = zero.String
	}

	err = c("store.Add", &zero.Args{
		Key:   zero.Key(args[0]),
		Value: v,
		Type:  kind,
	}, nil)

	ln = "OK"
	return
}

func push(c call, args []string) (ln string, err error) {
	if len(args) <= 1 {
		return "", fmt.Errorf("must provide a key and value")
	}

	var reply int
	v := strings.Join(args[1:], " ")

	if err = c("store.Push", &zero.Args{
		Key:   zero.Key(args[0]),
		Value: v,
	}, &reply); err != nil {
		return
	}

	ln = fmt.Sprintf("%v", reply)
	return
}

func pop(c call, args []string) (ln string, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("must provide a key")
	}

	var reply string
	if err = c("store.Pop", &zero.Args{
		Key: zero.Key(args[0]),
	}, &reply); err != nil {
		return
	}

	ln = fmt.Sprintf("%v", reply)
	return
}

func show(c call, args []string) (ln string, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("must provide a key")
	}

	var reply string
	if err = c("store.Show", &zero.Args{
		Key: zero.Key(args[0]),
	}, &reply); err != nil {
		return
	}

	if args[0] != "*" {
		ln = fmt.Sprintf("%v", reply)
		return
	}

	if reply == "" {
		return
	}

	data := [][]string{}
	for _, r := range strings.Split(reply, zero.RowSplitter) {
		data = append(data, strings.Split(r, zero.ColumnSplitter))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value", "Type", "Time"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	return
}

func keys(c call, args []string) (ln string, err error) {
	var keys []zero.Key
	err = c("store.Keys", &zero.Args{}, &keys)
	if err != nil {
		return
	}

	ln = fmt.Sprintf("%v", keys)
	return
}

func count(c call, args []string) (ln string, err error) {
	var cnt int
	if err = c("store.Count", &zero.Args{}, &cnt); err != nil {
		return
	}

	ln = fmt.Sprintf("%v", cnt)
	return
}

func del(c call, args []string) (ln string, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("must provide a key")
	}
	if err = c("store.Del", &zero.Args{
		Key: zero.Key(args[0]),
	}, nil); err != nil {
		return
	}

	ln = "OK"
	return

}

func main() {

	client, err := rpc.DialHTTP("tcp", ":7161")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		var f fn

		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		cmd, args := parts[0], parts[1:]

		switch strings.ToUpper(cmd) {
		case Add:
			f = add
			break
		case Push:
			f = push
			break
		case Pop:
			f = pop
			break
		case Show:
			f = show
			break
		case Keys:
			f = keys
			break
		case Count:
			f = count
			break
		case Del:
			f = del
			break
		case Exit:
			os.Exit(0)
		default:
			fmt.Println("err: invalid command")
			continue
		}

		ln, err := f(client.Call, args)
		if err != nil {
			fmt.Println(fmt.Sprintf("err: %s", err.Error()))
			continue
		}
		fmt.Println(ln)
	}
}
