# How to Add New Models to New-API

This guide explains how to add support for new AI models to the New-API system.

## Overview

Adding a new model involves three main steps:
1. Adding the model to the supported model list
2. Configuring pricing and billing ratios
3. Adding UI color mappings for the web interface

## Step 1: Add Model to the Model List

Navigate to the appropriate channel's constant file and add your model to the `ModelList` array.

### For OpenAI-compatible models:
File: `relay/channel/openai/constant.go`

```go
var ModelList = []string{
    // ... existing models
    "your-new-model-name",
    "your-new-model-variant",
}
```

### For other providers:
Find the corresponding constant file in `relay/channel/[provider]/constant.go` and add your models.

## Step 2: Configure Pricing and Billing

File: `setting/ratio_setting/model_ratio.go`

Add pricing configuration with detailed comments explaining the billing structure:

```go
// Your New Model - Pricing explanation
// Provider: [Provider Name]
// Official pricing: $X.XX per 1M tokens
// Billing structure: [Explain token types, input/output, etc.]
// Calculation: $X.XX / $0.002 = X ratio
"your-new-model": X.X, // $X.XX/1M tokens

// If there are variants with different pricing
"your-new-model-mini": X.X, // $X.XX/1M tokens
```

### Important Pricing Considerations:

1. **Base Rate**: New-API uses $0.002/1k tokens as the baseline (ratio = 1)
2. **Token Types**: Document whether the model uses:
   - Input/output tokens separately
   - Combined token counting
   - Special token types (audio, text, etc.)
3. **Provider Differences**: Note if pricing differs between direct API and cloud providers
4. **Billing Formula**:
   ```
   Ratio = (Model Price per 1M tokens) / (Base Price $2 per 1M tokens)
   ```

### Example: Audio Transcription Models

```go
// Azure OpenAI Transcription Models
// Provider: Azure OpenAI
// API Version: 2025-03-01-preview
// Official pricing: https://azure.microsoft.com/pricing/details/cognitive-services/openai-service/
//
// Token billing structure for transcription models:
// - input_tokens = audio_tokens + text_tokens (prompt/instructions)
// - output_tokens = 0 (transcription models don't generate, only transcribe)
// - Billing is based on audio duration converted to tokens
//
// GPT-4o-Transcribe: $6.00 per 1M audio tokens
// Calculation: $6.00 / $2.00 = 3.0 ratio
"gpt-4o-transcribe": 3,

// GPT-4o-Mini-Transcribe: $3.00 per 1M audio tokens
// Calculation: $3.00 / $2.00 = 1.5 ratio
"gpt-4o-mini-transcribe": 1.5,
```

## Step 3: Add UI Color Mappings

File: `web/src/helpers/render.jsx`

Add color mappings for the web interface:

```javascript
const modelColorMap = {
    // ... existing mappings
    'your-new-model': 'rgb(R,G,B)',
    'your-new-model-mini': 'rgb(R,G,B)',
};
```

### Color Guidelines:
- Use RGB format: `rgb(R,G,B)`
- Choose colors that are visually distinct from existing models
- Consider grouping related models with similar color families
- Ensure good contrast for readability

### Suggested Color Palette:
- Blue family: `rgb(173,216,230)` (light blue)
- Green family: `rgb(144,238,144)` (light green)
- Purple family: `rgb(221,160,221)` (plum)
- Orange family: `rgb(255,218,185)` (peach)

## Step 4: Build and Test

After making the changes:

1. **Build the application**:
   ```bash
   go build -ldflags "-s -w -X 'one-api/common.Version=$(cat VERSION)'" -o new-api
   ```

2. **Restart the service**:
   ```bash
   # Stop existing service (if running)
   pkill new-api

   # Start new service
   ./new-api &
   ```

3. **Test the new model**:
   - Check that the model appears in the web interface
   - Verify pricing calculations are correct
   - Test API calls with the new model

## Best Practices

### Documentation
- Always include comprehensive comments explaining pricing sources
- Document any special billing considerations
- Note API versions and provider-specific details
- Include calculation formulas for transparency

### Version Control
- Create feature branches for new model additions
- Write descriptive commit messages
- Include relevant links to official documentation

### Testing
- Verify model compatibility with existing API structure
- Test edge cases and error handling
- Confirm pricing calculations match official rates

## Example: Complete Model Addition

Here's a complete example of adding the GPT-4o transcription models:

### 1. Model List (`relay/channel/openai/constant.go`)
```go
var ModelList = []string{
    // ... existing models
    "gpt-4o-transcribe", "gpt-4o-mini-transcribe",
}
```

### 2. Pricing (`setting/ratio_setting/model_ratio.go`)
```go
// Azure OpenAI Transcription Models - Added 2025-09-19
// Provider: Azure OpenAI Service
// API Version: 2025-03-01-preview
// Documentation: https://learn.microsoft.com/azure/ai-foundry/openai/how-to/realtime-audio
//
// Token billing explanation:
// For transcription models, Azure OpenAI bills based on audio tokens:
// - input_tokens = audio_tokens + text_tokens (any prompt/instructions)
// - output_tokens = 0 (transcription models only convert speech to text, no generation)
// - Audio is processed as PCM16 24kHz format
// - Duration is converted to token equivalent for billing
//
// Official Azure OpenAI pricing (Pay-As-You-Go):
// - GPT-4o-Transcribe Audio Input: $6.00 per 1M tokens
// - GPT-4o-Mini-Transcribe Audio Input: $3.00 per 1M tokens
//
// Ratio calculation (baseline: $2.00 per 1M tokens):
"gpt-4o-transcribe": 3,     // $6.00 / $2.00 = 3.0 ratio
"gpt-4o-mini-transcribe": 1.5, // $3.00 / $2.00 = 1.5 ratio
```

### 3. UI Colors (`web/src/helpers/render.jsx`)
```javascript
'gpt-4o-transcribe': 'rgb(173,216,230)',        // Light blue
'gpt-4o-mini-transcribe': 'rgb(144,238,144)',   // Light green
```

This systematic approach ensures new models are properly integrated with accurate pricing and clear documentation.