package confy

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SetFields(v reflect.Value, ss []string) any {
	n := 0
	for i := 0; i < v.NumField(); i++ {
		s := ss[n]
		vf := v.Field(i)
		if vf.Kind().String() == "struct" {
			n = SetStruct(vf, n, ss)
		}
		if vf.Kind().String() == "map" {
			n = SetMap(vf, n, ss)
		} else {
			SetField(vf, s)
			n++
			if n == len(ss) {
				break
			}
		}
	}
	return v.Interface()
}

func NumIndt(s string) int {
	for i := 6; i > 0; i-- {
		ind := strings.Repeat(Indt, i)
		if strings.HasPrefix(s, ind) {
			return i
		}
	}
	return 0
}

func SetMap(v reflect.Value, n int, ss []string) int {
	n++
	m := reflect.Indirect(reflect.MakeMap(v.Type()))
	tk := v.Type().Key()
	tv := v.Type().Elem()
	for n < len(ss) {
		sls := strings.Split(ss[n], ":")
		if len(sls) < 2 || sls[1] == "" || sls[1] == "\n" || strings.HasPrefix(sls[1], " ") {
			break
		}
		key := reflect.Indirect(reflect.New(tk))
		SetField(key, sls[0])
		val := reflect.Indirect(reflect.New(tv))
		SetField(val, sls[1])
		m.SetMapIndex(key, val)
		n++
	}
	v.Set(m)
	return n
}

func SetStruct(v reflect.Value, n int, ss []string) int {
	n++
	for j := 0; j < v.NumField(); j++ {
		if v.Field(j).Kind().String() != "struct" {
			SetField(v.Field(j), ss[n])
			n++
		} else {
			n = SetStruct(v.Field(j), n, ss) + 1
		}
	}
	return n - 1
}

func SetField(v reflect.Value, s string) {
	switch v.Kind().String() {
	case "string":
		v.SetString(s)
	case "int", "int8", "int16", "int32", "int64":
		var n int64
		if v.Type() == reflect.TypeOf(time.Hour) {
			td, _ := time.ParseDuration(s)
			n = int64(td)
		} else {
			n, _ = strconv.ParseInt(s, 10, 64)
		}
		v.SetInt(n)
	case "uint8", "uint16", "uint32", "uint64":
		n, _ := strconv.ParseUint(s, 10, 64)
		v.SetUint(n)
	case "float32", "float64":
		f, _ := strconv.ParseFloat(s, 64)
		v.SetFloat(f)
	case "bool":
		b, _ := strconv.ParseBool(s)
		v.SetBool(b)
	case "slice":
		v.Set(SetSlice(v, s))
	}
}

func SetSlice(v reflect.Value, ss string) reflect.Value {
	t := v.Type().Elem()
	nv := reflect.New(reflect.SliceOf(t)).Elem()
	sls := strings.Trim(ss, "[]")
	if t.Kind().String() == "struct" {
		return SetStructElem(nv, sls)
	}
	for _, s := range strings.Split(sls, ",") {
		ne := reflect.New(t).Elem()
		SetField(ne, s)
		nv = reflect.Append(nv, ne)
	}
	return nv
}

func SetStructElem(v reflect.Value, s string) reflect.Value {
	sls := strings.Trim(s, "{}")
	t := v.Type().Elem()
	ss := strings.Split(sls, "},{")
	for i := 0; i < len(ss); i++ {
		ne := reflect.New(t).Elem()
		se := strings.Split(ss[i], ",")
		for j := 0; j < ne.NumField(); j++ {
			SetField(ne.Field(j), se[j])
		}
		v = reflect.Append(v, ne)
	}
	return v
}

func ReadConfig(path string) ([]string, error) {
	bts, err := os.ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	sf := strings.Split(string(bts), "\n")
	sv := []string{}
	for _, ln := range sf {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "#") {
			continue
		}
		if len(ln) < 2 {
			break
		}
		slc := strings.Fields(ln)
		val := ""
		if len(slc) == 1 {
			val = slc[0]
		} else {
			val = strings.Join(slc[2:], " ")
		}
		sv = append(sv, val)
	}

	return sv, nil
}

func LoadConfy(c any, path string) (any, error) {
	if !ConfigExists(path) {
		return nil, errors.New("config does not exist")
	}
	ss, err := ReadConfig(path)
	if err != nil {
		return nil, err
	}
	t := reflect.TypeOf(c)
	v := reflect.Indirect(reflect.New(t))
	inf := SetFields(v, ss)
	return inf, nil
}
