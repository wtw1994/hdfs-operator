package namenode

import (
	v1 "github.com/dataworkbench/hdfs-operator/api/v1"
	com "github.com/dataworkbench/hdfs-operator/common"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildStatefulSet(hdfs v1.HDFS) (appsv1.StatefulSet, error) {
	statefulSetName := com.GetName(hdfs.Name, hdfs.Spec.Namenode.Name)
	// ssetSelector is used to match the StatefulSet pods
	ssetSelector := com.NewStatefulSetLabels(com.ExtractNamespacedName(&hdfs), statefulSetName)
	// add default PVCs to the node spec if no user defined PVCs exist
	hdfs.Spec.Namenode.VolumeClaimTemplates = com.AppendDefaultPVCs(hdfs.Spec.Namenode.VolumeClaimTemplates, "metadatadir", hdfs.Spec.Namenode.StorageClass)
	// build pod template,associate PVCs to pod container
	podTemplate, err := BuildPodTemplateSpec(hdfs, ssetSelector)
	if err != nil {
		return appsv1.StatefulSet{}, err
	}

	sset := appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: hdfs.Namespace,
			Name:      statefulSetName,
			Labels:    ssetSelector,
			OwnerReferences: com.GetOwnerReference(hdfs),
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: statefulSetName,
			Selector: &metav1.LabelSelector{
				MatchLabels: ssetSelector,
			},
			Replicas:             &hdfs.Spec.Namenode.Replicas,
			VolumeClaimTemplates: hdfs.Spec.Namenode.VolumeClaimTemplates,
			Template:             podTemplate,
		},
	}
	return sset, nil
}
