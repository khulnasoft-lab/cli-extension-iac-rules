// © 2023 Khulnasoft Limited All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package forms

import (
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/rs/zerolog"
	"github.com/khulnasoft-lab/cli-extension-iac-rules/internal/project"
	"github.com/khulnasoft/policy-engine/pkg/input"
)

type (
	RuleFields struct {
		RuleID      string
		Title       string
		Severity    string
		Description string
		Product     []string
		InputType   string
		SubForm     Form
	}

	RuleForm struct {
		Project *project.Project
		Fields  RuleFields
		Logger  *zerolog.Logger
	}
)

func (f *RuleForm) Run() error {
	if err := f.promptRuleID(); err != nil {
		return err
	}
	if err := f.promptTitle(); err != nil {
		return err
	}
	if err := f.promptSeverity(); err != nil {
		return err
	}
	if err := f.promptDescription(); err != nil {
		return err
	}
	if err := f.promptProduct(); err != nil {
		return err
	}
	if err := f.promptInputType(); err != nil {
		return err
	}
	if err := f.runSubForm(); err != nil {
		return err
	}
	return nil
}

func (f *RuleForm) promptRuleID() error {
	if f.Fields.RuleID != "" {
		return nil
	}

	var existingIDs []string
	metadata, err := f.Project.RuleMetadata()
	if err == nil {
		for id := range metadata {
			existingIDs = append(existingIDs, id)
		}
	}
	existingDirs := f.Project.ListRules()
	prompt := textinput.New("Rule ID:")
	prompt.Placeholder = "ACMECORP_001"
	prompt.CharLimit = ruleIDMaxLength
	prompt.Validate = ruleIDValidator(existingIDs, existingDirs)
	prompt.Template = verboseValidationTemplate
	ruleID, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	f.Fields.RuleID = ruleID
	return nil
}

func (f *RuleForm) promptSeverity() error {
	if f.Fields.Severity != "" {
		return nil
	}

	prompt := selection.New("Severity:", []string{
		"critical",
		"high",
		"medium",
		"low",
		"informational",
	})
	severity, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	f.Fields.Severity = severity
	return nil
}

func (f *RuleForm) promptTitle() error {
	if f.Fields.Title != "" {
		return nil
	}

	prompt := textinput.New("Title:")
	prompt.Placeholder = "S3 bucket is public"
	prompt.CharLimit = 256
	title, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	f.Fields.Title = title
	return nil
}

func (f *RuleForm) promptDescription() error {
	if f.Fields.Description != "" {
		return nil
	}

	prompt := textinput.New("Description:")
	prompt.Placeholder = "Public S3 buckets are open to unauthorized access"
	prompt.CharLimit = 1024
	description, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	f.Fields.Description = description
	return nil
}

func (f *RuleForm) promptProduct() error {
	if len(f.Fields.Product) > 0 {
		return nil
	}

	prompt := selection.New("Is this rule intended for Vulnmap IaC, Vulnmap Cloud, or both?", []string{
		"iac",
		"cloud",
		"both",
	})
	choice, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	switch choice {
	case "iac":
		f.Fields.Product = []string{choice}
	case "cloud":
		f.Fields.Product = []string{choice}
		f.Fields.InputType = input.CloudScan.Name
	case "both":
		f.Fields.Product = []string{"iac", "cloud"}
	}
	return nil
}

func (f *RuleForm) promptInputType() error {
	if f.Fields.InputType != "" {
		return nil
	}

	choices := allInputTypes()
	if len(f.Fields.Product) == 1 && f.Fields.Product[0] == "iac" {
		choices = iacInputTypes()
	}
	prompt := selection.New("Input type:", choices)
	inputType, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	f.Fields.InputType = inputType
	return nil
}

func (f *RuleForm) runSubForm() error {
	if f.Fields.SubForm != nil {
		return nil
	}

	prompt := confirmation.New("Does this rule need more than one resource type?", confirmation.No)
	choice, err := prompt.RunPrompt()
	if err != nil {
		return err
	}

	metadata := &project.RuleMetadata{
		ID:          f.Fields.RuleID,
		Severity:    f.Fields.Severity,
		Title:       f.Fields.Title,
		Description: f.Fields.Description,
		Product:     f.Fields.Product,
	}
	if choice {
		f.Fields.SubForm = &MultiResourceRuleForm{
			Project:   f.Project,
			RuleID:    f.Fields.RuleID,
			InputType: f.Fields.InputType,
			Metadata:  metadata,
			Logger:    f.Logger,
		}
	} else {
		f.Fields.SubForm = &SingleResourceRuleForm{
			Project:   f.Project,
			RuleID:    f.Fields.RuleID,
			InputType: f.Fields.InputType,
			Metadata:  metadata,
			Logger:    f.Logger,
		}
	}
	return f.Fields.SubForm.Run()
}
