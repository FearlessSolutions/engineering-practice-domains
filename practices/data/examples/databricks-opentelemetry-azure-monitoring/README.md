# Databricks OpenTelemetry Azure Monitoring

This project demonstrates how to instrument Databricks notebooks with OpenTelemetry and integrate with Azure Application Insights for monitoring and observability.

## Overview

This directory provides a framework and examples for adding OpenTelemetry instrumentation to Databricks notebooks, with a focus on ETL pipeline monitoring. It shows how to track spans, record metrics, and send telemetry data to Azure Application Insights.

## Features

- Flexible OpenTelemetryHelper class that encapsulates OpenTelemetry functionality
- Integration with Azure Application Insights for monitoring and alerting
- Comprehensive tracing for ETL pipeline stages (extraction, transformation, loading)
- Parent-child notebook workflow monitoring with two approaches:
  - Monitoring without modifying child notebooks
  - Alternative approach for directly instrumenting child notebooks with passed context
- Custom metrics collection and visualization
- Span attributes for detailed monitoring and troubleshooting
- Function decorators for automatic tracing with minimal code changes
- Comprehensive documentation for setup, usage, and monitoring
- Multiple installation options for different use cases
- Ready-to-use examples of instrumented ETL pipelines and notebook workflows

## Getting Started

### Prerequisites

- Azure Databricks workspace
- Azure Application Insights instance
- Python 3.6+

### Installation

There are three ways to install and use this package:

#### Option 1: Direct File Upload (Recommended)

The simplest approach for Databricks users:

1. Download the `otel_helper.py` file from the `databricks_opentelemetry_azure_monitoring` directory
2. In your Databricks workspace, use the "Create" > "File" option to upload the file
3. Import the helper class in your notebooks:

```python
from otel_helper import OpenTelemetryHelper
```

4. Configure your Azure Application Insights connection string in your Databricks environment

#### Option 2: Install as a Python package

For more advanced use cases:

1. Navigate to this directory within the parent repository
2. Install the package in development mode:

```bash
pip install -e databricks-opentelemetry-azure-monitoring
```

3. Import the helper class in your code:

```python
from databricks_opentelemetry_azure_monitoring import OpenTelemetryHelper
```

4. Configure your Azure Application Insights connection string in your Databricks environment

#### Option 3: Use the files with manual dependency installation

If you need more control over dependencies:

1. Navigate to this directory within the parent repository
2. Install the required dependencies:

```bash
pip install -r databricks-opentelemetry-azure-monitoring/requirements.txt
```

3. Copy the helper module to your project or upload it to your Databricks workspace
4. Configure your Azure Application Insights connection string in your Databricks environment

### Usage

See the [usage documentation](docs/usage.md) for detailed instructions on how to use this library in your Databricks notebooks.

## Examples

The `databricks-opentelemetry-azure-monitoring/examples` directory contains sample notebooks that demonstrate how to use the OpenTelemetry instrumentation:

- `etl_simulation_before.py`: A basic ETL pipeline without instrumentation
- `etl_simulation_with_otel.py`: The same ETL pipeline with OpenTelemetry instrumentation
- `parent_notebook_with_otel.py`: A parent notebook that uses OpenTelemetry to instrument child notebook executions
  - `child_notebook_1.py`: A child notebook that performs data validation (no OpenTelemetry instrumentation)
  - `child_notebook_2.py`: A child notebook that performs data aggregation (no OpenTelemetry instrumentation)

## Documentation

- [Setup Guide](docs/setup.md): Detailed setup instructions
- [Usage Guide](docs/usage.md): How to use the library in your notebooks
- [ETL Simulation Guide](docs/etl_simulation.md): Information about the simulated ETL with and without OpenTelemetry
- [Parent-Child Notebooks Guide](docs/parent_child_notebooks.md): How to instrument parent notebooks that call child notebooks
- [Tracing Guide](docs/tracing.md): Details about tracing attributes, components, and the OpenTelemetry Span Summary
- [Metrics Guide](docs/metrics.md): Information about metric values and instrumentation
- [Azure Monitoring Guide](docs/azure_monitoring.md): Querying data in Azure Application Insights, viewing metrics, and setting up visualizations and alerts

## Development Notes

### AI-Assisted Development

This project began with manual experimentation and initial development. To accelerate progress, refine implementation, and enhance documentation clarity, AI-assisted development tools were subsequently utilized:

- [ChatGPT](https://chat.openai.com/)
- [Anthropic Claude 3.7 Sonnet](https://www.anthropic.com/claude-3) via the [VSCode Cline plugin](https://marketplace.visualstudio.com/items?itemName=cline.codeline)

All AI-generated code and documentation were carefully reviewed, tested, and adjusted manually to ensure accuracy, reliability, and alignment with the project's standards and best practices.

This hybrid approach has also supported my personal initiative to enhance skills in leveraging AI tools effectively within software development workflows.

Contributors using similar AI tools are encouraged to transparently disclose their use in pull requests, following the example below:

```markdown
### AI Assistance Disclosure
This contribution utilized AI tools (e.g., ChatGPT, Claude 3.7 Sonnet via VSCode Cline). All outputs were manually reviewed and tested to ensure adherence to project standards.
```

## License

This project is part of a larger repository and is licensed under the terms of the parent repository's license.
