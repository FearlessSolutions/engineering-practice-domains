# Setup Guide

This guide provides detailed instructions for setting up the Databricks OpenTelemetry Azure Monitoring integration.

## Quick Start: Direct File Upload (Recommended)

The simplest way to use this integration in Databricks is to upload the `otel_helper.py` file directly to your workspace:

1. **Download the helper file**:
   - Navigate to the `databricks-opentelemetry-azure-monitoring/databricks_opentelemetry_azure_monitoring` directory
   - Download the `otel_helper.py` file to your local machine

2. **Upload to Databricks workspace**:
   - In your Databricks workspace, click the "Create" button in the sidebar
   - Select "File" from the dropdown menu
   - Upload the `otel_helper.py` file to your workspace

3. **Import in your notebooks**:
   ```python
   from otel_helper import OpenTelemetryHelper
   ```

4. **Configure Azure Monitor connection string** (see Azure Application Insights Setup section below)

This approach is recommended for most users as it's simpler and doesn't require package installation.

## Azure Application Insights Setup

1. **Create an Azure Application Insights resource**:
   - Go to the Azure Portal
   - Create a new Application Insights resource
   - Note the Instrumentation Key and Connection String

2. **Configure Databricks Secret Scope**:

   **Option A: Using Databricks CLI**
   - Create a secret scope in your Databricks workspace
   - Add your Application Insights connection string as a secret:
     ```
     databricks secrets create-scope --scope azure-monitor
     databricks secrets put --scope azure-monitor --key connection-string
     ```

   **Option B: Using Databricks Cluster Environment Variables**
   - Navigate to your Databricks workspace
   - Go to the "Compute" section in the sidebar
   - Select your cluster (or create a new one)
   - Click on "Edit" to modify the cluster configuration
   - Expand the "Advanced options" section
   - Select the "Environment variables" tab
   - Add a new environment variable:
     - Key: `AZURE_MONITOR_CONNECTION_STRING`
     - Value: Paste your Application Insights connection string directly here
   - Click "Save" to update your cluster configuration
   - Restart your cluster for the changes to take effect

## Databricks Cluster Setup

1. **Install Required Libraries**:
   - Create a new Databricks cluster or use an existing one
   - Install the following libraries:
     - `azure-monitor-opentelemetry`
     - `opentelemetry-sdk`
     - `azure-core`
     - `opentelemetry-api`
     - `opentelemetry-instrumentation`

2. **Configure Cluster Environment Variables**:
   
   **Option A: Using Secret References**
   - Add the following environment variable to your cluster configuration:
     ```
     AZURE_MONITOR_CONNECTION_STRING = {{secrets/azure-monitor/connection-string}}
     ```
   
   **Option B: Using Direct Environment Variable**
   - If you've set up the connection string directly as an environment variable (as described in Option B above), 
     this step is already complete.

## Project Setup Options

There are multiple ways to set up and use this package:

### Option 1: Direct File Upload (Recommended)

As described in the Quick Start section above, this is the simplest approach:

1. **Upload the helper file**:
   - Upload the `otel_helper.py` file to your Databricks workspace using the "Create" > "File" option

2. **Import in your notebooks**:
   ```python
   from otel_helper import OpenTelemetryHelper
   ```

### Option 2: Install as a Python package

For more advanced use cases or when you need to modify the package:

1. **Access the Project**:
   - Navigate to the `databricks-opentelemetry-azure-monitoring` directory within the parent repository

2. **Install the Package**:
   ```bash
   pip install -e databricks-opentelemetry-azure-monitoring
   ```

3. **Import the Helper Module**:
   - Import it in your notebooks:
     ```python
     from databricks_opentelemetry_azure_monitoring import OpenTelemetryHelper
     ```

### Option 3: Use the files with manual dependency installation

If you need more control over dependencies:

1. **Access the Project**:
   - Navigate to the `databricks-opentelemetry-azure-monitoring` directory within the parent repository

2. **Install Dependencies**:
   ```bash
   pip install -r databricks-opentelemetry-azure-monitoring/requirements.txt
   ```

3. **Import the Helper Module**:
   - Upload the `databricks-opentelemetry-azure-monitoring/databricks_opentelemetry_azure_monitoring/otel_helper.py` file to your Databricks workspace
   - Import it in your notebooks:
     ```python
     from otel_helper import OpenTelemetryHelper
     ```

## Verification

To verify your setup is working correctly:

1. Run the example notebook `etl_simulation_with_otel.py` in your Databricks workspace
2. Check your Azure Application Insights resource for telemetry data
3. Verify that spans, metrics, and events are being properly recorded

The example notebook is already configured to use the direct import approach (`from otel_helper import OpenTelemetryHelper`), so it should work seamlessly with your uploaded file.

## Troubleshooting

- **No telemetry data in Application Insights**:
  - Verify your connection string is correct
  - Check that the OpenTelemetry exporters are properly initialized
  - Ensure your cluster has internet access to send data to Azure

- **Missing spans or metrics**:
  - Verify that spans are properly started and ended
  - Check for any exceptions in the notebook execution

- **Performance issues**:
  - Consider using batch processing for span exports
  - Adjust the metric export interval based on your needs
