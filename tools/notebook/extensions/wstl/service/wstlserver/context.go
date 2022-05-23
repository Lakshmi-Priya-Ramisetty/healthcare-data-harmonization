// Copyright 2020 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wstlserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/transform"    /* copybara-comment: transform */
	"github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util/gcsutil" /* copybara-comment: gcsutil */
	"github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util/ioutil"  /* copybara-comment: ioutil */
	"google.golang.org/grpc/codes"                                                                 /* copybara-comment: codes */
	"google.golang.org/grpc/status"                                                                /* copybara-comment: status */

	dhpb "github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto"                 /* copybara-comment: data_harmonization_go_proto */
	hpb "github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto"                  /* copybara-comment: harmonization_go_proto */
	httppb "github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto"               /* copybara-comment: http_go_proto */
	lpb "github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto"                  /* copybara-comment: library_go_proto */
	wspb "github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/tools/notebook/extensions/wstl/proto" /* copybara-comment: wstlservice_go_proto */
)

// Context is a single transformation context that loosely corresponds to a Jupyter notebook cell.
type Context struct {
	// Whistle incremental transformer context.
	incrementalTransformer transform.Transformer

	// Google Cloud Storage client.
	storageClient gcsutil.StorageClient
}

// NewContext instantiates a transformation context.
// TODO () Modify the funtion to only take in non-nil StorageClient.
func NewContext(c gcsutil.StorageClient) (*Context, error) {
	if c == nil {
		return nil, fmt.Errorf("the argument to NewContext is nil")
	}
	return &Context{storageClient: c}, nil
}

// EvaluateIncrementalTransformation evaluates incremental updates to the whistle script and outputs
// the result.
func (c *Context) EvaluateIncrementalTransformation(request *wspb.IncrementalTransformRequest) ([]*wspb.TransformedRecords, error) {
	if request == nil {
		return nil, errors.New("empty request")
	}
	if request.GetWstl() == "" {
		return nil, errors.New("missing wstl script from session")
	}
	config := newHarmonizationConfig(request.GetWstl(), request.GetLibraryConfig(), request.GetCodeConfig(), request.GetUnitConfig())
	trans, err := transform.NewTransformer(context.Background(), config, transform.TransformationConfig{}, transform.GCSClient(c.storageClient))
	if err != nil {
		return nil, err
	}

	inputs := request.GetInput()
	if len(inputs) == 0 {
		inputs = append(inputs, &wspb.Location{Location: &wspb.Location_InlineJson{InlineJson: "{}"}})
	}
	return executeTransformation(trans, inputs), nil
}

// TODO (): move to wstlserver level.
func newHarmonizationConfig(wstl string, libraryConfigs []*wspb.Location, codeConfigs []*wspb.Location, unitConfig *wspb.Location) *dhpb.DataHarmonizationConfig {
	libConfig := []*lpb.LibraryConfig{}
	if len(libraryConfigs) > 0 {
		for _, libraryConfig := range libraryConfigs {
			switch l := libraryConfig.GetLocation().(type) {
			case *wspb.Location_LocalPath:
				names := ioutil.MustReadGlob(l.LocalPath, "library_config")
				for _, f := range names {
					if !strings.HasSuffix(f, ".wstl") {
						continue
					}
					tLibConfig := &lpb.LibraryConfig{}
					tLibConfig.UserLibraries = []*lpb.UserLibrary{
						&lpb.UserLibrary{
							Type: hpb.MappingType_MAPPING_LANGUAGE,
							Path: &httppb.Location{
								Location: &httppb.Location_LocalPath{
									LocalPath: f,
								},
							},
						},
					}
					libConfig = append(libConfig, tLibConfig)
				}
			case *wspb.Location_GcsLocation:
				tLibConfig := &lpb.LibraryConfig{}
				tLibConfig.UserLibraries = []*lpb.UserLibrary{
					&lpb.UserLibrary{
						Type: hpb.MappingType_MAPPING_LANGUAGE,
						Path: &httppb.Location{
							Location: &httppb.Location_GcsLocation{
								GcsLocation: l.GcsLocation,
							},
						},
					},
				}
				libConfig = append(libConfig, tLibConfig)
			}
		}
	}

	var codeConfigLocs []*httppb.Location
	if len(codeConfigs) > 0 {
		for _, cConfig := range codeConfigs {
			switch l := cConfig.GetLocation().(type) {
			case *wspb.Location_LocalPath:
				names := ioutil.MustReadGlob(l.LocalPath, "code_config")
				for _, f := range names {
					if !strings.HasSuffix(f, ".json") {
						continue
					}
					codeConfigLocs = append(codeConfigLocs, &httppb.Location{Location: &httppb.Location_LocalPath{LocalPath: f}})
				}
			case *wspb.Location_GcsLocation:
				codeConfigLocs = append(codeConfigLocs, &httppb.Location{Location: &httppb.Location_GcsLocation{GcsLocation: l.GcsLocation}})
			}
		}
	}

	uConfig := &hpb.UnitHarmonizationConfig{}
	if unitConfig != nil {
		switch l := unitConfig.GetLocation().(type) {
		case *wspb.Location_LocalPath:
			uConfig.UnitConversion = &httppb.Location{Location: &httppb.Location_LocalPath{LocalPath: l.LocalPath}}
		case *wspb.Location_GcsLocation:
			uConfig.UnitConversion = &httppb.Location{Location: &httppb.Location_GcsLocation{GcsLocation: l.GcsLocation}}
		}
	}

	return &dhpb.DataHarmonizationConfig{
		StructureMappingConfig: &hpb.StructureMappingConfig{
			Mapping: &hpb.StructureMappingConfig_MappingLanguageString{
				MappingLanguageString: wstl,
			},
		},
		LibraryConfig:           libConfig,
		HarmonizationConfig:     &hpb.CodeHarmonizationConfig{CodeLookup: codeConfigLocs},
		UnitHarmonizationConfig: uConfig,
	}
}

func executeTransformation(trans transform.Transformer, inputs []*wspb.Location) []*wspb.TransformedRecords {
	results := []*wspb.TransformedRecords{}
	for _, input := range inputs {
		tRecord := &wspb.TransformedRecords{}
		var source json.RawMessage
		switch l := input.GetLocation().(type) {
		case *wspb.Location_InlineJson:
			source = json.RawMessage(l.InlineJson)
		case *wspb.Location_GcsLocation:
			var err error
			if source, err = gcsutil.ReadFromGcs(context.Background(), l.GcsLocation); err != nil {
				tRecord.Record = &wspb.TransformedRecords_Error{
					Error: status.New(codes.InvalidArgument, err.Error()).Proto(),
				}
				results = append(results, tRecord)
				continue
			}
		default:
			tRecord.Record = &wspb.TransformedRecords_Error{
				Error: status.New(codes.InvalidArgument, "unsupported input type").Proto(),
			}
			results = append(results, tRecord)
			continue
		}

		output, err := trans.JSONtoJSON(source)

		if err != nil {
			tRecord.Record = &wspb.TransformedRecords_Error{
				Error: status.New(codes.InvalidArgument, err.Error()).Proto(),
			}
		} else {
			tRecord.Record = &wspb.TransformedRecords_Output{
				Output: string(output),
			}
		}
		results = append(results, tRecord)
	}
	return results
}
