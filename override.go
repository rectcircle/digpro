package digpro

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rectcircle/digpro/internal"
	"go.uber.org/dig"
)

// Override a registered provider and only support digpro high level api (support *digpro.ContainerWrapper and digglobal).
// if Container not exist provider, the option will return error: no provider to override was found.
// for example
//   c := digpro.New()
//   _ = c.Supply(1) // please handle error in production
//   _ = c.Supply(1, digpro.Override())
//   // _ = c.Supply("a", digpro.Override())  // has error
//   i, _ := c.Extract(0)
//   fmt.Println(i.(int) == 1)
//   // Output: true
func Override() dig.ProvideOption {
	return overrideProvideOption{}
}

func overrideProvideMiddleware(pc *provideContext) error {

	opts, digproOptResult := filterProvideOptionAndGetDigproOptions(pc.opts, overrideProvideOptionType)
	hasOverrideOpt := digproOptResult.enableOverride
	pc.opts = opts

	// get ProviderInfo
	c := dig.New()
	info := internal.ProvideInfosWrapper{}
	err := c.Provide(pc.constructor, append([]dig.ProvideOption{dig.FillProvideInfo(&info.ProvideInfo)}, pc.opts...)...)
	if err != nil {
		return pc.next()
	}

	hasGroupOpt := false
	outputs := info.ExportedOutputs()
	if len(outputs) >= 1 && outputs[0].Group != "" {
		hasGroupOpt = true
	}

	if hasGroupOpt && hasOverrideOpt {
		return errors.New("cannot use digpro.Override() with value groups")
	}
	if !hasOverrideOpt {
		return pc.next()
	}

	// check and remove conflict provider
	recoverOld, err := removeOldConflictProvideOutputs(pc.c, outputs)
	if err != nil {
		return err
	}
	err = pc.next()
	if err != nil {
		recoverOld()
	}

	return err
}

func removeOldConflictProvideOutputs(c *ContainerWrapper, outputs []internal.ProvideOutput) (recoverOld func(), err error) {
	// check has conflict? if not will return error
	index := -1
	for i, info := range c.provideInfos {
		if internal.EqualsProvideOutputs(info.ExportedOutputs(), outputs) {
			index = i
			break
		}
	}
	if index == -1 {
		err = errors.New("no provider to override was found")
		return
	}

	containerValue := reflect.ValueOf(&c.Container).Elem()

	providersValue := internal.EnsureValueExported(containerValue.FieldByName("providers")) // map[dig.key][]*dig.node
	nodesValue := internal.EnsureValueExported(containerValue.FieldByName("nodes"))         // []*node

	// create all keys
	keyType := providersValue.Type().Key()
	keys := []reflect.Value{} // dig.key
	for _, output := range outputs {
		keyValue := reflect.New(keyType).Elem()
		internal.EnsureValueExported(keyValue.FieldByName("t")).Set(reflect.ValueOf(output.Type))
		internal.EnsureValueExported(keyValue.FieldByName("name")).Set(reflect.ValueOf(output.Name))
		internal.EnsureValueExported(keyValue.FieldByName("group")).Set(reflect.ValueOf(output.Group))
		keys = append(keys, keyValue)
	}
	// get all keyNodes
	keyNodes := []reflect.Value{} // []*dig.node
	for i, key := range keys {
		node := providersValue.MapIndex(key)
		if !node.IsValid() {
			// dead code
			err = fmt.Errorf("no provider to override was found: [%d].%s", i, outputs[i])
			return
		}
		keyNodes = append(keyNodes, node)
	}
	// check old outputs from same node
	var finalNodes *reflect.Value // []*dig.node
	for i, node := range keyNodes {
		if finalNodes == nil {
			finalNodes = &node
		} else {
			if !reflect.DeepEqual(finalNodes.Interface(), node.Interface()) {
				// dead code
				err = fmt.Errorf("the registered provider of [%d].%s is different from other outputs", i, outputs[i])
				return
			}
		}
	}
	if finalNodes == nil {
		// dead code
		err = errors.New("no provider to override was found")
		return
	}
	if finalNodes.Len() != 1 {
		// dead code
		err = errors.New("unknown error: len(finalNode) != 1")
		return
	}
	finalNode := finalNodes.Index(0).Elem() // dig.node
	// node not allow called
	if internal.EnsureValueExported(finalNode.FieldByName("called")).Interface().(bool) {
		err = fmt.Errorf("the old provider has called, digpro.Override only use before call Invoke()")
		return
	}

	// delete nodes
	oldNodes := reflect.MakeSlice(nodesValue.Type(), 0, nodesValue.Len())
	newNodes := reflect.MakeSlice(nodesValue.Type(), 0, nodesValue.Len()-1)
	for i := 0; i < nodesValue.Len(); i++ {
		if nodesValue.Index(i).Elem() != finalNode {
			reflect.Append(newNodes, nodesValue.Index(i))
		} else {
			reflect.Append(oldNodes, nodesValue.Index(i))
		}
	}
	nodesValue.Set(newNodes)
	// delete container
	for _, key := range keys {
		providersValue.SetMapIndex(key, reflect.Value{})
	}
	recoverOld = func() {
		// recover node
		nodesValue.Set(oldNodes)
		// recover container
		for _, key := range keys {
			providersValue.SetMapIndex(key, *finalNodes)
		}
	}
	return
}
