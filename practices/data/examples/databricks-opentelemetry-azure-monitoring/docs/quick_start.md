# Quick Start Guide

This guide provides a streamlined approach to get up and running with the Databricks OpenTelemetry Azure Monitoring solution quickly.

## Table of Contents

- [Prerequisites](#prerequisites)
- [5-Minute Setup](#5-minute-setup)
- [Basic ETL Pipeline Instrumentation](#basic-etl-pipeline-instrumentation)
- [Parent-Child Notebook Instrumentation](#parent-child-notebook-instrumentation)
- [Viewing Results in Azure Application Insights](#viewing-results-in-azure-application-insights)
- [Next Steps](#next-steps)

## Prerequisites

Before you begin, ensure you have:

- An Azure Databricks workspace
- An Azure Application Insights instance
- Python 3.6+

## 5-Minute Setup

Follow these steps to quickly set up the OpenTelemetry instrumentation:

1. **Download the helper file**:
   - Download the `otel_helper.py` file from the `databricks_opentelemetry_azure_monitoring` directory

2. **Upload to Databricks workspace**:
   - In your Databricks workspace, click "Create" > "File"
   - Upload the `otel_helper.py` file

3. **Configure Azure Monitor connection string**:
   - In your Databricks cluster configuration, add an environment variable:
     - Key: `AZURE_MONITOR_CONNECTION_STRING`
     - Value: Your Application Insights connection string

4. **Install required packages**:
   - Add the following to a notebook cell and run it:
     ```python
     %pip install --upgrade azure-monitor-opentelemetry opentelemetry-sdk azure-core opentelemetry-api opentelemetry-sdk opentelemetry-instrumentation
     
     # Restart Python interpreter
     %restart_python
     dbutils.library.restartPython()
     ```

5. **Import the helper class**:
   ```python
   from otel_helper import OpenTelemetryHelper
   ```

That's it! You're now ready to instrument your Databricks notebooks.

## Basic ETL Pipeline Instrumentation

Here's a minimal example to instrument an ETL pipeline:

```python
import uuid
from otel_helper import OpenTelemetryHelper

# Generate a unique ID for this ETL run
etl_pipeline_id = str(uuid.uuid4())

# Initialize OpenTelemetry with a parent span
otel_helper = OpenTelemetryHelper(
    span_name="ETL_Pipeline",
    etl_pipeline_id=etl_pipeline_id,
    span_attributes={"data_source": "My Data Source", "etl_type": "Incremental Load"}
)

try:
    # Data Extraction
    otel_helper.start_tracing("DataExtraction")
    # ... your extraction code here ...
    records_extracted = 1000  # Replace with actual count
    otel_helper.set_span_attribute("DataExtraction", "records_extracted", records_extracted)
    otel_helper.end_tracing("DataExtraction")
    
    # Data Transformation
    otel_helper.start_tracing("DataTransformation")
    # ... your transformation code here ...
    records_transformed = 950  # Replace with actual count
    otel_helper.set_span_attribute("DataTransformation", "records_transformed", records_transformed)
    otel_helper.end_tracing("DataTransformation")
    
    # Data Loading
    @otel_helper.trace_function("DataLoading", {
        "records_written": "records_written"
    })
    def perform_data_loading():
        # ... your loading code here ...
        return {"records_written": 950}  # Replace with actual count
    
    loading_result = perform_data_loading()
    
finally:
    # Always end the parent span
    otel_helper.end_tracing("ETL_Pipeline")
```

## Parent-Child Notebook Instrumentation

Here's a minimal example to instrument a parent notebook that calls child notebooks:

### Parent Notebook

```python
import uuid
import json
from otel_helper import OpenTelemetryHelper

# Generate a unique ID for this workflow
workflow_id = str(uuid.uuid4())

# Initialize OpenTelemetry with a parent span
workflow_otel_helper = OpenTelemetryHelper(
    span_name="Notebook_Workflow",
    etl_pipeline_id=workflow_id,
    span_attributes={"workflow_type": "Data Processing"}
)

try:
    # Execute Child Notebook 1 with automatic tracing
    child1_result = workflow_otel_helper.run_notebook_with_tracing(
        notebook_path="./child_notebook_1",
        span_name="Child_Notebook_1",
        timeout_seconds=600,
        etl_pipeline_id=workflow_id
    )
    
    # Execute Child Notebook 2 with automatic tracing
    child2_result = workflow_otel_helper.run_notebook_with_tracing(
        notebook_path="./child_notebook_2",
        span_name="Child_Notebook_2",
        timeout_seconds=600,
        etl_pipeline_id=workflow_id
    )
    
    # Set overall workflow attributes
    workflow_otel_helper.set_span_attribute("Notebook_Workflow", "total_records", child1_result["total_records"])
    workflow_otel_helper.set_span_attribute("Notebook_Workflow", "total_errors", child1_result["validation_errors"])
    
finally:
    # Always end the parent span
    workflow_otel_helper.end_tracing("Notebook_Workflow")
```

### Child Notebook

```python
import json

# ... your notebook code here ...

# Return structured results to the parent notebook
result = {
    "status": "Success",
    "total_records": 1000,
    "validation_errors": 5,
    "processing_time_sec": 2.5
}

# Return the result to the parent notebook
dbutils.notebook.exit(json.dumps(result))
```

## Viewing Results in Azure Application Insights

After running your instrumented notebooks, you can view the results in Azure Application Insights:

1. **View traces**:
   - Navigate to your Azure Application Insights resource
   - Select "Transaction Search" from the left menu
   - Filter by operation name (e.g., "ETL_Pipeline" or "Notebook_Workflow")
   - Click on a trace to view the detailed span hierarchy and attributes

2. **View metrics**:
   - Select "Metrics" from the left menu
   - Choose "customMetrics" as the metric namespace
   - Select the metrics you want to visualize (e.g., "records_extracted")

3. **Run queries**:
   - Select "Logs" from the left menu
   - Use the following query to view ETL pipeline executions:
     ```sql
     dependencies
     | where name == "ETL_Pipeline"
     | extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"])
     | project timestamp, etl_pipeline_id, duration
     | order by timestamp desc
     ```

## Next Steps

- Explore the [Usage Guide](usage.md) for more detailed instructions
- See the [ETL Simulation Guide](etl_simulation.md) for a complete example
- Check out the [Parent-Child Notebooks Guide](parent_child_notebooks.md) for more advanced scenarios
- Review the [Tracing Guide](tracing.md) and [Metrics Guide](metrics.md) for details on the telemetry data
- Learn how to set up dashboards and alerts in the [Azure Monitoring Guide](azure_monitoring.md)
