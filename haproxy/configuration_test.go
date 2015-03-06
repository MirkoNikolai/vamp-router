package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const (
	TEMPLATE_FILE         = "../configuration/templates/haproxy_config.template"
	CONFIG_FILE           = "/tmp/haproxy_test.cfg"
	PREFILLED_CONFIG_FILE = "../test/haproxy_test.cfg"
	CFG_JSON              = "../test/test_config1.json"
	CFG_WRONG_JSON        = "../test/test_wrong_config1.json"
	BACKEND_JSON          = "../test/test_backend1.json"
	JSON_FILE             = "/tmp/vamp_lb_test.json"
	PID_FILE              = "/tmp/vamp_lb_test.pid"
)

var (
	haConfig = Config{TemplateFile: TEMPLATE_FILE, ConfigFile: CONFIG_FILE, JsonFile: JSON_FILE, PidFile: PID_FILE}
)

func TestConfiguration_GetConfigFromDisk(t *testing.T) {

	if haConfig.GetConfigFromDisk(CFG_JSON) != nil {
		t.Errorf("Failed to load configuration from disk")
	}

	// wait for https://github.com/magneticio/vamp/issues/119
	// if haConfig.GetConfigFromDisk(CFG_WRONG_JSON) == nil {
	// 	t.Errorf("Expected an error when loading malformend JSON")
	// }

	if haConfig.GetConfigFromDisk("/this_is_really_something_wrong") == nil {
		t.Errorf("Expected an error when loading non existent path")
	}

}

func TestConfiguration_SetWeight(t *testing.T) {
	err := haConfig.SetWeight("test_be_1", "test_be_1_a", 20)
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

// Frontends

func TestConfiguration_FrontendExists(t *testing.T) {

	if haConfig.FrontendExists("non_existent_frontent") {
		t.Errorf("Should return false on non existent frontend")
	}

	if !haConfig.FrontendExists("test_fe_1") {
		t.Errorf("Should return true")
	}
}

func TestConfiguration_GetFrontends(t *testing.T) {
	result := haConfig.GetFrontends()
	if result[0].Name != "test_fe_1" {
		t.Errorf("Failed to get frontends array")
	}
}

func TestConfiguration_GetFrontend(t *testing.T) {
	if _, err := haConfig.GetFrontend("test_fe_1"); err != nil {
		t.Errorf("Failed to get frontend")
	}
}

func TestConfiguration_AddFrontend(t *testing.T) {

	fe := Frontend{Name: "my_test_frontend", Mode: "http", DefaultBackend: "test_be_1"}
	err := haConfig.AddFrontend(&fe)
	if err != nil {
		t.Errorf("Failed to add frontend")
	}
	if haConfig.Frontends[3].Name != "my_test_frontend" {
		t.Errorf("Failed to add frontend")
	}
}

func TestConfiguration_DeleteFrontend(t *testing.T) {

	if err := haConfig.DeleteFrontend("test_fe_2"); err != nil {
		t.Errorf("Failed to remove frontend")
	}

	if err := haConfig.DeleteFrontend("non_existing_frontend"); err == nil {
		t.Errorf("Frontend should not be removed")
	}
}

func TestConfiguration_GetFilters(t *testing.T) {

	filters := haConfig.GetFilters("test_fe_1")
	if filters[0].Name != "uses_internetexplorer" {
		t.Errorf("Could not retrieve Filter")
	}
}

func TestConfiguration_AddFilter(t *testing.T) {

	filter := Filter{Name: "uses_firefox", Condition: "hdr_sub(user-agent) Mozilla", Destination: "test_be_1_b"}
	err := haConfig.AddFilter("test_fe_1", &filter)
	if err != nil {
		t.Errorf("Could not add Filter")
	}
	if haConfig.Frontends[0].Filters[1].Name != "uses_firefox" {
		t.Errorf("Could not add Filter")
	}
}

func TestConfiguration_DeleteFilter(t *testing.T) {

	if err := haConfig.DeleteFilter("test_fe_1", "uses_firefox"); err != nil {
		t.Errorf("Could not add filter")
	}

	if err := haConfig.DeleteFilter("test_fe_1", "non_existent_filter"); err == nil {
		t.Errorf("Should return error on non existent filter")
	}
}

// Backends

func TestConfiguration_BackendUsed(t *testing.T) {

	if err := haConfig.BackendUsed("non_existent_backend"); err != nil {
		t.Errorf("Should not return error on non existent backend")
	}

	if err := haConfig.BackendUsed("test_be_1"); err == nil {
		t.Errorf("Should return error on backend still used by frontend")
	}

	if err := haConfig.BackendUsed("test_be_1_b"); err == nil {
		t.Errorf("Should return error on backend still used by filter")
	}
}

func TestConfiguration_GetBackends(t *testing.T) {
	result := haConfig.GetBackends()
	if result[0].Name != "test_be_1" {
		t.Errorf("Failed to get backends array")
	}
}

func TestConfiguration_GetBackend(t *testing.T) {

	if _, err := haConfig.GetBackend("test_be_1_a"); err != nil {
		t.Errorf("Failed to get backend")
	}

	if _, err := haConfig.GetBackend("non_existent_backend"); err == nil {
		t.Errorf("Should return error on non existent backend")
	}
}

func TestConfiguration_AddBackend(t *testing.T) {
	j, _ := ioutil.ReadFile(BACKEND_JSON)
	var backend *Backend
	_ = json.Unmarshal(j, &backend)

	if haConfig.AddBackend(backend) != nil {
		t.Errorf("Failed to add Backend")
	}

	if haConfig.AddBackend(backend) == nil {
		t.Errorf("Adding should fail when a backend already exists")
	}
}

func TestConfiguration_DeleteBackend(t *testing.T) {

	if err := haConfig.DeleteBackend("test_be_1"); err == nil {
		t.Errorf("Backend should not be removed because it is still in use")
	}

	if err := haConfig.DeleteBackend("deletable_backend"); err != nil {
		t.Errorf("Could not delete backend that should be deletable")
	}

	if err := haConfig.DeleteBackend("non_existing_backend"); err == nil {
		t.Errorf("Backend should not be removed")
	}
}

func TestConfiguration_BackendExists(t *testing.T) {

	if haConfig.BackendExists("non_existent_backend") {
		t.Errorf("Should return false on non existent backend")
	}

	if !haConfig.BackendExists("test_be_1") {
		t.Errorf("Should return true")
	}
}

// Server

func TestConfiguration_GetServers(t *testing.T) {

	if _, err := haConfig.GetServers("test_be_1"); err != nil {
		t.Errorf("Failed to get server array")
	}

	if _, err := haConfig.GetServers("non_existent_backend"); err == nil {
		t.Errorf("Should return false on non existent backend")
	}

}

func TestConfiguration_GetServer(t *testing.T) {

	if _, err := haConfig.GetServer("test_be_1", "test_be_1_a"); err != nil {
		t.Errorf("Failed to get server")
	}

	if _, err := haConfig.GetServer("non_existent_backend", "test_be_1"); err == nil {
		t.Errorf("Should return error on non existent backend")
	}
}

func TestConfiguration_AddServer(t *testing.T) {

	server := &ServerDetail{Name: "add_server", Host: "192.168.0.1", Port: 12345, Weight: 10}

	if err := haConfig.AddServer("test_be_1", server); err != nil {
		t.Errorf("Failed to add server")
	}

	if err := haConfig.AddServer("non_existent_backend", server); err == nil {
		t.Errorf("Should return false on non existent backend")
	}
}

func TestConfiguration_DeleteServer(t *testing.T) {

	if err := haConfig.DeleteServer("test_be_1", "deletable_server"); err != nil {
		t.Errorf("Failed to delete server")
	}

	if err := haConfig.DeleteServer("test_be_1", "non_existent_server"); err == nil {
		t.Errorf("Should return false on non existent server")
	}
}

// Namers

func TestConfiguration_ServiceName(t *testing.T) {
	if ServiceName("a", "b") == "a.b." {
		t.Errorf("Service name not well formed")
	}
}

func TestConfiguration_RouteName(t *testing.T) {
	if RouteName("a", "b") == "a.b." {
		t.Errorf("Route name not well formed")
	}
}

// Rendering & Persisting

func TestConfiguration_Render(t *testing.T) {
	err := haConfig.Render()
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestConfiguration_Persist(t *testing.T) {
	err := haConfig.Persist()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}

func TestConfiguration_RenderAndPersist(t *testing.T) {
	err := haConfig.RenderAndPersist()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}
