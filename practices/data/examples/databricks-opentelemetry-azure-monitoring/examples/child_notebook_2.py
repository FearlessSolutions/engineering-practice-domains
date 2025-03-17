# Databricks notebook source
# MAGIC %md
# MAGIC # Child Notebook 2: Data Aggregation
# MAGIC 
# MAGIC This notebook simulates a data aggregation process. It does not contain any OpenTelemetry instrumentation itself.
# MAGIC The parent notebook will track this execution using OpenTelemetry.

# COMMAND ----------

import time
import random
import json
import psutil  # For memory usage simulation

# COMMAND ----------

# MAGIC %md
# MAGIC ## Simulate Data Aggregation Process

# COMMAND ----------

# Simulate processing time
start_time = time.time()

# Simulate data aggregation
num_aggregations = random.randint(10, 50)
print(f"Performing {num_aggregations} data aggregations...")

for i in range(num_aggregations):
    # Simulate processing delay for this aggregation
    agg_time = random.uniform(0.1, 0.5)
    time.sleep(agg_time)
    
    print(f"Completed aggregation {i+1}/{num_aggregations}")

# Simulate total processing time
execution_time = random.uniform(3, 7)
time.sleep(max(0, execution_time - (time.time() - start_time)))  # Ensure minimum execution time

# Calculate actual execution time
actual_execution_time = time.time() - start_time

# Simulate memory usage (in MB)
memory_usage = random.uniform(200, 500)

# Simulate status code
status_code = 200 if random.random() < 0.9 else random.choice([400, 500])

print(f"Aggregation complete with status code: {status_code}")
print(f"Performed {num_aggregations} aggregations in {actual_execution_time:.2f} seconds")
print(f"Simulated memory usage: {memory_usage:.2f} MB")

# COMMAND ----------

# MAGIC %md
# MAGIC ## Return Results to Parent Notebook

# COMMAND ----------

# Create result object to return to parent notebook
result = {
    "task": "data_aggregation",
    "status_code": status_code,
    "num_aggregations": num_aggregations,
    "execution_time_sec": round(actual_execution_time, 2),
    "memory_usage_mb": round(memory_usage, 2)
}

# Convert to JSON string to return to parent
result_json = json.dumps(result)

# Return the result to the parent notebook
dbutils.notebook.exit(result_json)
