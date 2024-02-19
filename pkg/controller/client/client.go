package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/scale"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"

	qdrantClient "github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/clientset/versioned"
)

var scaleConverter = scale.NewScaleConverter()
var codecs = serializer.NewCodecFactory(scaleConverter.Scheme())

type Client struct {
	ctrl.Client
	kubernetes.Interface

	Qdrant qdrantClient.Interface

	Discovery discovery.DiscoveryInterface

	scheme *runtime.Scheme
	config *rest.Config
	rest   rest.Interface
}

func NewClient(cfg *rest.Config, scheme *runtime.Scheme, cc ctrl.Client) (*Client, error) {

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	restClient, err := NewRESTClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	qClient, err := qdrantClient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	c := Client{
		Client:    cc,
		Interface: kubeClient,
		Qdrant:    qClient,
		Discovery: discoveryClient,
		scheme:    scheme,
		config:    cfg,
		rest:      restClient,
	}

	return &c, nil
}

func NewRESTClientForConfig(config *rest.Config) (*rest.RESTClient, error) {
	cfg := rest.CopyConfig(config)
	// so that the RESTClientFor doesn't complain
	cfg.GroupVersion = &schema.GroupVersion{}
	cfg.NegotiatedSerializer = codecs.WithoutConversion()
	if len(cfg.UserAgent) == 0 {
		cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return rest.RESTClientFor(cfg)
}
