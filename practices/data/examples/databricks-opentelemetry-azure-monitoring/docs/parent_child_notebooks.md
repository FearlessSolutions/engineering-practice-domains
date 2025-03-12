# Parent-Child Notebooks Guide

This guide explains the parent-child notebook example provided in this project, demonstrating how to use OpenTelemetry to instrument a parent notebook that calls child notebooks using `dbutils.notebook.run()`.

## Overview

The project includes three example notebooks that demonstrate a parent-child notebook workflow:

1. `parent_notebook_with_otel.py`: A parent notebook with OpenTelemetry instrumentation that calls two child notebooks
2. `child_notebook_1.py`: A child notebook that performs data validation (no OpenTelemetry instrumentation)
3. `child_notebook_2.py`: A child notebook that performs data aggregation (no OpenTelemetry instrumentation)

This example shows how to add observability to Databricks notebook workflows without requiring OpenTelemetry instrumentation in every notebook. The parent notebook creates spans and records metrics for the entire workflow, while the child notebooks remain uninstrumented but return structured data that can be used as span attributes and metrics.

## Key Concepts

### Parent-Child Notebook Architecture

In Databricks, you can create modular workflows by having a parent notebook call child notebooks using the `dbutils.notebook.run()` function. This approach offers several benefits:

- **Modularity**: Break down complex workflows into smaller, reusable components
- **Separation of concerns**: Each notebook can focus on a specific task
- **Reusability**: Child notebooks can be called from multiple parent notebooks
- **Maintainability**: Easier to maintain and update individual components

### OpenTelemetry Instrumentation Approach

The example demonstrates a pragmatic approach to instrumenting notebook workflows:

- **Parent notebook**: Contains all OpenTelemetry instrumentation
- **Child notebooks**: Remain uninstrumented but return structured data
- **Data flow**: Child notebooks return JSON data that the parent notebook uses for span attributes and metrics

This approach allows you to add comprehensive observability to your notebook workflows without modifying every notebook in your workspace.

## Notebook Cell Summaries

The parent notebook (`parent_notebook_with_otel.py`) contains the following key sections:

1. **Title & Introduction**  
   Introduces the notebook by describing its purpose: demonstrating how to use OpenTelemetry to instrument a parent notebook that calls child notebooks.

2. **Installation of Required Packages**  
   Provides pip commands to install the necessary libraries for OpenTelemetry and Azure Monitor integration.
   ```python
   %pip install --upgrade azure-monitor-opentelemetry opentelemetry-sdk azure-core opentelemetry-api opentelemetry-sdk opentelemetry-instrumentation
   ```

3. **Import OpenTelemetryHelper**  
   Imports the helper class from the uploaded `otel_helper.py` file.
   ```python
   from otel_helper import OpenTelemetryHelper
   ```

4. **Initialize OpenTelemetry**  
   Sets up the OpenTelemetry helper with a parent span and metrics configuration.
   ```python
   workflow_id = str(uuid.uuid4())
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

5. **Execute Child Notebook 1**  
   - Starts a span for the child notebook execution
   - Adds a span event to mark the start of the child notebook
   - Executes the child notebook using dbutils.notebook.run()
   - Parses the JSON result returned by the child notebook
   - Sets span attributes based on the child notebook results
   - Records metrics based on the child notebook results
   - Adds a span event to mark the completion of the child notebook
   - Handles potential errors
   - Ends the child notebook span

6. **Execute Child Notebook 2**  
   - Similar to the execution of Child Notebook 1, but for a different task
   - Demonstrates handling complex return values (like dictionaries)

7. **Finalization**  
   - Calculates overall workflow metrics
   - Sets span attributes for the parent span
   - Records metrics for the entire workflow
   - Adds a completion event
   - Ends the parent span

## Child Notebook Implementation

The child notebooks are simple Databricks notebooks that perform specific tasks and return structured data to the parent notebook.

### Child Notebook 1: Data Validation

This notebook simulates a data validation process:

1. Generates random data for validation
2. Simulates processing time
3. Calculates validation metrics
4. Returns a JSON object with the results:
   ```python
   result = {
       "task": "data_validation",
       "status": validation_status,
       "total_records": total_records,
       "validation_errors": validation_errors,
       "processing_time_sec": round(processing_time, 2),
       "records_per_second": round(records_per_second, 2),
       "timestamp": time.time()
   }
   
   # Return the result to the parent notebook
   dbutils.notebook.exit(json.dumps(result))
   ```

### Child Notebook 2: Data Aggregation

This notebook simulates a data aggregation process:

1. Performs multiple aggregation operations
2. Tracks execution time and memory usage
3. Returns a JSON object with the results:
   ```python
   result = {
       "task": "data_aggregation",
       "status_code": status_code,
       "num_aggregations": num_aggregations,
       "execution_time_sec": round(actual_execution_time, 2),
       "memory_usage_mb": round(memory_usage, 2),
       "aggregations_performed": aggregations_performed,
       "timestamp": time.time()
   }
   
   # Return the result to the parent notebook
   dbutils.notebook.exit(json.dumps(result))
   ```

## Span Structure

The parent-child notebook example creates the following span hierarchy:

- **Notebook_Workflow** (Parent Span)
  - **Child_Notebook_1** (Child Span)
  - **Child_Notebook_2** (Child Span)

Each span captures specific information relevant to its part of the workflow, providing a comprehensive view of the execution.

## OpenTelemetry Span Summary

| **Span**             | **Key Attributes**                                                                                    | **Description**                                            |
|----------------------|-------------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| `Notebook_Workflow`  | `workflow_id`, `total_records_processed`, `total_errors`, `workflow_success`, `validation_status`, `aggregation_status_code` | Overall workflow execution metrics and performance         |
| `Child_Notebook_1`   | `task`, `status`, `total_records`, `validation_errors`, `processing_time_sec`, `records_per_second`   | Data validation performance and error details              |
| `Child_Notebook_2`   | `task`, `status_code`, `num_aggregations`, `execution_time_sec`, `memory_usage_mb`, `aggregations_breakdown` | Data aggregation performance and processing details        |

## Span Events

The example demonstrates the use of span events to mark significant points in the workflow:

1. **Child Notebook Execution Started**: Added when a child notebook execution begins
   ```python
   workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Started", {
       "notebook_path": "child_notebook_1",
       "timestamp": time.time()
   })
   ```

2. **Child Notebook Execution Completed**: Added when a child notebook execution completes
   ```python
   workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Completed", {
       "status": child1_result["status"],
       "total_records": child1_result["total_records"],
       "validation_errors": child1_result["validation_errors"],
       "timestamp": time.time()
   })
   ```

3. **Workflow Completed**: Added when the entire workflow completes
   ```python
   workflow_otel_helper.add_span_event("Notebook_Workflow", "Workflow Completed", {
       "workflow_success": str(workflow_success),
       "total_records_processed": total_records_processed,
       "total_errors": total_errors,
       "total_execution_time_sec": round(total_execution_time, 2),
       "timestamp": time.time()
   })
   ```

## Benefits of This Approach

Using OpenTelemetry to instrument parent notebooks that call child notebooks provides several key benefits:

1. **Minimal Instrumentation**: Only the parent notebook needs OpenTelemetry instrumentation
2. **Comprehensive Monitoring**: Still captures detailed information about the entire workflow
3. **Flexibility**: Child notebooks can be used in multiple contexts without modification
4. **Simplified Maintenance**: Easier to maintain and update the instrumentation
5. **Standardized Return Values**: Encourages structured data return from child notebooks
6. **Correlation**: All spans are correlated with the same workflow ID
7. **Error Tracking**: Captures and correlates errors across the entire workflow

## Implementing in Your Own Workflows

To implement this approach in your own notebook workflows:

1. **Create a parent notebook** with OpenTelemetry instrumentation
2. **Ensure child notebooks return structured data** as JSON
3. **Use try/finally blocks** to properly end spans even if errors occur
4. **Capture return values** from child notebooks and use them as span attributes and metrics
5. **Add span events** to mark significant points in the workflow

## Alternative Approach: Instrumented Child Notebooks

While the primary approach described in this guide keeps OpenTelemetry instrumentation confined to the parent notebook, there may be scenarios where you want child notebooks to directly log their own span attributes and metrics. This section describes how to implement this alternative approach.

### Passing Context Information

Since you cannot directly pass Python objects (like the OpenTelemetryHelper instance) between notebooks, you need to pass context information that allows child notebooks to create properly correlated spans:

```python
# Parent notebook
workflow_id = str(uuid.uuid4())
parent_otel_helper = OpenTelemetryHelper(
    span_name="Notebook_Workflow",
    etl_pipeline_id=workflow_id,
    # Other configuration...
)

# Pass context to child notebook
child1_result_json = dbutils.notebook.run(
    "./child_notebook_1", 
    timeout_seconds=600,
    arguments={
        "workflow_id": workflow_id, 
        "parent_span": "Notebook_Workflow"
    }
)
```

### Initializing OpenTelemetryHelper in Child Notebooks

In the child notebook, retrieve the context information and initialize a new OpenTelemetryHelper instance:

```python
# Child notebook
# Retrieve context from parent
workflow_id = dbutils.notebook.getArgument("workflow_id")
parent_span = dbutils.notebook.getArgument("parent_span")

# Initialize OpenTelemetryHelper with the same workflow_id
child_otel_helper = OpenTelemetryHelper(
    span_name="Child_Notebook_1",
    etl_pipeline_id=workflow_id,
    span_metrics={
        "Child_Notebook_1": {
            "execution_time_sec": "histogram",
            "records_processed": "counter",
            "validation_errors": "counter"
        }
    },
    span_attributes={
        "parent_span": parent_span,  # Link to parent span
        "notebook_type": "child"
    }
)

# Use the helper to record metrics and set attributes directly
child_otel_helper.record_metric("Child_Notebook_1", "records_processed", total_records)
child_otel_helper.set_span_attribute("Child_Notebook_1", "validation_status", validation_status)

# Don't forget to end tracing before exiting
child_otel_helper.end_tracing("Child_Notebook_1")

# Still return structured data to parent
result = {
    "task": "data_validation",
    "status": validation_status,
    # Other result data...
}
dbutils.notebook.exit(json.dumps(result))
```

### Ensuring Proper Correlation

To ensure proper correlation between parent and child spans:

1. **Use the same workflow_id**: Pass the workflow_id from parent to child to ensure all spans are part of the same trace
2. **Reference the parent span**: Set a span attribute in the child that references the parent span name
3. **Use consistent naming**: Use a consistent naming convention for spans across notebooks
4. **Configure the same exporter**: Ensure all notebooks use the same Azure Monitor connection string

### Benefits and Trade-offs

**Benefits:**
- More detailed instrumentation: Child notebooks can add span attributes and events based on their internal processing
- Real-time metrics: Child notebooks can record metrics as they process data, rather than only at completion
- Finer-grained error tracking: Capture and report errors directly where they occur

**Trade-offs:**
- Increased complexity: Each notebook requires OpenTelemetry instrumentation
- Duplication: Configuration and setup code is duplicated across notebooks
- Maintenance overhead: Changes to instrumentation approach require updates to multiple notebooks
- Potential for inconsistency: Different notebooks might implement instrumentation differently

### When to Use This Approach

Consider using instrumented child notebooks when:
- Child notebooks contain complex logic that benefits from detailed internal instrumentation
- You need to capture metrics at specific points during child notebook execution
- Child notebooks are long-running and you want to track progress in real-time
- You need to capture detailed error information within child notebooks

## Next Steps

- Review the [Azure Monitoring Guide](azure_monitoring.md) for instructions on querying and visualizing the parent-child notebook telemetry data
- Explore the [Metrics Guide](metrics.md) for information on the metrics collected in the parent-child notebook example
- See the [Tracing Guide](tracing.md) for details on the trace attributes and components used in the parent-child notebook example
