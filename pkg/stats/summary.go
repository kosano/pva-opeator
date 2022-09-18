package stats

import (
	"context"
	"encoding/json"
	"log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	statsapi "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Summary struct {
	client      client.Client
	NodeSummary map[string]*statsapi.Summary
	VolumeStats map[string]*statsapi.VolumeStats
}

func NewSummary(client client.Client) *Summary {
	return &Summary{
		client:      client,
		NodeSummary: make(map[string]*statsapi.Summary),
		VolumeStats: make(map[string]*statsapi.VolumeStats),
	}
}

func (s *Summary) Run(ctx context.Context) {
	s.generateNodesSummary(ctx)
	s.generateVolumeStats()
}

func (s *Summary) generateVolumeStats() {
	for _, ns := range s.NodeSummary {
		for _, ps := range ns.Pods {
			for _, vs := range ps.VolumeStats {
				if vs.PVCRef != nil {
					s.VolumeStats[vs.PVCRef.Name] = &vs
				}
			}
		}
	}
}

func (s *Summary) generateNodesSummary(ctx context.Context) {
	nodes := &v1.NodeList{}
	s.client.List(ctx, nodes)
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	restCall, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err.Error())
	}
	for _, node := range nodes.Items {
		nodeSummaryBytes, err := restCall.Get().
			AbsPath("api/v1/nodes", node.Name, "proxy/stats/summary").
			Do(context.TODO()).
			Raw()
		if err != nil {
			log.Fatal(err)
			continue
		}
		var summary statsapi.Summary
		json.Unmarshal(nodeSummaryBytes, &summary)
		s.NodeSummary[node.Name] = &summary
	}
}
