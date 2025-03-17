# Metrics Guide

This guide provides detailed information about the OpenTelemetry metrics instrumentation used in this project, including metric types, values, and how to use them for monitoring your ETL pipelines.

## Table of Contents

- [Metrics Instrumentation Overview](#metrics-instrumentation-overview)
- [Metric Types](#metric-types)
- [Metrics Configuration](#metrics-configuration)
- [Metrics Table](#metrics-table)
- [Recording Metrics](#recording-metrics)
- [Metric Attributes](#metric-attributes)
- [Implementing Metrics in Your Code](#implementing-metrics-in-your-code)
- [Best Practices for Metrics](#best-practices-for-metrics)
- [Viewing Metrics in Azure Monitor](#viewing-metrics-in-azure-monitor)
- [Parent-Child Notebook Metrics](#parent-child-notebook-metrics)
  - [Metrics Configuration for Parent-Child Notebooks](#metrics-configuration-for-parent-child-notebooks)
  - [Parent-Child Notebook Metrics Table](#parent-child-notebook-metrics-table)
  - [Recording Metrics from Child Notebook Results](#recording-metrics-from-child-notebook-results)
  - [Best Practices for Parent-Child Notebook Metrics](#best-practices-for-parent-child-notebook-metrics)
- [Next Steps](#next-steps)

## Metrics Instrumentation Overview

In addition to tracing, the ETL simulation is instrumented with OpenTelemetry metrics to provide real-time performance insights. These custom metric instruments are exported to Azure Monitor Application Insights, allowing you to:

- Track API request durations and success rates
- Monitor transformation processing times and error counts
- Analyze data loading performance, including throughput and write durations
- Assess overall ETL efficiency and detect performance bottlenecks

## Metric Types

OpenTelemetry supports different types of metrics, each suited for specific use cases:

1. **Counters**: Used for values that only increase, such as the number of records processed or errors encountered.
2. **Histograms**: Used for measuring the distribution of values, such as request durations or processing times.

In the ETL simulation, both types are used to provide comprehensive monitoring.

## Metrics Configuration

Metrics are configured during the initialization of the OpenTelemetryHelper:

```python
etl_otel_helper = OpenTelemetryHelper(
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

This configuration creates the following metrics:
- For `DataExtraction`: a counter for `records_extracted` and a histogram for `http_request_duration_sec`
- For `DataTransformation`: a counter for `records_transformed` and a histogram for `transformation_duration_sec`
- For `DataLoading`: a counter for `records_written` and a histogram for `write_duration_sec`

## Metrics Table

The following table provides details on all metrics used in the ETL simulation:

| **Metric Name**                     | **Type**   | **Stage**           | **Description**                                                                 |
|-------------------------------------|------------|---------------------|---------------------------------------------------------------------------------|
| `records_extracted`                 | Counter    | Data Extraction     | Number of records extracted from the API.                                       |
| `http_request_duration_sec`         | Histogram  | Data Extraction     | Duration (in seconds) of the HTTP request during data extraction.               |
| `records_transformed`               | Counter    | Data Transformation | Number of records successfully transformed.                                     |
| `records_failed_transformation`     | Counter    | Data Transformation | Number of records that failed during transformation.                            |
| `transformation_duration_sec`       | Histogram  | Data Transformation | Time taken (in seconds) to transform the data.                                  |
| `records_written`                   | Counter    | Data Loading        | Number of records successfully written to storage.                              |
| `write_duration_sec`                | Histogram  | Data Loading        | Duration (in seconds) of the data loading process.                              |
| `etl_total_duration_sec`            | Histogram  | Pipeline            | Total duration of the ETL pipeline execution.                                   |
| `etl_efficiency_rate`               | Histogram  | Pipeline            | ETL efficiency percentage, calculated as `(records_written / records_extracted) * 100`. |
| `retry_attempts`                    | Counter    | Pipeline            | Number of retry attempts made during the ETL process.                           |
| `etl_error_count`                   | Counter    | Pipeline            | Total count of errors encountered during the ETL execution.                     |

## Recording Metrics

Metrics are recorded during the execution of the ETL pipeline using the `record_metric` method:

```python
etl_otel_helper.record_metric("DataExtraction", "records_extracted", records_extracted)
etl_otel_helper.record_metric("DataExtraction", "http_request_duration_sec", duration_sec)
```

For metrics associated with functions decorated with `trace_function`, the metrics are automatically recorded if the function returns a dictionary with keys matching the metric names:

```python
@etl_otel_helper.trace_function("DataLoading", {
    "etl_pipeline_id": "etl_pipeline_id",
    "records_written": "records_written",
    "write_duration_sec": "write_duration_sec",
    # other attributes...
})
def perform_data_loading(records_to_write):
    # Function implementation
    return {
        "records_written": records_written,
        "write_duration_sec": write_duration_sec,
        # other return values...
    }
```

## Metric Attributes

Metrics can include attributes to provide additional context. By default, the `etl_pipeline_id` is included as an attribute for all metrics, allowing you to correlate metrics with specific ETL runs:

```python
metric_attrs = {"etl_pipeline_id": self.etl_pipeline_id}
```

You can also provide custom attributes when recording metrics:

```python
etl_otel_helper.record_metric("DataExtraction", "records_extracted", records_extracted, {
    "source": "API",
    "batch_id": batch_id
})
```

## Implementing Metrics in Your Code

To add metrics to your own ETL pipelines, follow these steps:

1. **Define metrics during initialization**:
   ```python
   etl_otel_helper = OpenTelemetryHelper(
       span_name="ETL_Pipeline",
       etl_pipeline_id=etl_pipeline_id,
       span_metrics={
           "YourStage": {"your_counter": "counter", "your_histogram": "histogram"},
       }
   )
   ```

2. **Record metric values during execution**:
   ```python
   # Record a counter
   etl_otel_helper.record_metric("YourStage", "your_counter", value)
   
   # Record a histogram
   etl_otel_helper.record_metric("YourStage", "your_histogram", duration)
   ```

3. **Use the function decorator for automatic metric recording**:
   ```python
   @etl_otel_helper.trace_function("YourStage", {
       "your_counter": "result_counter",
       "your_histogram": "result_duration"
   })
   def your_function():
       # Function implementation
       return {
           "result_counter": value,
           "result_duration": duration
       }
   ```

## Best Practices for Metrics

1. **Choose appropriate metric types** based on what you're measuring:
   - Use counters for cumulative values (records processed, errors)
   - Use histograms for distributions (durations, sizes)

2. **Include meaningful attributes** to provide context and enable filtering.

3. **Use consistent naming conventions** for metrics across your application.

4. **Focus on actionable metrics** that provide insights into performance and reliability.

5. **Consider cardinality** when adding attributes to avoid excessive metric combinations.

6. **Correlate metrics with traces** using the `etl_pipeline_id` to enable comprehensive analysis.

## Viewing Metrics in Azure Monitor

Once exported to Azure Monitor, metrics can be:

- Visualized in dashboards
- Used for alerting
- Analyzed for trends and anomalies

See the [Azure Monitoring Guide](azure_monitoring.md) for detailed instructions on working with metrics in Azure Monitor.

## Parent-Child Notebook Metrics

In addition to the ETL simulation metrics, the parent-child notebook example demonstrates how to use metrics to monitor notebook workflows. These metrics provide insights into the performance and reliability of notebook executions.

### Metrics Configuration for Parent-Child Notebooks

Metrics for the parent-child notebook example are configured during the initialization of the OpenTelemetryHelper:

```python
workflow_otel_helper = OpenTelemetryHelper(
    span_name="Notebook_Workflow",
    etl_pipeline_id=workflow_id,
    span_metrics={
        "Notebook_Workflow": {
            "total_execution_time_sec": "histogram",
            "total_records_processed": "counter",
            "error_count": "counter"
        },
        "Child_Notebook_1": {
            "execution_time_sec": "histogram",
            "records_processed": "counter",
            "validation_errors": "counter"
        },
        "Child_Notebook_2": {
            "execution_time_sec": "histogram",
            "aggregations_performed": "counter",
            "memory_usage_mb": "histogram"
        }
    },
    span_attributes={
        "workflow_type": "Data Processing Pipeline",
        "environment": "Development"
    }
)
```

### Parent-Child Notebook Metrics Table

The following table provides details on all metrics used in the parent-child notebook example:

| **Metric Name**                | **Type**   | **Stage**           | **Description**                                                                 |
|--------------------------------|------------|---------------------|---------------------------------------------------------------------------------|
| `total_execution_time_sec`     | Histogram  | Notebook_Workflow   | Total duration (in seconds) of the entire notebook workflow.                    |
| `total_records_processed`      | Counter    | Notebook_Workflow   | Total number of records processed across all child notebooks.                   |
| `error_count`                  | Counter    | Notebook_Workflow   | Total count of errors encountered during the workflow execution.                |
| `execution_time_sec`           | Histogram  | Child_Notebook_1    | Duration (in seconds) of the data validation notebook execution.                |
| `records_processed`            | Counter    | Child_Notebook_1    | Number of records processed in the data validation notebook.                    |
| `validation_errors`            | Counter    | Child_Notebook_1    | Number of validation errors found during data validation.                       |
| `execution_time_sec`           | Histogram  | Child_Notebook_2    | Duration (in seconds) of the data aggregation notebook execution.               |
| `aggregations_performed`       | Counter    | Child_Notebook_2    | Number of aggregation operations performed in the data aggregation notebook.    |
| `memory_usage_mb`              | Histogram  | Child_Notebook_2    | Memory usage (in MB) during the data aggregation notebook execution.            |

### Recording Metrics from Child Notebook Results

A key feature of the parent-child notebook example is the ability to record metrics based on the results returned by child notebooks. This is done by parsing the JSON result returned by the child notebook and using the values to record metrics:

```python
# Execute Child Notebook 1
child1_result_json = dbutils.notebook.run("./child_notebook_1", timeout_seconds=600)
child1_result = json.loads(child1_result_json)

# Record metrics based on child notebook results
workflow_otel_helper.record_metric("Child_Notebook_1", "execution_time_sec", child1_result["processing_time_sec"])
workflow_otel_helper.record_metric("Child_Notebook_1", "records_processed", child1_result["total_records"])
workflow_otel_helper.record_metric("Child_Notebook_1", "validation_errors", child1_result["validation_errors"])
```

This approach allows you to collect metrics from child notebooks without requiring OpenTelemetry instrumentation in each notebook.

### Best Practices for Parent-Child Notebook Metrics

1. **Define consistent return structures** for child notebooks to ensure reliable metric recording.

2. **Use try/catch blocks** when parsing child notebook results to handle unexpected return values.

3. **Record overall workflow metrics** that aggregate results from all child notebooks.

4. **Include error counts** to track the reliability of notebook executions.

5. **Monitor execution times** to identify performance bottlenecks in specific notebooks.

6. **Use the workflow ID** to correlate metrics across all notebooks in the workflow.

## Next Steps

- Review the [Tracing Guide](tracing.md) for details on the trace attributes and components
- See the [Azure Monitoring Guide](azure_monitoring.md) for instructions on querying and visualizing the metric data
- Explore the [ETL Simulation Guide](etl_simulation.md) for a complete example of metrics implementation
- Check out the [Parent-Child Notebooks Guide](parent_child_notebooks.md) for details on instrumenting notebook workflows
