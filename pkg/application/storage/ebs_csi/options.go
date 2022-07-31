package ebs_csi

import (
	"context"
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EbsCsiOptions struct {
	application.ApplicationOptions

	DefaultGp3 bool
}

const IsDefaultStorageClassAnnotation = "storageclass.kubernetes.io/is-default-class"

func newOptions() (options *EbsCsiOptions, flags cmd.Flags) {
	options = &EbsCsiOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "ebs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.9.0",
				Latest:        "v1.10.0",
				PreviousChart: "2.6.11",
				Previous:      "v1.7.0",
			},
		},
		DefaultGp3: false,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "default-gp3",
				Description: "set gp3 StorageClass as default",
			},
			Option: &options.DefaultGp3,
		},
	}
	return
}

func (o *EbsCsiOptions) PreInstall() error {
	if !o.DefaultGp3 {
		return nil
	}

	if o.DryRun {
		fmt.Println("\nPreInstall Dry Run:")
		fmt.Println("Mark the current default StorageClass as non-default")
		return nil
	}

	fmt.Println("Checking for default StorageClass")
	k8sclient, err := kubernetes.Client(o.KubeContext())
	if err != nil {
		return err
	}

	scList, err := k8sclient.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, sc := range scList.Items {
		if sc.Annotations[IsDefaultStorageClassAnnotation] != "true" {
			continue
		}
		fmt.Printf("Marking StorageClass %q as non-default...", sc.Name)

		sc.Annotations[IsDefaultStorageClassAnnotation] = "false"
		k8sclient.StorageV1().StorageClasses().Update(context.Background(), &sc, metav1.UpdateOptions{})

		fmt.Println("done")
	}

	return nil
}
