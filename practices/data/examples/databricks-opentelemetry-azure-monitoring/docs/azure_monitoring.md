# Azure Monitoring Guide

This guide provides detailed information on how to query, visualize, and set up alerts for the OpenTelemetry data exported to Azure Application Insights.

## Overview

The OpenTelemetry instrumentation in this project exports both traces and metrics to Azure Application Insights, allowing you to:

- Query trace data to analyze ETL pipeline execution
- Visualize metrics to monitor performance
- Set up alerts for proactive notification of issues
- Create dashboards for continuous monitoring

## Querying Data in Azure Application Insights

Azure Application Insights uses the Kusto Query Language (KQL) to query telemetry data. Below are example queries for each span in the ETL pipeline.

### Accessing the Query Editor

1. Navigate to your Azure Application Insights resource in the Azure Portal
2. Select "Logs" from the left menu
3. Use the query editor to run KQL queries

### Querying Trace Data

#### Query for `ETL_Pipeline` (Parent Span)

```sql
dependencies
| where name == "ETL_Pipeline"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         records_extracted = toint(customDimensions["records_extracted"]),
         records_transformed = toint(customDimensions["records_transformed"]),
         records_written = toint(customDimensions["records_written"]),
         records_failed = toint(customDimensions["records_failed"]),
         duration_transformation_sec = todouble(customDimensions["duration_transformation_sec"]),
         duration_write_sec = todouble(customDimensions["duration_write_sec"]),
         rate_etl_efficiency_percent = todouble(customDimensions["rate_etl_efficiency_percent"]),
         etl_status = tostring(customDimensions["etl_status"]),
         etl_retry_attempts = toint(customDimensions["etl_retry_attempts"]),
         etl_error_count = toint(customDimensions["etl_error_count"])
| project timestamp, etl_pipeline_id, records_extracted, records_transformed, records_written, 
          records_failed, duration_transformation_sec, duration_write_sec, 
          rate_etl_efficiency_percent, etl_status, etl_retry_attempts, etl_error_count
| order by timestamp desc
```

#### Query for `DataExtraction` (API Call Span)

```sql
dependencies
| where name == "DataExtraction"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         function_name = tostring(customDimensions["function_name"]),
         records_extracted = toint(customDimensions["records_extracted"]),
         http_status_code = tostring(customDimensions["http_status_code"]),
         http_response_size = toint(customDimensions["http_response_size"]),
         http_request_duration_sec = todouble(customDimensions["http_request_duration_sec"]),
         error = tostring(customDimensions["error"]),
         error_message = tostring(customDimensions["error_message"])
| project timestamp, etl_pipeline_id, function_name, records_extracted, http_status_code, http_response_size, http_request_duration_sec, error, error_message
| order by timestamp desc
```

#### Query for `DataTransformation` (Processing Span)

```sql
dependencies
| where name == "DataTransformation"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         function_name = tostring(customDimensions["function_name"]),
         records_input = toint(customDimensions["records_input"]),
         records_transformed = toint(customDimensions["records_transformed"]),
         records_failed = toint(customDimensions["records_failed"]),
         transformation_duration_sec = todouble(customDimensions["transformation_duration_sec"]),
         records_transformed_rate = todouble(customDimensions["records_transformed_rate"]),
         transformation_type = tostring(customDimensions["transformation_type"]),
         processing_engine = tostring(customDimensions["processing_engine"]),
         error = tostring(customDimensions["error"]),
         error_message = tostring(customDimensions["error_message"])
| project timestamp, etl_pipeline_id, function_name, records_input, records_transformed, records_failed, transformation_duration_sec, 
          records_transformed_rate, transformation_type, processing_engine, error, error_message
| order by timestamp desc
```

#### Query for `DataLoading` (Storage Span)

```sql
dependencies
| where name == "DataLoading"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         function_name = tostring(customDimensions["function_name"]),
         records_written = toint(customDimensions["records_written"]),
         storage_path = tostring(customDimensions["storage_path"]),
         storage_type = tostring(customDimensions["storage_type"]),
         file_size_bytes = toint(customDimensions["file_size_bytes"]),
         compression_type = tostring(customDimensions["compression_type"]),
         partition_count = toint(customDimensions["partition_count"]),
         batch_size = toint(customDimensions["batch_size"]),
         write_duration_sec = todouble(customDimensions["write_duration_sec"]),
         records_write_rate = todouble(customDimensions["records_write_rate"]),
         storage_status = tostring(customDimensions["storage_status"]),
         error = tostring(customDimensions["error"]),
         error_message = tostring(customDimensions["error_msg"])
| project timestamp, etl_pipeline_id, function_name, records_written, storage_path, storage_type, file_size_bytes, compression_type, 
          partition_count, batch_size, write_duration_sec, records_write_rate, storage_status, error, error_message
| order by timestamp desc
```

### Querying Metrics Data

#### Query for Extraction Metrics

```sql
customMetrics
| where name in ("records_extracted", "http_request_duration_sec")
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, etl_pipeline_id, timestamp
| order by timestamp desc
```

#### Query for Transformation Metrics

```sql
customMetrics
| where name in ("records_transformed", "transformation_duration_sec")
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, etl_pipeline_id, timestamp
| order by timestamp desc
```

#### Query for Loading Metrics

```sql
customMetrics
| where name in ("records_written", "write_duration_sec")
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, etl_pipeline_id, timestamp
| order by timestamp desc
```

#### Query for Overall ETL Metrics

```sql
customMetrics
| where name in ("etl_total_duration_sec", "etl_efficiency_rate", "retry_attempts", "etl_error_count")
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, etl_pipeline_id, timestamp
| order by timestamp desc
```

## Advanced Queries

### Finding Failed ETL Runs

```sql
dependencies
| where name == "ETL_Pipeline"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         etl_status = tostring(customDimensions["etl_status"]),
         etl_error_count = toint(customDimensions["etl_error_count"])
| where etl_status == "Failed" or etl_error_count > 0
| project timestamp, etl_pipeline_id, etl_status, etl_error_count
| order by timestamp desc
```

### Analyzing ETL Performance Trends

```sql
dependencies
| where name == "ETL_Pipeline"
| extend etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         records_extracted = toint(customDimensions["records_extracted"]),
         records_written = toint(customDimensions["records_written"]),
         duration_sec = todouble(customDimensions["duration_sec"])
| summarize avg_duration = avg(duration_sec), 
            avg_records = avg(records_extracted), 
            avg_efficiency = avg(todouble(records_written) / todouble(records_extracted) * 100) 
            by bin(timestamp, 1d)
| render timechart
```

### Correlating Spans for a Specific ETL Run

```sql
dependencies
| where customDimensions.etl_pipeline_id == "your-etl-pipeline-id"
| extend span_name = name,
         etl_pipeline_id = tostring(customDimensions["etl_pipeline_id"]),
         function_name = tostring(customDimensions["function_name"]),
         duration_sec = todouble(customDimensions["duration_sec"])
| project timestamp, span_name, function_name, duration_sec, operation_Id
| order by timestamp asc
```

## Viewing Metrics in Azure Monitor

### Using the Metrics Explorer

1. Navigate to your Azure Application Insights resource in the Azure Portal
2. Select "Metrics" from the left menu
3. In the Metrics Explorer:
   - Select "customMetrics" as the metric namespace
   - Choose the specific metric you want to visualize (e.g., "records_extracted")
   - Select the aggregation method (e.g., Sum, Average, Max)
   - Set the time range and granularity
   - Add filters if needed (e.g., to focus on a specific ETL pipeline ID)

### Creating Custom Charts

You can create custom charts to visualize multiple metrics together:

1. In the Metrics Explorer, click "Add metric" to add additional metrics to the chart
2. Use the "Split by" option to segment the data (e.g., by ETL pipeline ID)
3. Experiment with different chart types (Line, Bar, Area) to find the most informative visualization
4. Save the chart to a dashboard for ongoing monitoring

## Metrics Visualization & Alerts

### Creating a Dashboard

1. In the Azure Portal, navigate to "Dashboard" and click "New dashboard"
2. Give your dashboard a name (e.g., "ETL Pipeline Monitoring")
3. Add tiles to your dashboard:
   - Pin charts from the Metrics Explorer
   - Pin query results from the Logs view
   - Add custom text and markdown for documentation
4. Arrange the tiles to create a comprehensive view of your ETL pipeline performance
5. Save the dashboard and share it with your team

### Setting Up Alerts

You can set up alerts to be notified when certain conditions are met:

1. Navigate to your Azure Application Insights resource
2. Select "Alerts" from the left menu
3. Click "Create" to create a new alert rule
4. Define the condition:
   - Select the signal type (e.g., Log search)
   - Configure the query (e.g., to detect failed ETL runs)
   - Set the threshold and evaluation frequency
5. Configure the action group:
   - Create a new action group or select an existing one
   - Define the notification methods (email, SMS, webhook, etc.)
6. Add details like alert name, description, and severity
7. Review and create the alert rule

### Example Alert Queries

#### Alert for Failed ETL Runs

```sql
dependencies
| where name == "ETL_Pipeline"
| extend etl_status = tostring(customDimensions["etl_status"]),
         etl_error_count = toint(customDimensions["etl_error_count"])
| where etl_status == "Failed" or etl_error_count > 0
| count
```

#### Alert for Slow Transformation

```sql
dependencies
| where name == "DataTransformation"
| extend transformation_duration_sec = todouble(customDimensions["transformation_duration_sec"])
| where transformation_duration_sec > 5 // Alert if transformation takes more than 5 seconds
| count
```

#### Alert for Low ETL Efficiency

```sql
dependencies
| where name == "ETL_Pipeline"
| extend rate_etl_efficiency_percent = todouble(customDimensions["rate_etl_efficiency_percent"])
| where rate_etl_efficiency_percent < 90 // Alert if efficiency is below 90%
| count
```

## Best Practices for Azure Monitoring

1. **Start with key metrics**: Focus on the most important metrics first, such as overall ETL duration, record counts, and error rates.

2. **Create a comprehensive dashboard**: Build a dashboard that provides a complete view of your ETL pipeline performance.

3. **Set up proactive alerts**: Configure alerts for critical conditions to be notified of issues before they impact users.

4. **Use log-based metrics**: For complex metrics that can't be captured directly, use log queries to derive them.

5. **Correlate traces and metrics**: Use the `etl_pipeline_id` to correlate traces and metrics for a complete picture of each ETL run.

6. **Regularly review and refine**: Periodically review your monitoring setup and refine it based on your evolving needs.

7. **Document your queries**: Save and document your most useful queries for future reference and team knowledge sharing.

## Next Steps

- Review the [Tracing Guide](tracing.md) for details on the trace attributes and components
- Explore the [Metrics Guide](metrics.md) for information on the metrics collected
- See the [ETL Simulation Guide](etl_simulation.md) for a complete example of OpenTelemetry implementation
