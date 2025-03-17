# Databricks notebook source
# MAGIC %md
# MAGIC # ETL Pipeline Simulation BEFORE OpenTelemetry 
# MAGIC
# MAGIC This notebook simulates an ETL pipeline. It will be later used to demonstrate both trace and metric instrumentation using OpenTelemetry and Azure Monitor.
# MAGIC
# MAGIC The instrumentation will be added to this notebook in a way that is modular and has the least code changes to the original notebook cells.
# MAGIC
# MAGIC FUTURE OpenTelemetry Parent Notebook/Span name:  ETL_Pipeline
# MAGIC

# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL Helper Functions Setup
# MAGIC
# MAGIC Part of the original ETL, these functions simulate sleep time and error generation.
# MAGIC

# COMMAND ----------

import random
import uuid


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

print("--- PIPELINE-SPAN (parent): ETL Pipeline")

# Unique ETL run identifier
etl_pipeline_id = str(uuid.uuid4())

# Initialize tracking values
etl_retry_attempts = random.randint(0, 2)
etl_status = "Success"
etl_error_count = 0
etl_error_details = {}

print(f"    ETL_PIPELINE_ID:  {etl_pipeline_id} ")



# COMMAND ----------

# MAGIC %md
# MAGIC ## ETL-STAGE #1: Data Extraction
# MAGIC
# MAGIC Call an external web API to retrieve and extract raw data.
# MAGIC
# MAGIC FUTURE OpenTelemetry Span/Stage Name:  DataExtraction

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

# MAGIC %md
# MAGIC ## ETL-STAGE #2: Data Transformation
# MAGIC
# MAGIC Clean and transform raw data using PySpark.
# MAGIC
# MAGIC FUTURE OpenTelemetry Span/Stage Name:  DataTransformation

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

# MAGIC %md
# MAGIC ## ETL-STAGE #3: Data Load
# MAGIC
# MAGIC Write the transformed data to a storage location.
# MAGIC
# MAGIC FUTURE OpenTelemetry Span/Stage Name:  DataLoading

# COMMAND ----------

# -------------------
# Data Loading
# -------------------
print("--- PIPELINE-SPAN (stage):  DataLoading")

error = False
error_msg = ""

duration_sec = sleep_func()

records_written = records_transformed
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
    storage_msg = "Data loading encountered an error"

print(f"    DataLoading - records_written: {records_written}, "
      f"storage_path: {storage_path}, storage_type: {storage_type}, "
      f"write_duration_sec: {write_duration_sec}, records_write_rate: {round(records_write_rate, 2)} records/sec")




# COMMAND ----------

# MAGIC %md
# MAGIC # ETL Pipeline Simulation COMPLETE
# MAGIC
# MAGIC This is the final section of the original notebook. This would have stuff like cleanup and final reporting.
# MAGIC
# MAGIC FUTURE OpenTelemetry Parent Notebook/Span name:  ETL_Pipeline
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
      f"Output Records: {records_written}, Failed Records: {records_failed}, ")


# COMMAND ----------

# MAGIC %md
# MAGIC