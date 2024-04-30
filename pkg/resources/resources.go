package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ToUnstructured(s *runtime.Scheme, obj runtime.Object) (*unstructured.Unstructured, error) {
	switch ot := obj.(type) {
	case *unstructured.Unstructured:
		return ot, nil
	default:
		var err error
		var u unstructured.Unstructured

		u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to unstructured: %w", err)
		}

		gvk := u.GroupVersionKind()
		if gvk.Group == "" || gvk.Kind == "" {
			gvks, _, err := s.ObjectKinds(obj)
			if err != nil {
				return nil, fmt.Errorf("failed to convert to unstructured - unable to get GVK %w", err)
			}
			apiv, k := gvks[0].ToAPIVersionAndKind()

			u.SetAPIVersion(apiv)
			u.SetKind(k)
		}

		return &u, nil
	}
}
