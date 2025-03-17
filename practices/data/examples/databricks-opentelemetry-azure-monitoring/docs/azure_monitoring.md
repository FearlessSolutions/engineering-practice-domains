# Azure Monitoring Guide

This guide provides detailed information on how to query, visualize, and set up alerts for the OpenTelemetry data exported to Azure Application Insights.

## Table of Contents

- [Overview](#overview)
- [Querying Data in Azure Application Insights](#querying-data-in-azure-application-insights)
  - [Accessing the Query Editor](#accessing-the-query-editor)
  - [Querying Trace Data](#querying-trace-data)
  - [Querying Metrics Data](#querying-metrics-data)
- [Advanced Queries](#advanced-queries)
  - [Finding Failed ETL Runs](#finding-failed-etl-runs)
  - [Analyzing ETL Performance Trends](#analyzing-etl-performance-trends)
  - [Correlating Spans for a Specific ETL Run](#correlating-spans-for-a-specific-etl-run)
- [Viewing Metrics in Azure Monitor](#viewing-metrics-in-azure-monitor)
  - [Using the Metrics Explorer](#using-the-metrics-explorer)
  - [Creating Custom Charts](#creating-custom-charts)
- [Metrics Visualization & Alerts](#metrics-visualization--alerts)
  - [Creating a Dashboard](#creating-a-dashboard)
  - [Setting Up Alerts](#setting-up-alerts)
  - [Example Alert Queries](#example-alert-queries)
- [Best Practices for Azure Monitoring](#best-practices-for-azure-monitoring)
- [Parent-Child Notebook Queries](#parent-child-notebook-queries)
  - [Querying Trace Data for Parent-Child Notebooks](#querying-trace-data-for-parent-child-notebooks)
  - [Querying Metrics Data for Parent-Child Notebooks](#querying-metrics-data-for-parent-child-notebooks)
  - [Advanced Queries for Parent-Child Notebooks](#advanced-queries-for-parent-child-notebooks)
  - [Visualizing Parent-Child Notebook Metrics](#visualizing-parent-child-notebook-metrics)
  - [Alert Examples for Parent-Child Notebooks](#alert-examples-for-parent-child-notebooks)
- [Next Steps](#next-steps)

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

## Parent-Child Notebook Queries

The following queries are specific to the parent-child notebook example and can be used to analyze the execution of notebook workflows.

### Querying Trace Data for Parent-Child Notebooks

#### Query for `Notebook_Workflow` (Parent Span)

```sql
dependencies
| where name == "Notebook_Workflow"
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         total_records_processed = toint(customDimensions["total_records_processed"]),
         total_errors = toint(customDimensions["total_errors"]),
         workflow_success = tostring(customDimensions["workflow_success"]),
         validation_status = tostring(customDimensions["validation_status"]),
         aggregation_status_code = toint(customDimensions["aggregation_status_code"]),
         validation_time_sec = todouble(customDimensions["validation_time_sec"]),
         aggregation_time_sec = todouble(customDimensions["aggregation_time_sec"])
| project timestamp, workflow_id, total_records_processed, total_errors, 
          workflow_success, validation_status, aggregation_status_code,
          validation_time_sec, aggregation_time_sec, duration
| order by timestamp desc
```

#### Query for `Child_Notebook_1` (Data Validation Span)

```sql
dependencies
| where name == "Child_Notebook_1"
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         notebook_path = tostring(customDimensions["notebook_path"]),
         task = tostring(customDimensions["task"]),
         status = tostring(customDimensions["status"]),
         total_records = toint(customDimensions["total_records"]),
         validation_errors = toint(customDimensions["validation_errors"]),
         processing_time_sec = todouble(customDimensions["processing_time_sec"]),
         records_per_second = todouble(customDimensions["records_per_second"]),
         error = tostring(customDimensions["error"]),
         error_message = tostring(customDimensions["error_message"])
| project timestamp, workflow_id, notebook_path, task, status, total_records, 
          validation_errors, processing_time_sec, records_per_second, error, error_message, duration
| order by timestamp desc
```

#### Query for `Child_Notebook_2` (Data Aggregation Span)

```sql
dependencies
| where name == "Child_Notebook_2"
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         notebook_path = tostring(customDimensions["notebook_path"]),
         task = tostring(customDimensions["task"]),
         status_code = toint(customDimensions["status_code"]),
         num_aggregations = toint(customDimensions["num_aggregations"]),
         execution_time_sec = todouble(customDimensions["execution_time_sec"]),
         memory_usage_mb = todouble(customDimensions["memory_usage_mb"]),
         aggregations_breakdown = tostring(customDimensions["aggregations_breakdown"]),
         error = tostring(customDimensions["error"]),
         error_message = tostring(customDimensions["error_message"])
| project timestamp, workflow_id, notebook_path, task, status_code, num_aggregations, 
          execution_time_sec, memory_usage_mb, aggregations_breakdown, error, error_message, duration
| order by timestamp desc
```

### Querying Metrics Data for Parent-Child Notebooks

#### Query for Notebook Workflow Metrics

```sql
customMetrics
| where name in ("total_execution_time_sec", "total_records_processed", "error_count")
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, workflow_id, timestamp
| order by timestamp desc
```

#### Query for Child Notebook 1 Metrics

```sql
customMetrics
| where name in ("execution_time_sec", "records_processed", "validation_errors")
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, workflow_id, timestamp
| order by timestamp desc
```

#### Query for Child Notebook 2 Metrics

```sql
customMetrics
| where name in ("execution_time_sec", "aggregations_performed", "memory_usage_mb")
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"])
| summarize count = sum(valueCount), avg_value = avg(value), max_value = max(value), min_value = min(value) by name, workflow_id, timestamp
| order by timestamp desc
```

### Advanced Queries for Parent-Child Notebooks

#### Finding Failed Notebook Workflows

```sql
dependencies
| where name == "Notebook_Workflow"
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         workflow_success = tostring(customDimensions["workflow_success"]),
         total_errors = toint(customDimensions["total_errors"])
| where workflow_success == "false" or total_errors > 0
| project timestamp, workflow_id, workflow_success, total_errors
| order by timestamp desc
```

#### Analyzing Child Notebook Performance

```sql
dependencies
| where name in ("Child_Notebook_1", "Child_Notebook_2")
| extend workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         notebook_name = name
| extend execution_time = case(
           notebook_name == "Child_Notebook_1", todouble(customDimensions["processing_time_sec"]),
           notebook_name == "Child_Notebook_2", todouble(customDimensions["execution_time_sec"]),
           0.0
         )
| summarize avg_execution_time = avg(execution_time), max_execution_time = max(execution_time) by notebook_name, bin(timestamp, 1d)
| render timechart
```

#### Correlating Spans for a Specific Workflow

```sql
dependencies
| where customDimensions.etl_pipeline_id == "your-workflow-id"
| extend span_name = name,
         workflow_id = tostring(customDimensions["etl_pipeline_id"]),
         duration_sec = todouble(duration) / 1000
| project timestamp, span_name, duration_sec, operation_Id
| order by timestamp asc
```

### Visualizing Parent-Child Notebook Metrics

When creating visualizations for parent-child notebook workflows, consider the following approaches:

1. **Workflow Success Rate**: Create a chart showing the percentage of successful workflows over time
   ```sql
   dependencies
   | where name == "Notebook_Workflow"
   | extend workflow_success = tostring(customDimensions["workflow_success"])
   | summarize success_count = countif(workflow_success == "true"), 
              total_count = count() 
              by bin(timestamp, 1d)
   | extend success_rate = 100.0 * success_count / total_count
   | project timestamp, success_rate
   | render timechart
   ```

2. **Child Notebook Execution Time Comparison**: Compare the execution time of different child notebooks
   ```sql
   dependencies
   | where name in ("Child_Notebook_1", "Child_Notebook_2")
   | extend execution_time = case(
     name == "Child_Notebook_1", todouble(customDimensions["processing_time_sec"]),
     name == "Child_Notebook_2", todouble(customDimensions["execution_time_sec"]),
     0.0
   )
   | summarize avg_execution_time = avg(execution_time), max_execution_time = max(execution_time) by name, bin(timestamp, 1d)
   | render timechart
   ```

3. **Validation Error Trends**: Track validation errors over time
   ```sql
   dependencies
   | where name == "Child_Notebook_1"
   | extend validation_errors = toint(customDimensions["validation_errors"]),
            total_records = toint(customDimensions["total_records"])
   | summarize avg_errors = avg(validation_errors), 
              avg_error_rate = 100.0 * avg(todouble(validation_errors) / todouble(total_records)) 
              by bin(timestamp, 1d)
   | render timechart
   ```

### Alert Examples for Parent-Child Notebooks

#### Alert for Failed Notebook Workflows

```sql
dependencies
| where name == "Notebook_Workflow"
| extend workflow_success = tostring(customDimensions["workflow_success"])
| where workflow_success == "false"
| count
```

#### Alert for High Validation Error Rate

```sql
dependencies
| where name == "Child_Notebook_1"
| extend validation_errors = toint(customDimensions["validation_errors"]),
         total_records = toint(customDimensions["total_records"])
| extend error_rate = (todouble(validation_errors) / todouble(total_records)) * 100
| where error_rate > 5
```

#### Alert for Slow Child Notebook Execution

```sql
dependencies
| where name == "Child_Notebook_2"
| extend execution_time_sec = todouble(customDimensions["execution_time_sec"])
| where execution_time_sec > 10 // Alert if execution takes more than 10 seconds
| count
```

## Next Steps

- Review the [Tracing Guide](tracing.md) for details on the trace attributes and components
- Explore the [Metrics Guide](metrics.md) for information on the metrics collected
- See the [ETL Simulation Guide](etl_simulation.md) for a complete example of OpenTelemetry implementation
- Check out the [Parent-Child Notebooks Guide](parent_child_notebooks.md) for details on instrumenting notebook workflows
