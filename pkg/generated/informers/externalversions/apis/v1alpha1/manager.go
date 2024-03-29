/*
Copyright 2022.

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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	apisv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	versioned "github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/janekbaraniewski/kubeserial/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/generated/listers/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ManagerInformer provides access to a shared informer and lister for
// Managers.
type ManagerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ManagerLister
}

type managerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewManagerInformer constructs a new informer for Manager type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewManagerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredManagerInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredManagerInformer constructs a new informer for Manager type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredManagerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppV1alpha1().Managers().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AppV1alpha1().Managers().Watch(context.TODO(), options)
			},
		},
		&apisv1alpha1.Manager{},
		resyncPeriod,
		indexers,
	)
}

func (f *managerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredManagerInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *managerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apisv1alpha1.Manager{}, f.defaultInformer)
}

func (f *managerInformer) Lister() v1alpha1.ManagerLister {
	return v1alpha1.NewManagerLister(f.Informer().GetIndexer())
}
