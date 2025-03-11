# Tracing Guide

This guide provides detailed information about the OpenTelemetry tracing instrumentation used in this project, including span attributes, components, and the overall span structure.

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
| `DataExtraction`    | `records_extracted`, `http_status_code`, `http_response_size`, `http_request_duration_sec`             | API extraction performance and error details               |
| `DataTransformation`| `records_input`, `records_transformed`, `records_failed`, `transformation_duration_sec`, `records_transformed_rate` | Data transformation performance and processing speed       |
| `DataLoading`       | `records_written`, `write_duration_sec`, `storage_path`, `storage_status`                              | Data loading performance and error tracking                |

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

## Next Steps

- Explore the [Metrics Guide](metrics.md) for information on the metrics collected
- See the [Azure Monitoring Guide](azure_monitoring.md) for instructions on querying and visualizing the trace data
- Review the [ETL Simulation Guide](etl_simulation.md) for a complete example of tracing implementation
