# Usage Guide

This guide provides instructions and examples for using the OpenTelemetry instrumentation in your Databricks notebooks.

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

## Additional Documentation

For more detailed information about using OpenTelemetry with this project, refer to these guides:

- [ETL Simulation Guide](etl_simulation.md): Information about the simulated ETL with and without OpenTelemetry
- [Tracing Guide](tracing.md): Details about tracing attributes, components, and the OpenTelemetry Span Summary
- [Metrics Guide](metrics.md): Information about metric values and instrumentation
- [Azure Monitoring Guide](azure_monitoring.md): Querying data in Azure Application Insights, viewing metrics, and setting up visualizations and alerts
