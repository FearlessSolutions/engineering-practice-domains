# Databricks OpenTelemetry Azure Monitoring

This project demonstrates how to instrument Databricks notebooks with OpenTelemetry and integrate with Azure Application Insights for monitoring and observability.

## Overview

This project provides a comprehensive framework and examples for adding OpenTelemetry instrumentation to Databricks notebooks, with a focus on ETL pipeline monitoring and parent-child notebook workflows. It enables data engineers and data scientists to track spans, record metrics, and send telemetry data to Azure Application Insights, providing end-to-end visibility into their Databricks workflows.

By implementing this instrumentation, teams can monitor performance, detect issues, and troubleshoot problems more effectively, leading to more reliable and efficient data processing pipelines.

## Business Value

### Enhanced Observability
This project enables comprehensive monitoring of Databricks workflows, providing real-time visibility into ETL processes and notebook executions. This observability helps teams:

- **Reduce Mean Time to Resolution (MTTR)** by quickly identifying the root cause of failures
- **Improve Performance** by identifying bottlenecks in data processing pipelines
- **Increase Reliability** through proactive monitoring and alerting
- **Optimize Resource Usage** by tracking efficiency metrics across pipeline stages

### Use Cases
- **Data Engineering Teams**: Monitor complex ETL workflows and quickly troubleshoot failures
- **Data Science Teams**: Track notebook execution performance and dependencies
- **Operations Teams**: Set up alerts for critical pipeline failures and performance degradation
- **Business Stakeholders**: Access dashboards showing data processing volumes and success rates

## Visual Diagrams

### Architecture Overview
The following diagram illustrates the architecture of the Databricks OpenTelemetry Azure Monitoring solution:

[View Architecture Diagram](images/architecture_diagram.md)

### ETL Pipeline Visualization
This diagram shows the flow of data and tracing through the ETL pipeline:

[View ETL Pipeline Diagram](images/etl_pipeline_visualization.md)

### Parent-Child Notebook Workflow
This diagram illustrates how spans are created and correlated in parent-child notebook workflows:

[View Parent-Child Workflow Diagram](images/parent_child_workflow_diagram.md)

## Features

### Core Functionality
- Flexible OpenTelemetryHelper class that encapsulates OpenTelemetry functionality
- Integration with Azure Application Insights for monitoring and alerting
- Comprehensive tracing for ETL pipeline stages (extraction, transformation, loading)
- Custom metrics collection and visualization
- Span attributes for detailed monitoring and troubleshooting

### ETL Pipeline Instrumentation
- Function decorators for automatic tracing with minimal code changes
- Detailed performance metrics for each ETL stage
- Error tracking and correlation across pipeline components
- Efficiency and throughput measurements

### Parent-Child Notebook Workflow Monitoring
- Two approaches for notebook workflow instrumentation:
  - Monitoring without modifying child notebooks
  - Alternative approach for directly instrumenting child notebooks with passed context
- Simplified instrumentation with `run_notebook_with_tracing` method
- Function wrapping with `instrument_function` for custom tracing needs
- Automatic correlation of parent and child notebook executions

### Azure Monitor Integration
- Pre-configured KQL queries for data analysis
- Dashboard templates for visualization
- Alert configuration examples
- Performance trend analysis

## Documentation Guide

This project includes comprehensive documentation to help you get started and make the most of the OpenTelemetry instrumentation:

- **[Quick Start Guide](docs/quick_start.md)**: Get up and running in 5 minutes
- **[Setup Guide](docs/setup.md)**: Detailed installation and configuration instructions
- **[Usage Guide](docs/usage.md)**: How to use the library in your notebooks
- **[ETL Simulation Guide](docs/etl_simulation.md)**: Example ETL pipeline with OpenTelemetry
- **[Parent-Child Notebooks Guide](docs/parent_child_notebooks.md)**: Instrumenting notebook workflows
- **[Tracing Guide](docs/tracing.md)**: Details about span attributes and components
- **[Metrics Guide](docs/metrics.md)**: Information about metrics collection
- **[Azure Monitoring Guide](docs/azure_monitoring.md)**: Querying and visualizing telemetry data
- **[Glossary](docs/glossary.md)**: Definitions of technical terms used throughout the documentation

### For Business Users
Start with the **Business Value** section above and then explore the **ETL Simulation Guide** for practical examples. The **Glossary** can help with understanding technical terms.

### For Developers
Begin with the **Quick Start Guide** for rapid implementation, or the **Setup Guide** and **Usage Guide** for more detailed instructions. Then explore the specific guides relevant to your implementation needs.

## Getting Started

### Prerequisites

- Azure Databricks workspace
- Azure Application Insights instance
- Python 3.6+

### Quick Setup (Recommended)

The simplest approach for Databricks users:

1. Download the `otel_helper.py` file from the `databricks_opentelemetry_azure_monitoring` directory
2. In your Databricks workspace, use the "Create" > "File" option to upload the file
3. Import the helper class in your notebooks:

```python
from otel_helper import OpenTelemetryHelper
```

4. Configure your Azure Application Insights connection string in your Databricks environment

For alternative installation methods and more detailed setup instructions, see the [Quick Start Guide](docs/quick_start.md).

### Next Steps

- Follow the [Quick Start Guide](docs/quick_start.md) for step-by-step instructions and basic examples
- See the [Usage Guide](docs/usage.md) for detailed instructions on using the library in your notebooks
- Explore the [ETL Simulation](docs/etl_simulation.md) and [Parent-Child Notebooks](docs/parent_child_notebooks.md) guides for complete examples

## Examples

The `examples` directory contains sample notebooks that demonstrate how to use the OpenTelemetry instrumentation:

- `etl_simulation_before.py`: A basic ETL pipeline without instrumentation
- `etl_simulation_with_otel.py`: The same ETL pipeline with OpenTelemetry instrumentation
- `parent_notebook_with_otel.py`: A parent notebook that uses OpenTelemetry to instrument child notebook executions
  - `child_notebook_1.py`: A child notebook that performs data validation (no OpenTelemetry instrumentation)
  - `child_notebook_2.py`: A child notebook that performs data aggregation (no OpenTelemetry instrumentation)

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
