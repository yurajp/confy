package confy

import (
	"fmt"
	"reflect"
	"strings"
	"os"
)

var (
	Path = "config/conf.ini"
	Indt = "   "
)

type Vall struct {
	Val reflect.Value
	Level int
}

func SetPath(p string) {
	Path = p
}

func SetIndent(n int) {
	i := strings.Repeat(" ", n)
	Indt = i
}

func Indent(s string, l int) string {
	indent := strings.Repeat(Indt, l)
	return indent + s
}

func NumStruct(sl []string) int {
	i := 0
	for _, ln := range sl {
		if ln == "@" {
			i++
		}
	}
	return i
}

func SliceFmt(a reflect.Value) string {
	str := "["
	view := fmt.Sprintf("%v", a.Interface())
	
	for _, e := range strings.Fields(strings.Trim(view, "[]")) {
		str += fmt.Sprintf("%s,", e)
	}
	str = strings.TrimSuffix(str, ",") + "]"
	return str
}

func MapToStruct(v Vall) ([]string, []Vall) {
	Info := []string{}
	Vals := []Vall{}
	it := v.Val.MapRange()
	for it.Next() {
		ln := fmt.Sprintf("%v:%v\n", it.Key(), it.Value())
		Info = append(Info, Indent(ln, v.Level))
	}
	return Info, Vals
}


func ValToStruct(v Vall) ([]string, []Vall) {
	if v.Val.Kind() == reflect.Map {
		return MapToStruct(v)
	}
	Info := []string{}
	Vals := []Vall{}
	for i := 0; i < v.Val.NumField(); i++ {
		line := ""
		iv := ""
	  name := v.Val.Type().Field(i).Name + " "
	  ftype := v.Val.Field(i).Type()
	  kind := ftype.Kind().String()
	  if kind == "struct" {
		  tp := " " + strings.TrimPrefix(ftype.String(), "main.") + ":\n"
		  line = name + kind + tp
		  Info = append(Info, Indent(line, v.Level))
		  Info = append(Info, "@")
		  nval := reflect.Indirect(reflect.ValueOf(v.Val.Field(i).Interface()))
		  Vals = append(Vals, Vall{nval, v.Level+1})
		  continue
	  }
	  if kind == "map" {
	  	mt := " " + ftype.String() + "\n"
	  	line = name + kind + mt
	  	Info = append(Info, Indent(line, v.Level))
	  	Info = append(Info, "@")
	  	nval := reflect.Indirect(reflect.ValueOf(v.Val.Field(i).Interface()))
	  	Vals = append(Vals, Vall{nval, v.Level + 1})
	  	continue
	  }
	  if kind == "slice" {
	  	iv = SliceFmt(v.Val.Field(i)) + "\n"
	  	sk := strings.TrimPrefix(ftype.String(), "[]")
	  	sk = strings.TrimPrefix(sk, "main.")
	  	kind = "[]" + sk
	  } else {
	  	iv = fmt.Sprintf("%v \n", v.Val.Field(i).Interface())
	  }
		line = name + kind + ": "
		line += iv
	  Info = append(Info, Indent(line, v.Level))
	}
	return Info, Vals
}


func BuildString(v Vall) string {
 	Info, Vals := ValToStruct(v)
 	m := len(Info)
  lens := []int{0, m}
  for len(Vals) > 0 {
   	j := 0
   	for _, v := range Vals {
   	  inf, vls := ValToStruct(v)
   	  Info = append(Info, inf...)
   	  Vals = append(Vals, vls...)
   	  j++
   	  lens = append(lens, len(Info))
   	}
   	Vals = Vals[j:]
  }
  structs := [][]string{}
  for i := 0; i < len(lens) - 1; i++ {
  	sl := Info[lens[i]:lens[i+1]]
  	structs = append(structs, sl)
  }
  for i := len(structs) - 2; i >= 0; i-- {
  	st := structs[i]
  	ns := NumStruct(st)
  	if ns == 0 {
  		continue
  	}
  	for j := len(st)-1; j >= 0; j-- {
  		if st[j] == "@" {
		    str := strings.Join(structs[i + ns], "")
		    structs[i][j] = str
		    structs[i + ns] = []string{}
		    ns--
  		}
  	}
  }
  res := ""
  for i := 0; i < len(structs); i++ {
  	content := strings.Join(structs[i], "")
  	res += content
  }
	return res
}

func WriteConfig(s string) error {
	par := strings.Split(Path, "/")
	if len(par) > 1 {
		par = par[:len(par)-1]
		fold := strings.Join(par, "/")
	  if _, err := os.Stat(fold); os.IsNotExist(err) {
		  os.Mkdir(fold, 0750)
	  }
	}
	f, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err)
	}
	defer f.Close()
	header := "# github.com/yurajp/confy\n"
	_, err = f.WriteString(header + s)
	if err != nil {
		return fmt.Errorf("Cannot write file: %s", err)
	}
	return nil
}

func ConfigExists() bool {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func WriteConfy(c any) error {
	v := reflect.ValueOf(c)
	vl := Vall{v, 0}
	s := BuildString(vl)
	err := WriteConfig(s)
	if err != nil {
		return err
	}
	return nil
}
