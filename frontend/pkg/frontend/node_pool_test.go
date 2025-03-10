package frontend

import (
	// This will invoke the init() function in each
	// API version package so it can register itself.
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Azure/ARO-HCP/internal/api"
	"github.com/Azure/ARO-HCP/internal/api/arm"
	_ "github.com/Azure/ARO-HCP/internal/api/v20240610preview"
	"github.com/Azure/ARO-HCP/internal/api/v20240610preview/generated"
	"github.com/Azure/ARO-HCP/internal/database"
	"github.com/Azure/ARO-HCP/internal/ocm"
)

const dummyTenantId = "dummy-tenant-id"
const dummySubscrtiptionId = "00000000-0000-0000-0000-000000000000"
const dummyResourceGroupId = "dummy_resource_group_name"
const dummyClusterName = "dev-test-cluster"
const dummyNodePoolName = "dev-nodepool"

const dummyClusterID = ("/subscriptions/" + dummySubscrtiptionId + "/resourcegroups/" + dummyResourceGroupId +
	"/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/" + dummyClusterName)
const dummyNodePoolID = dummyClusterID + "/nodePools/" + dummyNodePoolName

var dummyClusterHREF = ocm.GenerateClusterHREF(dummyClusterName)
var dummyNodePoolHREF = ocm.GenerateNodePoolHREF(dummyClusterHREF, dummyNodePoolName)

var dummyLocation = "Spain"
var dummyVMSize = "Big"
var dummyChannelGroup = "dummyChannelGroup"
var dummyVersionID = "dummy"

func TestCreateNodePool(t *testing.T) {
	clusterResouceID, _ := arm.ParseResourceID(dummyClusterID)
	clusterDoc := database.NewResourceDocument(clusterResouceID)
	clusterDoc.InternalID, _ = ocm.NewInternalID(dummyClusterHREF)

	nodePoolResouceID, _ := arm.ParseResourceID(dummyNodePoolID)
	nodePoolDoc := database.NewResourceDocument(nodePoolResouceID)
	nodePoolDoc.InternalID, _ = ocm.NewInternalID(dummyNodePoolHREF)

	requestBody := generated.HcpOpenShiftClusterNodePoolResource{
		Location:   &dummyLocation,
		Properties: &generated.NodePoolProperties{Spec: &generated.NodePoolSpec{Platform: &generated.NodePoolPlatformProfile{VMSize: &dummyVMSize}, Version: &generated.VersionProfile{ID: &dummyVersionID, ChannelGroup: &dummyChannelGroup}}},
	}
	tests := []struct {
		name               string
		urlPath            string
		subscription       *arm.Subscription
		systemData         *arm.SystemData
		subDoc             *database.SubscriptionDocument
		clusterDoc         *database.ResourceDocument
		nodePoolDoc        *database.ResourceDocument
		expectedStatusCode int
	}{
		{
			name:    "PUT Node Pool - Create a new Node Pool",
			urlPath: dummyNodePoolID + "?api-version=2024-06-10-preview",
			subDoc: &database.SubscriptionDocument{
				BaseDocument: database.BaseDocument{
					ID: dummySubscrtiptionId,
				},
				Subscription: &arm.Subscription{
					State:            arm.SubscriptionStateRegistered,
					RegistrationDate: api.Ptr(time.Now().String()),
					Properties:       nil,
				},
			},
			clusterDoc:         clusterDoc,
			nodePoolDoc:        nodePoolDoc,
			systemData:         &arm.SystemData{},
			expectedStatusCode: http.StatusCreated,
		},
	}
	mockCSClient := ocm.NewMockClusterServiceClient()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &Frontend{
				dbClient:             database.NewCache(),
				metrics:              NewPrometheusEmitter(prometheus.NewRegistry()),
				clusterServiceClient: &mockCSClient,
			}
			hcpCluster := api.NewDefaultHCPOpenShiftCluster()

			requestHeader := make(http.Header)
			requestHeader.Add(arm.HeaderNameHomeTenantID, dummyTenantId)

			hcpCluster.Name = dummyClusterName
			csCluster, _ := f.BuildCSCluster(clusterResouceID, requestHeader, hcpCluster, false)

			if test.subDoc != nil {
				err := f.dbClient.CreateSubscriptionDoc(context.TODO(), test.subDoc)
				if err != nil {
					t.Fatal(err)
				}
			}

			if test.clusterDoc != nil {
				err := f.dbClient.CreateResourceDoc(context.TODO(), test.clusterDoc)
				if err != nil {
					t.Fatal(err)
				}
				_, err = f.clusterServiceClient.PostCSCluster(context.TODO(), csCluster)
				if err != nil {
					t.Fatal(err)
				}
			}

			body, _ := json.Marshal(requestBody)

			ts := httptest.NewServer(f.routes())
			ts.Config.BaseContext = func(net.Listener) context.Context {
				ctx := context.Background()
				ctx = ContextWithLogger(ctx, testLogger) // defined in frontend_test.go
				ctx = ContextWithDBClient(ctx, f.dbClient)
				ctx = ContextWithSystemData(ctx, test.systemData)

				return ctx
			}

			req, err := http.NewRequest(http.MethodPut, ts.URL+test.urlPath, bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rs, err := ts.Client().Do(req)
			t.Log(rs)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if rs.StatusCode != test.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", test.expectedStatusCode, rs.StatusCode)
			}
		})
	}
}

// TODO: Fix the update logic for this test.

// func TestUpdateNodePool(t *testing.T) {
// 	clusterResouceID, _ := arm.ParseResourceID(dummyClusterID)
// 	clusterDoc := database.NewResourceDocument(clusterResouceID)
// 	clusterDoc.InternalID, _ = ocm.NewInternalID(dummyClusterHREF)

// 	nodePoolResouceID, _ := arm.ParseResourceID(dummyNodePoolID)
// 	nodePoolDoc := database.NewResourceDocument(nodePoolResouceID)
// 	nodePoolDoc.InternalID, _ = ocm.NewInternalID(dummyNodePoolHREF)

// 	var dummyReplicas int32 = 2
// 	requestBody := generated.HcpOpenShiftClusterNodePoolResource{
// 		Location: &dummyLocation,
// 		Properties: &generated.NodePoolProperties{
// 			Spec: &generated.NodePoolSpec{
// 				Replicas: &dummyReplicas,
// 				Version: &generated.VersionProfile{
// 					ID: &dummyVersionID, ChannelGroup: &dummyChannelGroup,
// 				},
// 			},
// 		},
// 	}

// 	tests := []struct {
// 		name               string
// 		urlPath            string
// 		subscription       *arm.Subscription
// 		systemData         *arm.SystemData
// 		subDoc             *database.SubscriptionDocument
// 		clusterDoc         *database.ResourceDocument
// 		nodePoolDoc        *database.ResourceDocument
// 		expectedStatusCode int
// 	}{
// 		{
// 			name:    "PUT Node Pool - Update an existing Node Pool",
// 			urlPath: dummyNodePoolID + "?api-version=2024-06-10-preview",
// 			subDoc: &database.SubscriptionDocument{
// 				BaseDocument: database.BaseDocument{
// 					ID: dummySubscrtiptionId,
// 				},
// 				Subscription: &arm.Subscription{
// 					State:            arm.SubscriptionStateRegistered,
// 					RegistrationDate: api.Ptr(time.Now().String()),
// 					Properties:       nil,
// 				},
// 			},
// 			clusterDoc:         clusterDoc,
// 			nodePoolDoc:        nodePoolDoc,
// 			systemData:         &arm.SystemData{},
// 			expectedStatusCode: http.StatusAccepted,
// 		},
// 	}
// 	mockCSClient := ocm.NewMockClusterServiceClient()

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			f := &Frontend{
// 				dbClient:             database.NewCache(),
// 				logger:               slog.New(slog.NewTextHandler(io.Discard, nil)),
// 				metrics:              NewPrometheusEmitter(),
// 				clusterServiceClient: &mockCSClient,
// 			}
// 			hcpCluster := api.NewDefaultHCPOpenShiftCluster()
// 			hcpCluster.Name = dummyCluster
// 			csCluster, _ := f.BuildCSCluster(clusterResouceID, dummyTenantId, hcpCluster, false)

// 			hcpNodePool := api.NewDefaultHCPOpenShiftClusterNodePool()
// 			hcpNodePool.Name = dummyNodePool
// 			csNodePool, _ := f.BuildCSNodePool(context.TODO(), hcpNodePool, false)

// 			if test.subDoc != nil {
// 				err := f.dbClient.CreateSubscriptionDoc(context.TODO(), test.subDoc)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}

// 			if test.clusterDoc != nil {
// 				err := f.dbClient.CreateResourceDoc(context.TODO(), test.clusterDoc)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				_, err = f.clusterServiceClient.PostCSCluster(context.TODO(), csCluster)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}

// 			if test.nodePoolDoc != nil {
// 				err := f.dbClient.CreateResourceDoc(context.TODO(), test.nodePoolDoc)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				_, err = f.clusterServiceClient.PostCSNodePool(context.TODO(), clusterDoc.InternalID, csNodePool)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}
// 			body, _ := json.Marshal(requestBody)

// 			ts := httptest.NewServer(f.routes())
// 			ts.Config.BaseContext = func(net.Listener) context.Context {
// 				ctx := context.Background()
// 				ctx = ContextWithLogger(ctx, f.logger)
// 				ctx = ContextWithDBClient(ctx, f.dbClient)
// 				ctx = ContextWithSystemData(ctx, test.systemData)

// 				return ctx
// 			}

// 			req, err := http.NewRequest(http.MethodPatch, ts.URL+test.urlPath, bytes.NewReader(body))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			rs, err := ts.Client().Do(req)
// 			t.Log(rs)
// 			if err != nil {
// 				t.Log(err)
// 				t.Fatal(err)
// 			}

// 			if rs.StatusCode != test.expectedStatusCode {
// 				t.Errorf("expected status code %d, got %d", test.expectedStatusCode, rs.StatusCode)
// 			}
// 		})
// 	}
// }
