# Databricks notebook source
# MAGIC %md
# MAGIC # Child Notebook 1: Data Validation
# MAGIC 
# MAGIC This notebook simulates a data validation process. It does not contain any OpenTelemetry instrumentation itself.
# MAGIC The parent notebook will track this execution using OpenTelemetry.

# COMMAND ----------

import time
import random
import json

# COMMAND ----------

# MAGIC %md
# MAGIC ## Simulate Data Validation Process

# COMMAND ----------

# Simulate processing time
start_time = time.time()

# Simulate data validation
total_records = random.randint(5000, 10000)
print(f"Validating {total_records} records...")

# Simulate processing delay
processing_time = random.uniform(2, 5)
time.sleep(processing_time)

# Simulate validation results
validation_errors = random.randint(0, int(total_records * 0.05))  # Up to 5% error rate
validation_status = "Success" if validation_errors < 100 else "Warning" if validation_errors < 500 else "Failed"

# Calculate processing rate
records_per_second = total_records / processing_time

print(f"Validation complete with status: {validation_status}")
print(f"Processed {total_records} records in {processing_time:.2f} seconds ({records_per_second:.2f} records/sec)")
print(f"Found {validation_errors} validation errors")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Return Results to Parent Notebook

# COMMAND ----------

# Create result object to return to parent notebook
result = {
    "task": "data_validation",
    "status": validation_status,
    "total_records": total_records,
    "validation_errors": validation_errors,
    "processing_time_sec": round(processing_time, 2),
    "records_per_second": round(records_per_second, 2),
    "timestamp": time.time()
}

# Convert to JSON string to return to parent
result_json = json.dumps(result)

# Return the result to the parent notebook
dbutils.notebook.exit(result_json)
