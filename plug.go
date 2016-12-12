package zero

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tj/go-debug"
)

type Plug struct {
	str *Str
	num *Num
	arr *Arr

	debug debug.DebugFunction
}

func (p *Plug) Add(args *Args, reply *int) (err error) {
	p.debug("adding value of type %s", args.Type)
	switch args.Type {
	case String:
		p.debug("parsing for string value")
		if v, ok := args.Value.(string); ok {
			p.str.Add(args.Key, v)
			add(args.Key, String)
			break
		}

		return fmt.Errorf("invalid value")
	case Number:
		p.debug("parsing for number value")
		n, err := strconv.ParseInt(fmt.Sprintf("%v", args.Value), 10, 0)
		if err != nil {
			return err
		}

		p.num.Add(args.Key, n)
		add(args.Key, Number)
		break
	case Array:
		p.debug("parsing for array values")
		v, ok := args.Value.(string)
		if !ok {
			return fmt.Errorf("invalid value")
		}

		raw := strings.TrimSpace(v)
		raw = raw[1 : len(raw)-1]

		arr := []interface{}{}
		for _, s := range strings.Split(raw, ",") {
			arr = append(arr, interface{}(strings.TrimSpace(s)))
		}
		if err = p.arr.Add(args.Key, arr); err != nil {
			p.debug("error while adding to array store")
			return err
		}
		add(args.Key, Array)
		break
	default:
		return fmt.Errorf("store not found")
	}
	return
}

func (p *Plug) Push(args *Args, reply *int) (err error) {
	kind, err := which(args.Key)
	if err != nil {
		return
	}

	if kind != Array {
		return fmt.Errorf("push only supported for array type")
	}

	l, err := p.arr.Push(args.Key, args.Value)
	if err != nil {
		return
	}

	*reply = l
	return
}

func (p *Plug) Pop(args *Args, reply *string) (err error) {
	kind, err := which(args.Key)
	if err != nil {
		return
	}

	if kind != Array {
		return fmt.Errorf("pop only supported for array type")
	}

	val, err := p.arr.Pop(args.Key)
	if err != nil {
		return
	}

	*reply = fmt.Sprintf("%v", val)
	return
}

func (p *Plug) Show(args *Args, reply *string) (err error) {
	if args.Key == "*" {
		num := p.num.All()
		str := p.str.All()
		arr := p.arr.All()

		all := make([][]string, len(num)+len(str)+len(arr))
		all = append(all, num...)
		all = append(all, str...)
		all = append(all, arr...)

		var norm []string
		for _, r := range all {
			if len(r) > 0 {
				p.debug("row: %v", r)
				norm = append(norm, strings.Join(r, ColumnSplitter))
			}
		}

		*reply = strings.Join(norm, RowSplitter)
		return
	}

	kind, err := which(args.Key)
	if err != nil {
		return
	}
	switch kind {
	case String:
		r, err := p.str.Show(args.Key)
		if err != nil {
			return err
		}

		*reply = r
		break
	case Number:
		r, err := p.num.Show(args.Key)
		if err != nil {
			return err
		}
		*reply = fmt.Sprintf("%v", r)
		break
	case Array:
		r, err := p.arr.Show(args.Key)
		if err != nil {
			return err
		}
		*reply = fmt.Sprintf("%v", r)
		break
	default:
		return fmt.Errorf("store not found")
	}
	return
}

func (p *Plug) Count(args *Args, reply *int) (err error) {
	c := 0
	c += p.num.Count()
	c += p.str.Count()
	c += p.arr.Count()

	*reply = c
	return
}

func (p *Plug) Keys(args *Args, reply *[]Key) (err error) {
	keys := []Key{}
	keys = append(keys, p.str.Keys()...)
	keys = append(keys, p.num.Keys()...)
	keys = append(keys, p.arr.Keys()...)

	*reply = keys
	return
}

func (p *Plug) Del(args *Args, reply *int) (err error) {
	kind, err := which(args.Key)
	if err != nil {
		return
	}
	switch kind {
	case String:
		err = p.str.Del(args.Key)
		if err != nil {
			return
		}
		break
	case Number:
		err = p.num.Del(args.Key)
		if err != nil {
			return
		}
		break
	case Array:
		err = p.arr.Del(args.Key)
		if err != nil {
			return
		}
		break
	default:
		return fmt.Errorf("store not found")
	}
	delete(index, args.Key)
	return
}

func NewPlug() *Plug {
	return &Plug{
		str: NewStr(),
		num: NewNum(),
		arr: NewArr(),

		debug: debug.Debug("zero:plug"),
	}
}
