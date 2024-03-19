package extractors

import (
	"github.com/CheckmarxDev/containers-resolver/internal/logger"
	"github.com/CheckmarxDev/containers-resolver/internal/types"
	"testing"
)

func TestExtractImagesFromHelmFiles(t *testing.T) {
	l := logger.NewLogger(false)

	t.Run("ValidHelmFiles", func(t *testing.T) {
		helmCharts := []types.HelmChartInfo{
			{Directory: "../../test_files/imageExtraction/helm"},
		}

		images, err := ExtractImagesFromHelmFiles(l, helmCharts)
		if err != nil {
			t.Errorf("Error extracting images: %v", err)
		}

		expectedImages := map[string]types.ImageLocation{
			"checkmarx.jfrog.io/ast-docker/containers-worker:b201b1f": {Origin: types.HelmFileOrigin, Path: "containers/templates/containers-worker.yaml"},
			"checkmarx.jfrog.io/ast-docker/image-insights:f4b507b":    {Origin: types.HelmFileOrigin, Path: "containers/templates/image-insights.yaml"},
		}

		checkHelmResult(t, images, expectedImages)
	})

	t.Run("NoHelmFilesFound", func(t *testing.T) {
		helmCharts := []types.HelmChartInfo{}

		images, err := ExtractImagesFromHelmFiles(l, helmCharts)
		if err != nil {
			t.Errorf("Error extracting images: %v", err)
		}

		if len(images) != 0 {
			t.Errorf("Expected 0 images, but got %d", len(images))
		}
	})

	t.Run("OneValidOneInvalidHelmFiles", func(t *testing.T) {
		helmCharts := []types.HelmChartInfo{
			{Directory: "../../test_files/imageExtraction/helm/"},
			{Directory: "../../test_files/imageExtraction/helm2/"},
		}

		images, err := ExtractImagesFromHelmFiles(l, helmCharts)
		if err != nil {
			t.Errorf("Error extracting images: %v", err)
		}

		expectedImages := map[string]types.ImageLocation{
			"checkmarx.jfrog.io/ast-docker/containers-worker:b201b1f": {Origin: types.HelmFileOrigin, Path: "containers/templates/containers-worker.yaml"},
			"checkmarx.jfrog.io/ast-docker/image-insights:f4b507b":    {Origin: types.HelmFileOrigin, Path: "containers/templates/image-insights.yaml"},
		}

		checkHelmResult(t, images, expectedImages)
	})
}

func TestExtractImageInfo(t *testing.T) {
	t.Run("ValidYAMLString", func(t *testing.T) {
		yamlString := `---
# Source: containers/templates/image-insights.yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: release-name-containers-image-insights
  labels:
    helm.sh/chart: containers-0.0.133
    app.kubernetes.io/name: containers
    app.kubernetes.io/instance: release-name
    app.kubernetes.io/version: "0.0.133"
    app.kubernetes.io/managed-by: Helm
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get"]
  - apiGroups: [""]
    resources: ["pods/exec"]
    verbs: ["create"]
  - apiGroups: ["batch"]
    resources: ["jobs"]
    verbs: ["get", "create"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list"]
---
# Source: containers/templates/image-insights.yaml
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: release-name-containers-image-insights
  labels:
    helm.sh/chart: containers-0.0.133
    app.kubernetes.io/name: containers
    app.kubernetes.io/instance: release-name
    app.kubernetes.io/version: "0.0.133"
    app.kubernetes.io/managed-by: Helm
subjects:
  - kind: ServiceAccount
    name: release-name-containers-image-insights
    namespace: default
roleRef:
  kind: Role
  name: release-name-containers-image-insights
  apiGroup: rbac.authorization.k8s.io
---
# Source: containers/templates/containers-image-risks.yaml
apiVersion: ast.checkmarx.com/v1
kind: Microservice
metadata:
  name: release-name-containers-containers-image-risks
  labels:
    helm.sh/chart: containers-0.0.133
    app.kubernetes.io/name: containers
    app.kubernetes.io/instance: release-name
    app.kubernetes.io/version: "0.0.133"
    app.kubernetes.io/managed-by: Helm


spec:
  component: "containers"

  image:
    registry: 
    name: nginx
    pullPolicy: IfNotPresent # Overrides the image tag whose default is the chart appVersion.
    tag: latest
    imagePullSecrets: [ ]

  scale:
    static:
      count: 1
    auto:
      enabled: false
      minReplicas: 1
      maxReplicas: 20
      targetCPUUtilizationPercentage: 70
      targetMemoryUtilizationPercentage: 70

  envirmentVariablesUnencrypted:
    - key: BucketName
      value: "containers-image-risks-bucket"
    - key: REDISSHAREDS3BUCKET
      value: "containers-image-risks-bucket"
    - key: MaxReconnectAttempts
      value: "10"
    - key: ReconnectDelayInSeconds
      value: "3"
    - key: RabbitMqProtocol
      value: "AMQP"
    - key: RabbitMqSendTimeout
      value: "30000"
    - key: TopicExchangeName
      value: "containers.topic"
    - key: RoutingKey
      value: "containers.scan"
    - key: ActivateImageRisksQueueName
      value: "activate-image-risks"
    - key: FinishedImageRisksQueueName
      value: "finished-image-risks"
    - key: ImageCorrelationsUrl
      value: /image-correlations
    - key: VulnerabilitiesServiceUrl
      value: /vulnerabilities
    - key: ImageInsightsGrpcUrl
      value: http://image-insights:50051
    - key: CacheObjectTimeToLiveInMinutes
      value: "60"
    - key: GrpcPort
      value: 50051

  persistent:
    minio:
      enabled: true
      includeSchema: true
      envirmentVariablesMap:
        address:
          - "LocalServerUrl"
        region:
          - "LocalRegion"
        accessKey:
          - "LocalAccessKey"
        accessSecret:
          - "LocalSecretKey"
        tls:
          enabled:
            - "OBJECT_STORE_STORAGE_TLS_ENABLED"
          skipVerify:
            - "OBJECT_STORE_STORAGE_TLS_SKIP_VERIFY"
    redis:
      enabled: true
      envirmentVariablesMap:
        isCluster: "REDIS_IS_CLUSTER_MODE"
        address: "RedisAddresses"
        password: "RedisPassword"
        tls:
          enabled: "RedisSSLEnabled"
    postgres:
      enabled: true
      liquibase:
        enabled: false
        definitionsDirInImage: "/app/db"
      envirmentVariablesMap:
        read:
          connection_strings:
          - "DATABASE_READ_URL"
        readWrite:
          host: "DatabaseHost"
          port: "DatabasePort"
          db: "DatabaseName"
          username: "DatabaseUser"
          password: "DatabasePassword"
          connection_strings:
          - "DATABASE_WRITE_URL"

  messaging:
    rabbitMQ:
      enabled: true
      envirmentVariablesMap:
        tls:
          enabled: ""
          skipVerify: "RABBIT_TLS_SKIP_VERIFY"
        uri: "RabbitMqUrl"

  internalNetworking:
    addionalServiceName: "image-risks"
    ports:
      - port: 80
        name: "rest"
      - port: 50051
        name: "containers-image-risks"

  livenessProbe:
    enabled: true
    type: "httpGet"
    httpGet:
      path: "/health"
      port: 80
  readinessProbe:
    enabled: true
    type: "httpGet"
    httpGet:
      path: "/health"
      port: 80`
		images, err := extractImageInfo(yamlString)
		if err != nil {
			t.Errorf("Error extracting images: %v", err)
		}

		expectedImages := map[string]types.ImageLocation{
			"nginx:latest": {Origin: types.HelmFileOrigin, Path: "containers/templates/containers-image-risks.yaml"},
		}

		checkHelmResult(t, images, expectedImages)
	})

	t.Run("InvalidYAMLString", func(t *testing.T) {
		yamlString := `invalid yaml string`

		_, err := extractImageInfo(yamlString)
		if err == nil {
			t.Errorf("Expected error extracting images from invalid YAML string, but got none")
		}
	})
}

func checkHelmResult(t *testing.T, images []types.ImageModel, expectedImages map[string]types.ImageLocation) {
	for _, image := range images {
		// Check if the image name exists in the expected images map
		expectedLocation, ok := expectedImages[image.Name]
		if !ok {
			t.Errorf("Unexpected image found: %s", image.Name)
			continue
		}

		// Check if the file path matches the expected file path
		if len(image.ImageLocations) != 1 {
			t.Errorf("Expected image %s to have exactly one location, but got %d", image.Name, len(image.ImageLocations))
			continue
		}

		if image.ImageLocations[0].Path != expectedLocation.Path {
			t.Errorf("Expected image %s to have path %s, but got %s", image.Name, expectedLocation.Path, image.ImageLocations[0].Path)
		}

		if image.ImageLocations[0].Origin != expectedLocation.Origin {
			t.Errorf("Expected image %s to have origin %s, but got %s", image.Name, expectedLocation.Origin, image.ImageLocations[0].Origin)
		}

		// Remove the checked image from the expected images map
		delete(expectedImages, image.Name)
	}

	// Check if any expected images are left unchecked
	for imageName, expectedLocation := range expectedImages {
		t.Errorf("Expected image %s not found (Origin: %s, Path: %s)", imageName, expectedLocation.Origin, expectedLocation.Path)
	}
}
