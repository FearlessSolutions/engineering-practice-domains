# Usage Guide

This guide provides instructions and examples for using the OpenTelemetry instrumentation in your Databricks notebooks.

## Table of Contents

- [Basic Usage](#basic-usage)
  - [Importing the Helper Class](#importing-the-helper-class)
  - [Initializing OpenTelemetry](#initializing-opentelemetry)
  - [Configuring Metrics](#configuring-metrics)
- [Creating and Managing Spans](#creating-and-managing-spans)
  - [Starting a Span](#starting-a-span)
  - [Setting Span Attributes](#setting-span-attributes)
  - [Recording Metrics](#recording-metrics)
  - [Adding Span Events](#adding-span-events)
  - [Ending a Span](#ending-a-span)
- [Using the Function Decorator](#using-the-function-decorator)
- [Complete Example](#complete-example)
- [Best Practices](#best-practices)
- [Parent-Child Notebook Instrumentation](#parent-child-notebook-instrumentation)
  - [Basic Approach](#basic-approach)
  - [Example](#example)
  - [Child Notebook Return Values](#child-notebook-return-values)
  - [Complete Example](#complete-example-1)
- [Additional Documentation](#additional-documentation)

## Basic Usage

### Importing the Helper Class

First, import the OpenTelemetryHelper class in your notebook:

```python
# If you've uploaded the file to your Databricks workspace (recommended)
from otel_helper import OpenTelemetryHelper

# Or if installed as a package
# from databricks_opentelemetry_azure_monitoring import OpenTelemetryHelper
```

### Initializing OpenTelemetry

Initialize the OpenTelemetryHelper with a parent span name and a unique identifier for your ETL pipeline:

```python
import uuid

# Generate a unique ID for this ETL run
etl_pipeline_id = str(uuid.uuid4())

# Initialize OpenTelemetry with a parent span
otel_helper = OpenTelemetryHelper(
    span_name="ETL_Pipeline",
    etl_pipeline_id=etl_pipeline_id,
    span_attributes={"data_source": "External API", "etl_type": "Full Load"}
)
```

### Configuring Metrics

You can define metrics for different spans during initialization:

```python
otel_helper = OpenTelemetryHelper(
    span_name="ETL_Pipeline",
    etl_pipeline_id=etl_pipeline_id,
    span_metrics={
        "DataExtraction": {"records_extracted": "counter", "http_request_duration_sec": "histogram"},
        "DataTransformation": {"records_transformed": "counter", "transformation_duration_sec": "histogram"},
        "DataLoading": {"records_written": "counter", "write_duration_sec": "histogram"},
    },
    span_attributes={"data_source": "External API", "etl_type": "Full Load"}
)
```

## Creating and Managing Spans

### Starting a Span

To start a new span for a specific stage of your ETL pipeline:

```python
otel_helper.start_tracing("DataExtraction", {
    "etl_pipeline_id": etl_pipeline_id,
    "additional_attribute": "value"
})
```

### Setting Span Attributes

Add attributes to a span to provide context and metadata:

```python
otel_helper.set_span_attribute("DataExtraction", "records_extracted", 10000)
otel_helper.set_span_attribute("DataExtraction", "http_status_code", 200)
```

### Recording Metrics

Record metric values for a specific span:

```python
otel_helper.record_metric("DataExtraction", "records_extracted", 10000)
otel_helper.record_metric("DataExtraction", "http_request_duration_sec", 1.5)
```

### Adding Span Events

Add events to a span to mark significant occurrences:

```python
otel_helper.add_span_event("ETL_Pipeline", "Data Validation Complete", {
    "records_validated": 10000,
    "validation_status": "Success"
})
```

### Ending a Span

End a span when the operation is complete:

```python
otel_helper.end_tracing("DataExtraction")
```

## Using the Function Decorator

You can use the decorator to automatically trace function execution:

```python
@otel_helper.trace_function("DataLoading", {
    "records_written": "records_written",
    "storage_path": "storage_path"
})
def perform_data_loading(records_to_write):
    # Function implementation
    return {
        "records_written": records_to_write,
        "storage_path": "/path/to/data"
    }
```

The `trace_function` decorator automatically adds the function name as a span attribute called `function_name`. This allows you to identify which specific function was executed when analyzing trace data. For example, in the code above, the span will include an attribute `function_name` with the value `"perform_data_loading"`.

This attribute is particularly useful when:
- Multiple functions are decorated with the same span name
- You want to filter or group trace data by specific function implementations
- You need to correlate performance issues with specific code functions

## Complete Example

Here's a complete example of instrumenting an ETL pipeline:

```python
import uuid
# Import the helper class
from otel_helper import OpenTelemetryHelper

# Initialize
etl_pipeline_id = str(uuid.uuid4())
otel_helper = OpenTelemetryHelper(
    span_name="ETL_Pipeline",
    etl_pipeline_id=etl_pipeline_id,
    span_metrics={
        "DataExtraction": {"records_extracted": "counter"},
        "DataTransformation": {"records_transformed": "counter"},
        "DataLoading": {"records_written": "counter"},
    }
)

try:
    # Data Extraction
    otel_helper.start_tracing("DataExtraction")
    # ... extraction code ...
    records_extracted = 10000
    otel_helper.set_span_attribute("DataExtraction", "records_extracted", records_extracted)
    otel_helper.record_metric("DataExtraction", "records_extracted", records_extracted)
    otel_helper.end_tracing("DataExtraction")
    
    # Data Transformation
    otel_helper.start_tracing("DataTransformation")
    # ... transformation code ...
    records_transformed = 9500
    otel_helper.set_span_attribute("DataTransformation", "records_transformed", records_transformed)
    otel_helper.record_metric("DataTransformation", "records_transformed", records_transformed)
    otel_helper.end_tracing("DataTransformation")
    
    # Data Loading
    @otel_helper.trace_function("DataLoading", {
        "records_written": "records_written"
    })
    def perform_data_loading():
        # ... loading code ...
        return {"records_written": records_transformed}
    
    loading_result = perform_data_loading()
    
    # Add final event
    otel_helper.add_span_event("ETL_Pipeline", "ETL Pipeline Completed", {
        "records_processed": loading_result["records_written"]
    })
    
finally:
    # Always end the parent span
    otel_helper.end_tracing("ETL_Pipeline")
```

## Best Practices

1. **Always end your spans**: Use try/finally blocks to ensure spans are properly ended, even if exceptions occur.

2. **Use meaningful span names**: Choose descriptive names that reflect the operation being performed.

3. **Add relevant attributes**: Include attributes that provide context and help with troubleshooting.

4. **Record appropriate metrics**: Focus on metrics that are useful for monitoring and alerting.

5. **Use the function decorator**: For operations that can be encapsulated in functions, use the decorator for cleaner code.

6. **Structure your spans hierarchically**: Create a parent span for the overall operation and child spans for specific stages.

7. **Include error information**: When errors occur, set error attributes to make troubleshooting easier.

## Parent-Child Notebook Instrumentation

You can use OpenTelemetry to instrument parent notebooks that call child notebooks using `dbutils.notebook.run()`. This approach allows you to monitor the execution of child notebooks without adding OpenTelemetry instrumentation to each child notebook.

### Basic Approach

1. Create a parent span in the parent notebook
2. Create child spans for each child notebook execution
3. Execute child notebooks using `dbutils.notebook.run()`
4. Capture return values from child notebooks
5. Use these return values as span attributes and metrics

### Example

Here's a simplified example of instrumenting a parent notebook that calls child notebooks:

```python
import json
from otel_helper import OpenTelemetryHelper

# Initialize OpenTelemetry
workflow_id = str(uuid.uuid4())
otel_helper = OpenTelemetryHelper(
    span_name="Parent_Notebook",
    etl_pipeline_id=workflow_id
)

# Start a span for the child notebook execution
otel_helper.start_tracing("Child_Notebook_1")

try:
    # Execute the child notebook
    child_result_json = dbutils.notebook.run("./child_notebook_1", timeout_seconds=600)
    
    # Parse the JSON result
    child_result = json.loads(child_result_json)
    
    # Set span attributes based on child notebook results
    for key, value in child_result.items():
        if isinstance(value, (str, int, float, bool)):
            otel_helper.set_span_attribute("Child_Notebook_1", key, value)
    
    # Add a span event for the completion
    otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Completed", {
        "status": child_result["status"]
    })
    
finally:
    # End the child notebook span
    otel_helper.end_tracing("Child_Notebook_1")

# End the parent span when the workflow is complete
otel_helper.end_tracing("Parent_Notebook")
```

### Child Notebook Return Values

Child notebooks should return structured data that can be used as span attributes or metrics. For example:

```python
# In the child notebook
result = {
    "status": "Success",
    "records_processed": 10000,
    "processing_time_sec": 5.2
}

# Return the result to the parent notebook
dbutils.notebook.exit(json.dumps(result))
```

### Complete Example

For a complete example of parent-child notebook instrumentation, see the `examples/parent_notebook_with_otel.py` notebook in this project. This example demonstrates:

- Creating a parent span for the entire workflow
- Executing multiple child notebooks
- Capturing return values from child notebooks
- Adding span attributes and metrics based on child notebook results
- Using span events to mark significant points in the workflow

## Additional Documentation

For more detailed information about using OpenTelemetry with this project, refer to these guides:

- [ETL Simulation Guide](etl_simulation.md): Information about the simulated ETL with and without OpenTelemetry
- [Tracing Guide](tracing.md): Details about tracing attributes, components, and the OpenTelemetry Span Summary
- [Metrics Guide](metrics.md): Information about metric values and instrumentation
- [Azure Monitoring Guide](azure_monitoring.md): Querying data in Azure Application Insights, viewing metrics, and setting up visualizations and alerts
