// Copyright 2019 Google LLC
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

package transpiler

import (
	"github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language/parser" /* copybara-comment: parser */
)

func (t *transpiler) VisitProjectorDef(ctx *parser.ProjectorDefContext) interface{} {
	var aliases []string
	var requiredArgs []string
	for i := range ctx.AllArgAlias() {
		aliases = append(aliases, ctx.ArgAlias(i).(*parser.ArgAliasContext).TOKEN().GetText())
		if ctx.ArgAlias(i).(*parser.ArgAliasContext).REQUIRED() != nil {
			requiredArgs = append(requiredArgs, getTokenText(ctx.ArgAlias(i).(*parser.ArgAliasContext).TOKEN()))
		}
	}

	// Create a new environment for each projector.
	t.pushEnv(t.environment.newChild(getTokenText(ctx.TOKEN()), aliases, requiredArgs))

	ctx.Block().Accept(t)

	proj := t.environment.generateProjector()

	t.popEnv()

	if t.includeSourcePositions {
		proj.Meta = makeSourcePositionMeta(ctx, proj.Meta)
	}

	return proj
}
