package qgen

var Install *installer

func init() {
	Install = &installer{instructions: []DBInstallInstruction{}}
}

type DBInstallInstruction struct {
	Table    string
	Contents string
	Type     string
}

// TODO: Add methods to this to construct it OO-like
type DBInstallTable struct {
	Name      string
	Charset   string
	Collation string
	Columns   []DBTableColumn
	Keys      []DBTableKey
}

// A set of wrappers around the generator methods, so we can use this in the installer
// TODO: Re-implement the query generation, query builder and installer adapters as layers on-top of a query text adapter
type installer struct {
	adapter      Adapter
	instructions []DBInstallInstruction
	tables       []*DBInstallTable // TODO: Use this in Record() in the next commit to allow us to auto-migrate settings rather than manually patching them in on upgrade
	plugins      []QueryPlugin
}

func (install *installer) SetAdapter(name string) error {
	adap, err := GetAdapter(name)
	if err != nil {
		return err
	}
	install.SetAdapterInstance(adap)
	return nil
}

func (install *installer) SetAdapterInstance(adapter Adapter) {
	install.adapter = adapter
	install.instructions = []DBInstallInstruction{}
}

func (install *installer) AddPlugins(plugins ...QueryPlugin) {
	install.plugins = append(install.plugins, plugins...)
}

func (install *installer) CreateTable(table string, charset string, collation string, columns []DBTableColumn, keys []DBTableKey) error {
	tableStruct := &DBInstallTable{table, charset, collation, columns, keys}
	err := install.RunHook("CreateTableStart", tableStruct)
	if err != nil {
		return err
	}
	res, err := install.adapter.CreateTable("", table, charset, collation, columns, keys)
	if err != nil {
		return err
	}
	err = install.RunHook("CreateTableAfter", tableStruct)
	if err != nil {
		return err
	}
	install.instructions = append(install.instructions, DBInstallInstruction{table, res, "create-table"})
	install.tables = append(install.tables, tableStruct)
	return nil
}

// TODO: Let plugins manipulate the parameters like in CreateTable
func (install *installer) AddIndex(table string, iname string, colname string) error {
	err := install.RunHook("AddIndexStart", table, iname, colname)
	if err != nil {
		return err
	}
	res, err := install.adapter.AddIndex("", table, iname, colname)
	if err != nil {
		return err
	}
	err = install.RunHook("AddIndexAfter", table, iname, colname)
	if err != nil {
		return err
	}
	install.instructions = append(install.instructions, DBInstallInstruction{table, res, "index"})
	return nil
}

// TODO: Let plugins manipulate the parameters like in CreateTable
func (install *installer) SimpleInsert(table string, columns string, fields string) error {
	err := install.RunHook("SimpleInsertStart", table, columns, fields)
	if err != nil {
		return err
	}
	res, err := install.adapter.SimpleInsert("", table, columns, fields)
	if err != nil {
		return err
	}
	err = install.RunHook("SimpleInsertAfter", table, columns, fields, res)
	if err != nil {
		return err
	}
	install.instructions = append(install.instructions, DBInstallInstruction{table, res, "insert"})
	return nil
}

func (install *installer) RunHook(name string, args ...interface{}) error {
	for _, plugin := range install.plugins {
		err := plugin.Hook(name, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (install *installer) Write() error {
	var inserts string
	// We can't escape backticks, so we have to dump it out a file at a time
	for _, instr := range install.instructions {
		if instr.Type == "create-table" {
			err := writeFile("./schema/"+install.adapter.GetName()+"/query_"+instr.Table+".sql", instr.Contents)
			if err != nil {
				return err
			}
		} else {
			inserts += instr.Contents + ";\n"
		}
	}

	err := writeFile("./schema/"+install.adapter.GetName()+"/inserts.sql", inserts)
	if err != nil {
		return err
	}

	for _, plugin := range install.plugins {
		err := plugin.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
