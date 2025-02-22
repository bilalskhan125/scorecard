// Copyright 2023 OpenSSF Scorecard Authors
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

// nolint:stylecheck
package securityPolicyPresent

import (
	"embed"
	"fmt"

	"github.com/ossf/scorecard/v4/checker"
	"github.com/ossf/scorecard/v4/finding"
	"github.com/ossf/scorecard/v4/probes/internal/utils"
)

//go:embed *.yml
var fs embed.FS

const Probe = "securityPolicyPresent"

func Run(raw *checker.RawResults) ([]finding.Finding, string, error) {
	if raw == nil {
		return nil, "", fmt.Errorf("%w: raw", utils.ErrNil)
	}
	var files []checker.File
	for i := range raw.SecurityPolicyResults.PolicyFiles {
		files = append(files, raw.SecurityPolicyResults.PolicyFiles[i].File)
	}

	var findings []finding.Finding
	for i := range files {
		file := &files[i]
		f, err := finding.NewWith(fs, Probe, "security policy file detected",
			file.Location(), finding.OutcomePositive)
		if err != nil {
			return nil, Probe, fmt.Errorf("create finding: %w", err)
		}
		f = f.WithRemediationMetadata(raw.Metadata.Metadata)
		findings = append(findings, *f)
	}

	// No file found.
	if len(findings) == 0 {
		f, err := finding.NewWith(fs, Probe, "no security policy file detected",
			nil, finding.OutcomeNegative)
		if err != nil {
			return nil, Probe, fmt.Errorf("create finding: %w", err)
		}
		f = f.WithRemediationMetadata(raw.Metadata.Metadata)
		findings = append(findings, *f)
	}

	return findings, Probe, nil
}
