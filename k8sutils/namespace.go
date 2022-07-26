package k8sutils

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func MakesureNamespaceExists(ctx context.Context, cliset *kubernetes.Clientset, ns string) (err error) {
	nsCli := cliset.CoreV1().Namespaces()
	_, err = nsCli.Get(ctx, ns, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		return errors.Wrapf(err, "failed to get namespace %s", ns)
	}
	if isNotFound {
		_, err = nsCli.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				err = nil
			} else {
				err = errors.Wrapf(err, "failed to create namespace %s", ns)
				return
			}
		}
	}
	return
}
