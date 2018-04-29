// Package ast implements utility functions for generating abstract syntax
// trees from Go BNF.

package ast

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/shivansh/gogo/src/utils"
)

type symTabType map[string][]string

type SymInfo struct {
	varSymTab symTabType
	parent    *SymInfo
}

// DeferStackItem is an individual item stored when a call to defer is made. It
// contains the code for the corresponding function call which is placed at the
// end of function body.
type DeferStackItem []string

var (
	tmpIndex   int
	labelIndex int
	varIndex   int
	// funcSymtabCreated keeps track whether a symbol table corresponding to a
	// function declaration has to be instantiated. This is because usually
	// a new symbol table is created when the corresponding block begins.
	// However, in case of functions the arguments also need to be added to
	// the symbol table. Thus the symbol table is instantiated when the
	// production rule corresponding to the arguments is reached and not
	// when the block begins.
	funcSymtabCreated bool
	forSymtabCreated  bool
	symTab            symTabType // symbol table for temporaries ; TODO: Update this.
	// currSymTab keeps track of the currently active symbol table
	// depending on scope.
	currSymTab *SymInfo
	// globalSymTab keeps track of the global struct and function declarations.
	// NOTE: structs and functions can only be declared globally.
	globalSymTab symTabType
	// deferStack stores the deferred function calls which are then called
	// when the surrounding function block ends.
	deferStack *utils.Stack
	re         *regexp.Regexp
)

func init() {
	symTab = make(symTabType)
	globalSymTab = make(symTabType)
	// currSymTab now allocates space for a global symbol table for variable
	// declarations. Ideally, instead of creating a new symbol table it
	// should've been global symbol table, but since the type of
	// currSymTab.varSymTab is not pointer, updates made elsewhere will not
	// be reflected globally.
	// TODO: Update type of currSymTab.varSymTab to a pointer.
	currSymTab = &SymInfo{make(symTabType), nil}
	deferStack = utils.CreateStack()
	re = regexp.MustCompile("(^-?[0-9]+$)") // integers
}

// SearchInScope returns the symbol table entry for a given variable in the
// current scope. If not found, the parent symbol table is looked up until the
// topmost symbol table is reached. If not found in all these symbol tables,
// then the global symbol table is looked up which contains the entries
// corresponding to structs and functions.
func SearchInScope(v string) ([]string, bool) {
	currScope := currSymTab
	for currScope != nil {
		symTabEntry, ok := currScope.varSymTab[v]
		if ok {
			return symTabEntry, true
		} else {
			currScope = currScope.parent
		}
	}
	// Lookup in global scope in case the variable corresponds to a struct
	// or a function name.
	for k, symTabEntry := range globalSymTab {
		if k == v {
			return symTabEntry, true
		}
	}
	return []string{}, false
}

// GetRealName extracts the original name of variable from its renamed version.
func GetRealName(s string) string {
	realName := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			break
		} else {
			realName = realName + string(s[i])
		}
	}
	return realName
}

// NewTmp generates a unique temporary variable.
func NewTmp() string {
	t := fmt.Sprintf("t%d", tmpIndex)
	tmpIndex++
	return t
}

// NewLabel generates a unique label name.
func NewLabel() string {
	l := fmt.Sprintf("l%d", labelIndex)
	labelIndex++
	return l
}

// Node represents a node in the AST of a given program.
type Node struct {
	// If the node represents an expression, then place stores the name of
	// the variable storing the value of the expression.
	Place string
	Code  []string // IR instructions
}

type Attrib interface{}

// NewVar generates a unique variable name used for renaming. A variable named
// var will be renamed to 'var.int_lit' where int_lit is an integer. Since
// variable names cannot contain a '.', this will not result in a naming
// conflict with an existing variable. The renamed variable will only occur in
// the IR (there is no constraint on variable names in IR as of now).
func RenameVariable(v string) string {
	ret := fmt.Sprintf("%s.%d", v, varIndex)
	varIndex++
	return ret
}

func PrintIR(src interface{}) (interface{}, error) {
	re := regexp.MustCompile("\n(\n)*")
	c := src.(Node).Code
	for _, v := range c {
		v := strings.TrimSpace(v)
		// Compress multiple newlines within IR statements into
		// a single newline.
		v = re.ReplaceAllString(v, "\n")
		if v != "" {
			fmt.Println(v)
		}
	}
	return nil, nil
}

// --- [ Variable declarations ] -----------------------------------------------

// NewVarSpec creates a new variable specification.
// The accepted arguments (args) in their order are -
// 	- IdentifierList
// 	- Type
// 	- ExpressionList
// typ determines the index of the production rule invoked starting from top.
func NewVarSpec(typ int, args ...interface{}) (Attrib, error) {
	n := Node{"", []string{}}
	expr := []string{}
	// Add the IR instructions for ExpressionList.
	switch typ {
	case 1:
		n.Code = args[2].(Node).Code
		expr = utils.SplitAndSanitize(args[2].(Node).Place, ",")
	case 2:
		n.Code = args[1].(Node).Code
		expr = utils.SplitAndSanitize(args[1].(Node).Place, ",")
	case 3:
		return n, nil
	}
	for k, v := range args[0].(Node).Code {
		renamedVar := RenameVariable(v)
		// TODO: Handle other types
		currSymTab.varSymTab[v] = []string{renamedVar, "int"}
		if typ == 0 {
			n.Code = append(n.Code, fmt.Sprintf("=, %s, 0", renamedVar))
		} else if typ == 1 || typ == 2 {
			n.Code = append(n.Code, fmt.Sprintf("=, %s, %s", renamedVar, expr[k]))
		}
	}
	return n, nil
}

// --- [ Type declarations ] ---------------------------------------------------
func NewTypeDecl(args ...interface{}) (Attrib, error) {
	typeInfo := utils.SplitAndSanitize(args[1].(Node).Place, ",")
	structName := strings.TrimSpace(typeInfo[1])
	switch strings.TrimSpace(typeInfo[0]) {
	case "struct":
		// Create a global symbol table entry.
		// NOTE: The symbol table entry of a struct is of the form -
		//      structName : []{"struct", memberName1, memberType1, ...}
		globalSymTab[structName] = []string{"struct"}
		globalSymTab[structName] = append(globalSymTab[structName], args[1].(Node).Code...)
	default: // TODO: Add remaining types.
		return nil, errors.New("Unknown type")
	}
	// TODO: Member initialization will be done when a new object is
	// instantiated.
	return Node{"", []string{}}, nil
}

// --- [ Constant declarations ] -----------------------------------------------
func NewConstSpec(typ int, args ...interface{}) (Node, error) {
	n := Node{"", []string{}}
	expr := []string{}
	if typ == 1 {
		n.Code = append(n.Code, args[1].(Node).Code...)
		expr = utils.SplitAndSanitize(args[1].(Node).Place, ",")
	}
	for k, v := range args[0].(Node).Code {
		renamedVar := RenameVariable(v)
		currSymTab.varSymTab[v] = []string{renamedVar, "int"}
		switch typ {
		case 0:
			n.Code = append(n.Code, fmt.Sprintf("=, %s, 0", renamedVar))
		case 1:
			n.Code = append(n.Code, fmt.Sprintf("=, %s, %s", renamedVar, expr[k]))
		}
	}
	return n, nil
}
