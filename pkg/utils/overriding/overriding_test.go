package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	workspaces "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/json"
	yamlMachinery "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

func TestBasicToplevelOverriding(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Commands: []workspaces.Command{
			{
				Id: "command-with-type-changed",
				CommandUnion: workspaces.CommandUnion{
					Exec: &workspaces.ExecCommand{},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnion: workspaces.CommandUnion{
					Exec: &workspaces.ExecCommand{
						Env: []workspaces.EnvVar{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceOriginalValue",
							},
							{
								Name:  "envVarNotChanged",
								Value: "envVarNotChangedOriginalValue",
							},
						},
					},
				},
			},
			{
				Id: "command-not-changed",
				CommandUnion: workspaces.CommandUnion{
					Exec: &workspaces.ExecCommand{
						LabeledCommand: workspaces.LabeledCommand{
							Label: "commandNotChangedLabel",
						},
					},
				},
			},
		},
	}

	patch := workspaces.ParentOverrides{
		Commands: []workspaces.CommandParentOverride{
			{
				Id: "command-with-type-changed",
				CommandUnionParentOverride: workspaces.CommandUnionParentOverride{
					Apply: &workspaces.ApplyCommandParentOverride{
						Component: "mycomponent",
					},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnionParentOverride: workspaces.CommandUnionParentOverride{
					Exec: &workspaces.ExecCommandParentOverride{
						Env: []workspaces.EnvVarParentOverride{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceNewValue",
							},
							{
								Name:  "endVarToAdd",
								Value: "endVarToAddValue",
							},
						},
					},
				},
			},
		},
	}

	expected := &workspaces.DevWorkspaceTemplateSpecContent{
		Commands: []workspaces.Command{
			{
				Id: "command-with-type-changed",
				CommandUnion: workspaces.CommandUnion{
					Apply: &workspaces.ApplyCommand{
						Component: "mycomponent",
					},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnion: workspaces.CommandUnion{
					Exec: &workspaces.ExecCommand{
						Env: []workspaces.EnvVar{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceNewValue",
							},
							{
								Name:  "endVarToAdd",
								Value: "endVarToAddValue",
							},
							{
								Name:  "envVarNotChanged",
								Value: "envVarNotChangedOriginalValue",
							},
						},
					},
				},
			},
			{
				Id: "command-not-changed",
				CommandUnion: workspaces.CommandUnion{
					Exec: &workspaces.ExecCommand{
						LabeledCommand: workspaces.LabeledCommand{
							Label: "commandNotChangedLabel",
						},
					},
				},
			},
		},
	}

	result, err := OverrideDevWorkspaceTemplateSpec(&original, &patch)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected, result, "The two values should be the same.")
}

func overridingPatchTest(original, patch, expected []byte, expectedError string) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := OverrideDevWorkspaceTemplateSpecBytes(original, patch)
		if err != nil {
			compareErrorMessages(t, expectedError, err.Error(), "wrong error")
			return
		}
		if expectedError != "" {
			t.Error("Expected error but did not get one")
			return
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		resultYaml, err := yaml.JSONToYAML(resultJson)
		if err != nil {
			t.Error(err)
		}

		expectedJson, err := yamlMachinery.ToJSON(expected)
		if err != nil {
			t.Error(err)
		}
		expectedYaml, err := yaml.JSONToYAML(expectedJson)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(expectedYaml), string(resultYaml), "The two values should be the same.")
	}
}

func TestOverridingPatches(t *testing.T) {
	filepath.Walk("test-fixtures/patches", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "original.yaml" {
			if err != nil {
				t.Error(err)
				return nil
			}
			original, err := ioutil.ReadFile(path)
			if err != nil {
				t.Error(err)
				return nil
			}
			dirPath := filepath.Dir(path)
			patch, err := ioutil.ReadFile(filepath.Join(dirPath, "patch.yaml"))
			if err != nil {
				t.Error(err)
				return nil
			}
			result := []byte{}
			resultError := ""
			errorFile := filepath.Join(dirPath, "result-error.txt")
			if _, err = os.Stat(errorFile); err == nil {
				resultErrorBytes, err := ioutil.ReadFile(errorFile)
				if err != nil {
					t.Error(err)
					return nil
				}
				resultError = string(resultErrorBytes)
			} else {
				result, err = ioutil.ReadFile(filepath.Join(dirPath, "result.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, overridingPatchTest(original, patch, result, resultError))
		}
		return nil
	})
}

func mergingPatchTest(main, parent, expected []byte, expectedError string, plugins ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := MergeDevWorkspaceTemplateSpecBytes(main, parent, plugins...)
		if err != nil {
			compareErrorMessages(t, expectedError, err.Error(), "wrong error")
			return
		}
		if expectedError != "" {
			t.Error("Expected error but did not get one")
			return
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		resultYaml, err := yaml.JSONToYAML(resultJson)
		if err != nil {
			t.Error(err)
		}

		expectedJson, err := yamlMachinery.ToJSON(expected)
		if err != nil {
			t.Error(err)
		}
		expectedYaml, err := yaml.JSONToYAML(expectedJson)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(expectedYaml), string(resultYaml), "The two values should be the same.")
	}
}

func TestMerging(t *testing.T) {
	filepath.Walk("test-fixtures/merges", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "main.yaml" {
			if err != nil {
				t.Error(err)
				return nil
			}
			main, err := ioutil.ReadFile(path)
			if err != nil {
				t.Error(err)
				return nil
			}
			dirPath := filepath.Dir(path)
			parent := []byte{}
			parentFile := filepath.Join(dirPath, "parent.yaml")
			if _, err = os.Stat(parentFile); err == nil {
				parent, err = ioutil.ReadFile(parentFile)
				if err != nil {
					t.Error(err)
					return nil
				}
			}

			plugins := [][]byte{}
			pluginFile := filepath.Join(dirPath, "plugin.yaml")
			if _, err = os.Stat(pluginFile); err == nil {
				plugin, err := ioutil.ReadFile(filepath.Join(dirPath, "plugin.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
				plugins = append(plugins, plugin)
			}
			result := []byte{}
			resultError := ""
			errorFile := filepath.Join(dirPath, "result-error.txt")
			if _, err = os.Stat(errorFile); err == nil {
				resultErrorBytes, err := ioutil.ReadFile(errorFile)
				if err != nil {
					t.Error(err)
					return nil
				}
				resultError = string(resultErrorBytes)
			} else {
				result, err = ioutil.ReadFile(filepath.Join(dirPath, "result.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, mergingPatchTest(main, parent, result, resultError, plugins...))
		}
		return nil
	})
}

func TestPluginOverrides(t *testing.T) {
	originalFile := "test-fixtures/patches/override-just-plugin/original.yaml"
	patchFile := "test-fixtures/patches/override-just-plugin/patch.yaml"
	resultFile := "test-fixtures/patches/override-just-plugin/result.yaml"

	originalDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	patch := workspaces.PluginOverrides{}
	expectedDWT := workspaces.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, originalFile, &originalDWT)
	readFileToStruct(t, patchFile, &patch)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := OverrideDevWorkspaceTemplateSpec(&originalDWT, patch)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}

func TestMergingOnlyPlugins(t *testing.T) {
	baseFile := "test-fixtures/merges/no-parent/main.yaml"
	pluginFile := "test-fixtures/merges/no-parent/plugin.yaml"
	resultFile := "test-fixtures/merges/no-parent/result.yaml"

	baseDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	pluginDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	expectedDWT := workspaces.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, baseFile, &baseDWT)
	readFileToStruct(t, pluginFile, &pluginDWT)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := MergeDevWorkspaceTemplateSpec(&baseDWT, nil, &pluginDWT)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}

func TestMergingOnlyParent(t *testing.T) {
	// Reuse only plugin case since it's compatible
	baseFile := "test-fixtures/merges/no-parent/main.yaml"
	parentFile := "test-fixtures/merges/no-parent/plugin.yaml"
	resultFile := "test-fixtures/merges/no-parent/result.yaml"

	baseDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	parentDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	expectedDWT := workspaces.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, baseFile, &baseDWT)
	readFileToStruct(t, parentFile, &parentDWT)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := MergeDevWorkspaceTemplateSpec(&baseDWT, &parentDWT)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}

func readFileToStruct(t *testing.T, path string, into interface{}) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file from %s: %s", path, err.Error())
	}
	err = yaml.Unmarshal(bytes, into)
	if err != nil {
		t.Fatalf("Failed to unmarshal file into struct: %s", err.Error())
	}
}

// Since order of error message lines is not deterministic, it's necessary to compare
// in a weaker way than asserting string equality.
func compareErrorMessages(t *testing.T, expected, actual string, failReason string) {
	if expected == "" {
		t.Error("Received error but did not expect one: " + actual)
		return
	}
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	assert.ElementsMatch(t, expectedLines, actualLines, failReason)
}
