package controllers

import (
	"github.com/dataworkbench/hdfs-operator/api/v1"
	com "github.com/dataworkbench/hdfs-operator/common"
	dn "github.com/dataworkbench/hdfs-operator/controllers/datanode"
	jn "github.com/dataworkbench/hdfs-operator/controllers/journalnode"
	nn "github.com/dataworkbench/hdfs-operator/controllers/namenode"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type HdfsResources struct {
	StatefulSets  []appsv1.StatefulSet
	Datanode      appsv1.StatefulSet
	ConfigMaps    []corev1.ConfigMap
	Services      []corev1.Service
}

func BuildExpectedResources(hdfs v1.HDFS) (HdfsResources, error) {

	configs, err := BuildConfigMaps(hdfs)
	if err != nil {
		return HdfsResources{}, err
	}

	services, err := BuildServices(hdfs)
	if err != nil {
		return HdfsResources{}, err
	}

	statefulSets, err := BuildStatefulSets(hdfs)
	if err != nil {
		return HdfsResources{}, err
	}

	dnSet, err := dn.BuildStatefulSet(hdfs)
	if err != nil {
		return HdfsResources{}, err
	}

	return HdfsResources{
		StatefulSets: statefulSets,
		Datanode:     dnSet,
		ConfigMaps:   configs ,
		Services:     services,
	}, nil
}

func BuildConfigMaps(hdfs v1.HDFS) (c []corev1.ConfigMap,err error) {

	config, err := com.BuildHdfsConfig(hdfs, com.GetName(hdfs.Name, com.CommonConfigName))
	if err != nil {
		return c, err
	}
	nnScripts := nn.BuildConfigMap(hdfs)
	dnScripts := dn.BuildConfigMap(hdfs)

	return append(c, config, nnScripts, dnScripts), nil
}

func BuildServices(hdfs v1.HDFS) (svc []corev1.Service,err error) {

	nnSvc := com.HeadlessService(hdfs,
		com.GetName(hdfs.Name, hdfs.Spec.Namenode.Name),
		nn.GetDefaultServicePorts())
	jnSvc := com.HeadlessService(hdfs, com.GetName(hdfs.Name,
		hdfs.Spec.Journalnode.Name),
		jn.GetDefaultServicePorts())

	return append(svc, nnSvc, jnSvc ), nil
}

func BuildStatefulSets(hdfs v1.HDFS) (s []appsv1.StatefulSet,err error) {

	nnStatefulSet, err := nn.BuildStatefulSet(hdfs)
	if err != nil {
		return s, err
	}
	jnStatefulSet, err := jn.BuildStatefulSet(hdfs)
	if err != nil {
		return s, err
	}

	//if hdfs.Spec.ZkQuorum == "" {
	//	zkStatefulSet, err := zk.BuildStatefulSet(hdfs)
	//	if err != nil {
	//		return s, err
	//	}
	//	s = append(s, zkStatefulSet)
	//}

	return append(s, nnStatefulSet, jnStatefulSet ), nil
}