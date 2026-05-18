# mse validate

Validate analysis JSON files before running analysis.

## Usage

```bash
mse validate [file] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `file` | Path to analysis JSON file (required) |

## What Gets Validated

### Required Fields

- Analysis must have `id`, `name`, and `market`
- Market must have `id`, `name`, and at least one category
- Categories must have `id`, `name`, and at least one capability
- Capabilities must have `id` and `name`
- Segments must have `id` and `name`
- Vendors must have `id` and `name`

### ID Uniqueness

- All capability IDs must be unique
- All segment IDs must be unique
- All vendor IDs must be unique

### Reference Integrity

- Weights must reference valid segment and capability IDs
- Scores must reference valid vendor and capability IDs
- `focusVendorId` (if set) must reference a valid vendor ID

### Weight Rules

- Weights for each segment must sum to approximately 1.0 (tolerance: 0.001)
- Individual weights must be between 0.0 and 1.0

### Score Rules

- Scores must be between 0.0 and 10.0
- Each vendor must have at least one score
- Warning for missing vendor/capability combinations

## Output

### Valid File

```bash
mse validate analysis.json
```

```
Validating: analysis.json

Analysis: CRM Market Analysis
Market: CRM Software
Vendors: 3
Segments: 3
Capabilities: 8
Scores: 24
Weights: 24

✓ Validation passed
```

### Invalid File

```bash
mse validate broken.json
```

```
Validating: broken.json

✗ Validation failed with 3 error(s)

Errors:
  - weights for segment "smb" sum to 0.900, must equal 1.0
  - score references unknown capability "invalid-cap"
  - focusVendorId "unknown" not found in vendors

Warnings:
  - missing score for vendor "competitor" on capability "mobile"
```

## Examples

Validate a single file:

```bash
mse validate crm-analysis.json
```

Validate multiple files:

```bash
for f in *.json; do mse validate "$f"; done
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Validation passed |
| 1 | Validation failed or file error |

## Integration with Analysis

The `mse analyze` command runs validation automatically before analysis. Use `--skip-validate` to bypass:

```bash
# Skip validation (for debugging only)
mse analyze broken.json --skip-validate
```
