/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"k8s.io/code-generator/cmd/validation-gen/validators"
	"k8s.io/gengo/v2"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
)

var (
	fieldPkg      = "k8s.io/apimachinery/pkg/util/validation/field"
	errorListType = types.Name{Package: fieldPkg, Name: "ErrorList"}
	fieldPathType = types.Name{Package: fieldPkg, Name: "Path"}
	errorfType    = types.Name{Package: "fmt", Name: "Errorf"}
	runtimePkg    = "k8s.io/apimachinery/pkg/runtime"
	schemeType    = types.Name{Package: runtimePkg, Name: "Scheme"}
)

// genValidations produces a file with autogenerated validations.
type genValidations struct {
	generator.GoGenerator
	outputPackage       string
	inputToPkg          map[string]string // Maps input packages to generated validation packages
	rootTypes           []*types.Type
	typeNodes           map[*types.Type]*typeNode
	imports             namer.ImportTracker
	validator           validators.DeclarativeValidator
	hasValidationsCache map[*types.Type]bool
}

// NewGenValidations cretes a new generator for the specified package.
func NewGenValidations(outputFilename, outputPackage string, rootTypes []*types.Type, typeNodes map[*types.Type]*typeNode, inputToPkg map[string]string, validator validators.DeclarativeValidator) generator.Generator {
	return &genValidations{
		GoGenerator: generator.GoGenerator{
			OutputFilename: outputFilename,
		},
		outputPackage:       outputPackage,
		inputToPkg:          inputToPkg,
		rootTypes:           rootTypes,
		typeNodes:           typeNodes,
		imports:             generator.NewImportTrackerForPackage(outputPackage),
		validator:           validator,
		hasValidationsCache: map[*types.Type]bool{},
	}
}

func (g *genValidations) Namers(_ *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *genValidations) Filter(_ *generator.Context, t *types.Type) bool {
	_, ok := g.typeNodes[t]
	return ok
}

func (g *genValidations) Imports(_ *generator.Context) (imports []string) {
	var importLines []string
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func (g *genValidations) isOtherPackage(pkg string) bool {
	if pkg == g.outputPackage {
		return false
	}
	if strings.HasSuffix(pkg, `"`+g.outputPackage+`"`) {
		return false
	}
	return true
}

func (g *genValidations) Init(c *generator.Context, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	scheme := c.Universe.Type(schemeType)
	schemePtr := &types.Type{
		Kind: types.Pointer,
		Elem: scheme,
	}
	sw.Do("func init() { localSchemeBuilder.Register(RegisterValidations)}\n\n", nil)

	sw.Do("// RegisterValidations adds validation functions to the given scheme.\n", nil)
	sw.Do("// Public to allow building arbitrary schemes.\n", nil)
	sw.Do("func RegisterValidations(scheme $.|raw$) error {\n", schemePtr)
	for _, t := range g.rootTypes {
		tn, ok := g.typeNodes[t]
		if !ok {
			continue
		}
		if tn == nil {
			// Should never happen.
			klog.Fatalf("found nil typeNode for type %v", t)
		}

		// TODO: It would be nice if these were not hard-coded.
		var statusType *types.Type
		var statusField string
		if status := tn.lookupField("status"); status != nil {
			statusType = status.underlyingType
			statusField = status.name
		}

		targs := generator.Args{
			"rootType":    t,
			"statusType":  statusType,
			"statusField": statusField,
			"errorList":   c.Universe.Type(errorListType),
			"fieldPath":   c.Universe.Type(fieldPathType),
			"fmtErrorf":   c.Universe.Type(errorfType),
		}
		//TODO: can this be (*$.rootType|raw$)(nil) ?
		sw.Do("scheme.AddValidationFunc(new($.rootType|raw$), func(obj, oldObj interface{}, subresources ...string) $.errorList|raw$ {\n", targs)
		sw.Do("  if len(subresources) == 0 {\n", targs)
		sw.Do("    return $.rootType|objectvalidationfn$(obj.(*$.rootType|raw$), nil)\n", targs)
		sw.Do("  }\n", targs)

		if statusType != nil {
			sw.Do("  if len(subresources) == 1 && subresources[0] == \"status\" {\n", targs)
			if g.hasValidations(statusType) {
				sw.Do("    root := obj.(*$.rootType|raw$)\n", targs)
				sw.Do("    return $.statusType|objectvalidationfn$(&root.$.statusField$, nil)\n", targs)
			} else {
				sw.Do("    return nil // $.statusType|raw$ has no validation\n", targs)
			}
			sw.Do("  }\n", targs)
		}
		sw.Do("  return $.errorList|raw${field.InternalError(nil, $.fmtErrorf|raw$(\"no validation found for %T, subresources: %v\", obj, subresources))}\n", targs)
		sw.Do("})\n", targs)

		// TODO: Support update validations
		//       This will require correlating old object.
	}
	sw.Do("return nil\n", nil)
	sw.Do("}\n\n", nil)
	return sw.Error()
}

func (g *genValidations) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	klog.V(5).Infof("generating for type %v", t)

	var errs []error

	isRoot := false
	for _, rt := range g.rootTypes {
		if rt == t {
			isRoot = true
			break
		}
	}
	if !isRoot && !g.hasValidations(t) {
		return nil
	}

	sw := generator.NewSnippetWriter(w, c, "$", "$")
	g.emitValidationFunction(c, t, sw)
	if err := sw.Error(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (g *genValidations) hasValidations(t *types.Type) bool {
	if result, found := g.hasValidationsCache[t]; found {
		return result
	}
	r := g.hasValidationsMiss(t)
	g.hasValidationsCache[t] = r
	return r
}

// Called in case of a cache miss.
func (g *genValidations) hasValidationsMiss(t *types.Type) bool {
	tn := g.typeNodes[t]
	if len(tn.validations) > 0 {
		return true
	}
	allChildren := tn.children
	if tn.key != nil {
		allChildren = append(allChildren, tn.key)
	}
	if tn.elem != nil {
		allChildren = append(allChildren, tn.elem)
	}
	for _, cn := range allChildren {
		if len(cn.validations)+len(cn.eachKey)+len(cn.eachVal) > 0 {
			return true
		}
		if g.hasValidations(cn.underlyingType) {
			return true
		}
	}
	return false
}

func (g *genValidations) emitValidationFunction(c *generator.Context, t *types.Type, sw *generator.SnippetWriter) {
	targs := generator.Args{
		"inType":    t,
		"errorList": c.Universe.Type(errorListType),
		"fieldPath": c.Universe.Type(fieldPathType),
	}

	sw.Do("func $.inType|objectvalidationfn$(obj *$.inType|raw$, fldPath *$.fieldPath|raw$) (errs $.errorList|raw$) {\n", targs)
	g.emitValidationForType(c, t, true, sw, nil, nil)
	sw.Do("return errs\n", nil)
	sw.Do("}\n\n", nil)
}

// typeDiscoverer contains fields necessary to build a tree of types.
type typeDiscoverer struct {
	validator  validators.DeclarativeValidator
	inputToPkg map[string]string
	inProgress map[*types.Type]bool
	knownTypes map[*types.Type]*typeNode
}

// discoverTypes walks the type graph and populates the result map.
func discoverTypes(validator validators.DeclarativeValidator, inputToPkg map[string]string, t *types.Type, results map[*types.Type]*typeNode) error {
	td := &typeDiscoverer{
		validator:  validator,
		inputToPkg: inputToPkg,
		knownTypes: results,
	}
	return td.discover(t)
}

// typeNode carries validation informatiuon for a single type.
type typeNode struct {
	underlyingType *types.Type
	validations    []validators.FunctionGen
	children       []*childNode // populated when parent is a Struct
	elem           *childNode   // populated when parent is a list
	key            *childNode   // populated when parent is a map
	funcName       types.Name
}

func (n typeNode) lookupField(jsonName string) *childNode {
	for _, c := range n.children {
		if c.jsonName == jsonName {
			return c
		}
	}
	return nil
}

// childNode represents a field in a struct.
type childNode struct {
	name           string
	jsonName       string
	underlyingType *types.Type
	validations    []validators.FunctionGen

	// iterated validation has to be tracked separately from field's validations.
	eachKey, eachVal []validators.FunctionGen
}

const (
	eachKeyTag = "eachKey"
	eachValTag = "eachVal"
)

// discover walks the type graph, starting at t, and registers all types into
// knownTypes.  The specified comments represent the parent context for this
// type - the type comments for a type definition or the field comments for a
// field.
func (td *typeDiscoverer) discover(t *types.Type) error {
	// If we already know this type, we are done.
	if _, ok := td.knownTypes[t]; ok {
		return nil
	}

	klog.V(5).InfoS("discovering", "type", t)

	thisNode := &typeNode{
		underlyingType: t,
	}

	// Publish it right away in case we hit it recursively.
	td.knownTypes[t] = thisNode

	// Extract any type-attached validation rules.
	if validations, err := td.validator.ExtractValidations(t.Name.Name, t, t.CommentLines); err != nil {
		return err
	} else {
		if len(validations) > 0 {
			klog.V(5).InfoS("  found type-attached validations", "n", len(validations))
			thisNode.validations = validations
		}
	}

	switch t.Kind {
	case types.Builtin:
		// Nothing more to do.
	case types.Pointer:
		klog.V(5).InfoS("  type is a pointer", "type", t.Elem)
		if t.Elem.Kind == types.Pointer {
			klog.Fatalf("type %v: pointers to pointers are not supported", t)
		}
		if err := td.discover(t.Elem); err != nil {
			return err
		}
	case types.Slice, types.Array:
		klog.V(5).InfoS("  type is a list", "type", t.Elem)
		if err := td.discover(t.Elem); err != nil {
			return err
		}
		thisNode.elem = &childNode{
			underlyingType: t.Elem,
		}
	case types.Map:
		klog.V(5).InfoS("  type is a map", "type", t.Elem)
		if err := td.discover(t.Key); err != nil {
			return err
		}
		thisNode.key = &childNode{
			underlyingType: t.Elem,
		}

		if err := td.discover(t.Elem); err != nil {
			return err
		}
		thisNode.elem = &childNode{
			underlyingType: t.Elem,
		}
	case types.Struct:
		klog.V(5).InfoS("  type is a struct")
		fn, ok := td.getValidationFunctionName(t)
		if !ok {
			//FIXME: this seems like an error, but is it?  Or just "opaque from here"
			return nil
		}
		thisNode.funcName = fn

		for _, field := range t.Members {
			name := field.Name
			if len(name) == 0 {
				// embedded fields
				if field.Type.Kind == types.Pointer {
					name = field.Type.Elem.Name.Name
				} else {
					name = field.Type.Name.Name
				}
			}
			// If we try to emit code for this field and find no JSON name, we
			// will abort.
			jsonName := ""
			if tags, ok := lookupJSONTags(field); ok {
				jsonName = tags.name
			}
			// Only do exported fields.
			if unicode.IsLower([]rune(field.Name)[0]) {
				continue
			}
			klog.V(5).InfoS("  field", "name", name)

			if err := td.discover(field.Type); err != nil {
				return err
			}

			child := &childNode{
				name:           name,
				jsonName:       jsonName,
				underlyingType: field.Type,
			}

			switch field.Type.Kind {
			case types.Map:
				//TODO: also support +k8s:eachKey
				if tagVals, found := gengo.ExtractCommentTags("+", field.CommentLines)[eachKeyTag]; found {
					for _, tagVal := range tagVals {
						fakeComments := []string{tagVal}
						// Extract any embedded key-validation rules.
						if validations, err := td.validator.ExtractValidations(fmt.Sprintf("%s[keys]", field.Name), field.Type.Key, fakeComments); err != nil {
							return err
						} else {
							if len(validations) > 0 {
								klog.V(5).InfoS("  found key-validations", "n", len(validations))
								child.eachKey = append(child.eachKey, validations...)
							}
						}
					}
				}
				//TODO: also support +k8s:eachVal
				if tagVals, found := gengo.ExtractCommentTags("+", field.CommentLines)[eachValTag]; found {
					for _, tagVal := range tagVals {
						fakeComments := []string{tagVal}
						// Extract any embedded list-validation rules.
						if validations, err := td.validator.ExtractValidations(fmt.Sprintf("%s[vals]", field.Name), field.Type.Elem, fakeComments); err != nil {
							return err
						} else {
							if len(validations) > 0 {
								klog.V(5).InfoS("  found list-validations", "n", len(validations))
								child.eachVal = append(child.eachVal, validations...)
							}
						}
					}
				}
			case types.Slice, types.Array:
				//TODO: also support +k8s:eachVal
				if tagVals, found := gengo.ExtractCommentTags("+", field.CommentLines)[eachValTag]; found {
					for _, tagVal := range tagVals {
						fakeComments := []string{tagVal}
						// Extract any embedded list-validation rules.
						if validations, err := td.validator.ExtractValidations(fmt.Sprintf("%s[vals]", field.Name), field.Type.Elem, fakeComments); err != nil {
							return err
						} else {
							if len(validations) > 0 {
								klog.V(5).InfoS("  found list-validations", "n", len(validations))
								child.eachVal = append(child.eachVal, validations...)
							}
						}
					}
				}
			}

			// Extract any field-attached validation rules.
			if validations, err := td.validator.ExtractValidations(name, field.Type, field.CommentLines); err != nil {
				return err
			} else {
				if len(validations) > 0 {
					klog.V(5).InfoS("  found field-attached value-validations", "n", len(validations))
					child.validations = append(child.validations, validations...)
				}
			}
			thisNode.children = append(thisNode.children, child)
		}
	case types.Alias:
		klog.V(5).InfoS("  type is an alias", "type", t.Underlying)
		// Note: By the language definition, what gengo calls "Aliases" (really
		// just "type definitions") have underlying types of the type literal.
		// In other words, if we define `type T1 string` and `type T2 T1`, the
		// underlying type of T2 is string, not T1.  This means that:
		//    1) We will emit code for both underlying types. If the underlying
		//       type is a struct with many fields, we will emit two identical
		//       functions.
		//    2) Validating a field of type T2 will NOT call any validation
		//       defined on the type T1.
		//    3) In the case of a type definition whose RHS is a struct which
		//       has fields with validation tags, the validation for those fields
		//       WILL be called from the generated for for the new type.
		if t.Underlying.Kind == types.Pointer {
			klog.Fatalf("type %v: aliases to pointers are not supported", t)
		}
		if err := td.discover(t.Underlying); err != nil {
			return err
		}
		fn, ok := td.getValidationFunctionName(t)
		if !ok {
			//FIXME: this seems like an error, but is it?  Or just "opaque from here"
			return nil
		}
		thisNode.funcName = fn
	default:
		klog.Fatalf("unhandled type: %v (%v)", t, t.Kind)
	}

	return nil
}

func (td *typeDiscoverer) getValidationFunctionName(t *types.Type) (types.Name, bool) {
	pkg, ok := td.inputToPkg[t.Name.Package]
	if !ok {
		return types.Name{}, false
	}
	return types.Name{Package: pkg, Name: "Validate_" + t.Name.Name}, true
}

// emitValidationForType writes code for inType, calling type-attached
// validations and then descending into the type (e.g. struct fields).
// inType is always a value type, with pointerness removed, and isVarPtr
// accomodates for that.
func (g *genValidations) emitValidationForType(c *generator.Context, inType *types.Type, isVarPtr bool, sw *generator.SnippetWriter, eachKey, eachVal []validators.FunctionGen) {
	if inType.Kind == types.Pointer {
		klog.Fatalf("unexpected pointer: %v", inType)
	}

	targs := generator.Args{
		"inType":    inType,
		"errorList": c.Universe.Type(errorListType),
		"fieldPath": c.Universe.Type(fieldPathType),
	}

	didSome := false // for prettier output later

	// Emit code for type-attached validations.
	tn := g.typeNodes[inType]
	if len(tn.validations) > 0 {
		sw.Do("// type $.inType|raw$\n", targs)
		g.emitCallsToValidators(c, tn.validations, isVarPtr, sw)
		sw.Do("\n", nil)
		didSome = true
	}

	// Descend into the type.
	switch inType.Kind {
	case types.Builtin:
		// Nothing further.
	case types.Alias:
		// Nothing further.
	case types.Struct:
		for _, child := range tn.children {
			if len(child.name) == 0 {
				klog.Fatalf("missing child name for field in %v", inType)
			}
			// Missing JSON name is checked iff we have code to emit.

			targs := targs.WithArgs(generator.Args{
				"fieldName": child.name,
				"fieldJSON": child.jsonName,
				"fieldType": child.underlyingType,
			})

			childIsPtr := child.underlyingType.Kind == types.Pointer

			// Accumulate into a buffer so we don't emit empty functions.
			buf := bytes.NewBuffer(nil)
			bufsw := sw.Dup(buf)

			if len(child.validations) > 0 {
				// When calling registered validators, we always pass the
				// underlying value-type.  E.g. if the field's type is string,
				// we pass string, and if the field's type is *string, we also
				// pass string (checking for nil, first).  This means those
				// validators don't have to know the difference, but it also
				// means that large structs will be passed by value.  If this
				// turns out to be a real problem, we could change this to pass
				// everything by pointer.
				g.emitCallsToValidators(c, child.validations, childIsPtr, bufsw)
			}

			// Get to the real type.
			t := child.underlyingType
			if t.Kind == types.Pointer {
				t = t.Elem
			}

			if t.Kind == types.Struct || t.Kind == types.Alias {
				// If this field is another type, call its validation function.
				// Checking for nil is handled inside this call.
				g.emitCallToOtherTypeFunc(c, t, childIsPtr, bufsw)
			} else {
				// Descend into this field.
				g.emitValidationForType(c, t, childIsPtr, bufsw, child.eachKey, child.eachVal)
			}

			if buf.Len() > 0 {
				if len(child.jsonName) == 0 {
					klog.Fatalf("missing child JSON name for field %v.%s", inType, child.name)
				}

				if didSome {
					sw.Do("\n", nil)
				}
				sw.Do("// field $.inType|raw$.$.fieldName$\n", targs)
				//TODO: pass val first or fldpath first?  validators do fldpath, why?
				sw.Do("errs = append(errs,\n", targs)
				sw.Do("  func(obj $.fieldType|raw$, fldPath *$.fieldPath|raw$) (errs $.errorList|raw$) {\n", targs)
				sw.Append(buf)
				sw.Do("    return\n", targs)
				sw.Do("  }(obj.$.fieldName$, fldPath.Child(\"$.fieldJSON$\"))...)\n", targs)
				sw.Do("\n", nil)
			} else {
				sw.Do("// field $.inType|raw$.$.fieldName$ has no validation\n", targs)
			}
			didSome = true
		}
	case types.Slice, types.Array:
		//FIXME: figure out if we can make this a wrapper-function and do it in one call to validate.ValuesInSlice()
		targs := targs.WithArgs(generator.Args{
			"elemType": inType.Elem,
		})

		elemIsPtr := inType.Elem.Kind == types.Pointer

		// Accumulate into a buffer so we don't emit empty functions.
		elemBuf := bytes.NewBuffer(nil)
		elemSW := sw.Dup(elemBuf)

		// Validate each value.
		validations := tn.elem.validations
		validations = append(validations, eachVal...)
		if len(validations) > 0 {
			// When calling registered validators, we always pass the
			// underlying value-type.  E.g. if the field's type is string,
			// we pass string, and if the field's type is *string, we also
			// pass string (checking for nil, first).  This means those
			// validators don't have to know the difference, but it also
			// means that large structs will be passed by value.  If this
			// turns out to be a real problem, we could change this to pass
			// everything by pointer.
			g.emitCallsToValidators(c, validations, elemIsPtr, elemSW)
		}

		// Get to the real type.
		t := inType.Elem
		if t.Kind == types.Pointer {
			t = t.Elem
		}

		if t.Kind == types.Struct || t.Kind == types.Alias {
			// If this field is another type, call its validation function.
			// Checking for nil is handled inside this call.
			g.emitCallToOtherTypeFunc(c, t, elemIsPtr, elemSW)
		} else {
			// No need to go further.  Struct- or alias-typed fields might have
			// validations attached to the type, but anything else (e.g.
			// string) can't, and we already emitted code for the field
			// validations.
		}

		if elemBuf.Len() > 0 {
			sw.Do("for i, val := range obj {\n", targs)
			sw.Do("  errs = append(errs,\n", targs)
			sw.Do("    func(obj $.elemType|raw$, fldPath *$.fieldPath|raw$) (errs $.errorList|raw$) {\n", targs)
			sw.Append(elemBuf)
			sw.Do("      return\n", targs)
			sw.Do("    }(val, fldPath.Index(i))...)\n", targs)
			sw.Do("}\n", nil)
		}
	case types.Map:
		targs := targs.WithArgs(generator.Args{
			"keyType": inType.Key,
			"valType": inType.Elem,
		})

		keyIsPtr := inType.Key.Kind == types.Pointer
		valIsPtr := inType.Elem.Kind == types.Pointer

		// Accumulate into a buffer so we don't emit empty functions.
		keyBuf := bytes.NewBuffer(nil)
		keySW := sw.Dup(keyBuf)

		// Validate each key.
		keyValidations := tn.key.validations
		keyValidations = append(keyValidations, eachKey...)
		if len(keyValidations) > 0 {
			// When calling registered validators, we always pass the
			// underlying value-type.  E.g. if the field's type is string,
			// we pass string, and if the field's type is *string, we also
			// pass string (checking for nil, first).  This means those
			// validators don't have to know the difference, but it also
			// means that large structs will be passed by value.  If this
			// turns out to be a real problem, we could change this to pass
			// everything by pointer.
			g.emitCallsToValidators(c, keyValidations, keyIsPtr, keySW)
		}

		// Get to the real type.
		t := inType.Key
		if t.Kind == types.Pointer {
			t = t.Elem
		}

		if t.Kind == types.Struct || t.Kind == types.Alias {
			// If this field is another type, call its validation function.
			// Checking for nil is handled inside this call.
			g.emitCallToOtherTypeFunc(c, t, keyIsPtr, keySW)
		} else {
			// No need to go further.  Struct- or alias-typed fields might have
			// validations attached to the type, but anything else (e.g.
			// string) can't, and we already emitted code for the field
			// validations.
		}

		// Accumulate into a buffer so we don't emit empty functions.
		valBuf := bytes.NewBuffer(nil)
		valSW := sw.Dup(valBuf)

		// Validate each value.
		valValidations := tn.elem.validations
		valValidations = append(valValidations, eachVal...)
		if len(valValidations) > 0 {
			// When calling registered validators, we always pass the
			// underlying value-type.  E.g. if the field's type is string,
			// we pass string, and if the field's type is *string, we also
			// pass string (checking for nil, first).  This means those
			// validators don't have to know the difference, but it also
			// means that large structs will be passed by value.  If this
			// turns out to be a real problem, we could change this to pass
			// everything by pointer.
			g.emitCallsToValidators(c, valValidations, valIsPtr, valSW)
		}

		// Get to the real type.
		t = inType.Elem
		if t.Kind == types.Pointer {
			t = t.Elem
		}

		if t.Kind == types.Struct || t.Kind == types.Alias {
			// If this field is another type, call its validation function.
			// Checking for nil is handled inside this call.
			g.emitCallToOtherTypeFunc(c, t, valIsPtr, valSW)
		} else {
			// No need to go further.  Struct- or alias-typed fields might have
			// validations attached to the type, but anything else (e.g.
			// string) can't, and we already emitted code for the field
			// validations.
		}

		kName, vName := "_", "_"
		if keyBuf.Len() > 0 {
			kName = "key"
		}
		if valBuf.Len() > 0 {
			vName = "val"
		}
		if keyBuf.Len()+valBuf.Len() > 0 {
			sw.Do("for $.key$, $.val$ := range obj {\n", targs.WithArgs(generator.Args{"key": kName, "val": vName}))
			if keyBuf.Len() > 0 {
				sw.Do("  errs = append(errs,\n", targs)
				sw.Do("    func(obj $.keyType|raw$, fldPath *$.fieldPath|raw$) (errs $.errorList|raw$) {\n", targs)
				sw.Append(keyBuf)
				sw.Do("      return\n", targs)
				sw.Do("    }(key, fldPath)...)\n", targs) // TODO: we need a way to denote "invalid key"
			}
			if valBuf.Len() > 0 {
				sw.Do("  errs = append(errs,\n", targs)
				sw.Do("    func(obj $.valType|raw$, fldPath *$.fieldPath|raw$) (errs $.errorList|raw$) {\n", targs)
				sw.Append(valBuf)
				sw.Do("      return\n", targs)
				sw.Do("    }(val, fldPath.Key(key))...)\n", nil) // TODO: what if the key is not a string?
			}
			sw.Do("}\n", nil)
		}
	default:
		klog.Fatalf("unhandled type: %v (%s)", inType, inType.Kind)
	}
}

// emitCallToOtherTypeFunc generates a call to a different generated validation
// function for a field in some parent context.  inType is the value type
// being validated with pointerness removed.  isVarPtr indicates that the value
// was a pointer in the parent context.  The variable in question is always
// named "obj" and the field path is always "fldPath".
func (g *genValidations) emitCallToOtherTypeFunc(c *generator.Context, inType *types.Type, isVarPtr bool, sw *generator.SnippetWriter) {
	// If this type has no validations (transitively) then we don't need to do
	// anything.
	if !g.hasValidations(inType) {
		return
	}

	addr := "" // adjusted below if needed
	if isVarPtr {
		sw.Do("if obj != nil {\n", nil)
		defer func() {
			sw.Do("}\n", nil)
		}()
	} else {
		addr = "&"
	}

	tn := g.typeNodes[inType]
	targs := generator.Args{
		"addr":     addr,
		"funcName": c.Universe.Type(tn.funcName),
	}
	sw.Do("errs = append(errs, $.funcName|raw$($.addr$obj, fldPath)...)\n", targs)
}

// emitCallsToValidators generates calls to a list of validation functions for
// a single field or type. validations is a list of functions to call, with
// arguments.  The name of this value is always "obj" and the field path is
// "fldPath".  isVarPtr indicates that the value  was a pointer in the parent
// context.
func (g *genValidations) emitCallsToValidators(c *generator.Context, validations []validators.FunctionGen, isVarPtr bool, sw *generator.SnippetWriter) {
	// Helper func
	sort := func(in []validators.FunctionGen) []validators.FunctionGen {
		fatal := make([]validators.FunctionGen, 0, len(in))
		fatalPtr := make([]validators.FunctionGen, 0, len(in))
		nonfatal := make([]validators.FunctionGen, 0, len(in))
		nonfatalPtr := make([]validators.FunctionGen, 0, len(in))

		for _, fg := range in {
			isFatal := (fg.Flags()&validators.IsFatal != 0)
			isPtrOK := (fg.Flags()&validators.PtrOK != 0)

			if isFatal {
				if isPtrOK {
					fatalPtr = append(fatalPtr, fg)
				} else {
					fatal = append(fatal, fg)
				}
			} else {
				if isPtrOK {
					nonfatalPtr = append(nonfatalPtr, fg)
				} else {
					nonfatal = append(nonfatal, fg)
				}
			}
		}
		result := fatalPtr
		result = append(result, fatal...)
		result = append(result, nonfatalPtr...)
		result = append(result, nonfatal...)
		return result
	}

	validations = sort(validations)

	insideNilCheck := false
	for _, v := range validations {
		ptrOK := (v.Flags()&validators.PtrOK != 0)
		isFatal := (v.Flags()&validators.IsFatal != 0)

		fn, extraArgs := v.SignatureAndArgs()
		targs := generator.Args{
			"funcName": c.Universe.Type(fn),
			"deref":    "", // updated below if needed
		}
		if isVarPtr && !ptrOK {
			if !insideNilCheck {
				sw.Do("if obj != nil {\n", targs)
				insideNilCheck = true
			}
			targs["deref"] = "*"
		} else {
			if insideNilCheck {
				sw.Do("}\n", nil)
				insideNilCheck = false
			}
		}

		emitCall := func() {
			sw.Do("$.funcName|raw$(fldPath, $.deref$obj", targs)
			for _, arg := range extraArgs {
				sw.Do(", "+toGolangSourceDataLiteral(arg), nil)
			}
			sw.Do(")", targs)
		}

		if isFatal {
			sw.Do("if e := ", nil)
			emitCall()
			sw.Do("; len(e) != 0 {\n", nil)
			sw.Do("errs = append(errs, e...)\n", nil)
			sw.Do("    return // fatal\n", nil)
			sw.Do("}\n", nil)
		} else {
			sw.Do("errs = append(errs, ", nil)
			emitCall()
			sw.Do("...)\n", nil)
		}
	}
	if insideNilCheck {
		sw.Do("}\n", nil)
	}
}

func toGolangSourceDataLiteral(value any) string {
	// For safety, be strict in what values we output to visited source, and ensure strings
	// are quoted.
	switch value.(type) {
	case uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, bool:
		return fmt.Sprintf("%v", value)
	case string:
		// If the incoming string was quoted, we still do it ourselves, JIC.
		str := value.(string)
		if s, err := strconv.Unquote(str); err == nil {
			str = s
		}
		return fmt.Sprintf("%q", str)
	}
	klog.Fatalf("Unsupported extraArg type: %T", value)
	return ""
}
