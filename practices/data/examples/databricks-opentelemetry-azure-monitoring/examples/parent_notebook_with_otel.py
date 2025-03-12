# Databricks notebook source
# MAGIC %md
# MAGIC # Parent Notebook with OpenTelemetry Instrumentation
# MAGIC 
# MAGIC This notebook demonstrates how to use OpenTelemetry to instrument a parent notebook that calls child notebooks using `dbutils.notebook.run()`. The parent notebook creates spans and records metrics for the entire workflow, while the child notebooks remain uninstrumented.
# MAGIC 
# MAGIC Key features demonstrated:
# MAGIC - Creating a parent span for the entire workflow
# MAGIC - Executing child notebooks using dbutils.notebook.run()
# MAGIC - Capturing return values from child notebooks
# MAGIC - Adding span attributes and metrics based on child notebook results
# MAGIC - Using span events to mark significant points in the workflow
# MAGIC 
# MAGIC OpenTelemetry Parent Notebook/Span name: Notebook_Workflow

# COMMAND ----------

# MAGIC %md
# MAGIC ## OpenTelemetry Setup: Install prerequisite packages

# COMMAND ----------

# MAGIC %pip install --upgrade azure-monitor-opentelemetry opentelemetry-sdk azure-core opentelemetry-api opentelemetry-sdk opentelemetry-instrumentation
# MAGIC 
# MAGIC # Restart Python interpreter
# MAGIC %restart_python
# MAGIC dbutils.library.restartPython()

# COMMAND ----------

# MAGIC %md
# MAGIC ## OpenTelemetry Setup: Import OpenTelemetryHelper class
# MAGIC 
# MAGIC Import the OpenTelemetryHelper class from the uploaded otel_helper.py file.
# MAGIC 
# MAGIC ### Setup in Databricks
# MAGIC 
# MAGIC To use this in Databricks:
# MAGIC 
# MAGIC 1. **Upload the otel_helper.py file**:
# MAGIC    - In your Databricks workspace, use the "Create" button with the "File" option
# MAGIC    - Upload the otel_helper.py file to your workspace
# MAGIC 
# MAGIC 2. **Import the helper class**:
# MAGIC    ```python
# MAGIC    from otel_helper import OpenTelemetryHelper
# MAGIC    ```
# MAGIC 
# MAGIC 3. **Restart the Python interpreter if needed**:
# MAGIC    ```
# MAGIC    %restart_python
# MAGIC    ```

# COMMAND ----------

import os
import time
import uuid
import json
from otel_helper import OpenTelemetryHelper

# COMMAND ----------

# MAGIC %md
# MAGIC ## Initialize OpenTelemetryHelper for the Notebook Workflow
# MAGIC 
# MAGIC - Create a unique workflow ID
# MAGIC - Initialize the OpenTelemetryHelper with a parent span
# MAGIC - Define metrics for the workflow and child notebook executions

# COMMAND ----------

print("--- PARENT-SPAN: Notebook_Workflow")

# Create a unique workflow ID
workflow_id = str(uuid.uuid4())
print(f"    WORKFLOW_ID: {workflow_id}")

# Initialize OpenTelemetry with a parent span and metrics
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

print(f"    OpenTelemetryHelper instance created with workflow_id: {workflow_id}")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Execute Child Notebook 1: Data Validation
# MAGIC 
# MAGIC - Start a span for the child notebook execution
# MAGIC - Execute the child notebook using dbutils.notebook.run()
# MAGIC - Parse the return value from the child notebook
# MAGIC - Add span attributes and metrics based on the child notebook results

# COMMAND ----------

# Start a span for Child Notebook 1 execution
print("Starting Child_Notebook_1 span...")
workflow_otel_helper.start_tracing("Child_Notebook_1", {
    "etl_pipeline_id": workflow_id,
    "notebook_path": "child_notebook_1"
})

# Add a span event to mark the start of the child notebook
workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Started", {
    "notebook_path": "child_notebook_1",
    "timestamp": time.time()
})

try:
    # Execute Child Notebook 1
    print("Executing Child Notebook 1...")
    child1_result_json = dbutils.notebook.run("./child_notebook_1", timeout_seconds=600)
    
    # Parse the JSON result
    child1_result = json.loads(child1_result_json)
    print(f"Child Notebook 1 completed with status: {child1_result['status']}")
    
    # Set span attributes based on child notebook results
    for key, value in child1_result.items():
        if isinstance(value, (str, int, float, bool)):
            workflow_otel_helper.set_span_attribute("Child_Notebook_1", key, value)
    
    # Record metrics
    workflow_otel_helper.record_metric("Child_Notebook_1", "execution_time_sec", child1_result["processing_time_sec"])
    workflow_otel_helper.record_metric("Child_Notebook_1", "records_processed", child1_result["total_records"])
    workflow_otel_helper.record_metric("Child_Notebook_1", "validation_errors", child1_result["validation_errors"])
    
    # Add a span event for the completion
    workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Completed", {
        "status": child1_result["status"],
        "total_records": child1_result["total_records"],
        "validation_errors": child1_result["validation_errors"],
        "timestamp": time.time()
    })
    
except Exception as e:
    error_msg = str(e)
    print(f"Error executing Child Notebook 1: {error_msg}")
    
    # Set error attributes
    workflow_otel_helper.set_span_attribute("Child_Notebook_1", "error", "true")
    workflow_otel_helper.set_span_attribute("Child_Notebook_1", "error_message", error_msg)
    
    # Add an error event
    workflow_otel_helper.add_span_event("Child_Notebook_1", "Child Notebook Execution Failed", {
        "error_message": error_msg,
        "timestamp": time.time()
    })
    
    # Raise the exception to stop execution
    raise
    
finally:
    # End the Child_Notebook_1 span
    workflow_otel_helper.end_tracing("Child_Notebook_1")
    print("Child_Notebook_1 span completed")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Execute Child Notebook 2: Data Aggregation
# MAGIC 
# MAGIC - Start a span for the child notebook execution
# MAGIC - Execute the child notebook using dbutils.notebook.run()
# MAGIC - Parse the return value from the child notebook
# MAGIC - Add span attributes and metrics based on the child notebook results

# COMMAND ----------

# Start a span for Child Notebook 2 execution
print("Starting Child_Notebook_2 span...")
workflow_otel_helper.start_tracing("Child_Notebook_2", {
    "etl_pipeline_id": workflow_id,
    "notebook_path": "child_notebook_2"
})

# Add a span event to mark the start of the child notebook
workflow_otel_helper.add_span_event("Child_Notebook_2", "Child Notebook Execution Started", {
    "notebook_path": "child_notebook_2",
    "timestamp": time.time()
})

try:
    # Execute Child Notebook 2
    print("Executing Child Notebook 2...")
    child2_result_json = dbutils.notebook.run("./child_notebook_2", timeout_seconds=600)
    
    # Parse the JSON result
    child2_result = json.loads(child2_result_json)
    print(f"Child Notebook 2 completed with status code: {child2_result['status_code']}")
    
    # Set span attributes based on child notebook results
    for key, value in child2_result.items():
        if isinstance(value, (str, int, float, bool)):
            workflow_otel_helper.set_span_attribute("Child_Notebook_2", key, value)
    
    # For complex objects like dictionaries, we can convert to string or extract specific values
    if "aggregations_performed" in child2_result:
        workflow_otel_helper.set_span_attribute(
            "Child_Notebook_2", 
            "aggregations_breakdown", 
            json.dumps(child2_result["aggregations_performed"])
        )
    
    # Record metrics
    workflow_otel_helper.record_metric("Child_Notebook_2", "execution_time_sec", child2_result["execution_time_sec"])
    workflow_otel_helper.record_metric("Child_Notebook_2", "aggregations_performed", child2_result["num_aggregations"])
    workflow_otel_helper.record_metric("Child_Notebook_2", "memory_usage_mb", child2_result["memory_usage_mb"])
    
    # Add a span event for the completion
    workflow_otel_helper.add_span_event("Child_Notebook_2", "Child Notebook Execution Completed", {
        "status_code": child2_result["status_code"],
        "num_aggregations": child2_result["num_aggregations"],
        "execution_time_sec": child2_result["execution_time_sec"],
        "timestamp": time.time()
    })
    
except Exception as e:
    error_msg = str(e)
    print(f"Error executing Child Notebook 2: {error_msg}")
    
    # Set error attributes
    workflow_otel_helper.set_span_attribute("Child_Notebook_2", "error", "true")
    workflow_otel_helper.set_span_attribute("Child_Notebook_2", "error_message", error_msg)
    
    # Add an error event
    workflow_otel_helper.add_span_event("Child_Notebook_2", "Child Notebook Execution Failed", {
        "error_message": error_msg,
        "timestamp": time.time()
    })
    
    # Raise the exception to stop execution
    raise
    
finally:
    # End the Child_Notebook_2 span
    workflow_otel_helper.end_tracing("Child_Notebook_2")
    print("Child_Notebook_2 span completed")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Finalize Workflow and Record Overall Metrics
# MAGIC 
# MAGIC - Calculate overall workflow metrics
# MAGIC - Set span attributes for the parent span
# MAGIC - Record metrics for the entire workflow
# MAGIC - Add a completion event
# MAGIC - End the parent span

# COMMAND ----------

# Calculate overall workflow metrics
total_records_processed = child1_result["total_records"]
total_errors = child1_result["validation_errors"] + (1 if child2_result["status_code"] != 200 else 0)
workflow_success = child1_result["status"] == "Success" and child2_result["status_code"] == 200

# Set span attributes for the parent span
print("Setting attributes for the Notebook_Workflow span...")
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "total_records_processed", total_records_processed)
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "total_errors", total_errors)
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "workflow_success", str(workflow_success).lower())

# Add child notebook results as attributes to the parent span
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "validation_status", child1_result["status"])
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "aggregation_status_code", child2_result["status_code"])
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "validation_time_sec", child1_result["processing_time_sec"])
workflow_otel_helper.set_span_attribute("Notebook_Workflow", "aggregation_time_sec", child2_result["execution_time_sec"])

# Record metrics for the entire workflow
total_execution_time = time.time() - workflow_otel_helper.span_start_times["Notebook_Workflow"]
workflow_otel_helper.record_metric("Notebook_Workflow", "total_execution_time_sec", total_execution_time)
workflow_otel_helper.record_metric("Notebook_Workflow", "total_records_processed", total_records_processed)
workflow_otel_helper.record_metric("Notebook_Workflow", "error_count", total_errors)

# Add a completion event to the parent span
workflow_otel_helper.add_span_event("Notebook_Workflow", "Workflow Completed", {
    "workflow_success": str(workflow_success),
    "total_records_processed": total_records_processed,
    "total_errors": total_errors,
    "total_execution_time_sec": round(total_execution_time, 2),
    "timestamp": time.time()
})

# End the parent span
workflow_otel_helper.end_tracing("Notebook_Workflow")
print("Notebook_Workflow trace ended")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Workflow Summary
# MAGIC 
# MAGIC This notebook has demonstrated how to use OpenTelemetry to instrument a parent notebook that calls child notebooks. The key points are:
# MAGIC 
# MAGIC 1. The parent notebook creates spans and records metrics for the entire workflow
# MAGIC 2. Child notebooks remain uninstrumented but return structured data
# MAGIC 3. The parent notebook captures return values from child notebooks and uses them as span attributes and metrics
# MAGIC 4. Span events are used to mark significant points in the workflow
# MAGIC 
# MAGIC This approach allows for comprehensive monitoring of notebook workflows without requiring OpenTelemetry instrumentation in every notebook.
