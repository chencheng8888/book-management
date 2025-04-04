package configs

import (
	"os"
	"testing"
)

func TestAppConfig_LoadConfig(t *testing.T) {
	file, err := os.Create("config_test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			t.Error(err)
		}
		err = os.Remove("./config_test.yaml")
		if err != nil {
			t.Error(err)
		}
	}()

	_, err = file.Write([]byte(`
db:
  addr: "123"
cache:
  addr: "123"
`))
	if err != nil {
		t.Fatal(err)
	}
	var appConf AppConfig
	err = appConf.LoadConfig("./config_test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(appConf)
}
