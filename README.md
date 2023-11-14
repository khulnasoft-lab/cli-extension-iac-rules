<p align="center">
  <img src="https://vulnmap.khulnasoft.com/style/asset/logo/vulnmap-print.svg" />
</p>

# Vulnmap IaC Rules CLI Extension

## Overview

This repository contains an extension to the Vulnmap CLI that provides workflows to
author and manage custom rules for Vulnmap IaC.

## Usage

This repository produces a standalone binary for debugging purposes. This
extension is also built into the [Vulnmap CLI](https://github.com/khulnasoft-lab/vulnmap).
Outside of debugging and development, we advise to use the Vulnmap CLI instead of
the standalone binary.

## Workflows

- `vulnmap iac rules push`
  - Builds and pushes a custom rules project to the Vulnmap API
  - Can also be used to delete a custom rules project from the Vulnmap API
- `vulnmap iac rules init`
  - Prompts to initialize a custom rules project, relation, rule, or spec
- `vulnmap iac test`
  - Tests all rules in the project against their specs
  - Also used to generate the expected output for specs
