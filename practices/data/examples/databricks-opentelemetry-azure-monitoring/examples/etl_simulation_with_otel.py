# Databricks notebook source
# MAGIC %md
# MAGIC # ETL Pipeline Simulation WITH OpenTelemetry 
# MAGIC
# MAGIC This notebook simulates an ETL pipeline. It will be later used to demonstrate both trace and metric instrumentation using OpenTelemetry and Azure Monitor.
# MAGIC
# MAGIC The instrumentation will be added to this notebook in a way that is modular and has the least code changes to the original notebook cells.
# MAGIC
# MAGIC OpenTelemetry Parent Notebook/Span name:  ETL_Pipeline
# MAGIC

# COMMAND ----------

# MAGIC %md
# MAGIC
# MAGIC ## OpenTelemetry Setup For Notebook: Install prereq packages
# MAGIC

# COMMAND ----------

# MAGIC %pip install --upgrade azure-monitor-opentelemetry opentelemetry-sdk azure-core opentelemetry-api opentelemetry-sdk opentelemetry-instrumentation
# MAGIC
# MAGIC # Restart Python interpreter
# MAGIC
# MAGIC %restart_python
# MAGIC dbutils.library.restartPython()
# MAGIC
# MAGIC

# COMMAND ----------

# MAGIC %md
# MAGIC ## OpenTelemetry Setup For Notebook: Import OpenTelemetryHelper class
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
from otel_helper import OpenTelemetryHelper


# COMMAND ----------

# MAGIC %md
# MAGIC ## OpenTelemetry Setup For Notebook: Initialize helper class for ETL_Pipeline span
# MAGIC
# MAGIC - Starts of the parent ETL_Pipeline span (or entire notebooks span).
# MAGIC - Creates a Unique ETL run identifier
# MAGIC - Init tracking values for span
# MAGIC - Initialize OpenTelemetry with a parent span and metrics
# MAGIC

# COMMAND ----------


import random
import uuid
import os

print("--- PIPELINE-SPAN (parent): ETL Pipeline")

# Unique ETL run identifier
etl_pipeline_id = str(uuid.uuid4())

# Initialize tracking values
etl_retry_attempts = random.randint(0, 2)
etl_status = "Success"
etl_error_count = 0
etl_error_details = {}

print(f"    ETL_PIPELINE_ID:  {etl_pipeline_id} ")

# Initialize OpenTelemetry with a parent span and metrics
etl_otel_helper = OpenTelemetryHelper(
    span_name="ETL_Pipeline",
    etl_pipeline_id=etl_pipeline_id,
    span_metrics={
        "DataExtraction": {"records_extracted": "counter", "http_request_duration_sec": "histogram"},
        "DataTransformation": {"records_transformed": "counter", "transformation_duration_sec": "histogram"},
        "DataLoading": {"records_written": "counter", "write_duration_sec": "histogram"},
        "ETL_Pipeline": {
            "etl_total_duration_sec": "histogram",
            "etl_efficiency_rate": "histogram", 
            "retry_attempts": "counter",
            "etl_error_count": "counter"
        }
    },
    span_attributes={"data_source": "External API", "etl_type": "Full Load"}
)

print(f"     OpenTelemetryHelper instance created for:  etl_otel_helper.etl_pipeline_id: {etl_otel_helper.etl_pipeline_id}")

# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL Helper Functions Setup
# MAGIC
# MAGIC Part of the original ETL, these functions simulate sleep time and error generation.
# MAGIC

# COMMAND ----------

import random



def sleep_func():
    start_time = time.time()
    sleep_time = random.uniform(0, 5)
    print("    ...Sleeping for ", sleep_time)
    time.sleep(sleep_time)
    end_time = time.time()
    duration_sec = end_time - start_time
    return round(duration_sec, 3)

def select_random_value(input_list):
    """Randomly selects and returns a value from the input list."""
    if not input_list:
        return None  # Handle case where the list is empty
    return random.choice(input_list)

def risky_function():
    """Randomly raises an exception 1 out of 5 times."""
    if random.randint(1, 5) == 1:  # 1 in 5 chance
        raise Exception("Oops! Random failure occurred.")
    return "Success! No exception this time."



# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL-STAGE #1: Data Extraction
# MAGIC
# MAGIC Call an external web API to retrieve and extract raw data.
# MAGIC
# MAGIC OpenTelemetry Span/Stage Name:  DataExtraction

# COMMAND ----------

# Start the DataExtraction span
print("Starting DataExtraction span...")
etl_otel_helper.start_tracing("DataExtraction", {
    "etl_pipeline_id": etl_pipeline_id
})

# COMMAND ----------

import time


# -------------------
# Data Extraction - DataExtraction
# -------------------
print("--- PIPELINE-SPAN (stage):  DataExtraction")

extraction_failed = False
err_msg = ""

# Simulate data extraction
duration_sec = sleep_func()

records_extracted = random.randint(10000, 20000)
response_size = records_extracted * 50  # Simulated response size

response_status = 200

# Simulate error handling with a 10% chance
if random.random() < 0.1:
    response_status = select_random_value([400, 401, 403, 404, 500])
    err_msg = f"API call failed with response_code={response_status}"
    extraction_failed = True
    print("    EXTRACTION FAILED!")

print(f"    STATUS: records_extracted: {records_extracted}, "
      f"http.status_code: {response_status}, "
      f"http.response_size: {response_size} bytes")



# COMMAND ----------

# Set span attributes for DataExtraction
print("Setting DataExtraction span attributes...")
etl_otel_helper.set_span_attribute("DataExtraction", "etl_pipeline_id", etl_pipeline_id)
etl_otel_helper.set_span_attribute("DataExtraction", "http_request_duration_sec", duration_sec)
etl_otel_helper.set_span_attribute("DataExtraction", "records_extracted", records_extracted)
etl_otel_helper.set_span_attribute("DataExtraction", "http_response_size", response_size)
etl_otel_helper.set_span_attribute("DataExtraction", "extraction_status", "Failed" if extraction_failed else "Success")
etl_otel_helper.set_span_attribute("DataExtraction", "http_status_code", response_status)
etl_otel_helper.set_span_attribute("DataExtraction", "error", str(extraction_failed).lower())
etl_otel_helper.set_span_attribute("DataExtraction", "error_message", err_msg if extraction_failed else "")

# Record metrics
etl_otel_helper.record_metric("DataExtraction", "records_extracted", records_extracted)
etl_otel_helper.record_metric("DataExtraction", "http_request_duration_sec", duration_sec)

# End the DataExtraction span
etl_otel_helper.end_tracing("DataExtraction")
print("DataExtraction span completed")

# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL-STAGE #2: Data Transformation
# MAGIC
# MAGIC Clean and transform raw data using PySpark.
# MAGIC
# MAGIC OpenTelemetry Span/Stage Name:  DataTransformation

# COMMAND ----------

# Start the DataTransformation span
print("Starting DataTransformation span...")
etl_otel_helper.start_tracing("DataTransformation", {
    "etl_pipeline_id": etl_pipeline_id
})

# COMMAND ----------

# -------------------
# Data Transformation
# -------------------

print("--- PIPELINE-SPAN (stage):  DataTransformation")

duration_sec = sleep_func()

records_input = records_extracted

if random.random() < 0.1:
    records_failed = random.randint(0, 500)
else:
    records_failed = 0

records_transformed = records_input - records_failed
transformation_duration_sec = duration_sec
records_transformed_rate = records_transformed / transformation_duration_sec if transformation_duration_sec > 0 else 0  

transformation_type = "Data Cleaning"
processing_engine = "PySpark"

if records_failed > 0:
    print("    TRANSFORM FAILED!")
    err_msg = f"{records_failed} records failed transformation"

print(f"    DataTransformation - records_input: {records_input}, "
      f"records_transformed: {records_transformed}, records_failed: {records_failed}, "
      f"transformation_duration_sec: {transformation_duration_sec}, "
      f"records_transformed_rate: {round(records_transformed_rate, 2)} records/sec, "
      f"transformation_type: {transformation_type}, processing_engine: {processing_engine}")



# COMMAND ----------

# Set span attributes for DataTransformation
print("Setting DataTransformation span attributes...")
etl_otel_helper.set_span_attribute("DataTransformation", "etl_pipeline_id", etl_pipeline_id)
etl_otel_helper.set_span_attribute("DataTransformation", "error", str(records_failed > 0).lower())
etl_otel_helper.set_span_attribute("DataTransformation", "records_input", records_input)
etl_otel_helper.set_span_attribute("DataTransformation", "records_transformed", records_transformed)
etl_otel_helper.set_span_attribute("DataTransformation", "records_failed", records_failed)
etl_otel_helper.set_span_attribute("DataTransformation", "transformation_duration_sec", transformation_duration_sec)
etl_otel_helper.set_span_attribute("DataTransformation", "records_transformed_rate", round(records_transformed_rate, 2))
etl_otel_helper.set_span_attribute("DataTransformation", "transformation_type", transformation_type)
etl_otel_helper.set_span_attribute("DataTransformation", "processing_engine", processing_engine)
etl_otel_helper.set_span_attribute("DataTransformation", "error_message", err_msg if records_failed > 0 and 'err_msg' in locals() else "")

# Record metrics
etl_otel_helper.record_metric("DataTransformation", "records_transformed", records_transformed)
etl_otel_helper.record_metric("DataTransformation", "transformation_duration_sec", transformation_duration_sec)

# End the DataTransformation span
etl_otel_helper.end_tracing("DataTransformation")
print("DataTransformation span completed")

# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL-STAGE #3: Data Load
# MAGIC
# MAGIC Write the transformed data to a storage location.
# MAGIC
# MAGIC OpenTelemetry Span/Stage Name:  DataLoading

# COMMAND ----------

# Define the data loading function with OpenTelemetry tracing
@etl_otel_helper.trace_function("DataLoading", {
    "etl_pipeline_id": "etl_pipeline_id",
    "records_written": "records_written",
    "storage_path": "storage_path",
    "storage_type": "storage_type",
    "file_size_bytes": "file_size_bytes",
    "compression_type": "compression_type",
    "partition_count": "partition_count",
    "batch_size": "batch_size",
    "write_duration_sec": "write_duration_sec",
    "records_write_rate": "records_write_rate",
    "storage_status": "storage_status",
    "error": "error",
    "error_message": "error_msg"
})
def perform_data_loading(records_to_write):
    """
    Perform data loading operation and return metrics and status information.
    
    :param records_to_write: Number of records to write to storage
    :return: Dictionary containing metrics and status information
    """
    # -------------------
    # Data Loading
    # -------------------
    print("--- PIPELINE-SPAN (stage):  DataLoading")

    error = False
    error_msg = ""
    storage_status = "Success"

    duration_sec = sleep_func()

    records_written = records_to_write
    storage_path = "/mnt/output/transformed_data"
    storage_type = "Azure Data Lake"
    file_size_bytes = records_written * 50
    compression_type = "Parquet"
    partition_count = random.randint(1, 10)
    batch_size = 1000
    write_duration_sec = duration_sec
    records_write_rate = records_written / write_duration_sec if write_duration_sec > 0 else 0  

    if random.random() < 0.1:
        print("    LOAD FAILED!")
        error = True
        error_msg = "Storage write failed: Permission denied"
        storage_status = "Failed"

    print(f"    DataLoading - records_written: {records_written}, "
          f"storage_path: {storage_path}, storage_type: {storage_type}, "
          f"write_duration_sec: {write_duration_sec}, records_write_rate: {round(records_write_rate, 2)} records/sec")
    
    # Return all metrics and status information
    return {
        "etl_pipeline_id": etl_pipeline_id,
        "error": str(error).lower(),
        "error_msg": error_msg,
        "records_written": records_written,
        "storage_path": storage_path,
        "storage_type": storage_type,
        "file_size_bytes": file_size_bytes,
        "compression_type": compression_type,
        "partition_count": partition_count,
        "batch_size": batch_size,
        "write_duration_sec": write_duration_sec,
        "records_write_rate": round(records_write_rate, 2),
        "storage_status": storage_status
    }

# Call the function and unpack the results
loading_results = perform_data_loading(records_transformed)

# Extract variables from the results for use in subsequent cells
error = loading_results["error"] == "true"
error_msg = loading_results["error_msg"]
records_written = loading_results["records_written"]
storage_path = loading_results["storage_path"]
storage_type = loading_results["storage_type"]
write_duration_sec = loading_results["write_duration_sec"]
records_write_rate = loading_results["records_write_rate"]



# COMMAND ----------

# MAGIC %md
# MAGIC # ETL Pipeline Simulation COMPLETE
# MAGIC
# MAGIC This is the final section of the original notebook. This would have stuff like cleanup and final reporting.
# MAGIC
# MAGIC OpenTelemetry Parent Notebook/Span name:  ETL_Pipeline
# MAGIC
# MAGIC NOTE: We can add final stats for the ETL_Pipeline here.
# MAGIC

# COMMAND ----------

# -------------------
# Finalize ETL Pipeline
# -------------------
print("--- ETL Pipeline COMPLETE")

etl_efficiency_rate = (records_written / records_extracted) * 100 if records_extracted > 0 else 0

# Record pipeline-level metrics
print(f"    ETL_Pipeline - ID: {etl_pipeline_id}, Input Records: {records_extracted}, "
      f"Output Records: {records_written}, Failed Records: {records_failed}")

# Log trace attributes for global values from previous cell stages
print("    Adding trace attributes to ETL_Pipeline span...")

# Record counts
etl_otel_helper.set_span_attribute("ETL_Pipeline", "records_extracted", records_extracted)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "records_transformed", records_transformed)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "records_written", records_written)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "records_failed", records_failed)

# Performance metrics
etl_otel_helper.set_span_attribute("ETL_Pipeline", "duration_transformation_sec", transformation_duration_sec)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "duration_write_sec", write_duration_sec)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "rate_transformation_records_per_sec", round(records_transformed_rate, 2))
etl_otel_helper.set_span_attribute("ETL_Pipeline", "rate_write_records_per_sec", round(records_write_rate, 2))
etl_otel_helper.set_span_attribute("ETL_Pipeline", "rate_etl_efficiency_percent", round(etl_efficiency_rate, 2))

# Status information
etl_otel_helper.set_span_attribute("ETL_Pipeline", "extraction_failed", extraction_failed)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "extraction_response_status", response_status)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "loading_error", error if 'error' in locals() else False)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "etl_status", etl_status)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "etl_retry_attempts", etl_retry_attempts)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "etl_error_count", etl_error_count)

# Other metadata
etl_otel_helper.set_span_attribute("ETL_Pipeline", "storage_path", storage_path)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "storage_type", storage_type)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "transformation_type", transformation_type)
etl_otel_helper.set_span_attribute("ETL_Pipeline", "processing_engine", processing_engine)

# Add error details if any
if 'error_msg' in locals() and error_msg:
    etl_otel_helper.set_span_attribute("ETL_Pipeline", "error_message", error_msg)
elif 'err_msg' in locals() and err_msg:
    etl_otel_helper.set_span_attribute("ETL_Pipeline", "error_message", err_msg)

print("    Trace attributes added successfully")

# Record pipeline-level metrics for Azure Monitor
print("    Recording pipeline-level metrics...")
etl_total_duration_sec = time.time() - etl_otel_helper.span_start_times["ETL_Pipeline"]
etl_otel_helper.record_metric("ETL_Pipeline", "etl_total_duration_sec", etl_total_duration_sec)
etl_otel_helper.record_metric("ETL_Pipeline", "etl_efficiency_rate", etl_efficiency_rate)
etl_otel_helper.record_metric("ETL_Pipeline", "retry_attempts", etl_retry_attempts)
etl_otel_helper.record_metric("ETL_Pipeline", "etl_error_count", etl_error_count)
print("    Pipeline-level metrics recorded successfully")

# Add a completion event to the ETL_Pipeline span
event_attributes = {
    "records_processed": records_written,
    "efficiency_rate": round(etl_efficiency_rate, 2),
    "extraction_status": "Failed" if extraction_failed else "Success",
    "loading_status": "Failed" if ('error' in locals() and error) else "Success"
}
etl_otel_helper.add_span_event("ETL_Pipeline", "ETL Pipeline Completed", event_attributes)
print("    Added completion event to ETL_Pipeline span")

# End the ETL_Pipeline trace
etl_otel_helper.end_tracing("ETL_Pipeline")
print("    ETL_Pipeline trace ended")


# COMMAND ----------

# MAGIC %md
# MAGIC
