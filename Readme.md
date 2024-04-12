Confy is the small and simple library to store struct data such as config. It is available to write data to file and to load it back to struct. Supports nested structs, slices and maps (map cannot contain structs).
Similar to json you can view and edit confy file but more user friendly - file contains field types and indents.

go get github.com/yurajp/confy
...

Just two parameters you can change: path to file and indent size. E.g.:

confy.SetIndent(4)
confy.SetPath("conf/myconf.ini")

Defaults are "   " and "config/conf.ini".

You can use comments in writen ini. Use '#'.
To store your data struct in file:

err := confy.WriteConfy(<myconf>)
...
This will write variable 'myconf' of your type to file that has defined in Path.
Then you load data to interface variable and convert it to your type.

myconf := Config{}
iface, err := confy.LoadConfy(myconf)
...
mc, ok = iface.(Config)
if !ok {
	...
}
myconf = &mc

Good luck!