package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"strconv"
	"strings"
)

var errExpectedLiteral = errors.New("expected basic literal (ex. string or number)")

func literal(ex ast.Expr) (string, error) {
	lit, ok := ex.(*ast.BasicLit)
	if !ok {
		return "", errExpectedLiteral
	}

	return lit.Value, nil
}

func expr(e string) error {
	ex, err := parser.ParseExpr(e)
	if err != nil {
		return err
	}

	switch ex := ex.(type) {
	case *ast.CallExpr:
		switch fun := ex.Fun.(type) {
		case *ast.SelectorExpr:
			ident, ok := fun.X.(*ast.Ident)
			if !ok {
				return errors.New("expected identifier")
			}
			k := ident.Name
			if v, ok := workspace[k]; ok {
				switch fun.Sel.Name {
				case "delete":
					if len(ex.Args) != 1 {
						return errors.New("delete() expects exactly one argument")
					}
					lit, err := literal(ex.Args[0])
					if err != nil {
						return fmt.Errorf("invalid input to delete(): %w", err)
					}
					return v.Delete(lit)
				case "dedup":
					if len(ex.Args) != 0 {
						return errors.New("dedup() expects no arguments")
					}
					v.Dedup()
					return nil
				case "keep":
					if len(ex.Args) != 1 {
						return errors.New("delete() expects exactly one argument")
					}
					lit, err := literal(ex.Args[0])
					if err != nil {
						return fmt.Errorf("invalid input to keep(): %w", err)
					}
					return v.Keep(lit)
				case "save":
					if len(ex.Args) != 1 {
						return errors.New("save() expects exactly one argument")
					}
					lit, err := literal(ex.Args[0])
					if err != nil {
						return fmt.Errorf("invalid input to save(): %w", err)
					}
					fn := strings.Trim(lit, "\"")
					if _, err := os.Stat(fn); err == nil {
						pmpt := fmt.Sprintf("File %s already exists, overwrite it? [Y]/n: ", fn)
						var confirm string
						if confirm, err = line.Prompt(pmpt); err != nil {
							return err
						}
						confirm = strings.ToLower(strings.TrimSpace(confirm))
						if confirm != "y" && confirm != "" {
							return nil
						}
					}
					if err := v.Save(fn); err != nil {
						return err
					}
				case "search":
					var err error
					offset := 0
					limit := 10
					switch len(ex.Args) {
					case 0:
						return errors.New("search() expects at least one argument")
					case 1:
					case 2:
						offset, err = strconv.Atoi(ex.Args[1].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument %s", ex.Args[1])
						}
					case 3:
						offset, err = strconv.Atoi(ex.Args[1].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument 'offset' %s", ex.Args[1])
						}
						limit, err = strconv.Atoi(ex.Args[2].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument 'limit' %s", ex.Args[2])
						}
					default:
						return errors.New("search() expects at most three arguments")
					}
					lit, ok := ex.Args[0].(*ast.BasicLit)
					if !ok {
						return errors.New("search() expects a quoted expression")
					}
					v.Search(lit.Value, offset, limit)
					return nil
				case "show":
					var err error
					offset := 0
					limit := 10
					switch len(ex.Args) {
					case 0:
					case 1:
						offset, err = strconv.Atoi(ex.Args[0].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument %s", ex.Args[0])
						}
					case 2:
						offset, err = strconv.Atoi(ex.Args[0].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument 'offset' %s", ex.Args[0])
						}
						limit, err = strconv.Atoi(ex.Args[1].(*ast.BasicLit).Value)
						if err != nil {
							return fmt.Errorf("invalid argument 'limit' %s", ex.Args[1])
						}
					default:
						return errors.New("show() expects at least one and at most two arguments")
					}
					v.Show(offset, limit)
					return nil
				default:
					return fmt.Errorf("unknown instruction")
				}
			}
		default:
			return fmt.Errorf("unknown instruction")
		}
	case *ast.Ident:
		if v, ok := workspace[ex.String()]; ok {
			v.Summary()
		} else {
			return fmt.Errorf("variable %s not found in workspace", e)
		}
	default:
		return fmt.Errorf("unknown instruction")
	}

	return nil
}
