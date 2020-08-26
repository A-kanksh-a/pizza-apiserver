/*
Copyright 2018 The Kubernetes Authors.

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

package v1beta1

import (
	"github.com/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	err := scheme.AddConversionFuncs(
		Convert_v1beta1_PizzaSpec_To_restaurant_PizzaSpec,
		Convert_restaurant_PizzaSpec_To_v1beta1_PizzaSpec,
	)
	if err != nil {
		return err
	}

	return nil
}

func Convert_v1beta1_PizzaSpec_To_restaurant_PizzaSpec(in *PizzaSpec, out *restaurant.PizzaSpec, s conversion.Scope) error {

	for _, top := range in.Toppings {
		out.Toppings = append(out.Toppings, restaurant.PizzaTopping{
			Name:     top.Name,
			Quantity: top.Quantity,
		})
	}
	out.Size = "Medium"
	return nil
}

func Convert_restaurant_PizzaSpec_To_v1beta1_PizzaSpec(in *restaurant.PizzaSpec, out *PizzaSpec, s conversion.Scope) error {
	for _, top := range in.Toppings {
		out.Toppings = append(out.Toppings, PizzaTopping{
			Name:     top.Name,
			Quantity: top.Quantity,
		})
	}
	return nil
}
