// Cluster-wide permissions the controller manager needs regardless of which
// reconciler is running. Kept as free-floating package-level markers: markers
// attached directly to a method's doc comment are not collected by
// controller-gen, so they must not live on a Reconcile function.
//
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// controllers contains implementation of kubeserial controllers
package controllers
