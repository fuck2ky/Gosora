package tmpl

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	//"regexp"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"text/template/parse"
)

// TODO: Turn this file into a library
var textOverlapList = make(map[string]int)

// nolint
type VarItem struct {
	Name        string
	Destination string
	Type        string
}
type VarItemReflect struct {
	Name        string
	Destination string
	Value       reflect.Value
}
type CTemplateSet struct {
	tlist                map[string]*parse.Tree
	dir                  string
	funcMap              map[string]interface{}
	importMap            map[string]string
	Fragments            map[string]int
	FragmentCursor       map[string]int
	FragOut              string
	varList              map[string]VarItem
	localVars            map[string]map[string]VarItemReflect
	hasDispInt           bool
	localDispStructIndex int
	stats                map[string]int
	pVarList             string
	pVarPosition         int
	previousNode         parse.NodeType
	currentNode          parse.NodeType
	nextNode             parse.NodeType
	//tempVars map[string]string
	doImports   bool
	minify      bool
	debug       bool
	superDebug  bool
	skipHandles bool
	expectsInt  interface{}
}

func (c *CTemplateSet) Minify(on bool) {
	c.minify = on
}

func (c *CTemplateSet) Debug(on bool) {
	c.debug = on
}

func (c *CTemplateSet) SuperDebug(on bool) {
	c.superDebug = on
}

func (c *CTemplateSet) SkipHandles(on bool) {
	c.skipHandles = on
}

func (c *CTemplateSet) Compile(name string, dir string, expects string, expectsInt interface{}, varList map[string]VarItem, imports ...string) (out string, err error) {
	if c.debug {
		fmt.Println("Compiling template '" + name + "'")
	}

	c.dir = dir
	c.doImports = true
	c.funcMap = map[string]interface{}{
		"and":      "&&",
		"not":      "!",
		"or":       "||",
		"eq":       true,
		"ge":       true,
		"gt":       true,
		"le":       true,
		"lt":       true,
		"ne":       true,
		"add":      true,
		"subtract": true,
		"multiply": true,
		"divide":   true,
	}

	c.importMap = map[string]string{
		"net/http": "net/http",
		"./common": "./common",
	}
	if len(imports) > 0 {
		for _, importItem := range imports {
			c.importMap[importItem] = importItem
		}
	}

	c.varList = varList
	c.hasDispInt = false
	c.localDispStructIndex = 0
	//c.pVarList = ""
	//c.pVarPosition = 0
	c.stats = make(map[string]int)
	c.expectsInt = expectsInt
	holdreflect := reflect.ValueOf(expectsInt)

	res, err := ioutil.ReadFile(dir + name)
	if err != nil {
		return "", err
	}

	content := string(res)
	if c.minify {
		content = minify(content)
	}

	tree := parse.New(name, c.funcMap)
	var treeSet = make(map[string]*parse.Tree)
	tree, err = tree.Parse(content, "{{", "}}", treeSet, c.funcMap)
	if err != nil {
		return "", err
	}
	c.log(name)

	out = ""
	fname := strings.TrimSuffix(name, filepath.Ext(name))
	c.tlist = make(map[string]*parse.Tree)
	c.tlist[fname] = tree
	varholder := "tmpl_" + fname + "_vars"

	c.log(c.tlist)
	c.localVars = make(map[string]map[string]VarItemReflect)
	c.localVars[fname] = make(map[string]VarItemReflect)
	c.localVars[fname]["."] = VarItemReflect{".", varholder, holdreflect}
	if c.Fragments == nil {
		c.Fragments = make(map[string]int)
	}
	c.FragmentCursor = make(map[string]int)
	c.FragmentCursor[fname] = 0

	out += c.rootIterate(c.tlist[fname], varholder, holdreflect, fname)

	var importList string
	if c.doImports {
		for _, item := range c.importMap {
			importList += "import \"" + item + "\"\n"
		}
	}

	var varString string
	for _, varItem := range c.varList {
		varString += "var " + varItem.Name + " " + varItem.Type + " = " + varItem.Destination + "\n"
	}

	fout := "// +build !no_templategen\n\n// Code generated by Gosora. More below:\n/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */\n"

	fout += "package main\n" + importList + c.pVarList + "\n"
	fout += "// nolint\nfunc init() {\n"

	if !c.skipHandles {
		fout += "\tcommon.Template_" + fname + "_handle = Template_" + fname + "\n"

		fout += "\tcommon.Ctemplates = append(common.Ctemplates,\"" + fname + "\")\n\tcommon.TmplPtrMap[\"" + fname + "\"] = &common.Template_" + fname + "_handle\n"
	}

	fout += "\tcommon.TmplPtrMap[\"o_" + fname + "\"] = Template_" + fname + "\n}\n\n"

	fout += "// nolint\nfunc Template_" + fname + "(tmpl_" + fname + "_vars " + expects + ", w http.ResponseWriter) error {\n" + varString + out + "\treturn nil\n}\n"

	fout = strings.Replace(fout, `))
w.Write([]byte(`, " + ", -1)
	fout = strings.Replace(fout, "` + `", "", -1)
	//spstr := "`([:space:]*)`"
	//whitespaceWrites := regexp.MustCompile(`(?s)w.Write\(\[\]byte\(`+spstr+`\)\)`)
	//fout = whitespaceWrites.ReplaceAllString(fout,"")

	if c.debug {
		for index, count := range c.stats {
			fmt.Println(index+": ", strconv.Itoa(count))
		}
		fmt.Println(" ")
	}

	c.log("Output!")
	c.log(fout)
	return fout, nil
}

func (c *CTemplateSet) rootIterate(tree *parse.Tree, varholder string, holdreflect reflect.Value, fname string) (out string) {
	c.log(tree.Root)
	treeLength := len(tree.Root.Nodes)
	for index, node := range tree.Root.Nodes {
		c.log("Node:", node.String())
		c.previousNode = c.currentNode
		c.currentNode = node.Type()
		if treeLength != (index + 1) {
			c.nextNode = tree.Root.Nodes[index+1].Type()
		}
		out += c.compileSwitch(varholder, holdreflect, fname, node)
	}
	return out
}

func (c *CTemplateSet) compileSwitch(varholder string, holdreflect reflect.Value, templateName string, node parse.Node) (out string) {
	c.log("in compileSwitch")
	switch node := node.(type) {
	case *parse.ActionNode:
		c.log("Action Node")
		if node.Pipe == nil {
			break
		}
		for _, cmd := range node.Pipe.Cmds {
			out += c.compileSubswitch(varholder, holdreflect, templateName, cmd)
		}
	case *parse.IfNode:
		c.log("If Node:")
		c.log("node.Pipe", node.Pipe)

		var expr string
		for _, cmd := range node.Pipe.Cmds {
			c.log("If Node Bit:", cmd)
			c.log("If Node Bit Type:", reflect.ValueOf(cmd).Type().Name())
			expr += c.compileVarswitch(varholder, holdreflect, templateName, cmd)
			c.log("If Node Expression Step:", c.compileVarswitch(varholder, holdreflect, templateName, cmd))
		}

		c.log("If Node Expression:", expr)
		c.previousNode = c.currentNode
		c.currentNode = parse.NodeList
		c.nextNode = -1
		if node.ElseList == nil {
			c.log("Selected Branch 1")
			return "if " + expr + " {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.List) + "}\n"
		}

		c.log("Selected Branch 2")
		return "if " + expr + " {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.List) + "} else {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.ElseList) + "}\n"
	case *parse.ListNode:
		c.log("List Node")
		for _, subnode := range node.Nodes {
			out += c.compileSwitch(varholder, holdreflect, templateName, subnode)
		}
	case *parse.RangeNode:
		return c.compileRangeNode(varholder, holdreflect, templateName, node)
	case *parse.TemplateNode:
		return c.compileSubtemplate(varholder, holdreflect, node)
	case *parse.TextNode:
		c.previousNode = c.currentNode
		c.currentNode = node.Type()
		c.nextNode = 0
		tmpText := bytes.TrimSpace(node.Text)
		if len(tmpText) == 0 {
			return ""
		}

		fragmentName := templateName + "_" + strconv.Itoa(c.FragmentCursor[templateName])
		_, ok := c.Fragments[fragmentName]
		if !ok {
			c.Fragments[fragmentName] = len(node.Text)
			c.FragOut += "var " + fragmentName + " = []byte(`" + string(node.Text) + "`)\n"
		}
		c.FragmentCursor[templateName] = c.FragmentCursor[templateName] + 1
		return "w.Write(" + fragmentName + ")\n"
	default:
		return c.unknownNode(node)
	}
	return out
}

func (c *CTemplateSet) compileRangeNode(varholder string, holdreflect reflect.Value, templateName string, node *parse.RangeNode) (out string) {
	c.log("Range Node!")
	c.log(node.Pipe)

	var outVal reflect.Value
	for _, cmd := range node.Pipe.Cmds {
		c.log("Range Bit:", cmd)
		out, outVal = c.compileReflectSwitch(varholder, holdreflect, templateName, cmd)
	}
	c.log("Returned:", out)
	c.log("Range Kind Switch!")

	switch outVal.Kind() {
	case reflect.Map:
		var item reflect.Value
		for _, key := range outVal.MapKeys() {
			item = outVal.MapIndex(key)
		}
		if c.debug {
			fmt.Println("Range item:", item)
		}
		if !item.IsValid() {
			panic("item" + "^\n" + "Invalid map. Maybe, it doesn't have any entries for the template engine to analyse?")
		}

		if node.ElseList != nil {
			out = "if len(" + out + ") != 0 {\nfor _, item := range " + out + " {\n" + c.compileSwitch("item", item, templateName, node.List) + "}\n} else {\n" + c.compileSwitch("item", item, templateName, node.ElseList) + "}\n"
		} else {
			out = "if len(" + out + ") != 0 {\nfor _, item := range " + out + " {\n" + c.compileSwitch("item", item, templateName, node.List) + "}\n}"
		}
	case reflect.Slice:
		if outVal.Len() == 0 {
			panic("The sample data needs at-least one or more elements for the slices. We're looking into removing this requirement at some point!")
		}
		item := outVal.Index(0)
		out = "if len(" + out + ") != 0 {\nfor _, item := range " + out + " {\n" + c.compileSwitch("item", item, templateName, node.List) + "}\n}"
	case reflect.Invalid:
		return ""
	}

	if node.ElseList != nil {
		out += " else {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.ElseList) + "}"
	}
	return out + "\n"
}

func (c *CTemplateSet) compileSubswitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.log("in compileSubswitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		c.log("Field Node:", n.Ident)

		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Variable declarations are coming soon! */
		cur := holdreflect

		var varbit string
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			varbit += ".(" + cur.Type().Name() + ")"
		}

		// ! Might not work so well for non-struct pointers
		skipPointers := func(cur reflect.Value, id string) reflect.Value {
			if cur.Kind() == reflect.Ptr {
				c.log("Looping over pointer")
				for cur.Kind() == reflect.Ptr {
					cur = cur.Elem()
				}
				c.log("Data Kind:", cur.Kind().String())
				c.log("Field Bit:", id)
			}
			return cur
		}

		var assLines string
		var multiline = false
		for _, id := range n.Ident {
			c.log("Data Kind:", cur.Kind().String())
			c.log("Field Bit:", id)
			cur = skipPointers(cur, id)

			if !cur.IsValid() {
				c.error("Debug Data:")
				c.error("Holdreflect:", holdreflect)
				c.error("Holdreflect.Kind():", holdreflect.Kind())
				if !c.superDebug {
					c.error("cur.Kind():", cur.Kind().String())
				}
				c.error("")
				if !multiline {
					panic(varholder + varbit + "^\n" + "Invalid value. Maybe, it doesn't exist?")
				} else {
					panic(varbit + "^\n" + "Invalid value. Maybe, it doesn't exist?")
				}
			}

			c.log("in-loop varbit: " + varbit)
			if cur.Kind() == reflect.Map {
				cur = cur.MapIndex(reflect.ValueOf(id))
				varbit += "[\"" + id + "\"]"
				cur = skipPointers(cur, id)

				if cur.Kind() == reflect.Struct || cur.Kind() == reflect.Interface {
					// TODO: Move the newVarByte declaration to the top level or to the if level, if a dispInt is only used in a particular if statement
					var dispStr, newVarByte string
					if cur.Kind() == reflect.Interface {
						dispStr = "Int"
						if !c.hasDispInt {
							newVarByte = ":"
							c.hasDispInt = true
						}
					}
					// TODO: De-dupe identical struct types rather than allocating a variable for each one
					if cur.Kind() == reflect.Struct {
						dispStr = "Struct" + strconv.Itoa(c.localDispStructIndex)
						newVarByte = ":"
						c.localDispStructIndex++
					}
					varbit = "disp" + dispStr + " " + newVarByte + "= " + varholder + varbit + "\n"
					varholder = "disp" + dispStr
					multiline = true
				} else {
					continue
				}
			}
			if cur.Kind() != reflect.Interface {
				cur = cur.FieldByName(id)
				varbit += "." + id
			}

			// TODO: Handle deeply nested pointers mixed with interfaces mixed with pointers better
			if cur.Kind() == reflect.Interface {
				cur = cur.Elem()
				varbit += ".("
				// TODO: Surely, there's a better way of doing this?
				if cur.Type().PkgPath() != "main" && cur.Type().PkgPath() != "" {
					c.importMap["html/template"] = "html/template"
					varbit += strings.TrimPrefix(cur.Type().PkgPath(), "html/") + "."
				}
				varbit += cur.Type().Name() + ")"
			}
			c.log("End Cycle: ", varbit)
		}

		if multiline {
			assSplit := strings.Split(varbit, "\n")
			varbit = assSplit[len(assSplit)-1]
			assSplit = assSplit[:len(assSplit)-1]
			assLines = strings.Join(assSplit, "\n") + "\n"
		}
		varbit = varholder + varbit
		out = c.compileVarsub(varbit, cur, assLines)

		for _, varItem := range c.varList {
			if strings.HasPrefix(out, varItem.Destination) {
				out = strings.Replace(out, varItem.Destination, varItem.Name, 1)
			}
		}
		return out
	case *parse.DotNode:
		c.log("Dot Node:", node.String())
		return c.compileVarsub(varholder, holdreflect, "")
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	case *parse.VariableNode:
		c.log("Variable Node:", n.String())
		c.log(n.Ident)
		varname, reflectVal := c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
		return c.compileVarsub(varname, reflectVal, "")
	case *parse.StringNode:
		return n.Quoted
	case *parse.IdentifierNode:
		c.log("Identifier Node:", node)
		c.log("Identifier Node Args:", node.Args)
		out, outval := c.compileIdentSwitch(varholder, holdreflect, templateName, node)
		return c.compileVarsub(out, outval, "")
	default:
		return c.unknownNode(node)
	}
}

func (c *CTemplateSet) compileVarswitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.log("in compileVarswitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		if c.superDebug {
			fmt.Println("Field Node:", n.Ident)
			for _, id := range n.Ident {
				fmt.Println("Field Bit:", id)
			}
		}

		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Coming Soon. */
		return c.compileBoolsub(n.String(), varholder, templateName, holdreflect)
	case *parse.ChainNode:
		c.log("Chain Node:", n.Node)
		c.log("Chain Node Args:", node.Args)
	case *parse.IdentifierNode:
		c.log("Identifier Node:", node)
		c.log("Identifier Node Args:", node.Args)
		return c.compileIdentSwitchN(varholder, holdreflect, templateName, node)
	case *parse.DotNode:
		return varholder
	case *parse.VariableNode:
		c.log("Variable Node:", n.String())
		c.log("Variable Node Identifier:", n.Ident)
		out, _ = c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	case *parse.PipeNode:
		c.log("Pipe Node!")
		c.log(n)
		c.log("Args:", node.Args)
		out += c.compileIdentSwitchN(varholder, holdreflect, templateName, node)

		c.log("Out:", out)
	default:
		return c.unknownNode(firstWord)
	}
	return out
}

func (c *CTemplateSet) unknownNode(node parse.Node) (out string) {
	fmt.Println("Unknown Kind:", reflect.ValueOf(node).Elem().Kind())
	fmt.Println("Unknown Type:", reflect.ValueOf(node).Elem().Type().Name())
	panic("I don't know what node this is! Grr...")
}

func (c *CTemplateSet) compileIdentSwitchN(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.log("in compileIdentSwitchN")
	out, _ = c.compileIdentSwitch(varholder, holdreflect, templateName, node)
	return out
}

func (c *CTemplateSet) dumpSymbol(pos int, node *parse.CommandNode, symbol string) {
	c.log("symbol: ", symbol)
	c.log("node.Args[pos + 1]", node.Args[pos+1])
	c.log("node.Args[pos + 2]", node.Args[pos+2])
}

func (c *CTemplateSet) compareFunc(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, compare string) (out string) {
	c.dumpSymbol(pos, node, compare)
	return c.compileIfVarsubN(node.Args[pos+1].String(), varholder, templateName, holdreflect) + " " + compare + " " + c.compileIfVarsubN(node.Args[pos+2].String(), varholder, templateName, holdreflect)
}

func (c *CTemplateSet) simpleMath(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, symbol string) (out string, val reflect.Value) {
	leftParam, val2 := c.compileIfVarsub(node.Args[pos+1].String(), varholder, templateName, holdreflect)
	rightParam, val3 := c.compileIfVarsub(node.Args[pos+2].String(), varholder, templateName, holdreflect)

	if val2.IsValid() {
		val = val2
	} else if val3.IsValid() {
		val = val3
	} else {
		// TODO: What does this do?
		numSample := 1
		val = reflect.ValueOf(numSample)
	}

	c.dumpSymbol(pos, node, symbol)
	return leftParam + " " + symbol + " " + rightParam, val
}

func (c *CTemplateSet) compareJoin(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, symbol string) (pos2 int, out string) {
	c.log("Building " + symbol + " function")
	if pos == 0 {
		fmt.Println("pos:", pos)
		panic(symbol + " is missing a left operand")
	}
	if len(node.Args) <= pos {
		fmt.Println("post pos:", pos)
		fmt.Println("len(node.Args):", len(node.Args))
		panic(symbol + " is missing a right operand")
	}

	left := c.compileBoolsub(node.Args[pos-1].String(), varholder, templateName, holdreflect)
	_, funcExists := c.funcMap[node.Args[pos+1].String()]

	var right string
	if !funcExists {
		right = c.compileBoolsub(node.Args[pos+1].String(), varholder, templateName, holdreflect)
	}
	out = left + " " + symbol + " " + right

	c.log("Left operand:", node.Args[pos-1])
	c.log("Right operand:", node.Args[pos+1])
	if !funcExists {
		pos++
	}
	c.log("pos:", pos)
	c.log("len(node.Args):", len(node.Args))

	return pos, out
}

func (c *CTemplateSet) compileIdentSwitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string, val reflect.Value) {
	c.log("in compileIdentSwitch")
ArgLoop:
	for pos := 0; pos < len(node.Args); pos++ {
		id := node.Args[pos]
		c.log("pos:", pos)
		c.log("ID:", id)
		switch id.String() {
		case "not":
			out += "!"
		case "or", "and":
			var rout string
			pos, rout = c.compareJoin(varholder, holdreflect, templateName, pos, node, c.funcMap[id.String()].(string)) // TODO: Test this
			out += rout
		case "le": // TODO: Can we condense these comparison cases down into one?
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "<=")
			break ArgLoop
		case "lt":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "<")
			break ArgLoop
		case "gt":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, ">")
			break ArgLoop
		case "ge":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, ">=")
			break ArgLoop
		case "eq":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "==")
			break ArgLoop
		case "ne":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "!=")
			break ArgLoop
		case "add":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "+")
			out += rout
			val = rval
			break ArgLoop
		case "subtract":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "-")
			out += rout
			val = rval
			break ArgLoop
		case "divide":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "/")
			out += rout
			val = rval
			break ArgLoop
		case "multiply":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "*")
			out += rout
			val = rval
			break ArgLoop
		default:
			c.log("Variable!")
			if len(node.Args) > (pos + 1) {
				nextNode := node.Args[pos+1].String()
				if nextNode == "or" || nextNode == "and" {
					continue
				}
			}
			out += c.compileIfVarsubN(id.String(), varholder, templateName, holdreflect)
		}
	}
	return out, val
}

func (c *CTemplateSet) compileReflectSwitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string, outVal reflect.Value) {
	c.log("in compileReflectSwitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		if c.superDebug {
			fmt.Println("Field Node:", n.Ident)
			for _, id := range n.Ident {
				fmt.Println("Field Bit:", id)
			}
		}
		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Coming Soon. */
		return c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
	case *parse.ChainNode:
		c.log("Chain Node:", n.Node)
		c.log("node.Args:", node.Args)
	case *parse.DotNode:
		return varholder, holdreflect
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	default:
		//panic("I don't know what node this is")
	}
	return "", outVal
}

func (c *CTemplateSet) compileIfVarsubN(varname string, varholder string, templateName string, cur reflect.Value) (out string) {
	c.log("in compileIfVarsubN")
	out, _ = c.compileIfVarsub(varname, varholder, templateName, cur)
	return out
}

func (c *CTemplateSet) compileIfVarsub(varname string, varholder string, templateName string, cur reflect.Value) (out string, val reflect.Value) {
	c.log("in compileIfVarsub")
	if varname[0] != '.' && varname[0] != '$' {
		return varname, cur
	}

	bits := strings.Split(varname, ".")
	if varname[0] == '$' {
		var res VarItemReflect
		if varname[1] == '.' {
			res = c.localVars[templateName]["."]
		} else {
			res = c.localVars[templateName][strings.TrimPrefix(bits[0], "$")]
		}
		out += res.Destination
		cur = res.Value

		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
		}
	} else {
		out += varholder
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			out += ".(" + cur.Type().Name() + ")"
		}
	}
	bits[0] = strings.TrimPrefix(bits[0], "$")

	c.log("Cur Kind:", cur.Kind())
	c.log("Cur Type:", cur.Type().Name())
	for _, bit := range bits {
		c.log("Variable Field:", bit)
		if bit == "" {
			continue
		}

		// TODO: Fix this up so that it works for regular pointers and not just struct pointers. Ditto for the other cur.Kind() == reflect.Ptr we have in this file
		if cur.Kind() == reflect.Ptr {
			c.log("Looping over pointer")
			for cur.Kind() == reflect.Ptr {
				cur = cur.Elem()
			}
			c.log("Data Kind:", cur.Kind().String())
			c.log("Field Bit:", bit)
		}

		cur = cur.FieldByName(bit)
		out += "." + bit
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			out += ".(" + cur.Type().Name() + ")"
		}
		if !cur.IsValid() {
			panic(out + "^\n" + "Invalid value. Maybe, it doesn't exist?")
		}

		c.log("Data Kind:", cur.Kind())
		c.log("Data Type:", cur.Type().Name())
	}

	c.log("Out Value:", out)
	c.log("Out Kind:", cur.Kind())
	c.log("Out Type:", cur.Type().Name())

	for _, varItem := range c.varList {
		if strings.HasPrefix(out, varItem.Destination) {
			out = strings.Replace(out, varItem.Destination, varItem.Name, 1)
		}
	}

	c.log("Out Value:", out)
	c.log("Out Kind:", cur.Kind())
	c.log("Out Type:", cur.Type().Name())

	_, ok := c.stats[out]
	if ok {
		c.stats[out]++
	} else {
		c.stats[out] = 1
	}

	return out, cur
}

func (c *CTemplateSet) compileBoolsub(varname string, varholder string, templateName string, val reflect.Value) string {
	c.log("in compileBoolsub")
	out, val := c.compileIfVarsub(varname, varholder, templateName, val)
	// TODO: What if it's a pointer or an interface? I *think* we've got pointers handled somewhere, but not interfaces which we don't know the types of at compile time
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		out += " > 0"
	case reflect.Bool: // Do nothing
	case reflect.String:
		out += " != \"\""
	case reflect.Slice, reflect.Map:
		out = "len(" + out + ") != 0"
	default:
		fmt.Println("Variable Name:", varname)
		fmt.Println("Variable Holder:", varholder)
		fmt.Println("Variable Kind:", val.Kind())
		panic("I don't know what this variable's type is o.o\n")
	}
	return out
}

func (c *CTemplateSet) compileVarsub(varname string, val reflect.Value, assLines string) (out string) {
	c.log("in compileVarsub")
	for _, varItem := range c.varList {
		if strings.HasPrefix(varname, varItem.Destination) {
			varname = strings.Replace(varname, varItem.Destination, varItem.Name, 1)
		}
	}

	_, ok := c.stats[varname]
	if ok {
		c.stats[varname]++
	} else {
		c.stats[varname] = 1
	}

	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	c.log("varname: ", varname)
	c.log("assLines: ", assLines)
	switch val.Kind() {
	case reflect.Int:
		c.importMap["strconv"] = "strconv"
		out = "w.Write([]byte(strconv.Itoa(" + varname + ")))\n"
	case reflect.Bool:
		out = "if " + varname + " {\nw.Write([]byte(\"true\"))} else {\nw.Write([]byte(\"false\"))\n}\n"
	case reflect.String:
		if val.Type().Name() != "string" && !strings.HasPrefix(varname, "string(") {
			varname = "string(" + varname + ")"
		}
		out = "w.Write([]byte(" + varname + "))\n"
	case reflect.Int64:
		c.importMap["strconv"] = "strconv"
		out = "w.Write([]byte(strconv.FormatInt(" + varname + ", 10)))"
	default:
		if !val.IsValid() {
			panic(assLines + varname + "^\n" + "Invalid value. Maybe, it doesn't exist?")
		}
		fmt.Println("Unknown Variable Name:", varname)
		fmt.Println("Unknown Kind:", val.Kind())
		fmt.Println("Unknown Type:", val.Type().Name())
		panic("// I don't know what this variable's type is o.o\n")
	}
	c.log("out: ", out)
	return assLines + out
}

func (c *CTemplateSet) compileSubtemplate(pvarholder string, pholdreflect reflect.Value, node *parse.TemplateNode) (out string) {
	c.log("in compileSubtemplate")
	c.log("Template Node:", node.Name)

	fname := strings.TrimSuffix(node.Name, filepath.Ext(node.Name))
	varholder := "tmpl_" + fname + "_vars"
	var holdreflect reflect.Value
	if node.Pipe != nil {
		for _, cmd := range node.Pipe.Cmds {
			firstWord := cmd.Args[0]
			switch firstWord.(type) {
			case *parse.DotNode:
				varholder = pvarholder
				holdreflect = pholdreflect
			case *parse.NilNode:
				panic("Nil is not a command x.x")
			default:
				out = "var " + varholder + " := false\n"
				out += c.compileCommand(cmd)
			}
		}
	}

	// TODO: Cascade errors back up the tree to the caller?
	res, err := ioutil.ReadFile(c.dir + node.Name)
	if err != nil {
		log.Fatal(err)
	}

	content := string(res)
	if c.minify {
		content = minify(content)
	}

	tree := parse.New(node.Name, c.funcMap)
	var treeSet = make(map[string]*parse.Tree)
	tree, err = tree.Parse(content, "{{", "}}", treeSet, c.funcMap)
	if err != nil {
		log.Fatal(err)
	}

	c.tlist[fname] = tree
	subtree := c.tlist[fname]
	c.log("subtree.Root", subtree.Root)

	c.localVars[fname] = make(map[string]VarItemReflect)
	c.localVars[fname]["."] = VarItemReflect{".", varholder, holdreflect}
	c.FragmentCursor[fname] = 0

	out += c.rootIterate(subtree, varholder, holdreflect, fname)
	return out
}

func (c *CTemplateSet) log(args ...interface{}) {
	if c.superDebug {
		fmt.Println(args...)
	}
}

func (c *CTemplateSet) error(args ...interface{}) {
	if c.debug {
		fmt.Println(args...)
	}
}

func (c *CTemplateSet) compileCommand(*parse.CommandNode) (out string) {
	panic("Uh oh! Something went wrong!")
}

// TODO: Write unit tests for this
func minify(data string) string {
	data = strings.Replace(data, "\t", "", -1)
	data = strings.Replace(data, "\v", "", -1)
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)
	data = strings.Replace(data, "  ", " ", -1)
	return data
}

// TODO: Strip comments
// TODO: Handle CSS nested in <style> tags?
// TODO: Write unit tests for this
func minifyHTML(data string) string {
	return minify(data)
}

// TODO: Have static files use this
// TODO: Strip comments
// TODO: Convert the rgb()s to hex codes?
// TODO: Write unit tests for this
func minifyCSS(data string) string {
	return minify(data)
}

// TODO: Convert this to three character hex strings whenever possible?
// TODO: Write unit tests for this
// nolint
func rgbToHexstr(red int, green int, blue int) string {
	return strconv.FormatInt(int64(red), 16) + strconv.FormatInt(int64(green), 16) + strconv.FormatInt(int64(blue), 16)
}

/*
// TODO: Write unit tests for this
func hexstrToRgb(hexstr string) (red int, blue int, green int, err error) {
	// Strip the # at the start
	if hexstr[0] == '#' {
		hexstr = strings.TrimPrefix(hexstr,"#")
	}
	if len(hexstr) != 3 && len(hexstr) != 6 {
		return 0, 0, 0, errors.New("Hex colour codes may only be three or six characters long")
	}

	if len(hexstr) == 3 {
		hexstr = hexstr[0] + hexstr[0] + hexstr[1] + hexstr[1] + hexstr[2] + hexstr[2]
	}
}*/
