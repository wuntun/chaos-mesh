package clientpool

import (
	"encoding/json"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

const stressChaos = `
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: burn-cpu
  namespace: chaos-testing
spec:
  mode: one
  selector:
    namespaces:
      - tidb-cluster-demo
  stressors:
    cpu:
      workers: 1
  duration: '30s'
  scheduler:
    cron: '@every 2m'
`

type KubeClient struct {
	DynamicClient   dynamic.Interface
	DiscoveryClient *discovery.DiscoveryClient
}

var MyClient KubeClient

func abc() {
	//kubeconfig := "/Users/xiang/.kube/config"
	kubeconfig := `
apiVersion: v1
clusters:
- cluster:
	server: https://192.168.64.66:8443
  name: test
- cluster:
	certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeE1ERXhNakEzTVRBek0xb1hEVE14TURFeE1EQTNNVEF6TTFvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTTZQClZNVGFPU01tQTNJakdxM2tiVTM0dkFyWEFZdS9haDUvV2FaNHhuMVN2N0Q5YTh1Nk1zQjNab21sSGhkN1lzODYKSmFQWmx6SUthRS9HcHFGSXUrelhERk5nWVl3YlVpc2RrRGpEUWtSYmovbFYzcHdERjJDNzV4QjEwNC8yREZjRQpwaWRsZjdPSEZkRnZrYTVrdTd6c01UZFFDUGRDUUdtblBudCtENGlSc3pPNzBLcEEyckVpTFBmaFFQTmFUM0FEClplZGN1RTJGT0FFS3lvaWNFSlFJWjJWLy9XTW9vYjhyWkhGN3I4OWsrMzhkWWtHUS9TRG1lcEd3L1QrUnBJZDYKWXpiWC9sQWtJSkZSaC93WEtVOU9ZRkNtYWRPOEpEV3FWcm4xMWNrVVBrWWRONVc3OFZlcnVKOWFBa2JhOVZuMApBU1ZDMDZ5VmVCV2dmSjl0NTBNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFMUXI2YzFpTXB6dEgybS9NWnFhSjI4WnYyNFEKaHZXYkpyQlpFck1uWlpXbzJYUndCSzJQV0RjeS8yQWdwUGc3REltTW10cWRDUmYrc21rU0tITHBONGI0dWI0cwpoVC94aW9BU1MrMFdTNW1EUUVYT3lFS3hYa3hUSVNmS3ZZQmpzRzZsRkVmcHJRUDdacDlvR0VKL3RCaDkwNVRxCmtIbDFZZWthYlhxMElKaFh1d1RneFlvYTkydjFnNDFBN1lwN0w3d1lMQkpLY05LWEMrT2wrU0pDejRIY1ZwS3MKZzQ4S0VScVVNdytCcUJvejcxK0hlSkJBazNnQzZDb3FLT0V0QzJQZ1R0ZFV2YTkrRDVaNUR3UDFXcmFSeXR1eQpQSWZQWlkzV0RKdjVFK2NuUzhScC9qYmZkRjBDSGdva1hvN3BVWUV3QlNmWUw2VmppQWwvelI0MEtlYz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
	server: https://127.0.0.1:60133
  name: kind-kind
contexts:
- context:
	cluster: test
	user: test
  name: test
- context:
	cluster: kind-kind
	user: kind-kind
  name: kind-kind
current-context: kind-kind
kind: Config
preferences: {}
users:
- name: cluster-admin
  user:
	password: uXFGweU9l35qcif
	username: admin
- name: kind-kind
  user:
	client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM4akNDQWRxZ0F3SUJBZ0lJQmx5bnFmdXdBTHN3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TVRBeE1USXdOekV3TXpOYUZ3MHlNakF4TVRJd056RXdNemxhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXdjcGZnWjJ6b1J3MEc3WWsKZ2dQKzNsZy9xdE5YL2NrSWVPQzc5ZitGL3NVbSs2K1lkcW5oU0NoVVFuWUd6aERoOE5rYytIcjZkV2hQNXdNTQpaVVpxbFZ6YXRuN1dSd0JRcDkydDZJOGVPTVBsN0g1Q2FRMmxWMjRhcEo1aFdzNVY3cUpTU2l6eE9lM1A1RFNqCm1aNjFXQ0NLR3ZWZnRxb3o3MUNkVllBQk9NNjkwcXlZakdQaUl2OXFrb3BvaGxCSEw1eG1oQ1BmeFptSkt1aEsKaVJlRS9RS1hkY09pOXhLcXBZdStCKzJtU3VoaW9wQlhtVngra25namJxL0s5T05zOGh4SVJJZEY3M2Z6YXdaVQpwOWxhMXYzRGxQci9PLzh2U2UzaWFZdm50S1h1c2ZMdlJQSkFCNVUvOFIzT2lDL0xtZTlRTTlRSmd1ZStvdFcxCnNDMVhyd0lEQVFBQm95Y3dKVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFMRnlNWVZ3TkpNalhnNC9VSmMzTUR0SVdtVGFXZzRVeDNWZQpuRWtwTm4wZG5DRkFGSit1OFY5UjdiVTFlWU1tWWtyOUVQcncySjN3aU5VOXdWY0ZHRUpNSXA3SmNKZEM5c0duCk5UZ2Q1K3VtSzNLdHhNOXozU01TYVBIeUlpcEp4eWxuNVc0MUM1MXZ6Yi9zemxNU20wNVNYaG04dGxpQSsyWngKVVlUZGt5dENhMlFmZmlUVnJrcmMxR2hvc2dwUDVuZmloTlUyd085ejhPaTRrV3F5VUZXMktJOCtrUkNmNWU0egpNbktrYURBclFYaWxvaFNRWnlmakF0OUJ2WHhCRGFGdnQ2MEhOT1N6SVRWNzlEcWRaR1N1dkRyamJlNjdudE02Ci9BOEh1NFQ0WHNOTnppYTFRc1ozTjRmT3RGYmpoUGlRdWRoeWk0eHJvcktsb3J0a2tsOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
	client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBd2NwZmdaMnpvUncwRzdZa2dnUCszbGcvcXROWC9ja0llT0M3OWYrRi9zVW0rNitZCmRxbmhTQ2hVUW5ZR3poRGg4TmtjK0hyNmRXaFA1d01NWlVacWxWemF0bjdXUndCUXA5MnQ2SThlT01QbDdINUMKYVEybFYyNGFwSjVoV3M1VjdxSlNTaXp4T2UzUDVEU2ptWjYxV0NDS0d2VmZ0cW96NzFDZFZZQUJPTTY5MHF5WQpqR1BpSXY5cWtvcG9obEJITDV4bWhDUGZ4Wm1KS3VoS2lSZUUvUUtYZGNPaTl4S3FwWXUrQisybVN1aGlvcEJYCm1WeCtrbmdqYnEvSzlPTnM4aHhJUklkRjczZnphd1pVcDlsYTF2M0RsUHIvTy84dlNlM2lhWXZudEtYdXNmTHYKUlBKQUI1VS84UjNPaUMvTG1lOVFNOVFKZ3VlK290VzFzQzFYcndJREFRQUJBb0lCQUVqUmRJWE43bHVScjNyaQpQR0dtZ3JTbDBIYXVKNWd6WEQyZnBNRlJITmFZMm9ja2VsUE1qZHlCV3ZnR1JaUlUvN0Z5dzlJUzA5NGVMamdPCkN6QmEvMTNVb0ZLRzRwbVhZcmRSTXpINTVVOUxQVEJhV1RZRWJLYW50dGM0dERoYzVDbGhVUzZTS0txdDA3cGEKbFViNlBnWTVZK3V4WEIvVllPS0NGankwZHNFVE54VkdVR09POHJpeUdoU1FsVzFZWGp3R3p5S25zbUZSQ1R3QgpOS2wxQ3RzL28xNWxxUjE5M0xJS1p2N1N5N29pN2ZyYTV6czVibzFGZ2R6NFIrS3REcG1OODc2MWVKUG10UXU3CmlzVnVGVVFKTTVudHAwK0VyYUs5WjA2Yjd4VGdHRGJyeHpJSGlTV05KWjBXUHNpUTRjNFR1eEZlR2F5SGdacGsKM0VlZ1RRRUNnWUVBNGF3bmYxTWFoR3Vlc3RqMVM0MDY4THl1YjZnOXhhU2VKNU9hckFFeGxveXRJOXRxVjlOOApubHlNbU1TenUvRk5pM3ppZGF6T1FxYWJ1bmlZTWk0dEM3VzlmVGpBNkdMUVNDcS90RnVxeUNYOEwwNEdWRkpwCmJCZllSVXYzT1ppaUpEeFJSMGZvZmJuR0VyazZRWGthdFAvQWNQVFA4N3JrU3lQU3pVZ2dlRzhDZ1lFQTI5VmcKOTVzaTFsN1IrcXBGKzRxdllnRXhGaUpmTGFRZCtNWEprWWVSanVGVDhyMWUrWHNMT1d4cHhhSEQ1bGt6OXhXZQpqd0lSUXZEKzZVeWt6UE9uSUhIUmU5aVljZ2orSmxtbjF2MGEwSFhaejBxS0h4bng4QncrTjllZ2VNcjFpYVN1CmdQeUlFclpkUlcvOEMxQitsQUQ2a2htWnlsMmREYzdvRDBLRU5NRUNnWUJJdnZ3RWVUK0ZERVFlRnY0TG1yMHoKT1Q5cDB1d0d1Q2diVGVPQUt1cFhRNFhVbHpoU2syUUtrSDdxQ0E2QU9TcnNHaGZPSXlSaUs5N3JYMUNBYkk0cwp4aXNOSUt4ZXZPdXpOOFNRV1RSV1RKaGNqMlJPN2puNWxENHRLRzNMYlQ1bk8rSmZmZmlkL3JLdytuQ2pCbXpyCmg3MzdLNCtWVzl1WHRUYVE0ZjFGbFFLQmdRQ0hWWk4rbTVrVTFBYjlCRHBWSXduWmtkWkFSQ1RJR2xNQlJmSlQKajF4QzArRTBmUFN0TGR5NUcwNzZoaDN0LzFpSWNsek11WDhhOFBaZGRmdTIyUUU0YmhtQzN0THEwVEoxTlppbwpOK1Y4RkRGazlnU1dKUWpXd3V4aXdISmdLc0tDWEVtNXlyMlNsNFpRS3lMRmJHYWdnd0cvVTlkang4SGFNRmlZCm5HQjdnUUtCZ0FQaFRkNVlac1VHWWNDWU50eHIzMWlXY3pPTjNXejVXdnBSQTAweWlNUzFnbW9RUGhSUitwTkwKd0QrYjE0b2RYR29EVzY2N1Q1UEd6c0FjcnJEeFpqd2RqUEdNekR3ZFd0NFBWdWNPSytCU3c1a0NDcHdWTzNBVApXMGtsTFZWczkrNWs1blRpSFIzSW1ZRll6dmNTbTF1VGlUYkFvQmVpeDMweDhjOVZVQVNLCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
`
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		log.Println(err)
		return
	}

	dynamicClient, _ := dynamic.NewForConfig(config)
	MyClient.DynamicClient = dynamicClient

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)
	MyClient.DiscoveryClient = discoveryClient

	obj := &unstructured.Unstructured{}

	// decode YAML into unstructured.Unstructured
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, _ := dec.Decode([]byte(stressChaos), nil, obj)

	// show GVK
	log.Println("GVK", gvk.String())

	// encode back to JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	_ = enc.Encode(obj)

	unStructGVR := getGVR(obj)
	log.Println("GVR", unStructGVR)

	dynamicHandle := dynamicClient.Resource(*unStructGVR).Namespace(obj.GetNamespace())

	// create
	Create(dynamicHandle, obj)

	time.Sleep(5 * time.Second)

	// delete
	Delete(dynamicHandle, obj)
}

func getGVR(unStruct *unstructured.Unstructured) *schema.GroupVersionResource {
	gvk := unStruct.GroupVersionKind()
	kind := getKind(gvk.Kind)
	return &schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: kind,
	}
}

func getKind(kind string) string {
	group, resourceList := GetResource(MyClient.DiscoveryClient)
	log.Println("group", group)
	log.Println("resource list", resourceList)
	for _, list := range resourceList {
		for _, resource := range list.APIResources {
			if resource.Kind == kind {
				return resource.Name
			}
		}
	}
	return ""
}

func GetResource(handler *discovery.DiscoveryClient) ([]*metav1.APIGroup, []*metav1.APIResourceList) {
	group, source, _ := handler.ServerGroupsAndResources()
	return group, source
}

func Create(handler dynamic.ResourceInterface, unstructured *unstructured.Unstructured) {
	ret, _ := handler.Create(unstructured, metav1.CreateOptions{})
	log.Println(ret)
}

func Delete(handler dynamic.ResourceInterface, unstructured *unstructured.Unstructured) {
	_ = handler.Delete(unstructured.GetName(), &metav1.DeleteOptions{})
}
