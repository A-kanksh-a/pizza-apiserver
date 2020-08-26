/*
Copyright 2017 The Kubernetes Authors.

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

package pizzasize

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/apiserver/pkg/admission"

	"github.com/programming-kubernetes/pizza-apiserver/pkg/admission/custominitializer"
	"github.com/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	informers "github.com/programming-kubernetes/pizza-apiserver/pkg/generated/informers/externalversions"
	listers "github.com/programming-kubernetes/pizza-apiserver/pkg/generated/listers/restaurant/v2alpha1"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register("PizzaSize", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

type PizzaSizePlugin struct {
	*admission.Handler
	pizzaLister listers.PizzaLister
}

var _ = custominitializer.WantsRestaurantInformerFactory(&PizzaSizePlugin{})
var _ = admission.ValidationInterface(&PizzaSizePlugin{})

// Admit ensures that the object in-flight is of kind Pizza.
// In addition checks that the toppings are known.
func (d *PizzaSizePlugin) Validate(a admission.Attributes, _ admission.ObjectInterfaces) error {
	// we are only interested in pizzas
	if a.GetKind().GroupKind() != restaurant.Kind("Pizza") {
		return nil
	}

	if !d.WaitForReady() {
		return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	}

	obj := a.GetObject()
	pizza := obj.(*restaurant.Pizza)

	pizzaListing, _ := d.pizzaLister.List(labels.Everything())
	for _, piz := range pizzaListing {
		if piz.Spec.Size == pizza.Spec.Size {
			return admission.NewForbidden(
				a,
				fmt.Errorf("size already present: %s", piz.Spec.Size),
			)
		}
	}

	return nil
}

// SetRestaurantInformerFactory gets Lister from SharedInformerFactory.
// The lister knows how to lists Toppings.
func (d *PizzaSizePlugin) SetRestaurantInformerFactory(f informers.SharedInformerFactory) {
	d.pizzaLister = f.Restaurant().V2alpha1().Pizzas().Lister()
	d.SetReadyFunc(f.Restaurant().V2alpha1().Pizzas().Informer().HasSynced)
}

// ValidaValidateInitializationte checks whether the plugin was correctly initialized.
func (d *PizzaSizePlugin) ValidateInitialization() error {
	if d.pizzaLister == nil {
		return fmt.Errorf("missing policy lister")
	}
	return nil
}

// New creates a new ban pizza topping admission plugin
func New() (*PizzaSizePlugin, error) {
	return &PizzaSizePlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}
