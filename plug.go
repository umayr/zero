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

	debug debug.DebugFunction
}

func (p *Plug) Add(args *Args, reply *int) (err error) {
	p.debug("adding value of type %s", args.Type)

	switch args.Type {
	case String:
		if v, ok := args.Value.(string); ok {
			p.str.Add(args.Key, v)
			add(args.Key, String)
			break
		}

		return fmt.Errorf("invalid value")
	case Number:
		n, err := strconv.ParseInt(fmt.Sprintf("%v", args.Value), 10, 0)
		if err != nil {
			return err
		}

		p.num.Add(args.Key, n)
		add(args.Key, Number)
		break
	case Array:
		break
	default:
		return fmt.Errorf("store not found")
	}
	return
}

func (p *Plug) Show(args *Args, reply *string) (err error) {
	if args.Key == "*" {
		n := p.num.All()
		s := p.str.All()

		all := make([][]string, len(n) + len(s))
		all = append(all, n...)
		all = append(all, s...)

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

	*reply = c
	return
}

func (p *Plug) Keys(args *Args, reply *[]Key) (err error) {
	keys := []Key{}
	keys = append(keys, p.str.Keys()...)
	keys = append(keys, p.num.Keys()...)

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

		debug: debug.Debug("zero:plug"),
	}
}
