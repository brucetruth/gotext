package pos

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

const (
	startTag string = "<s>"
	endTag   string = "</s>"
)

// link struct stores information of two tokens and the relation between both.
type link struct {
	current     string
	previous    string // word, tag (emission) - tag, tag (transition)
	occurrences float64
	weight      float64
}

// List of links to sort it.
type links []*link

func (l links) Len() int           { return len(l) }
func (l links) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l links) Less(i, j int) bool { return l[i].weight < l[j].weight }

// getLink function search is called by links struct and search in whole links
// inside them to find a relation between current and previous provided.
func (l links) getLink(current, previous string) (*link, bool) {
	if len(l) > 0 {
		for _, sl := range l {
			if sl.current == current && sl.previous == previous {
				return sl, true
			}
		}
	}

	return &link{current: current, previous: previous, occurrences: 1}, false
}

// Model define a trained model based on tagged corpus. Contains emissions and
// transitions tables and tag list available.
type Model struct {
	tags        []string
	transitions links
	emissions   links
}

// LoadModel function returns model instance opening transitions and emissions
// tables based on path provided.
func LoadModel(p string) (m *Model, err error) {
	if l, err := filepath.Abs(p); err != nil {
		return m, err
	} else if _, err = os.Open(l); err != nil {
		return m, err
	}

	var (
		tp string = fmt.Sprintf("%s/transitions", p)
		ep string = fmt.Sprintf("%s/emissions", p)
	)

	m = &Model{}
	if err = m.loadTransitions(tp); err != nil {
		return m, err
	}

	if err = m.loadEmissions(ep); err != nil {
		return m, err
	}
	return m, err
}

// loadTransitions function opens transition table file associated to current
// model. Then parses and generates links each line.
func (m *Model) loadTransitions(p string) (err error) {
	var re *regexp.Regexp = regexp.MustCompile(`\t`)

	var tfd *os.File
	if tfd, err = os.Open(p); err != nil {
		return err
	}
	defer tfd.Close()

	var sc *bufio.Scanner = bufio.NewScanner(tfd)
	for sc.Scan() {
		var ln string = sc.Text()
		var data []string = re.Split(ln, -1)
		if len(data) == 3 {
			var w float64
			if w, err = strconv.ParseFloat(data[2], 64); err != nil {
				return err
			}

			m.transitions = append(m.transitions, &link{previous: data[0], current: data[1], weight: w})
		}
	}
	return nil
}

// loadTransitions function opens emission table file associated to current
// model. Then parses and generates links each line.
func (m *Model) loadEmissions(p string) (e error) {
	var re *regexp.Regexp = regexp.MustCompile(`\t`)
	var efd *os.File
	if efd, e = os.Open(p); e != nil {
		return e
	}
	defer efd.Close()

	var sc *bufio.Scanner = bufio.NewScanner(efd)
	for sc.Scan() {
		var line string = sc.Text()
		var data []string = re.Split(line, -1)
		if len(data) == 3 {
			var w float64
			if w, e = strconv.ParseFloat(data[2], 64); e != nil {
				return e
			}

			m.emissions = append(m.emissions, &link{data[1], data[0], 0, w})
		}
	}
	return nil
}

// probs function calculate word possibilities based on previous tag, with
// transmission and emission costs using Model provided. If model doesn't have
// emission record for current word, return proposed tag with '?' after.
func (m *Model) probs(cw, pt string) (ps map[string]float64, sg string) {
	var ts links
	for _, t := range m.transitions {
		if t.previous == pt {
			ts = append(ts, t)
		}
	}

	var es links
	for _, e := range m.emissions {
		if e.current == cw {
			es = append(es, e)
		}
	}

	ps = make(map[string]float64, len(ts))
	for _, e := range es {
		var s float64 = e.weight
		for _, t := range ts {
			if e.current == t.previous {
				s += t.weight
			}
		}
		ps[e.previous] = s
	}

	if len(ps) == 0 {
		var _t string = startTag
		var max float64
		for _, t := range ts {
			if t.weight > max {
				_t = t.current
				max = t.weight
			}
		}

		sg = fmt.Sprintf("%s?", _t)
	}

	return ps, sg
}

// Train function trains Model with corpus provided and generates transitions
// and emissions tables. Receives corpus path and return Model instance.
func Train(p string) (m *Model, err error) {
	if l, err := filepath.Abs(p); err != nil {
		return m, err
	} else if fd, err := os.Open(l); err != nil {
		return m, err
	} else {
		defer fd.Close()

		m = &Model{tags: []string{}}
		var (
			data []sentence
			rs   *regexp.Regexp = regexp.MustCompile(`\s|\t`)
			rtg  *regexp.Regexp = regexp.MustCompile(`(.+)/(.+)`)
		)

		var sc *bufio.Scanner = bufio.NewScanner(fd)
		for sc.Scan() {
			var ln string = sc.Text()
			var cdts []string = rs.Split(ln, -1)

			var s sentence
			for i, cdt := range cdts {
				if g := rtg.FindStringSubmatch(cdt); len(g) > 1 {
					var (
						r  string = g[1]
						tg string = g[2]
					)

					var in bool = false
					for _, t := range m.tags {
						in = in || t == tg
					}
					if !in {
						m.tags = append(m.tags, tg)
					}

					s = append(s, &token{i, r, tg})
				}
			}

			if len(s) > 0 {
				data = append(data, s)
			}
		}

		if err := sc.Err(); err != nil {
			return m, err
		}
		m.score(data)
	}
	return m, err
}

// score function calculates transitions and emissions for untrained corpus
// provided.
func (m *Model) score(data []sentence) {
	var (
		ts  links
		es  links
		ctx map[string]float64 = make(map[string]float64, len(m.tags)+2)
	)

	for _, s := range data {
		var prev string = startTag
		ctx[startTag]++

		sort.Sort(s)
		for _, t := range s {
			if t, ok := ts.getLink(t.tag, prev); ok {
				t.occurrences++
			} else {
				ts = append(ts, t)
			}

			if e, ok := es.getLink(t.raw, t.tag); ok {
				e.occurrences++
			} else {
				es = append(es, e)
			}

			ctx[t.tag]++
			prev = t.tag
		}

		if t, exists := ts.getLink(prev, endTag); exists {
			t.occurrences++
		} else {
			ts = append(ts, t)
		}
		ctx[endTag]++
	}

	// Normalize weights
	for _, t := range ts {
		t.weight = t.occurrences / ctx[t.previous]
	}
	m.transitions = ts

	for _, e := range es {
		e.weight = e.occurrences / ctx[e.previous]
	}
	m.emissions = es
}

// Store function saves trained Model locally. Creates tabbed separated file
// with transitions and emissions and each weight.
func (m *Model) Store(o string) (err error) {
	var l string
	if l, err = filepath.Abs(o); err == nil {
		if err = os.Mkdir(l, os.ModePerm); err != nil {
			return err
		}

		var (
			tp       string = fmt.Sprintf("%s/transitions", l)
			ep       string = fmt.Sprintf("%s/emissions", l)
			fdt, fde *os.File
		)

		if fdt, err = os.Create(tp); err == nil {
			defer fdt.Close()

			for _, t := range m.transitions {
				var ln string = fmt.Sprintf("%s\t%s\t%g\n", t.previous, t.current, t.weight)
				if _, err = fdt.WriteString(ln); err != nil {
					return err
				}
			}
		}

		if fde, err = os.Create(ep); err == nil {
			defer fde.Close()

			for _, e := range m.emissions {
				var ln string = fmt.Sprintf("%s\t%s\t%g\n", e.previous, e.current, e.weight)
				if _, err = fde.WriteString(ln); err != nil {
					return err
				}
			}
		}
	}
	return err
}
