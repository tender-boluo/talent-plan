package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})

	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

func URLCountMap(filename string, contents string) []KeyValue {
	urls := strings.Split(string(contents), "\n")
	urlCnt := make(map[string]int)
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if len(url) == 0 {
			continue
		}
		urlCnt[url]++
	}
	KeyValuePairs := make([]KeyValue, 0, len(urlCnt))
	for url, cnt := range urlCnt {
		KeyValuePairs = append(KeyValuePairs, KeyValue{Key: url, Value: strconv.Itoa(cnt)})
	}
	return KeyValuePairs
}

func URLCountReduce(key string, values []string) string {
	cnt := 0
	for _, value := range values {
		value, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		cnt += value
	}
	return fmt.Sprintf("%s %s\n", key, strconv.Itoa(cnt))
}

func URLTop10Map(filename string, contents string) []KeyValue {
	block := 10
	lines := strings.Split(string(contents), "\n")
	ucs := make([]*urlCount, 0, len(lines))
	KeyValuePairs := make([]KeyValue, 0, block)
	for _, value := range lines {
		tmp := strings.Split(value, " ")
		if tmp[0] != "" {
			cnt, err := strconv.Atoi(tmp[1])
			if err != nil {
				panic(err)
			}
			ucs = append(ucs, &urlCount{tmp[0], cnt})
		}
	}
	sort.Slice(ucs, func(i, j int) bool {
		if ucs[i].cnt == ucs[j].cnt {
			return ucs[i].url < ucs[j].url
		}
		return ucs[i].cnt > ucs[j].cnt
	})
	for i, u := range ucs {
		if i == block {
			break
		}
		KeyValuePairs = append(KeyValuePairs, KeyValue{Key: "", Value: fmt.Sprintf("%s %s", u.url, strconv.Itoa(u.cnt))})
	}
	return KeyValuePairs
}

func URLTop10Reduce(key string, values []string) string {
	cnts := make(map[string]int, len(values))
	for _, v := range values {
		v := strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		cnts[tmp[0]] = n
	}

	us, cs := TopN(cnts, 10)
	buf := new(bytes.Buffer)
	for i := range us {
		fmt.Fprintf(buf, "%s: %d\n", us[i], cs[i])
	}
	return buf.String()
}
