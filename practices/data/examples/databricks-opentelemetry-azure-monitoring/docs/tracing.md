# Tracing Guide

This guide provides detailed information about the OpenTelemetry tracing instrumentation used in this project, including span attributes, components, and the overall span structure.

## Table of Contents

- [Tracing Instrumentation Overview](#tracing-instrumentation-overview)
- [Span Structure](#span-structure)
- [OpenTelemetry Span Summary](#opentelemetry-span-summary)
- [Detailed Span Trace Attributes](#detailed-span-trace-attributes)
  - [ETL_Pipeline (Parent Span)](#etl_pipeline-parent-span)
  - [DataExtraction (API Call Span)](#dataextraction-api-call-span)
  - [DataTransformation (Processing Span)](#datatransformation-processing-span)
  - [DataLoading (Storage Span)](#dataloading-storage-span)
- [Span Events](#span-events)
- [Implementing Tracing in Your Code](#implementing-tracing-in-your-code)
- [Best Practices for Tracing](#best-practices-for-tracing)
- [Parent-Child Notebook Tracing](#parent-child-notebook-tracing)
  - [Span Structure for Parent-Child Notebooks](#span-structure-for-parent-child-notebooks)
  - [Detailed Span Trace Attributes for Parent-Child Notebooks](#detailed-span-trace-attributes-for-parent-child-notebooks)
  - [Span Events in Parent-Child Notebooks](#span-events-in-parent-child-notebooks)
  - [Capturing Child Notebook Results as Span Attributes](#capturing-child-notebook-results-as-span-attributes)
  - [Best Practices for Parent-Child Notebook Tracing](#best-practices-for-parent-child-notebook-tracing)
  - [Alternative Approach: Instrumented Child Notebooks](#alternative-approach-instrumented-child-notebooks)
- [Next Steps](#next-steps)

## Tracing Instrumentation Overview

The ETL simulation is instrumented with OpenTelemetry tracing to provide detailed insights into the execution flow and performance characteristics of each stage of the pipeline. Tracing creates a hierarchical representation of the ETL process, with a parent span for the overall pipeline and child spans for each stage.

## Span Structure

The tracing instrumentation creates the following span hierarchy:

- **ETL_Pipeline** (Parent Span)
  - **DataExtraction** (Child Span)
  - **DataTransformation** (Child Span)
  - **DataLoading** (Child Span)

Each span captures specific information relevant to its stage of the ETL process, providing a comprehensive view of the pipeline's execution.

## OpenTelemetry Span Summary

| **Span**            | **Key Attributes**                                                                                    | **Description**                                            |
|---------------------|-------------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| `ETL_Pipeline`      | `etl_pipeline_id`, `records_input_total`, `records_output_total`, `etl_total_duration_sec`, `etl_efficiency_rate`, `retry_attempts`, `etl_error_count` | Overall ETL execution metrics and performance              |
| `DataExtraction`    | `function_name`, `records_extracted`, `http_status_code`, `http_response_size`, `http_request_duration_sec`             | API extraction performance and error details               |
| `DataTransformation`| `function_name`, `records_input`, `records_transformed`, `records_failed`, `transformation_duration_sec`, `records_transformed_rate` | Data transformation performance and processing speed       |
| `DataLoading`       | `function_name`, `records_written`, `write_duration_sec`, `storage_path`, `storage_status`                              | Data loading performance and error tracking                |

## Detailed Span Trace Attributes

Below are detailed tables that break down the trace attributes for each span in the ETL simulation.

### `ETL_Pipeline` (Parent Span)

| **Attribute Name**          | **Example Value**                           | **Description** |
|-----------------------------|---------------------------------------------|-----------------|
| `etl_pipeline_id`           | `"etl_run_1234"`                            | Unique identifier for the ETL run. |
| `records_extracted`         | `10000`                                     | Number of records extracted from the source. |
| `records_transformed`       | `9800`                                      | Number of records successfully transformed. |
| `records_written`           | `9800`                                      | Number of records successfully written to storage. |
| `records_failed`            | `200`                                       | Count of records that failed during transformation. |
| `duration_transformation_sec` | `2.345`                                   | Duration of the transformation stage in seconds. |
| `duration_write_sec`        | `3.876`                                     | Duration of the data loading stage in seconds. |
| `rate_transformation_records_per_sec` | `4179.5`                          | Transformation processing rate (records per second). |
| `rate_write_records_per_sec` | `2528.4`                                   | Data loading rate (records per second). |
| `rate_etl_efficiency_percent` | `98.0`                                    | Efficiency percentage (output/input records). |
| `extraction_failed`         | `false`                                     | Indicates if the extraction stage failed. |
| `extraction_response_status` | `200`                                      | HTTP status code from the API call. |
| `loading_error`             | `false`                                     | Indicates if the loading stage encountered errors. |
| `etl_status`                | `"Success"`                                 | Overall ETL process status. |
| `etl_retry_attempts`        | `1`                                         | Number of retry attempts made during execution. |
| `etl_error_count`           | `0`                                         | Total number of errors encountered. |
| `storage_path`              | `"/mnt/output/transformed_data"`            | The file or table where data is stored. |
| `storage_type`              | `"Azure Data Lake"`                         | The storage system used. |
| `transformation_type`       | `"Data Cleaning"`                           | Type of transformation applied. |
| `processing_engine`         | `"PySpark"`                                 | Engine used for processing. |
| `error_message`             | `""`                                        | Error message if any stage failed. |
| `data_source`               | `"External API"`                            | Source of the data (set during initialization). |
| `etl_type`                  | `"Full Load"`                               | Type of ETL process (set during initialization). |

### `DataExtraction` (API Call Span)

| **Attribute Name**             | **Example Value**                          | **Description** |
|--------------------------------|--------------------------------------------|-----------------|
| `etl_pipeline_id`              | `"etl_run_1234"`                           | Unique identifier for the ETL run. |
| `function_name`                | `"extract_data_from_api"`                  | Name of the function that was decorated with the trace_function decorator. |
| `http_request_duration_sec`    | `1.234`                                    | Duration of the API request in seconds. |
| `records_extracted`            | `10000`                                    | Number of records fetched from the API. |
| `http_response_size`           | `500000`                                   | Simulated size of the API response in bytes. |
| `extraction_status`            | `"Success"`                                | Status of the extraction process. |
| `http_status_code`             | `200`                                      | HTTP status code from the API call. |
| `error`                        | `"false"`                                  | Indicates if an error occurred. |
| `error_message`                | `""`                                       | Error message if the extraction failed. |
| `duration_sec`                 | `1.234`                                    | Total duration of the extraction span. |

### `DataTransformation` (Processing Span)

| **Attribute Name**             | **Example Value**                          | **Description** |
|--------------------------------|--------------------------------------------|-----------------|
| `etl_pipeline_id`              | `"etl_run_1234"`                           | Unique identifier for the ETL run. |
| `function_name`                | `"transform_data"`                         | Name of the function that was decorated with the trace_function decorator. |
| `error`                        | `"false"`                                  | Indicates if an error occurred during transformation. |
| `records_input`                | `10000`                                    | Number of records received for transformation. |
| `records_transformed`          | `9800`                                     | Number of records successfully transformed. |
| `records_failed`               | `200`                                      | Number of records that failed during transformation. |
| `transformation_duration_sec`  | `2.345`                                    | Duration of the transformation process in seconds. |
| `records_transformed_rate`     | `4179.5`                                   | Transformation speed (records per second). |
| `transformation_type`          | `"Data Cleaning"`                          | Type of transformation applied. |
| `processing_engine`            | `"PySpark"`                                | Engine used for processing. |
| `error_message`                | `""`                                       | Error details if transformation failed. |
| `duration_sec`                 | `2.345`                                    | Total duration of the transformation span. |

### `DataLoading` (Storage Span)

| **Attribute Name**          | **Example Value**                | **Description** |
|-----------------------------|----------------------------------|-----------------|
| `etl_pipeline_id`           | `"etl_run_1234"`                 | Unique identifier for the ETL run. |
| `function_name`             | `"perform_data_loading"`         | Name of the function that was decorated with the trace_function decorator. |
| `error`                     | `"false"`                        | Indicates if the write process failed. |
| `error_msg`                 | `""`                             | Error details if loading failed. |
| `records_written`           | `9800`                           | Number of records successfully written to storage. |
| `storage_path`              | `"/mnt/output/transformed_data"` | The file or table where data is stored. |
| `storage_type`              | `"Azure Data Lake"`              | The storage system used. |
| `file_size_bytes`           | `490000`                         | Total size of the stored file in bytes. |
| `compression_type`          | `"Parquet"`                      | Compression format used. |
| `partition_count`           | `5`                              | Number of partitions used for storage. |
| `batch_size`                | `1000`                           | Number of records written per batch. |
| `write_duration_sec`        | `3.876`                          | Time taken to write the data. |
| `records_write_rate`        | `2528.4`                         | Write speed (records per second). |
| `storage_status`            | `"Success"`                      | Final status of data writing. |
| `duration_sec`              | `3.876`                          | Total duration of the loading span. |

## Span Events

In addition to attributes, spans can also include events that mark significant occurrences during the execution. For example, the `ETL_Pipeline` span includes a completion event:

```python
event_attributes = {
    "records_processed": records_written,
    "efficiency_rate": round(etl_efficiency_rate, 2),
    "extraction_status": "Failed" if extraction_failed else "Success",
    "loading_status": "Failed" if error else "Success"
}
etl_otel_helper.add_span_event("ETL_Pipeline", "ETL Pipeline Completed", event_attributes)
```

## Implementing Tracing in Your Code

To add tracing to your own ETL pipelines, follow these steps:

1. **Initialize the OpenTelemetryHelper**:
   ```python
   etl_pipeline_id = str(uuid.uuid4())
   etl_otel_helper = OpenTelemetryHelper(
       span_name="ETL_Pipeline",
       etl_pipeline_id=etl_pipeline_id,
       span_attributes={"data_source": "Your Data Source", "etl_type": "Your ETL Type"}
   )
   ```

2. **Create spans for each stage**:
   ```python
   # Start a span
   etl_otel_helper.start_tracing("YourStage", {
       "etl_pipeline_id": etl_pipeline_id,
       "additional_attribute": "value"
   })
   
   # Your stage code here
   
   # Set span attributes
   etl_otel_helper.set_span_attribute("YourStage", "attribute_name", attribute_value)
   
   # End the span
   etl_otel_helper.end_tracing("YourStage")
   ```

3. **Use the function decorator for encapsulated operations**:
   ```python
   @etl_otel_helper.trace_function("YourStage", {
       "attribute_name": "result_key"
   })
   def your_function():
       # Function implementation
       return {
           "result_key": "value"
       }
   
   # Note: The function_name attribute ("your_function") is automatically added to the span
   ```

4. **Add events for significant occurrences**:
   ```python
   etl_otel_helper.add_span_event("YourStage", "Event Name", {
       "event_attribute": "value"
   })
   ```

5. **End the parent span when the pipeline is complete**:
   ```python
   etl_otel_helper.end_tracing("ETL_Pipeline")
   ```

## Best Practices for Tracing

1. **Use meaningful span names** that reflect the operation being performed.
2. **Include the `etl_pipeline_id`** in all spans to correlate related operations.
3. **Set appropriate attributes** that provide context and help with troubleshooting.
4. **Record error information** when failures occur to make debugging easier.
5. **Use try/finally blocks** to ensure spans are properly ended, even if exceptions occur.
6. **Add events** to mark significant occurrences during execution.
7. **Keep span hierarchy consistent** with your application's logical structure.

## Parent-Child Notebook Tracing

The parent-child notebook example demonstrates how to use OpenTelemetry tracing to monitor the execution of notebook workflows. This approach allows you to track the execution of child notebooks without requiring OpenTelemetry instrumentation in each notebook.

### Span Structure for Parent-Child Notebooks

The parent-child notebook example creates the following span hierarchy:

- **Notebook_Workflow** (Parent Span)
  - **Child_Notebook_1** (Child Span)
  - **Child_Notebook_2** (Child Span)

Each span captures specific information relevant to its part of the workflow, providing a comprehensive view of the execution.

### Detailed Span Trace Attributes for Parent-Child Notebooks

Below are detailed tables that break down the trace attributes for each span in the parent-child notebook example.

#### `Notebook_Workflow` (Parent Span)

| **Attribute Name**          | **Example Value**                           | **Description** |
|-----------------------------|---------------------------------------------|-----------------|
| `etl_pipeline_id`           | `"workflow_1234"`                           | Unique identifier for the notebook workflow. |
| `workflow_type`             | `"Data Processing Pipeline"`                | Type of workflow being executed. |
| `environment`               | `"Development"`                             | Environment where the workflow is running. |
| `total_records_processed`   | `8500`                                      | Total number of records processed in the workflow. |
| `total_errors`              | `25`                                        | Total number of errors encountered during the workflow. |
| `workflow_success`          | `"true"`                                    | Indicates if the workflow completed successfully. |
| `validation_status`         | `"Success"`                                 | Status of the data validation process. |
| `aggregation_status_code`   | `200`                                       | Status code from the data aggregation process. |
| `validation_time_sec`       | `3.45`                                      | Duration of the data validation process in seconds. |
| `aggregation_time_sec`      | `5.67`                                      | Duration of the data aggregation process in seconds. |
| `duration_sec`              | `10.23`                                     | Total duration of the workflow in seconds. |

#### `Child_Notebook_1` (Data Validation Span)

| **Attribute Name**          | **Example Value**                           | **Description** |
|-----------------------------|---------------------------------------------|-----------------|
| `etl_pipeline_id`           | `"workflow_1234"`                           | Unique identifier for the notebook workflow. |
| `notebook_path`             | `"child_notebook_1"`                        | Path to the child notebook. |
| `task`                      | `"data_validation"`                         | Task performed by the child notebook. |
| `status`                    | `"Success"`                                 | Status of the child notebook execution. |
| `total_records`             | `8500`                                      | Number of records processed in the child notebook. |
| `validation_errors`         | `25`                                        | Number of validation errors found. |
| `processing_time_sec`       | `3.45`                                      | Duration of the child notebook execution in seconds. |
| `records_per_second`        | `2463.77`                                   | Processing rate in records per second. |
| `error`                     | `"false"`                                   | Indicates if an error occurred during execution. |
| `error_message`             | `""`                                        | Error message if an error occurred. |
| `duration_sec`              | `3.45`                                      | Total duration of the child notebook span. |

#### `Child_Notebook_2` (Data Aggregation Span)

| **Attribute Name**          | **Example Value**                           | **Description** |
|-----------------------------|---------------------------------------------|-----------------|
| `etl_pipeline_id`           | `"workflow_1234"`                           | Unique identifier for the notebook workflow. |
| `notebook_path`             | `"child_notebook_2"`                        | Path to the child notebook. |
| `task`                      | `"data_aggregation"`                        | Task performed by the child notebook. |
| `status_code`               | `200`                                       | Status code of the child notebook execution. |
| `num_aggregations`          | `35`                                        | Number of aggregation operations performed. |
| `execution_time_sec`        | `5.67`                                      | Duration of the child notebook execution in seconds. |
| `memory_usage_mb`           | `350.5`                                     | Memory usage during execution in MB. |
| `aggregations_breakdown`    | `{"sum": 10, "average": 8, "count": 17}`    | Breakdown of aggregation operations by type. |
| `error`                     | `"false"`                                   | Indicates if an error occurred during execution. |
| `error_message`             | `""`                                        | Error message if an error occurred. |
| `duration_sec`              | `5.67`                                      | Total duration of the child notebook span. |

### Span Events in Parent-Child Notebooks

The parent-child notebook example demonstrates the use of span events to mark significant points in the workflow:

#### Child Notebook Execution Started

```python
workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Started", {
    "notebook_path": "child_notebook_1",
    "timestamp": time.time()
})
```

This event marks the start of a child notebook execution and includes the notebook path and timestamp.

#### Child Notebook Execution Completed

```python
workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Completed", {
    "status": child1_result["status"],
    "total_records": child1_result["total_records"],
    "validation_errors": child1_result["validation_errors"],
    "timestamp": time.time()
})
```

This event marks the completion of a child notebook execution and includes key results from the child notebook.

#### Workflow Completed

```python
workflow_otel_helper.add_span_event("Notebook_Workflow", "Workflow Completed", {
    "workflow_success": str(workflow_success),
    "total_records_processed": total_records_processed,
    "total_errors": total_errors,
    "total_execution_time_sec": round(total_execution_time, 2),
    "timestamp": time.time()
})
```

This event marks the completion of the entire workflow and includes summary information about the execution.

### Capturing Child Notebook Results as Span Attributes

A key feature of the parent-child notebook tracing approach is the ability to capture the results returned by child notebooks as span attributes. This is done by parsing the JSON result returned by the child notebook and setting the values as span attributes:

```python
# Execute Child Notebook 1
child1_result_json = dbutils.notebook.run("./child_notebook_1", timeout_seconds=600)
child1_result = json.loads(child1_result_json)

# Set span attributes based on child notebook results
for key, value in child1_result.items():
    if isinstance(value, (str, int, float, bool)):
        workflow_otel_helper.set_span_attribute("Child_Notebook_1", key, value)
```

This approach allows you to capture detailed information about the execution of child notebooks without requiring OpenTelemetry instrumentation in each notebook.

### Best Practices for Parent-Child Notebook Tracing

1. **Create a span for each child notebook execution** to track the performance and results of each notebook.

2. **Use span events** to mark significant points in the workflow, such as the start and completion of child notebook executions.

3. **Capture child notebook results as span attributes** to provide detailed information about the execution.

4. **Use try/finally blocks** to ensure spans are properly ended, even if exceptions occur.

5. **Include error handling** to capture and report errors that occur during child notebook execution.

6. **Use consistent naming conventions** for spans, attributes, and events across all notebook workflows.

7. **Correlate spans using the workflow ID** to enable comprehensive analysis of the entire workflow.

## Next Steps

- Explore the [Metrics Guide](metrics.md) for information on the metrics collected
- See the [Azure Monitoring Guide](azure_monitoring.md) for instructions on querying and visualizing the trace data
- Review the [ETL Simulation Guide](etl_simulation.md) for a complete example of tracing implementation
- Check out the [Parent-Child Notebooks Guide](parent_child_notebooks.md) for details on instrumenting notebook workflows
