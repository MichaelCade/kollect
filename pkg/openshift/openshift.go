package openshift

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	buildv1 "github.com/openshift/api/build/v1"
	configv1 "github.com/openshift/api/config/v1"
	consolev1 "github.com/openshift/api/console/v1"
	imagev1 "github.com/openshift/api/image/v1"
	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	projectv1 "github.com/openshift/api/project/v1"
	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	buildclientset "github.com/openshift/client-go/build/clientset/versioned"
	configclientset "github.com/openshift/client-go/config/clientset/versioned"
	consoleclientset "github.com/openshift/client-go/console/clientset/versioned"
	imageclientset "github.com/openshift/client-go/image/clientset/versioned"
	machineconfigclientset "github.com/openshift/client-go/machineconfiguration/clientset/versioned"
	oauthclientset "github.com/openshift/client-go/oauth/clientset/versioned"
	operatorclientset "github.com/openshift/client-go/operator/clientset/versioned"
	projectclientset "github.com/openshift/client-go/project/clientset/versioned"
	routeclientset "github.com/openshift/client-go/route/clientset/versioned"
	securityclientset "github.com/openshift/client-go/security/clientset/versioned"
	operatorframeworkclientset "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenShiftData represents the collected OpenShift 4.x data
type OpenShiftData struct {
	// Cluster information
	ClusterInfo      ClusterInfo                `json:"clusterInfo,omitempty"`
	ClusterOperators []configv1.ClusterOperator `json:"clusterOperators,omitempty"`
	ClusterVersions  []configv1.ClusterVersion  `json:"clusterVersions,omitempty"`
	DNSes            []configv1.DNS             `json:"dnses,omitempty"`
	Infrastructures  []configv1.Infrastructure  `json:"infrastructures,omitempty"`
	Networks         []configv1.Network         `json:"networks,omitempty"`

	// Node configuration
	MachineConfigs     []machineconfigv1.MachineConfig     `json:"machineConfigs,omitempty"`
	MachineConfigPools []machineconfigv1.MachineConfigPool `json:"machineConfigPools,omitempty"`

	// Operators and OLM
	OperatorGroups []operatorsv1.OperatorGroup      `json:"operatorGroups,omitempty"`
	Subscriptions  []operatorsv1alpha1.Subscription `json:"subscriptions,omitempty"`
	InstallPlans   []operatorsv1alpha1.InstallPlan  `json:"installPlans,omitempty"`

	// Routes and Ingress
	Routes             []routev1.Route                `json:"routes,omitempty"`
	IngressControllers []operatorv1.IngressController `json:"ingressControllers,omitempty"`

	// Projects and Security
	Projects                   []projectv1.Project                     `json:"projects,omitempty"`
	SecurityContextConstraints []securityv1.SecurityContextConstraints `json:"securityContextConstraints,omitempty"`

	// Developer Experience
	BuildConfigs []buildv1.BuildConfig `json:"buildConfigs,omitempty"`
	ImageStreams []imagev1.ImageStream `json:"imageStreams,omitempty"`

	// Console Elements
	ConsoleLinks         []consolev1.ConsoleLink         `json:"consoleLinks,omitempty"`
	ConsoleNotifications []consolev1.ConsoleNotification `json:"consoleNotifications,omitempty"`
	ConsoleCLIDownloads  []consolev1.ConsoleCLIDownload  `json:"consoleCLIDownloads,omitempty"`

	// Authentication
	OAuthClients []oauthv1.OAuthClient `json:"oAuthClients,omitempty"`
}

// ClusterInfo contains information about the OpenShift 4.x cluster
type ClusterInfo struct {
	Name                 string   `json:"name"`
	ID                   string   `json:"id,omitempty"`
	Version              string   `json:"version"`
	Channel              string   `json:"channel,omitempty"`
	BaseDomain           string   `json:"baseDomain,omitempty"`
	Platform             string   `json:"platform"`
	Provider             string   `json:"provider"`
	Region               string   `json:"region,omitempty"`
	ControlPlaneTopology string   `json:"controlPlaneTopology,omitempty"`
	InfraTopology        string   `json:"infraTopology,omitempty"`
	APIServerURL         string   `json:"apiServerUrl"`
	ConsoleURL           string   `json:"consoleUrl,omitempty"`
	FeatureGates         []string `json:"featureGates,omitempty"`
}

// CheckCredentials checks if the OpenShift 4.x credentials are valid
func CheckCredentials(ctx context.Context, kubeconfigPath string, contextName ...string) (bool, error) {
	config, err := buildConfig(kubeconfigPath, contextName...)
	if err != nil {
		return false, fmt.Errorf("failed to build Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	_, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list namespaces: %v", err)
	}

	// Check if this is actually an OpenShift 4.x cluster
	ocpConfigClient, err := configclientset.NewForConfig(config)
	if err != nil {
		return false, fmt.Errorf("not an OpenShift cluster or missing permissions: %v", err)
	}

	clusterVersions, err := ocpConfigClient.ConfigV1().ClusterVersions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("not an OpenShift cluster or missing permissions: %v", err)
	}

	// Verify it's OpenShift 4.x
	if len(clusterVersions.Items) > 0 {
		version := clusterVersions.Items[0].Status.Desired.Version
		if version[0:1] != "4" {
			return false, fmt.Errorf("unsupported OpenShift version: %s. Only OpenShift 4.x is supported", version)
		}
		return true, nil
	}

	return false, fmt.Errorf("could not determine OpenShift version")
}

// CollectOpenShiftData collects data from an OpenShift 4.x cluster
func CollectOpenShiftData(ctx context.Context, kubeconfigPath string, contextName ...string) (interface{}, error) {
	data := OpenShiftData{}

	config, err := buildConfig(kubeconfigPath, contextName...)
	if err != nil {
		return nil, fmt.Errorf("failed to build Kubernetes config: %v", err)
	}

	// Standard Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Create OpenShift clients
	routeClient, err := routeclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create route client: %v", err)
	}

	buildClient, err := buildclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create build client: %v", err)
	}

	imageClient, err := imageclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create image client: %v", err)
	}

	securityClient, err := securityclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create security client: %v", err)
	}

	projectClient, err := projectclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create project client: %v", err)
	}

	configClient, err := configclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create config client: %v", err)
	}

	machineConfigClient, err := machineconfigclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create machine config client: %v", err)
	}

	operatorClient, err := operatorclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create operator client: %v", err)
	}

	operatorFrameworkClient, err := operatorframeworkclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create operator framework client: %v", err)
	}

	consoleClient, err := consoleclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create console client: %v", err)
	}

	oauthClient, err := oauthclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Warning: Failed to create OAuth client: %v", err)
	}

	// Collect cluster info
	clusterInfo, err := collectClusterInfo(ctx, configClient, clientset)
	if err != nil {
		log.Printf("Warning: Failed to collect cluster info: %v", err)
	}
	data.ClusterInfo = clusterInfo

	// Collect Route resources
	if routeClient != nil {
		routes, err := collectRoutes(ctx, routeClient)
		if err != nil {
			log.Printf("Warning: Failed to collect Routes: %v", err)
		} else {
			data.Routes = routes
		}
	}

	// Collect BuildConfig resources
	if buildClient != nil {
		buildConfigs, err := collectBuildConfigs(ctx, buildClient)
		if err != nil {
			log.Printf("Warning: Failed to collect BuildConfigs: %v", err)
		} else {
			data.BuildConfigs = buildConfigs
		}
	}

	// Collect Image resources
	if imageClient != nil {
		imageStreams, err := collectImageStreams(ctx, imageClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ImageStreams: %v", err)
		} else {
			data.ImageStreams = imageStreams
		}
	}

	// Collect SecurityContextConstraints
	if securityClient != nil {
		sccs, err := collectSecurityContextConstraints(ctx, securityClient)
		if err != nil {
			log.Printf("Warning: Failed to collect SecurityContextConstraints: %v", err)
		} else {
			data.SecurityContextConstraints = sccs
		}
	}

	// Collect Projects
	if projectClient != nil {
		projects, err := collectProjects(ctx, projectClient)
		if err != nil {
			log.Printf("Warning: Failed to collect Projects: %v", err)
		} else {
			data.Projects = projects
		}
	}

	// Collect Config resources
	if configClient != nil {
		clusterOperators, err := collectClusterOperators(ctx, configClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ClusterOperators: %v", err)
		} else {
			data.ClusterOperators = clusterOperators
		}

		clusterVersions, err := collectClusterVersions(ctx, configClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ClusterVersions: %v", err)
		} else {
			data.ClusterVersions = clusterVersions
		}

		dnses, err := collectDNSes(ctx, configClient)
		if err != nil {
			log.Printf("Warning: Failed to collect DNSes: %v", err)
		} else {
			data.DNSes = dnses
		}

		infrastructures, err := collectInfrastructures(ctx, configClient)
		if err != nil {
			log.Printf("Warning: Failed to collect Infrastructures: %v", err)
		} else {
			data.Infrastructures = infrastructures
		}

		networks, err := collectNetworks(ctx, configClient)
		if err != nil {
			log.Printf("Warning: Failed to collect Networks: %v", err)
		} else {
			data.Networks = networks
		}
	}

	// Collect MachineConfig and MachineConfigPool
	if machineConfigClient != nil {
		machineConfigs, err := collectMachineConfigs(ctx, machineConfigClient)
		if err != nil {
			log.Printf("Warning: Failed to collect MachineConfigs: %v", err)
		} else {
			data.MachineConfigs = machineConfigs
		}

		machineConfigPools, err := collectMachineConfigPools(ctx, machineConfigClient)
		if err != nil {
			log.Printf("Warning: Failed to collect MachineConfigPools: %v", err)
		} else {
			data.MachineConfigPools = machineConfigPools
		}
	}

	// Collect Operator Framework resources
	if operatorFrameworkClient != nil {
		operatorGroups, err := collectOperatorGroups(ctx, operatorFrameworkClient)
		if err != nil {
			log.Printf("Warning: Failed to collect OperatorGroups: %v", err)
		} else {
			data.OperatorGroups = operatorGroups
		}

		subscriptions, err := collectSubscriptions(ctx, operatorFrameworkClient)
		if err != nil {
			log.Printf("Warning: Failed to collect Subscriptions: %v", err)
		} else {
			data.Subscriptions = subscriptions
		}

		installPlans, err := collectInstallPlans(ctx, operatorFrameworkClient)
		if err != nil {
			log.Printf("Warning: Failed to collect InstallPlans: %v", err)
		} else {
			data.InstallPlans = installPlans
		}
	}

	// Collect Console resources
	if consoleClient != nil {
		consoleLinks, err := collectConsoleLinks(ctx, consoleClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ConsoleLinks: %v", err)
		} else {
			data.ConsoleLinks = consoleLinks
		}

		consoleNotifications, err := collectConsoleNotifications(ctx, consoleClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ConsoleNotifications: %v", err)
		} else {
			data.ConsoleNotifications = consoleNotifications
		}

		consoleCLIDownloads, err := collectConsoleCLIDownloads(ctx, consoleClient)
		if err != nil {
			log.Printf("Warning: Failed to collect ConsoleCLIDownloads: %v", err)
		} else {
			data.ConsoleCLIDownloads = consoleCLIDownloads
		}
	}

	// Collect OAuth resources
	if oauthClient != nil {
		oauthClients, err := collectOAuthClients(ctx, oauthClient)
		if err != nil {
			log.Printf("Warning: Failed to collect OAuthClients: %v", err)
		} else {
			data.OAuthClients = oauthClients
		}
	}

	// Collect IngressController resources
	if operatorClient != nil {
		ingressControllers, err := collectIngressControllers(ctx, operatorClient)
		if err != nil {
			log.Printf("Warning: Failed to collect IngressControllers: %v", err)
		} else {
			data.IngressControllers = ingressControllers
		}
	}

	return data, nil
}

// Helper functions for building configuration
func buildConfig(kubeconfigPath string, contextName ...string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeconfigPath
	configOverrides := &clientcmd.ConfigOverrides{}

	if len(contextName) > 0 && contextName[0] != "" {
		configOverrides.CurrentContext = contextName[0]
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return kubeConfig.ClientConfig()
}

// Collection functions for each OpenShift resource type
func collectClusterInfo(ctx context.Context, configClient *configclientset.Clientset, k8sClient *kubernetes.Clientset) (ClusterInfo, error) {
	clusterInfo := ClusterInfo{}

	// Get version and update channel
	versions, err := configClient.ConfigV1().ClusterVersions().List(ctx, metav1.ListOptions{})
	if err == nil && len(versions.Items) > 0 {
		version := versions.Items[0]
		clusterInfo.Name = version.ObjectMeta.Name
		clusterInfo.ID = string(version.Spec.ClusterID)
		clusterInfo.Version = version.Status.Desired.Version

		if len(version.Spec.Channel) > 0 {
			clusterInfo.Channel = version.Spec.Channel
		}

		// Add feature gates if any
		for _, gate := range version.Status.Capabilities.EnabledCapabilities {
			clusterInfo.FeatureGates = append(clusterInfo.FeatureGates, string(gate))
		}
	}

	// Extract provider information from infrastructure
	infras, err := configClient.ConfigV1().Infrastructures().List(ctx, metav1.ListOptions{})
	if err == nil && len(infras.Items) > 0 {
		infra := infras.Items[0]
		clusterInfo.Platform = string(infra.Status.PlatformStatus.Type)
		clusterInfo.Provider = string(infra.Status.PlatformStatus.Type)

		// Get regional information if available
		switch infra.Status.PlatformStatus.Type {
		case configv1.AWSPlatformType:
			if infra.Status.PlatformStatus.AWS != nil {
				clusterInfo.Region = infra.Status.PlatformStatus.AWS.Region
			}
		case configv1.AzurePlatformType:
			if infra.Status.PlatformStatus.Azure != nil {
				clusterInfo.Region = string(infra.Status.PlatformStatus.Azure.CloudName)
			}
		case configv1.GCPPlatformType:
			if infra.Status.PlatformStatus.GCP != nil {
				clusterInfo.Region = infra.Status.PlatformStatus.GCP.Region
			}
		}

		clusterInfo.ControlPlaneTopology = string(infra.Status.ControlPlaneTopology)
		clusterInfo.InfraTopology = string(infra.Status.InfrastructureTopology)

		if infra.Status.APIServerURL != "" {
			clusterInfo.APIServerURL = infra.Status.APIServerURL
		}
	}

	// Get base domain from DNS
	dnses, err := configClient.ConfigV1().DNSes().List(ctx, metav1.ListOptions{})
	if err == nil && len(dnses.Items) > 0 {
		dns := dnses.Items[0]
		clusterInfo.BaseDomain = dns.Spec.BaseDomain
	}

	return clusterInfo, nil
}

func collectRoutes(ctx context.Context, client *routeclientset.Clientset) ([]routev1.Route, error) {
	var allRoutes []routev1.Route

	routes, err := client.RouteV1().Routes("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allRoutes = append(allRoutes, routes.Items...)
	return allRoutes, nil
}

func collectBuildConfigs(ctx context.Context, client *buildclientset.Clientset) ([]buildv1.BuildConfig, error) {
	var allBuildConfigs []buildv1.BuildConfig

	buildConfigs, err := client.BuildV1().BuildConfigs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allBuildConfigs = append(allBuildConfigs, buildConfigs.Items...)
	return allBuildConfigs, nil
}

func collectImageStreams(ctx context.Context, client *imageclientset.Clientset) ([]imagev1.ImageStream, error) {
	var allImageStreams []imagev1.ImageStream

	imageStreams, err := client.ImageV1().ImageStreams("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allImageStreams = append(allImageStreams, imageStreams.Items...)
	return allImageStreams, nil
}

func collectSecurityContextConstraints(ctx context.Context, client *securityclientset.Clientset) ([]securityv1.SecurityContextConstraints, error) {
	var allSCCs []securityv1.SecurityContextConstraints

	sccs, err := client.SecurityV1().SecurityContextConstraints().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allSCCs = append(allSCCs, sccs.Items...)
	return allSCCs, nil
}

func collectProjects(ctx context.Context, client *projectclientset.Clientset) ([]projectv1.Project, error) {
	var allProjects []projectv1.Project

	projects, err := client.ProjectV1().Projects().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allProjects = append(allProjects, projects.Items...)
	return allProjects, nil
}

func collectClusterOperators(ctx context.Context, client *configclientset.Clientset) ([]configv1.ClusterOperator, error) {
	var allClusterOperators []configv1.ClusterOperator

	clusterOperators, err := client.ConfigV1().ClusterOperators().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allClusterOperators = append(allClusterOperators, clusterOperators.Items...)
	return allClusterOperators, nil
}

func collectClusterVersions(ctx context.Context, client *configclientset.Clientset) ([]configv1.ClusterVersion, error) {
	var allClusterVersions []configv1.ClusterVersion

	clusterVersions, err := client.ConfigV1().ClusterVersions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allClusterVersions = append(allClusterVersions, clusterVersions.Items...)
	return allClusterVersions, nil
}

func collectInfrastructures(ctx context.Context, client *configclientset.Clientset) ([]configv1.Infrastructure, error) {
	var allInfrastructures []configv1.Infrastructure

	infrastructures, err := client.ConfigV1().Infrastructures().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allInfrastructures = append(allInfrastructures, infrastructures.Items...)
	return allInfrastructures, nil
}

func collectNetworks(ctx context.Context, client *configclientset.Clientset) ([]configv1.Network, error) {
	var allNetworks []configv1.Network

	networks, err := client.ConfigV1().Networks().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allNetworks = append(allNetworks, networks.Items...)
	return allNetworks, nil
}

func collectMachineConfigs(ctx context.Context, client *machineconfigclientset.Clientset) ([]machineconfigv1.MachineConfig, error) {
	var allMachineConfigs []machineconfigv1.MachineConfig

	machineConfigs, err := client.MachineconfigurationV1().MachineConfigs().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allMachineConfigs = append(allMachineConfigs, machineConfigs.Items...)
	return allMachineConfigs, nil
}

func collectMachineConfigPools(ctx context.Context, client *machineconfigclientset.Clientset) ([]machineconfigv1.MachineConfigPool, error) {
	var allMachineConfigPools []machineconfigv1.MachineConfigPool

	machineConfigPools, err := client.MachineconfigurationV1().MachineConfigPools().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allMachineConfigPools = append(allMachineConfigPools, machineConfigPools.Items...)
	return allMachineConfigPools, nil
}

func collectOperatorGroups(ctx context.Context, client *operatorframeworkclientset.Clientset) ([]operatorsv1.OperatorGroup, error) {
	var allOperatorGroups []operatorsv1.OperatorGroup

	operatorGroups, err := client.OperatorsV1().OperatorGroups("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allOperatorGroups = append(allOperatorGroups, operatorGroups.Items...)
	return allOperatorGroups, nil
}

func collectSubscriptions(ctx context.Context, client *operatorframeworkclientset.Clientset) ([]operatorsv1alpha1.Subscription, error) {
	var allSubscriptions []operatorsv1alpha1.Subscription

	subscriptions, err := client.OperatorsV1alpha1().Subscriptions("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allSubscriptions = append(allSubscriptions, subscriptions.Items...)
	return allSubscriptions, nil
}

func collectInstallPlans(ctx context.Context, client *operatorframeworkclientset.Clientset) ([]operatorsv1alpha1.InstallPlan, error) {
	var allInstallPlans []operatorsv1alpha1.InstallPlan

	installPlans, err := client.OperatorsV1alpha1().InstallPlans("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allInstallPlans = append(allInstallPlans, installPlans.Items...)
	return allInstallPlans, nil
}

func collectConsoleLinks(ctx context.Context, client *consoleclientset.Clientset) ([]consolev1.ConsoleLink, error) {
	var allConsoleLinks []consolev1.ConsoleLink

	consoleLinks, err := client.ConsoleV1().ConsoleLinks().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allConsoleLinks = append(allConsoleLinks, consoleLinks.Items...)
	return allConsoleLinks, nil
}

func collectConsoleNotifications(ctx context.Context, client *consoleclientset.Clientset) ([]consolev1.ConsoleNotification, error) {
	var allConsoleNotifications []consolev1.ConsoleNotification

	consoleNotifications, err := client.ConsoleV1().ConsoleNotifications().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allConsoleNotifications = append(allConsoleNotifications, consoleNotifications.Items...)
	return allConsoleNotifications, nil
}

func collectConsoleCLIDownloads(ctx context.Context, client *consoleclientset.Clientset) ([]consolev1.ConsoleCLIDownload, error) {
	var allConsoleCLIDownloads []consolev1.ConsoleCLIDownload

	consoleCLIDownloads, err := client.ConsoleV1().ConsoleCLIDownloads().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allConsoleCLIDownloads = append(allConsoleCLIDownloads, consoleCLIDownloads.Items...)
	return allConsoleCLIDownloads, nil
}

func collectOAuthClients(ctx context.Context, client *oauthclientset.Clientset) ([]oauthv1.OAuthClient, error) {
	var allOAuthClients []oauthv1.OAuthClient

	oauthClients, err := client.OauthV1().OAuthClients().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allOAuthClients = append(allOAuthClients, oauthClients.Items...)
	return allOAuthClients, nil
}

func collectIngressControllers(ctx context.Context, client *operatorclientset.Clientset) ([]operatorv1.IngressController, error) {
	var allIngressControllers []operatorv1.IngressController

	ingressControllers, err := client.OperatorV1().IngressControllers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allIngressControllers = append(allIngressControllers, ingressControllers.Items...)
	return allIngressControllers, nil
}

func collectDNSes(ctx context.Context, client *configclientset.Clientset) ([]configv1.DNS, error) {
	var allDNSes []configv1.DNS

	dnses, err := client.ConfigV1().DNSes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	allDNSes = append(allDNSes, dnses.Items...)
	return allDNSes, nil
}
